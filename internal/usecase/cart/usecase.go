package cart

import (
	"context"
	"time"

	entities "domashka-backend/internal/entity/carts"
)

type UseCase struct {
	cartRepo cartRepo
}

func New(repo cartRepo) *UseCase {
	return &UseCase{
		cartRepo: repo,
	}
}

func (uc *UseCase) CreateCart(ctx context.Context, cart *entities.Cart) error {
	return uc.cartRepo.CreateCart(ctx, cart)
}

func (uc *UseCase) GetCart(ctx context.Context, userId string) (*entities.Cart, error) {
	return uc.cartRepo.GetCartByUserID(ctx, userId)
}

func (uc *UseCase) AddCartItem(ctx context.Context, cartItem *entities.CartItem) (itemID *string, err error) {
	// проверить на доступность блюдо
	return uc.cartRepo.AddCartItem(ctx, cartItem)
}

func (uc *UseCase) UpdateCartItem(ctx context.Context, cartItem *entities.CartItem) (itemID *string, err error) {
	cartItem.AddedAt = time.Now()
	return uc.cartRepo.UpdateCartItem(ctx, cartItem)
}

func (uc *UseCase) DeleteCartItem(ctx context.Context, itemID string) error {
	return uc.cartRepo.DeleteCartItem(ctx, itemID)
}

func (uc *UseCase) GetCartItems(ctx context.Context, userID string) ([]entities.CartItem, error) {
	return uc.cartRepo.GetCartItemsByUserID(ctx, userID)
}

func (uc *UseCase) ClearCart(ctx context.Context, userID string) error {
	return uc.cartRepo.ClearCartByUserID(ctx, userID)
}
