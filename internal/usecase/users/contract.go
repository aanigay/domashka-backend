package users

import (
	"context"
	cartEntities "domashka-backend/internal/entity/carts"

	usersEntity "domashka-backend/internal/entity/users"
)

type usersRepo interface {
	Create(ctx context.Context, user *usersEntity.User) (*string, error)
	GetByID(ctx context.Context, id string) (*usersEntity.User, error)
	Update(ctx context.Context, id string, user usersEntity.User) error
	Delete(ctx context.Context, id string) error
}

type cartRepo interface {
	CreateCart(ctx context.Context, cart *cartEntities.Cart) error
	DeleteCart(ctx context.Context, userID string) error
}
