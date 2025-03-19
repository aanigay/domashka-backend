package users

import (
	"context"
	cartEntities "domashka-backend/internal/entity/carts"

	usersEntity "domashka-backend/internal/entity/users"
)

type UseCase struct {
	usersRepo usersRepo
	cartRepo  cartRepo
}

func New(usersRepo usersRepo, cartRepo cartRepo) *UseCase {
	return &UseCase{
		usersRepo: usersRepo,
		cartRepo:  cartRepo,
	}
}

func (u *UseCase) Create(ctx context.Context, user *usersEntity.User) (*string, error) {
	userID, err := u.usersRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// При создании юзера сразу добавляем ему корзину
	emptyCart := &cartEntities.Cart{
		UserID: *userID,
	}
	err = u.cartRepo.CreateCart(ctx, emptyCart)
	if err != nil {
		return nil, err
	}
	return userID, nil
}

func (u *UseCase) GetByID(ctx context.Context, id string) (*usersEntity.User, error) {
	return u.usersRepo.GetByID(ctx, id)
}

func (u *UseCase) Update(ctx context.Context, id string, user usersEntity.User) error {
	return u.usersRepo.Update(ctx, id, user)
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	err := u.cartRepo.DeleteCart(ctx, id)
	if err != nil {
		return err
	}
	return u.usersRepo.Delete(ctx, id)
}
