package middleware

import (
	"net/http"
	"os"
	"strings"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет JWT токен
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}

		// Извлечение токена из заголовка "Bearer token"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		// Парсинг и валидация токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			c.Abort()
			return
		}

		// Извлечение данных пользователя из токена
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := uint(claims["user_id"].(float64))
			
			// Проверка существования пользователя в базе данных
			var user models.User
			if err := database.DB.First(&user, userID).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
				c.Abort()
				return
			}

			// Проверка активности пользователя
			if !user.IsActive {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Аккаунт заблокирован"})
				c.Abort()
				return
			}

			// Сохранение пользователя в контексте
			c.Set("user", user)
			c.Set("user_id", userID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware проверяет права администратора
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
			c.Abort()
			return
		}

		userModel := user.(models.User)
		if userModel.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав доступа"})
			c.Abort()
			return
		}

		c.Next()
	}
}