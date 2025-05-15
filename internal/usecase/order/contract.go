package order

import (
	"context"
	"domashka-backend/internal/entity/orders"

	cartentity "domashka-backend/internal/entity/cart"
	chefEntity "domashka-backend/internal/entity/chefs"
	entity "domashka-backend/internal/entity/dishes"
	addressentity "domashka-backend/internal/entity/geo"
	reviewEntity "domashka-backend/internal/entity/reviews"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type geoUsecase interface {
	GetLastUpdatedClientAddress(ctx context.Context, clientID int64) (*addressentity.Address, error)
}

type cartUsecase interface {
	GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error)
}

type shiftsRepo interface {
	GetActiveShiftIDByChefID(ctx context.Context, chefID int64) (int64, error)
	AddToTotalProfit(ctx context.Context, shiftID int64, profit float64) error
}

type ordersRepo interface {
	CreateOrder(
		ctx context.Context,
		chefID int64,
		shiftID int64,
		userID int64,
		clientAddressID int64,
		totalCost float32,
		leaveByTheDoor bool,
		callBeforehand bool,

	) (int64, error)
	AddCartItemToOrder(ctx context.Context, cartItem *cartentity.CartItem, orderID, userID int64) error
	GetOrdersByShiftID(ctx context.Context, shiftID int64) ([]orders.Order, error)
	GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error)
	GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error)
	SetStatus(ctx context.Context, orderID int64, status int32) error
	GetShiftIDByOrderID(ctx context.Context, orderID int64) (int64, error)
	GetOrderByID(ctx context.Context, orderID int64) (*orders.Order, error)
	GetOrderedDishesIDsAndChefsIDs(ctx context.Context, userID int64, dishesLimit, chefsLimit int) ([]int64, []int64, error)
	GetStatus(ctx context.Context, orderID int64) (int32, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]orders.Order, error)
	CountDishInOrders(ctx context.Context, dishID int64) (int, error)
}

type dishesUsecase interface {
	GetDishByID(ctx context.Context, dishID int64) (*entity.Dish, error)
	GetDishRatingByID(ctx context.Context, dishID int64) (*entity.Dish, error)
}

type chefsUsecase interface {
	GetChefByID(ctx context.Context, chefID int64) (*chefEntity.Chef, error)
}
type reviewUsecase interface {
	GetReviewByOrderAndUserID(ctx context.Context, chefID, userID int64) (*reviewEntity.Review, error)
}
