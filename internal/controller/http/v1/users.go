package v1

import (
	"errors"
	"net/http"
	"strconv"
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
		rg.GET("/:id", u.GetProfile)
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

func (u *usersHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	r := CreateRequest{}

	if err := c.ShouldBindJSON(&r); err != nil {
		u.log.Error(ctx, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.usecase.Create(c.Request.Context(), &usersEntity.User{
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

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (u *usersHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// TODO: handle error
	}
	user, err := u.usecase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
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

func (u *usersHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// TODO: handle error
	}

	r := UpdateRequest{}

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = u.usecase.Update(ctx, id, usersEntity.User{
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

func (u *usersHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// TODO: handle error
	}

	err = u.usecase.Delete(ctx, id)
	if err != nil {
		u.log.Error(ctx, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (u *usersHandler) GetProfile(c *gin.Context) {
	mock := `
{
  "status": "success",
  "data": {
    "user_id": 1,
    "name": "Ульяна",
    "phone_number": "+7(925) 040 12-34",
    "last_order": {
      "label": "Заказ в субботу, 14:30",
      "dishes": [
        {
          "dish_id": 301,
          "title": "Цезарь с курицей и салатом айсберг",
          "image_url": "https://example.com/images/caesar.jpg",
          "rating": 4.7,
          "details": "S, без грибов и лука и еще какой то хуйни",
          "price_from": {
            "value": "320",
            "currency": "RUB"
          }
        },
        {
          "dish_id": 302,
          "title": "Греческий салат",
          "image_url": "https://example.com/images/greek_salad.jpg",
          "rating": 4.5,
          "price_from": {
            "value": "2 x 250",
            "currency": "RUB"
          }
        }
      ],
      "chef": {
        "id": 1,
        "name": "Раиса Виноградова",
        "avatar_url": "https://example.com/images/chef_raisa.jpg",
        "rating": 4.7,
        "reviews_count": 34
      },
      "review": {
        "can_write_review": true,
		"rating": null, 
		"comment": null
      }
    },
    "favorite_chefs": [
      {
        "id": 2,
        "name": "Александр Иванов",
        "avatar_url": "https://example.com/images/chef_raisa.jpg",
        "rating": 4.5,
        "reviews_count": 12
      },
      {
        "id": 3,
        "name": "Елена Петрова",
        "avatar_url": "https://example.com/images/chef_raisa.jpg",
        "rating": 4.2,
        "reviews_count": 5
      }
    ],
    "stop_list_of_ingredients" : [
      {
        "id": 1,
        "name": "Морковь",
        "image_url": "/images/carrot.jpg",
        "stop": true
      },
      {
        "id": 2,
        "name": "Картошка",
        "image_url": "/images/potato.jpg",
        "stop": false
      },
      {
        "id": 3,
        "name": "Свекла",
        "image_url": "/images/cabbage.jpg",
        "stop": true
      },
      {
        "id": 4,
        "name": "Орехи",
        "image_url": "/images/carrot.jpg",
        "stop": true
      }
    ],
    "orders": [
      {
        "label": "Заказ в субботу, 14:30",
        "dishes": [
          {
            "dish_id": 301,
            "title": "Цезарь с курицей и салатом айсберг",
            "image_url": "https://example.com/images/caesar.jpg",
            "rating": 4.7,
            "details": "S, без грибов и лука и еще какой то хуйни",
            "price_from": {
              "value": "320",
              "currency": "RUB"
            }
          },
          {
            "dish_id": 302,
            "title": "Греческий салат",
            "image_url": "https://example.com/images/greek_salad.jpg",
            "details": null,
            "rating": 4.5,
            "price_from": {
              "value": "2 x 250",
              "currency": "RUB"
            }
          }
        ],
        "chef": {
          "id": 1,
          "name": "Раиса Виноградова",
          "avatar_url": "https://example.com/images/chef_raisa.jpg",
          "rating": 4.7,
          "reviews_count": 34
        },
        "review": {
          "can_write_review": false,
          "rating": 4,
          "comment": "Отличный блюдо"
        }
      },
      {
        "label": "Заказ в субботу, 14:30",
        "dishes": [
          {
            "dish_id": 301,
            "title": "Цезарь с курицей и салатом айсберг",
            "image_url": "/images/caesar.jpg",
            "rating": 4.7,
            "details": "S, без грибов и лука и еще какой то хуйни",
            "price_from": {
              "value": "320",
              "currency": "RUB"
            }
          },
          {
            "dish_id": 302,
            "title": "Греческий салат",
            "image_url": "images/greek_salad.jpg",
            "rating": 4.5,
            "details": null,
            "price_from": {
              "value": "2 x 250",
              "currency": "RUB"
            }
          }
        ],
        "chef": {
          "id": 1,
          "name": "Раиса Виноградова",
          "avatar_url": "/images/chef_raisa.jpg",
          "rating": 4.7,
          "reviews_count": 34
        },
        "review": {
          "can_write_review": false,
          "rating": 5,
          "comment": "Отличный повар, люблю его <3"
        }
      }
    ]
  }
}
`
	c.Data(200, "application/json; charset=utf-8", []byte(mock))
}
