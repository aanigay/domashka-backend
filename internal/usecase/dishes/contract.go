package dishes

import (
	"context"

	dishEntity "domashka-backend/internal/entity/dishes"
)

type dishRepo interface {
	GetDishByID(ctx context.Context, dishID int64) (*dishEntity.Dish, error)
	GetDishRatingByID(ctx context.Context, dishID int64) (*dishEntity.Dish, error)
	GetDishesByChefID(ctx context.Context, chefID int64) ([]dishEntity.Dish, error)
	GetNutritionByDishID(ctx context.Context, dishID int64) (*dishEntity.Nutrition, error)
	GetDishSizesByDishID(ctx context.Context, dishID int64) ([]dishEntity.Size, error)
	GetIngredientsByDishID(ctx context.Context, dishID int64) ([]dishEntity.Ingredient, error)
	GetTopDishes(ctx context.Context, limit int) ([]dishEntity.Dish, error)
}
