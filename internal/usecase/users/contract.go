package users

import (
	"context"

	chefsEntity "domashka-backend/internal/entity/chefs"
	dishesEntity "domashka-backend/internal/entity/dishes"
	usersEntity "domashka-backend/internal/entity/users"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type repo interface {
	Create(ctx context.Context, user *usersEntity.User) error
	GetByID(ctx context.Context, id int64) (*usersEntity.User, error)
	Update(ctx context.Context, id int64, user usersEntity.User) error
	Delete(ctx context.Context, id int64) error
	CheckIfUserIsChef(ctx context.Context, userID int64) (*int64, bool, error)
	GetFavoritesDishesByUserID(ctx context.Context, userID int64) ([]dishesEntity.Dish, error)
	GetFavoritesChefsByUserID(ctx context.Context, userID int64) ([]chefsEntity.Chef, error)
}
