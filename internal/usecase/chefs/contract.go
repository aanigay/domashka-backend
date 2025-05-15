package chefs

import (
	"context"
	"mime/multipart"

	entity "domashka-backend/internal/entity/chefs"
	dish "domashka-backend/internal/entity/dishes"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type chefRepo interface {
	GetChefByDishID(ctx context.Context, dishID int64) (*entity.Chef, error)
	GetChefRatingByChefID(ctx context.Context, chefID int64) (*entity.Chef, error)
	GetChefByID(ctx context.Context, chefID int64) (*entity.Chef, error)
	GetChefExperienceYears(ctx context.Context, chefID int64) (int, error)
	SaveChefAvatar(ctx context.Context, chefID int64, publicURL string) error
	GetTopChefs(ctx context.Context, limit int) ([]entity.Chef, error)
	GetChefAvatarURLByDishID(ctx context.Context, dishID int64) (string, error)
	GetChefAvatarURLByChefID(ctx context.Context, chefID int64) (string, error)
	GetChefCertifications(ctx context.Context, chefID int64) ([]entity.Certification, error)
	GetDishesByChefID(ctx context.Context, chefID int64) ([]dish.Dish, error)
	GetNearestChefs(ctx context.Context, lat, long float64, distance, limit int) ([]entity.Chef, error)
	SetSmallAvatar(ctx context.Context, chefID int64, publicURL string) error
	GetAll(ctx context.Context) ([]entity.Chef, error)
}

type geoRepo interface {
	GetDistanceToChef(ctx context.Context, lat, long float64, chefID int64) (float64, error)
}

type s3Client interface {
	UploadPicture(ctx context.Context, filePrefix string, fileHeader *multipart.FileHeader) (string, error)
}
