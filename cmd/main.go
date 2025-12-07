// Package main TexnoUsta Backend API
//
//	@title						TexnoUsta API
//	@version					1.0
//	@description				API для интернет-магазина техники TexnoUsta
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.email				admin@texnousta.com
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"log"
	"os"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/handlers"
	"texnousta-backend/internal/middleware"

	_ "texnousta-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются системные переменные")
	}

	// Инициализация базы данных
	database.Init()

	// Настройка Gin режима
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Создание роутера
	r := gin.Default()

	// Настройка CORS
	config := cors.DefaultConfig()
	
	// Настройка CORS для продакшена и разработки
	if gin.Mode() == gin.ReleaseMode {
		// Продакшен - разрешить только определенные домены
		corsOrigins := os.Getenv("CORS_ORIGINS")
		if corsOrigins != "" {
			config.AllowOrigins = []string{corsOrigins}
		} else {
			config.AllowOrigins = []string{"https://your-frontend-domain.com"}
		}
	} else {
		// Разработка - разрешить localhost
		config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000"}
	}
	
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Middleware для логирования
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Статические файлы
	r.Static("/uploads", "./uploads")

	// API роуты
	api := r.Group("/api/v1")
	{
		// Публичные роуты
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		
		// Роуты для товаров (публичные)
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.GET("/categories", handlers.GetCategories)
		
		// Контактная форма (публичная)
		api.POST("/contact", handlers.CreateContact)
		api.POST("/quick-contact", handlers.CreateQuickContact)
		api.POST("/phone-contact", handlers.CreatePhoneContact)
		
		// Защищенные роуты
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Пользовательские роуты
			protected.GET("/profile", handlers.GetProfile)
			protected.PUT("/profile", handlers.UpdateProfile)
			
			// Админские роуты
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// Управление товарами
				admin.POST("/products", handlers.CreateProduct)
				admin.PUT("/products/:id", handlers.UpdateProduct)
				admin.DELETE("/products/:id", handlers.DeleteProduct)
				
				// Управление категориями
				admin.POST("/categories", handlers.CreateCategory)
				admin.PUT("/categories/:id", handlers.UpdateCategory)
				admin.DELETE("/categories/:id", handlers.DeleteCategory)
				
				// Управление пользователями
				admin.GET("/users", handlers.GetUsers)
				admin.PUT("/users/:id", handlers.UpdateUser)
				admin.DELETE("/users/:id", handlers.DeleteUser)
				
				// Управление контактными обращениями
				admin.GET("/contacts", handlers.GetContacts)
				admin.GET("/contacts/:id", handlers.GetContact)
				admin.PUT("/contacts/:id/read", handlers.MarkContactAsRead)
				admin.DELETE("/contacts/:id", handlers.DeleteContact)
			}
		}
	}

	// Swagger документация
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Сервер запущен на порту %s", port)
	r.Run(":" + port)
}