package handlers

import (
	"net/http"
	"strconv"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// GetCategories получает список категорий
//
//	@Summary		Получить список категорий
//	@Description	Получение всех активных категорий товаров
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	
	if err := database.DB.Where("is_active = ?", true).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении категорий"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// CreateCategory создает новую категорию (только для админов)
//
//	@Summary		Создать категорию
//	@Description	Создание новой категории товаров (только для администраторов)
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			category	body		models.CategoryRequest	true	"Данные категории"
//	@Success		201			{object}	map[string]interface{}
//	@Failure		400			{object}	map[string]interface{}
//	@Failure		401			{object}	map[string]interface{}
//	@Failure		403			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/admin/categories [post]
func CreateCategory(c *gin.Context) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании категории"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Категория успешно создана",
		"category": category,
	})
}

// UpdateCategory обновляет категорию (только для админов)
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Категория не найдена"})
		return
	}

	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"is_active":   req.IsActive,
	}

	if err := database.DB.Model(&category).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении категории"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Категория успешно обновлена",
		"category": category,
	})
}

// DeleteCategory удаляет категорию (только для админов)
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Категория не найдена"})
		return
	}

	// Проверяем, есть ли товары в этой категории
	var productCount int64
	database.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Невозможно удалить категорию, в ней есть товары",
		})
		return
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении категории"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Категория успешно удалена"})
}

// GetUsers получает список пользователей (только для админов)
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	// Подсчет общего количества
	database.DB.Model(&models.User{}).Count(&total)

	// Получение пользователей с пагинацией
	if err := database.DB.Select("id, name, email, phone, role, is_active, created_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateUser обновляет пользователя (только для админов)
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	updates["is_active"] = req.IsActive

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь успешно обновлен",
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"phone":     user.Phone,
			"role":      user.Role,
			"is_active": user.IsActive,
		},
	})
}

// DeleteUser удаляет пользователя (только для админов)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Нельзя удалить себя
	currentUser, _ := c.Get("user")
	currentUserModel := currentUser.(models.User)
	if currentUserModel.ID == user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя удалить собственный аккаунт"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}