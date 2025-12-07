package handlers

import (
	"net/http"
	"strconv"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// GetProducts получает список товаров с фильтрацией и пагинацией
//
//	@Summary		Получить список товаров
//	@Description	Получение товаров с возможностью фильтрации по категории, поиску и пагинацией
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Номер страницы"		default(1)
//	@Param			limit		query		int		false	"Количество на странице"	default(12)
//	@Param			category	query		int		false	"ID категории"
//	@Param			search		query		string	false	"Поиск по названию"
//	@Param			featured	query		bool	false	"Только рекомендуемые"
//	@Param			sort		query		string	false	"Сортировка"			default(created_at)
//	@Param			order		query		string	false	"Порядок сортировки"		default(desc)
//	@Success		200			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/products [get]
func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	category := c.Query("category")
	search := c.Query("search")
	featured := c.Query("featured")
	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	offset := (page - 1) * limit

	query := database.DB.Model(&models.Product{}).Where("is_active = ?", true)

	// Применение фильтров
	if category != "" {
		query = query.Where("category_id = ?", category)
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if featured == "true" {
		query = query.Where("is_featured = ?", true)
	}

	// Подсчет общего количества
	var total int64
	query.Count(&total)

	// Получение товаров с пагинацией
	var products []models.Product
	if err := query.Preload("Category").
		Order(sortBy + " " + order).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении товаров"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetProduct получает товар по ID
//
//	@Summary		Получить товар по ID
//	@Description	Получение информации о конкретном товаре
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID товара"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/products/{id} [get]
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	
	var product models.Product
	if err := database.DB.Preload("Category").
		Where("id = ? AND is_active = ?", id, true).
		First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// CreateProduct создает новый товар (только для админов)
//
//	@Summary		Создать товар
//	@Description	Создание нового товара (только для администраторов)
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			product	body		models.ProductRequest	true	"Данные товара"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		403		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/admin/products [post]
func CreateProduct(c *gin.Context) {
	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования категории
	var category models.Category
	if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Категория не найдена"})
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		OldPrice:    req.OldPrice,
		CategoryID:  req.CategoryID,
		Brand:       req.Brand,
		Model:       req.Model,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
		IsFeatured:  req.IsFeatured,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании товара"})
		return
	}

	// Загрузка связанной категории
	database.DB.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Товар успешно создан",
		"product": product,
	})
}

// UpdateProduct обновляет товар (только для админов)
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования категории
	var category models.Category
	if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Категория не найдена"})
		return
	}

	// Обновление полей
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
		"old_price":   req.OldPrice,
		"category_id": req.CategoryID,
		"brand":       req.Brand,
		"model":       req.Model,
		"stock":       req.Stock,
		"is_active":   req.IsActive,
		"is_featured": req.IsFeatured,
	}

	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении товара"})
		return
	}

	// Загрузка обновленного товара с категорией
	database.DB.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Товар успешно обновлен",
		"product": product,
	})
}

// DeleteProduct удаляет товар (только для админов)
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении товара"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар успешно удален"})
}