package dishes

import (
	"context"
	"fmt"

	entity "domashka-backend/internal/entity/dishes"
)

type Usecase struct {
	dishRepo dishRepo
}

func New(dishRepo dishRepo) *Usecase {
	return &Usecase{
		dishRepo: dishRepo,
	}
}

func (u *Usecase) GetDishByID(ctx context.Context, dishID int64) (*entity.Dish, error) {
	dish, err := u.dishRepo.GetDishByID(ctx, dishID)
	if err != nil {
		return nil, err
	}
	dishRating, err := u.dishRepo.GetDishRatingByID(ctx, dishID)
	if err != nil {
		return nil, err
	}
	dish.Rating = dishRating.Rating
	dish.ReviewsCount = dishRating.ReviewsCount
	return dish, nil
}

func (u *Usecase) GetDishesByChefID(ctx context.Context, chefID int64) ([]entity.Dish, error) {
	dishes, err := u.dishRepo.GetDishesByChefID(ctx, chefID)
	if err != nil {
		return nil, err
	}
	for idx, dish := range dishes {
		dishRating, err := u.dishRepo.GetDishRatingByID(ctx, dish.ID)
		if err != nil {
			return nil, err
		}
		dishes[idx].Rating = dishRating.Rating
		dishes[idx].ReviewsCount = dishRating.ReviewsCount
	}

	return dishes, nil
}

func (u *Usecase) GetNutritionByDishID(ctx context.Context, dishID int64) (*entity.Nutrition, error) {
	return u.dishRepo.GetNutritionByDishID(ctx, dishID)
}

func (u *Usecase) GetDishSizesByDishID(ctx context.Context, dishID int64) ([]entity.Size, error) {
	return u.dishRepo.GetDishSizesByDishID(ctx, dishID)
}

func (u *Usecase) GetIngredientsByDishID(ctx context.Context, dishID int64) ([]entity.Ingredient, error) {
	return u.dishRepo.GetIngredientsByDishID(ctx, dishID)
}

func (u *Usecase) GetMinimalPriceByDishID(ctx context.Context, dishID int64) (*entity.Price, error) {
	sizes, err := u.GetDishSizesByDishID(ctx, dishID)
	if err != nil {
		return nil, err
	}
	if len(sizes) == 0 {
		return nil, fmt.Errorf("no sizes attached to a dish")
	}
	// берем первое значение как минимальное
	minPrice := &entity.Price{
		Value:    sizes[0].PriceValue,
		Currency: sizes[0].PriceCurrency,
	}

	for _, size := range sizes {
		if size.PriceValue < minPrice.Value {
			minPrice.Value = size.PriceValue
			minPrice.Currency = size.PriceCurrency
		}
	}
	return minPrice, nil
}

func (u *Usecase) GetTopDishes(ctx context.Context, limit int) ([]entity.Dish, error) {
	return u.dishRepo.GetTopDishes(ctx, limit)
}
