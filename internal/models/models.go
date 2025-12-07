package models

import (
	"time"
)

// User - модель пользователя
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Email     string    `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Role      string    `json:"role" gorm:"size:20;default:'user'"` // user, admin
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category - модель категории товаров
type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Image       string    `json:"image" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Связи
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// Product - модель товара
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:200;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"not null"`
	OldPrice    float64   `json:"old_price"`
	Image       string    `json:"image" gorm:"size:255"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	Brand       string    `json:"brand" gorm:"size:100"`
	Model       string    `json:"model" gorm:"size:100"`
	Stock       int       `json:"stock" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	IsFeatured  bool      `json:"is_featured" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Связи
	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// Order - модель заказа
type Order struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	UserID     uint        `json:"user_id" gorm:"not null"`
	Total      float64     `json:"total" gorm:"not null"`
	Status     string      `json:"status" gorm:"size:50;default:'pending'"` // pending, confirmed, shipped, delivered, cancelled
	PaymentStatus string   `json:"payment_status" gorm:"size:50;default:'pending'"` // pending, paid, failed
	ShippingAddress string `json:"shipping_address" gorm:"type:text;not null"`
	Phone      string      `json:"phone" gorm:"size:20;not null"`
	Notes      string      `json:"notes" gorm:"type:text"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	
	// Связи
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderItem - модель позиции заказа
type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	
	// Связи
	Order   Order   `json:"-" gorm:"foreignKey:OrderID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// ContactForm - модель формы обратной связи
type ContactForm struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Email     string    `json:"email" gorm:"size:100"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Subject   string    `json:"subject" gorm:"size:200;not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest - структура для запроса входа
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest - структура для запроса регистрации
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

// ProductRequest - структура для создания/обновления товара
type ProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	OldPrice    float64  `json:"old_price"`
	CategoryID  uint     `json:"category_id" binding:"required"`
	Brand       string   `json:"brand"`
	Model       string   `json:"model"`
	Stock       int      `json:"stock"`
	IsActive    bool     `json:"is_active"`
	IsFeatured  bool     `json:"is_featured"`
}

// CategoryRequest - структура для создания/обновления категории
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}