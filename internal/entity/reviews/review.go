package reviews

import (
	"time"
)

// Review представляет отзыв пользователя о шеф-поваре.
type Review struct {
	ID              int64     `db:"id"`
	ChefID          int64     `db:"chef_id"`           // Идентификатор шеф-повара
	UserID          int64     `db:"user_id"`           // Идентификатор пользователя
	Stars           int16     `db:"stars"`             // Оценка (1–5)
	Comment         string    `db:"comment"`           // Текст комментария
	IsVerified      bool      `db:"is_verified"`       // Верифицирован ли отзыв
	IncludeInRating bool      `db:"include_in_rating"` // Участвует ли в рейтинге
	IsDeleted       bool      `db:"is_deleted"`        // Мягкое удаление
	CreatedAt       time.Time `db:"created_at"`        // Когда создан
	UpdatedAt       time.Time `db:"updated_at"`        // Когда обновлён
	OrderID         int64     `db:"order_id"`          // Идентификатор заказа
}
type ReviewWithUser struct {
	Review
	UserName string
}
