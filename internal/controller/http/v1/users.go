package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"domashka-backend/internal/custom_errors"
	usersEntity "domashka-backend/internal/entity/users"
)

type usersHandler struct {
	log     logger
	usecase usersUsecase
}

func newUsersHandler(rg *gin.RouterGroup, log logger, usecase usersUsecase) {
	u := usersHandler{
		log:     log,
		usecase: usecase,
	}

	rg = rg.Group("/users")
	{
		rg.POST("/create", u.Create)
		rg.GET("/:id", u.GetByID)
		rg.PUT("/:id", u.Update)
		rg.DELETE("/:id", u.Delete)
	}
}

type CreateRequest struct {
	Username         string     `json:"username"`                // Уникальный идентификатор
	Alias            string     `json:"alias"`                   // Полное имя (ФИО)
	FirstName        string     `json:"first_name"`              // Имя
	SecondName       *string    `json:"second_name,omitempty"`   // Отчество (опционально)
	LastName         *string    `json:"last_name,omitempty"`     // Фамилия (опционально)
	Email            *string    `json:"email,omitempty"`         // Электронная почта (опционально)
	NumberPhone      *string    `json:"number_phone,omitempty"`  // Телефон (опционально)
	Status           int        `json:"status"`                  // Состояние (0 = не бан, 1 = бан)
	ExternalType     int        `json:"external_type"`           // Внешний статус (0, 1, 2)
	TelegramName     *string    `json:"telegram_name,omitempty"` // Telegram username (опционально)
	ExternalID       *string    `json:"external_id,omitempty"`   // Внешний ID (опционально)
	NotificationFlag int        `json:"notification_flag"`       // Уведомления (0 = включены, 1 = выключены)
	Role             string     `json:"role"`                    // Роль (client, chef, admin)
	Birthday         *time.Time `json:"birthday,omitempty"`      // День рождения (опционально)
}

// Create godoc
// @Summary      Создание пользователя
// @Description  Создание пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body      CreateRequest  true  "body"
// @Router       /users/create [post]
func (u *usersHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	r := CreateRequest{}

	if err := c.ShouldBindJSON(&r); err != nil {
		u.log.Error(ctx, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := u.usecase.Create(c.Request.Context(), &usersEntity.User{
		Username:         r.Username,
		Alias:            r.Alias,
		FirstName:        r.FirstName,
		SecondName:       r.SecondName,
		LastName:         r.LastName,
		Email:            r.Email,
		NumberPhone:      r.NumberPhone,
		Status:           r.Status,
		ExternalType:     r.ExternalType,
		TelegramName:     r.TelegramName,
		ExternalID:       r.ExternalID,
		NotificationFlag: r.NotificationFlag,
		Role:             r.Role,
		Birthday:         r.Birthday,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user_id": id})
}

// GetByID godoc
// @Summary      Получение пользователя
// @Description  Получение пользователя по ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Router       /users/{id} [get]
func (u *usersHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	user, err := u.usecase.GetByID(ctx, id)
	if err != nil {
		if err == custom_errors.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			u.log.Error(ctx, err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch user", "error": err, "user": user})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateRequest struct {
	Alias            string     `json:"alias"`                   // Полное имя (ФИО)
	FirstName        string     `json:"first_name"`              // Имя
	SecondName       *string    `json:"second_name,omitempty"`   // Отчество (опционально)
	LastName         *string    `json:"last_name,omitempty"`     // Фамилия (опционально)
	Email            *string    `json:"email,omitempty"`         // Электронная почта (опционально)
	NumberPhone      *string    `json:"number_phone,omitempty"`  // Телефон (опционально)
	Status           int        `json:"status"`                  // Состояние (0 = не бан, 1 = бан)
	ExternalType     int        `json:"external_type"`           // Внешний статус (0, 1, 2)
	TelegramName     *string    `json:"telegram_name,omitempty"` // Telegram username (опционально)
	ExternalID       *string    `json:"external_id,omitempty"`   // Внешний ID (опционально)
	NotificationFlag int        `json:"notification_flag"`       // Уведомления (0 = включены, 1 = выключены)
	Role             string     `json:"role"`                    // Роль (client, chef, admin)
	Birthday         *time.Time `json:"birthday,omitempty"`      // День рождения (опционально)

}

// Update godoc
// @Summary      Обновление пользователя
// @Description  Обновление данных пользователя по ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Param        request   body      UpdateRequest  true  "body"
// @Router       /users/{id} [put]
func (u *usersHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	r := UpdateRequest{}

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.usecase.Update(ctx, id, usersEntity.User{
		Alias:            r.Alias,
		FirstName:        r.FirstName,
		SecondName:       r.SecondName,
		LastName:         r.LastName,
		Email:            r.Email,
		NumberPhone:      r.NumberPhone,
		Status:           r.Status,
		ExternalType:     r.ExternalType,
		TelegramName:     r.TelegramName,
		ExternalID:       r.ExternalID,
		NotificationFlag: r.NotificationFlag,
		Role:             r.Role,
		Birthday:         r.Birthday,
	})

	if err != nil {
		u.log.Error(ctx, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Delete godoc
// @Summary      Удаление пользователя
// @Description  Удаление пользователя по ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Router       /users/{id} [delete]
func (u *usersHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	err := u.usecase.Delete(ctx, id)
	if err != nil {
		u.log.Error(ctx, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
