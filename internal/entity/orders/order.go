package orders

import (
	"time"

	cartentity "domashka-backend/internal/entity/cart"
	chefEntity "domashka-backend/internal/entity/chefs"
)

const (
	StatusRejected = iota
	StatusCreated
	StatusAccepted
	StatusCooked
	StatusInDelivery
	StatusDelivered
)

type Order struct {
	ID              int64     `db:"id"`
	ChefID          int64     `db:"chef_id"`
	ShiftID         int64     `db:"shift"`
	Status          int32     `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	TotalCost       float64   `db:"total_cost"`
	LeaveByTheDoor  bool      `db:"leave_by_the_door"`
	ClientAddressID int64     `db:"client_address_id"`
	UserID          int64     `db:"user_id"`
}

// ReviewDetail должен быть определён в том же пакете или импортирован
type ReviewDetail struct {
	CanWriteReview bool
	Rating         *int
	Comment        *string
}

// OrderProfile объединяет основные данные заказа и сопутствующую информацию
type OrderProfile struct {
	Order  *Order                // основная сущность заказа из этой же модели
	Items  []cartentity.CartItem // позиции в заказе
	Chef   *chefEntity.Chef      // данные шеф-повара
	Review *ReviewDetail         // отзыв пользователя по этому заказу
}
