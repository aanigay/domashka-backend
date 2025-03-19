package dishes

import (
	"context"

	entity "domashka-backend/internal/entity/dishes"
)

type dishesRepo interface {
	GetAllCategories(ctx context.Context) ([]entity.DishCategory, error)
	CreateDish(ctx context.Context, dish *entity.Dish) (*string, error)
	GetDishByID(ctx context.Context, id string) (*entity.Dish, error)
	GetDishesByChefID(ctx context.Context, chefID string) ([]entity.Dish, error)
	GetDishesByCategoryID(ctx context.Context, categoryID string) ([]entity.Dish, error)
	UpdateDish(ctx context.Context, dish *entity.Dish) error
	RemoveDish(ctx context.Context, dishID string) error
}
