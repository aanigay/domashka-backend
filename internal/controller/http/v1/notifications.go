package v1

import (
	notifEntity "domashka-backend/internal/entity/notifications"
	"domashka-backend/internal/utils/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type NotificationHandler struct {
	notificationUsecase notificationUsecase
}

func newNotificationHandler(rg *gin.RouterGroup, notification notificationUsecase) {
	h := &NotificationHandler{notificationUsecase: notification}

	rg = rg.Group("/notifications")
	{
		rg.POST("/", h.CreateNotification)
		rg.POST("/:id/resend", h.ResendNotification)
		rg.GET("/:id", h.GetNotificationByID)
		rg.GET("/", h.GetNotifications)
	}

}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	// Проверяем наличие JWT claims
	if _, exists := c.Get("claims"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized: authentication claims missing.",
		})
		return
	}

	var notification notifEntity.Notification
	// Привязываем JSON к структуре Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid JSON format: " + err.Error(),
		})
		return
	}

	// Проверяем, что поле Recipient не пустое
	if notification.Recipient == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Recipient is required.",
		})
		return
	}

	// Проверяем, что указан корректный сценарий
	allowedScenarios := []string{
		notifEntity.ScenarioSystem,
		notifEntity.ScenarioMarketing,
		notifEntity.ScenarioOrderStatus,
	}
	isValidScenario := false
	for _, s := range allowedScenarios {
		if notification.Scenario == s {
			isValidScenario = true
			break
		}
	}
	if !isValidScenario {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid scenario provided. Allowed scenarios: system, marketing, order_status.",
		})
		return
	}

	// Валидация и обработка в зависимости от канала
	switch notification.Channel {
	case notifEntity.ChannelEmail:
		if !validateEmail(notification.Recipient) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid email address.",
			})
			return
		}
		if !notification.Subject.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Email subject is required.",
			})
			return
		}
		if notification.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Email message is required.",
			})
			return
		}
		// Отправка email-уведомления
		if err := h.notificationUsecase.SendEmailNotification(c, notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to send email notification: " + err.Error(),
			})
			return
		}
	case notifEntity.ChannelSMS:
		if !validation.ValidatePhoneNumber(notification.Recipient) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid phone number format.",
			})
			return
		}
		// Логика отправки SMS-уведомления (если требуется) может быть добавлена здесь
	case notifEntity.ChannelPush:
		// Если push-уведомления ещё не поддерживаются, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Push notifications are not implemented yet.",
		})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Unsupported channel provided. Allowed channels: email, sms, push.",
		})
		return
	}

	// Создание уведомления в базе данных
	id, err := h.notificationUsecase.CreateNotification(c, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create notification: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":     id,
		"status": "created",
	})
}

func validateEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

func (h *NotificationHandler) ResendNotification(c *gin.Context) {
	// Проверка наличия "claims" в контексте
	if _, exists := c.Get("claims"); !exists {
		log.Printf("DEBUG: claims отсутствуют в контексте")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized: missing claims"})
		return
	}

	// Получение параметра id
	id := c.Param("id")
	if id == "" {
		log.Printf("DEBUG: Параметр id пустой")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request: missing id"})
		return
	}

	// Преобразование id в число
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("DEBUG: Не удалось преобразовать id '%s' в число: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid id format"})
		return
	}

	log.Printf("DEBUG: Попытка повторной отправки уведомления с id=%d", idInt)
	err = h.notificationUsecase.ResendNotification(c, idInt)
	if err != nil {
		log.Printf("DEBUG: Ошибка повторной отправки уведомления с id=%d: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	log.Printf("DEBUG: Уведомление с id=%d успешно отправлено повторно", idInt)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	if _, exists := c.Get("claims"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": ""})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	notification, err := h.notificationUsecase.GetNotificationByID(c, idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "notification": notification})
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	if _, exists := c.Get("claims"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": ""})
		return
	}

	filters := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		filters[k] = v[0]
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	notifications, total, err := h.notificationUsecase.GetNotifications(c, filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "notifications": notifications, "total": total})
}
