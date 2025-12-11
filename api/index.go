package api

import (
	"net/http"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/handlers"
	"texnousta-backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router *gin.Engine

func init() {
	// Загрузка переменных окружения для Vercel
	godotenv.Load()
	
	// Инициализация базы данных
	database.Init()
	
	// Настройка Gin в production режиме для Vercel
	gin.SetMode(gin.ReleaseMode)
	
	// Создание роутера
	router = gin.New()
	
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS конфигурация для Vercel
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// API маршруты
	api := router.Group("/api/v1")
	{
		// Публичные маршруты
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.GET("/categories", handlers.GetCategories)
		
		// Контактные формы
		api.POST("/contact", handlers.CreateContact)
		api.POST("/quick-contact", handlers.CreateQuickContact)
		api.POST("/phone-contact", handlers.CreatePhoneContact)
		
		// Аналитика
		api.POST("/track-visitor", handlers.TrackVisitor)
		api.POST("/track-phone-click", handlers.TrackPhoneClick)
		
		// Авторизованные маршруты
		auth := api.Group("/")
		auth.Use(middleware.AuthRequired())
		{
			auth.GET("/profile", handlers.GetProfile)
			auth.PUT("/profile", handlers.UpdateProfile)
		}

		// Админ маршруты
		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
		{
			// Контакты
			admin.GET("/contacts", handlers.GetContacts)
			admin.GET("/contacts/:id", handlers.GetContact)
			admin.PUT("/contacts/:id/read", handlers.MarkContactAsRead)
			admin.DELETE("/contacts/:id", handlers.DeleteContact)
			
			// Аналитика
			admin.GET("/visitor-stats", handlers.GetVisitorStats)
			admin.GET("/phone-click-stats", handlers.GetPhoneClickStats)
			admin.GET("/phone-contacts", handlers.GetPhoneContacts)
			admin.DELETE("/phone-contacts/:id", handlers.DeletePhoneContact)
		}
	}
}

// Handler для Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}