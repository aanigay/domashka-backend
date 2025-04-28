package users

import (
	"context"

	usersEntity "domashka-backend/internal/entity/users"
)

type repo interface {
	Create(ctx context.Context, user *usersEntity.User) error
	GetByID(ctx context.Context, id int64) (*usersEntity.User, error)
	Update(ctx context.Context, id int64, user usersEntity.User) error
	Delete(ctx context.Context, id int64) error
}
