package cart

import (
	"context"

	cartentity "domashka-backend/internal/entity/cart"
	dishentity "domashka-backend/internal/entity/dishes"
)

type Usecase struct {
	cartRepo CartRepository
}

func New(cartRepo CartRepository) *Usecase {
	return &Usecase{cartRepo: cartRepo}
}

func (u *Usecase) AddItem(
	ctx context.Context,
	userID int64,
	dish dishentity.Dish,
	sizeID int64,
	addedIngredients []int64,
	removedIngredients []int64,
	notes string,
) (int64, error) {
	return u.cartRepo.AddItem(ctx, userID, dish, sizeID, addedIngredients, removedIngredients, notes)
}

func (u *Usecase) RemoveItem(ctx context.Context, cartItemID int64) error {
	return u.cartRepo.RemoveItem(ctx, cartItemID)
}

func (u *Usecase) GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error) {
	return u.cartRepo.GetCartItems(ctx, userID)
}

func (u *Usecase) ClearCart(ctx context.Context, userID int64) error {
	return u.cartRepo.Clear(ctx, userID)
}

func (u *Usecase) IncrementCartItem(ctx context.Context, cartItemID int64) (newQuantity int32, err error) {
	return u.cartRepo.IncrementCartItemQuantity(ctx, cartItemID)
}

func (u *Usecase) DecrementCartItem(ctx context.Context, cartItemID int64) (newQuantity int32, err error) {
	return u.cartRepo.DecrementCartItemQuantity(ctx, cartItemID)
}

func (u *Usecase) GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error) {
	return u.cartRepo.GetCartItemsByOrderID(ctx, orderID)
}
