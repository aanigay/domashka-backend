package cart

import (
	"context"

	entities "domashka-backend/internal/entity/carts"
)

type cartRepo interface {
	CreateCart(ctx context.Context, cart *entities.Cart) error
	GetCartByUserID(ctx context.Context, userID string) (*entities.Cart, error)
	GetCartItemsByUserID(ctx context.Context, userID string) ([]entities.CartItem, error)
	ClearCartByUserID(ctx context.Context, userID string) error
	AddCartItem(ctx context.Context, cartItem *entities.CartItem) (*string, error)
	UpdateCartItem(ctx context.Context, cartItem *entities.CartItem) (*string, error)
	DeleteCartItem(ctx context.Context, cartItemID string) error
}
