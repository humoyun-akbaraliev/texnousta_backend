package handlers

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"texnousta-backend/internal/database"
	"texnousta-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// AdminLogin входе в админ панель
// @Summary Вход в админ панель
// @Description Авторизация для доступа к админ панели аналитики
// @Tags Admin
// @Accept json
// @Produce json
// @Param login body models.AdminLoginRequest true "Данные для входа"
// @Success 200 {object} models.AdminLoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/admin/login [post]
func AdminLogin(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Простая проверка логина и пароля (в реальном проекте используйте хеширование)
	if req.Username == "admin" && req.Password == "admin123" {
		// Генерируем простой токен (в реальном проекте используйте JWT)
		token := generateSimpleToken(req.Username)
		
		c.JSON(http.StatusOK, models.AdminLoginResponse{
			Token:   token,
			Message: "Успешный вход в админ панель",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
	}
}

// generateSimpleToken генерирует простой токен
func generateSimpleToken(username string) string {
	data := fmt.Sprintf("%s:%d", username, time.Now().Unix())
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// isValidAdminToken проверяет валидность токена админа
func isValidAdminToken(c *gin.Context) bool {
	token := c.GetHeader("Authorization")
	if token == "" {
		return false
	}
	
	// Убираем префикс "Bearer " если есть
	token = strings.TrimPrefix(token, "Bearer ")
	
	// В реальном проекте здесь должна быть полная проверка токена
	// Пока просто проверяем, что токен не пустой и имеет нужную длину
	return len(token) == 32 // MD5 hash length
}

// TrackVisitor регистрирует посещение сайта
// @Summary Отслеживание посетителей
// @Description Регистрирует посещение сайта по IP адресу
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /track-visitor [post]
func TrackVisitor(c *gin.Context) {
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	
	now := time.Now()
	date := now.Format("2006-01-02")
	month := now.Format("2006-01")
	
	// Проверяем, был ли уже такой посетитель сегодня
	var existing models.VisitorStat
	result := database.DB.Where("ip_address = ? AND date = ?", clientIP, date).First(&existing)
	
	if result.Error != nil {
		// Новый посетитель за сегодня - создаем запись
		visitor := models.VisitorStat{
			IPAddress: clientIP,
			UserAgent: userAgent,
			Date:      date,
			Month:     month,
		}
		
		if err := database.DB.Create(&visitor).Error; err != nil {
			log.Printf("❌ Ошибка сохранения посетителя: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения статистики"})
			return
		}
		log.Printf("✅ Посетитель зарегистрирован: IP=%s, дата=%s", clientIP, date)
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Посещение зарегистрировано"})
}

// GetVisitorStats возвращает статистику посетителей
// @Summary Получить статистику посетителей
// @Description Возвращает статистику посетителей по дням и месяцам
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Количество дней для отображения (по умолчанию 30)"
// @Success 200 {object} models.VisitorStatsResponse
// @Router /admin/analytics/visitors [get]
func GetVisitorStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}
	
	// Получаем дату начала периода
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	
	// Дневная статистика
	var dailyStats []models.DailyStat
	database.DB.Raw(`
		SELECT 
			date,
			COUNT(DISTINCT ip_address) as unique_views,
			COUNT(*) as total_views
		FROM visitor_stats 
		WHERE date >= ? 
		GROUP BY date 
		ORDER BY date DESC
	`, startDate).Scan(&dailyStats)
	
	// Месячная статистика за последние 12 месяцев
	var monthlyStats []models.MonthlyStat
	startMonth := time.Now().AddDate(0, -12, 0).Format("2006-01")
	
	database.DB.Raw(`
		SELECT 
			month,
			COUNT(DISTINCT ip_address) as unique_views,
			COUNT(*) as total_views
		FROM visitor_stats 
		WHERE month >= ? 
		GROUP BY month 
		ORDER BY month DESC
	`, startMonth).Scan(&monthlyStats)
	
	// Общее количество уникальных посетителей
	var totalUnique int64
	database.DB.Model(&models.VisitorStat{}).
		Distinct("ip_address").
		Count(&totalUnique)
	
	response := models.VisitorStatsResponse{
		DailyStats:   dailyStats,
		MonthlyStats: monthlyStats,
		TotalUnique:  totalUnique,
	}
	
	c.JSON(http.StatusOK, response)
}

// TrackPhoneClick регистрирует клик по кнопке телефона
// @Summary Отслеживание кликов по телефону
// @Description Регистрирует клик по кнопке с номером телефона
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /track-phone-click [post]
func TrackPhoneClick(c *gin.Context) {
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	date := time.Now().Format("2006-01-02")
	
	phoneClick := models.PhoneClickStat{
		IPAddress: clientIP,
		UserAgent: userAgent,
		Date:      date,
	}
	
	if err := database.DB.Create(&phoneClick).Error; err != nil {
		log.Printf("❌ Ошибка сохранения клика по телефону: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения статистики кликов"})
		return
	}
	
	log.Printf("✅ Клик по телефону зарегистрирован: IP=%s, дата=%s", clientIP, date)
	c.JSON(http.StatusOK, gin.H{"message": "Клик по телефону зарегистрирован"})
}

// GetPhoneClickStats возвращает статистику кликов по телефону
// @Summary Получить статистику кликов по телефону
// @Description Возвращает статистику кликов по кнопке телефона
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Количество дней для отображения (по умолчанию 30)"
// @Success 200 {object} models.PhoneClickStatsResponse
// @Router /admin/analytics/phone-clicks [get]
func GetPhoneClickStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}
	
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	
	// Общее количество кликов
	var totalClicks int64
	database.DB.Model(&models.PhoneClickStat{}).
		Where("date >= ?", startDate).
		Count(&totalClicks)
	
	// Уникальные клики (по IP)
	var uniqueClicks int64
	database.DB.Model(&models.PhoneClickStat{}).
		Where("date >= ?", startDate).
		Distinct("ip_address").
		Count(&uniqueClicks)
	
	// Дневная статистика кликов
	var dailyClicks []models.DailyStat
	database.DB.Raw(`
		SELECT 
			date,
			COUNT(DISTINCT ip_address) as unique_views,
			COUNT(*) as total_views
		FROM phone_click_stats 
		WHERE date >= ? 
		GROUP BY date 
		ORDER BY date DESC
	`, startDate).Scan(&dailyClicks)
	
	response := models.PhoneClickStatsResponse{
		TotalClicks:  totalClicks,
		UniqueClicks: uniqueClicks,
		DailyClicks:  dailyClicks,
	}
	
	c.JSON(http.StatusOK, response)
}

// GetPhoneContacts возвращает список всех телефонных контактов
// @Summary Получить список телефонных контактов
// @Description Возвращает список всех оставленных телефонных номеров
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество записей на странице (по умолчанию 20)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/phone-contacts [get]
func GetPhoneContacts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	offset := (page - 1) * limit
	
	var contacts []models.PhoneContact
	var total int64
	
	// Получаем общее количество
	database.DB.Model(&models.PhoneContact{}).Count(&total)
	
	// Получаем контакты с пагинацией
	database.DB.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&contacts)
	
	c.JSON(http.StatusOK, gin.H{
		"contacts": contacts,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + int64(limit) - 1) / int64(limit),
	})
}

// DeletePhoneContact удаляет телефонный контакт
// @Summary Удалить телефонный контакт
// @Description Удаляет телефонный номер из базы данных
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID контакта"
// @Success 200 {object} map[string]interface{}
// @Router /admin/phone-contacts/{id} [delete]
func DeletePhoneContact(c *gin.Context) {
	id := c.Param("id")
	
	result := database.DB.Delete(&models.PhoneContact{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления контакта"})
		return
	}
	
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Контакт не найден"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Контакт успешно удален"})
}

// GetDatabaseStatus проверяет статус базы данных
// @Summary Статус базы данных
// @Description Получить информацию о типе базы данных и количестве записей
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Admin Token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/admin/database-status [get]
func GetDatabaseStatus(c *gin.Context) {
	if !isValidAdminToken(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	// Определяем тип базы данных
	dbType := "SQLite"
	databaseURL := os.Getenv("DATABASE_URL")
	ginMode := os.Getenv("GIN_MODE")
	if databaseURL != "" || ginMode == "release" {
		dbType = "PostgreSQL"
	}

	// Считаем записи в основных таблицах
	var phoneContactsCount, visitorsCount, phoneClicksCount int64
	database.DB.Model(&models.PhoneContact{}).Count(&phoneContactsCount)
	database.DB.Model(&models.VisitorStat{}).Count(&visitorsCount)
	database.DB.Model(&models.PhoneClickStat{}).Count(&phoneClicksCount)

	c.JSON(http.StatusOK, gin.H{
		"database_type":      dbType,
		"database_url_set":   databaseURL != "",
		"gin_mode":          ginMode,
		"phone_contacts":    phoneContactsCount,
		"visitors":          visitorsCount,
		"phone_clicks":      phoneClicksCount,
		"timestamp":         time.Now(),
	})
}

