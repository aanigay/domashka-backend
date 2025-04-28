package order

import (
	"context"
	"domashka-backend/internal/entity/orders"

	cartentity "domashka-backend/internal/entity/cart"
	addressentity "domashka-backend/internal/entity/geo"
)

type geoUsecase interface {
	GetLastUpdatedClientAddress(ctx context.Context, clientID int64) (*addressentity.Address, error)
}

type cartUsecase interface {
	GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error)
}

type shiftsRepo interface {
	GetActiveShiftIDByChefID(ctx context.Context, chefID int64) (int64, error)
}

type ordersRepo interface {
	CreateOrder(
		ctx context.Context,
		chefID int64,
		shiftID int64,
		clientAddressID int64,
		totalCost float32,
		leaveByTheDoor bool,
		callBeforehand bool,
	) (int64, error)
	AddCartItemToOrder(ctx context.Context, cartItemID int64, orderID int64) error
	GetOrdersByShiftID(ctx context.Context, shiftID int64) ([]orders.Order, error)
}
