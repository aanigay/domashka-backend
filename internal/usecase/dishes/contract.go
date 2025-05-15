package dishes

import (
	"context"
	"mime/multipart"

	dishEntity "domashka-backend/internal/entity/dishes"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type dishRepo interface {
	GetDishByID(ctx context.Context, dishID int64) (*dishEntity.Dish, error)
	GetDishRatingByID(ctx context.Context, dishID int64) (*dishEntity.Dish, error)
	GetDishesByChefID(ctx context.Context, chefID int64, limit *int) ([]dishEntity.Dish, error)
	GetNutritionByDishID(ctx context.Context, dishID int64) (*dishEntity.Nutrition, error)
	GetDishSizesByDishID(ctx context.Context, dishID int64) ([]dishEntity.Size, error)
	GetIngredientsByDishID(ctx context.Context, dishID int64) ([]dishEntity.Ingredient, error)
	GetTopDishes(ctx context.Context, limit int) ([]dishEntity.Dish, error)
	SetDishImageURL(ctx context.Context, dishID int64, imageURL string) error
	SetIngredientImageURL(ctx context.Context, ingredientID int64, imageURL string) error
	GetCategoryTitleByDishID(ctx context.Context, dishID int64) (string, error)
	GetAll(ctx context.Context) ([]dishEntity.Dish, error)
	CreateDish(ctx context.Context, dish *dishEntity.Dish, nutrition *dishEntity.Nutrition, sizes []dishEntity.Size, ingredients []dishEntity.Ingredient) (int64, error)
	DeleteDish(ctx context.Context, dishID int64) error
	SetRating(ctx context.Context, dishID int64, rating float32, reviewsCount int64) error
	GetAllIngredients(ctx context.Context) ([]dishEntity.Ingredient, error)
	GetAllCategories(ctx context.Context) ([]dishEntity.Category, error)
}

type s3Client interface {
	UploadPicture(ctx context.Context, filePrefix string, fileHeader *multipart.FileHeader) (string, error)
}
