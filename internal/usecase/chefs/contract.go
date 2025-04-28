package chefs

import (
	"context"
	"mime/multipart"

	entity "domashka-backend/internal/entity/chefs"
	dish "domashka-backend/internal/entity/dishes"
)

type chefRepo interface {
	GetChefByDishID(ctx context.Context, dishID int64) (*entity.Chef, error)
	GetChefRatingByChefID(ctx context.Context, chefID int64) (*entity.Chef, error)
	GetChefByID(ctx context.Context, chefID int64) (*entity.Chef, error)
	GetChefExperienceYears(ctx context.Context, chefID int64) (int, error)
	SaveChefAvatar(ctx context.Context, chefID int64, fileHeader *multipart.FileHeader) (string, error)
	GetTopChefs(ctx context.Context, limit int) ([]entity.Chef, error)
	GetChefAvatarURLByDishID(ctx context.Context, dishID int64) (string, error)
	GetChefAvatarURLByChefID(ctx context.Context, chefID int64) (string, error)
	GetChefCertifications(ctx context.Context, chefID int64) ([]entity.Certification, error)
	GetDishesByChefID(ctx context.Context, chefID int64) ([]dish.Dish, error)
}
