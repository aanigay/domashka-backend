package dishes

import (
	"context"
	entity "domashka-backend/internal/entity/dishes"
)

type Usecase struct {
	repo dishesRepo
}

func New(repo dishesRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetDishByID(ctx context.Context, id string) (*entity.Dish, error) {
	return u.repo.GetDishByID(ctx, id)
}

func (u *Usecase) GetDishesByChefID(ctx context.Context, chefID string) ([]entity.Dish, error) {
	return u.repo.GetDishesByChefID(ctx, chefID)
}

func (u *Usecase) GetDishesByCategoryID(ctx context.Context, catID string) ([]entity.Dish, error) {
	return u.repo.GetDishesByCategoryID(ctx, catID)
}

func (u *Usecase) CreateDish(ctx context.Context, dish *entity.Dish) (*string, error) {
	return u.repo.CreateDish(ctx, dish)
}

func (u *Usecase) UpdateDish(ctx context.Context, dish *entity.Dish) error {
	return u.repo.UpdateDish(ctx, dish)
}

func (u *Usecase) RemoveDish(ctx context.Context, id string) error {
	return u.repo.RemoveDish(ctx, id)
}

func (u *Usecase) GetAllCategories(ctx context.Context) ([]entity.DishCategory, error) {
	return u.repo.GetAllCategories(ctx)
}
