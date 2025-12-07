package database

import (
	"log"
	"os"
	"texnousta-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init инициализирует подключение к базе данных
func Init() {
	var err error
	
	// Выбор базы данных в зависимости от окружения
	databaseURL := os.Getenv("DATABASE_URL")
	
	if databaseURL != "" {
		// Продакшен - используем PostgreSQL
		DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatal("Ошибка подключения к PostgreSQL:", err)
		}
		log.Println("Успешное подключение к базе данных PostgreSQL")
	} else {
		// Локальная разработка - используем SQLite
		DB, err = gorm.Open(sqlite.Open("texnousta.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("Ошибка подключения к SQLite:", err)
		}
		log.Println("Успешное подключение к базе данных SQLite")
	}

	// Автомиграция
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.ContactForm{},
	)

	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	log.Println("Миграция базы данных завершена")
	
	// Создание тестовых данных
	createSeedData()
}

// createSeedData создает начальные данные для тестирования
func createSeedData() {
	// Проверяем, есть ли уже данные
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	
	if userCount > 0 {
		return // Данные уже существуют
	}

	log.Println("Создание начальных данных...")

	// Создание админа
	admin := models.User{
		Name:     "Администратор",
		Email:    "admin@texnousta.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Role:     "admin",
		Phone:    "+998901234567",
		IsActive: true,
	}
	DB.Create(&admin)

	// Создание категорий
	categories := []models.Category{
		{Name: "Смартфоны", Description: "Мобильные телефоны и смартфоны", IsActive: true},
		{Name: "Ноутбуки", Description: "Портативные компьютеры", IsActive: true},
		{Name: "Телевизоры", Description: "ЖК и OLED телевизоры", IsActive: true},
		{Name: "Бытовая техника", Description: "Техника для дома", IsActive: true},
		{Name: "Аксессуары", Description: "Различные аксессуары", IsActive: true},
	}

	for _, category := range categories {
		DB.Create(&category)
	}

	// Создание товаров
	products := []models.Product{
		{
			Name:        "iPhone 15 Pro",
			Description: "Новейший флагманский смартфон от Apple",
			Price:       1200,
			OldPrice:    1300,
			CategoryID:  1,
			Brand:       "Apple",
			Model:       "iPhone 15 Pro",
			Stock:       50,
			IsActive:    true,
			IsFeatured:  true,
		},
		{
			Name:        "Samsung Galaxy S24",
			Description: "Флагманский Android смартфон",
			Price:       1000,
			CategoryID:  1,
			Brand:       "Samsung",
			Model:       "Galaxy S24",
			Stock:       30,
			IsActive:    true,
			IsFeatured:  true,
		},
		{
			Name:        "MacBook Pro 16\"",
			Description: "Профессиональный ноутбук для работы",
			Price:       2500,
			CategoryID:  2,
			Brand:       "Apple",
			Model:       "MacBook Pro 16",
			Stock:       20,
			IsActive:    true,
			IsFeatured:  true,
		},
		{
			Name:        "LG OLED TV 55\"",
			Description: "Большой OLED телевизор высокого качества",
			Price:       1800,
			OldPrice:    2000,
			CategoryID:  3,
			Brand:       "LG",
			Model:       "OLED55C3",
			Stock:       15,
			IsActive:    true,
		},
	}

	for _, product := range products {
		DB.Create(&product)
	}

	log.Println("Начальные данные созданы")
}