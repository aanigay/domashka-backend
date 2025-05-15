package cart

import (
	"context"

	cartentity "domashka-backend/internal/entity/cart"
	dishentity "domashka-backend/internal/entity/dishes"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type CartRepository interface {
	AddItem(
		ctx context.Context,
		userID int64,
		dish dishentity.Dish,
		sizeID int64,
		addedIngredients []int64,
		removedIngredients []int64,
		notes string,
	) (cartItemID int64, err error)
	RemoveItem(ctx context.Context, cartItemID int64) error
	GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error)
	GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error)
	Clear(ctx context.Context, userID int64) error
	IncrementCartItemQuantity(ctx context.Context, cartItemID int64) (int32, error)
	DecrementCartItemQuantity(ctx context.Context, cartItemID int64) (int32, error)
}
