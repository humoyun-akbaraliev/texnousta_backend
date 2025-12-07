package handlers

import (
	"net/http"
	"os"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register регистрирует нового пользователя
//
//	@Summary		Регистрация пользователя
//	@Description	Создание нового аккаунта пользователя
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.RegisterRequest	true	"Данные пользователя"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		409		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/register [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования пользователя
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	// Создание пользователя
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     "user",
		IsActive: true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь успешно зарегистрирован",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  user.Role,
		},
		"token": tokenString,
	})
}

// Login авторизует пользователя
//
//	@Summary		Авторизация пользователя
//	@Description	Вход в систему с получением JWT токена
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		models.LoginRequest	true	"Учетные данные"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{object}	map[string]interface{}
//	@Failure		401			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	// Проверка активности
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Аккаунт заблокирован"})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Успешная авторизация",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  user.Role,
		},
		"token": tokenString,
	})
}

// GetProfile получает профиль пользователя
//
//	@Summary		Получить профиль пользователя
//	@Description	Получение информации о текущем авторизованном пользователе
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/profile [get]
func GetProfile(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         userModel.ID,
			"name":       userModel.Name,
			"email":      userModel.Email,
			"phone":      userModel.Phone,
			"role":       userModel.Role,
			"created_at": userModel.CreatedAt,
		},
	})
}

// UpdateProfile обновляет профиль пользователя
//
//	@Summary		Обновить профиль пользователя
//	@Description	Обновление имени и телефона пользователя
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			profile	body		object	true	"Данные для обновления"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/profile [put]
func UpdateProfile(c *gin.Context) {
	user, _ := c.Get("user")
	userModel := user.(models.User)

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновление данных
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if err := database.DB.Model(&userModel).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении профиля"})
		return
	}

	// Получение обновленных данных
	var updatedUser models.User
	database.DB.First(&updatedUser, userModel.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Профиль успешно обновлен",
		"user": gin.H{
			"id":    updatedUser.ID,
			"name":  updatedUser.Name,
			"email": updatedUser.Email,
			"phone": updatedUser.Phone,
			"role":  updatedUser.Role,
		},
	})
}