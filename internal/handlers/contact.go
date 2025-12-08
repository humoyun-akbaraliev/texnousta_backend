package handlers

import (
	"net/http"
	"strconv"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// ContactRequest - структура для запроса контактной формы
type ContactRequest struct {
	Name    string `json:"name" binding:"required,min=2"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required,min=10"`
}

// CreateContact создает новое обращение через контактную форму
//
//	@Summary		Отправить контактную форму
//	@Description	Отправка полной контактной формы с темой и сообщением
//	@Tags			contact
//	@Accept			json
//	@Produce		json
//	@Param			contact	body		ContactRequest	true	"Данные обращения"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/contact [post]
func CreateContact(c *gin.Context) {
	var req ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание контактного обращения
	contact := models.ContactForm{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Subject: req.Subject,
		Message: req.Message,
		IsRead:  false,
	}

	if err := database.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении обращения"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ваше обращение успешно отправлено. Мы свяжемся с вами в ближайшее время.",
		"id":      contact.ID,
	})
}

// GetContacts получает список контактных обращений (только для админов)
//
//	@Summary		Получить список обращений
//	@Description	Получение списка контактных обращений с фильтрацией (только для администраторов)
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int		false	"Номер страницы"		default(1)
//	@Param			limit	query		int		false	"Количество на странице"	default(20)
//	@Param			unread	query		bool	false	"Только непрочитанные"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		403		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/admin/contacts [get]
func GetContacts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	onlyUnread := c.Query("unread") == "true"

	offset := (page - 1) * limit

	query := database.DB.Model(&models.ContactForm{})
	
	// Фильтр только непрочитанных
	if onlyUnread {
		query = query.Where("is_read = ?", false)
	}

	// Подсчет общего количества
	var total int64
	query.Count(&total)

	// Получение обращений с пагинацией
	var contacts []models.ContactForm
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении обращений"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contacts": contacts,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetContact получает конкретное обращение по ID (только для админов)
func GetContact(c *gin.Context) {
	id := c.Param("id")

	var contact models.ContactForm
	if err := database.DB.First(&contact, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Обращение не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contact": contact})
}

// MarkContactAsRead помечает обращение как прочитанное (только для админов)
func MarkContactAsRead(c *gin.Context) {
	id := c.Param("id")

	var contact models.ContactForm
	if err := database.DB.First(&contact, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Обращение не найдено"})
		return
	}

	if err := database.DB.Model(&contact).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении статуса"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Обращение помечено как прочитанное"})
}

// DeleteContact удаляет обращение (только для админов)
func DeleteContact(c *gin.Context) {
	id := c.Param("id")

	var contact models.ContactForm
	if err := database.DB.First(&contact, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Обращение не найдено"})
		return
	}

	if err := database.DB.Delete(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении обращения"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Обращение успешно удалено"})
}

// QuickContactRequest - упрощенная структура для быстрого контакта (только телефон)
type QuickContactRequest struct {
	Name  string `json:"name" binding:"required,min=2"`
	Phone string `json:"phone" binding:"required"`
}

// PhoneContactRequest - структура только для номера телефона
type PhoneContactRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// CreateQuickContact создает быстрое обращение (только имя и телефон)
//
//	@Summary		Быстрая заявка
//	@Description	Отправка быстрой заявки с именем и телефоном для обратного звонка
//	@Tags			contact
//	@Accept			json
//	@Produce		json
//	@Param			contact	body		QuickContactRequest	true	"Имя и телефон"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/quick-contact [post]
func CreateQuickContact(c *gin.Context) {
	var req QuickContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание быстрого контактного обращения
	contact := models.ContactForm{
		Name:    req.Name,
		Phone:   req.Phone,
		Subject: "Быстрая заявка",
		Message: "Клиент оставил заявку на обратный звонок",
		IsRead:  false,
	}

	if err := database.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении заявки"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Заявка принята! Мы перезвоним вам в течение 15 минут.",
		"id":      contact.ID,
	})
}

// CreatePhoneContact создает обращение только с номером телефона
//
//	@Summary		Оставить телефон
//	@Description	Сохранение только номера телефона для обратного звонка
//	@Tags			contact
//	@Accept			json
//	@Produce		json
//	@Param			phone	body		PhoneContactRequest	true	"Номер телефона"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/phone-contact [post]
func CreatePhoneContact(c *gin.Context) {
	var req PhoneContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание контактного обращения только с телефоном
	contact := models.ContactForm{
		Name:    "Не указано",
		Phone:   req.Phone,
		Subject: "Оставлен телефон",
		Message: "Клиент оставил только номер телефона для связи",
		IsRead:  false,
	}

	if err := database.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении номера"})
		return
	}

	// Также сохраняем в отдельную таблицу для аналитики
	phoneContact := models.PhoneContact{
		Phone: req.Phone,
	}
	database.DB.Create(&phoneContact)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Номер телефона сохранен! Мы свяжемся с вами.",
		"id":      contact.ID,
	})
}