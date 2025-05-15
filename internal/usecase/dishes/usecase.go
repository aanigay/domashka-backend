package dishes

import (
	"context"
	"fmt"
	"mime/multipart"

	entity "domashka-backend/internal/entity/dishes"
)

const (
	dishesFilePrefixTpl      = "images/dish/%d"
	ingredientsFilePrefixTpl = "images/ingredient/%d"
)

type Usecase struct {
	dishRepo dishRepo
	s3Client s3Client
}

func New(dishRepo dishRepo, s3client s3Client) *Usecase {
	return &Usecase{
		dishRepo: dishRepo,
		s3Client: s3client,
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

func (u *Usecase) GetDishRatingByID(ctx context.Context, dishID int64) (*entity.Dish, error) {
	return u.dishRepo.GetDishRatingByID(ctx, dishID)
}

func (u *Usecase) GetAll(ctx context.Context) ([]entity.Dish, error) {
	dishes, err := u.dishRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for idx, dish := range dishes {
		dishRating, err := u.dishRepo.GetDishRatingByID(ctx, dish.ID)
		if err != nil {
			return nil, err
		}
		dishes[idx].Rating = dishRating.Rating
	}

	return dishes, nil
}

func (u *Usecase) GetDishesByChefID(ctx context.Context, chefID int64, limit int) ([]entity.Dish, error) {
	dishes, err := u.dishRepo.GetDishesByChefID(ctx, chefID, &limit)
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

func (u *Usecase) GetAllDishesByChefID(ctx context.Context, chefID int64) ([]entity.Dish, error) {
	dishes, err := u.dishRepo.GetDishesByChefID(ctx, chefID, nil)
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

func (u *Usecase) SetDishImage(ctx context.Context, dishID int64, image *multipart.FileHeader) (string, error) {
	filePath := fmt.Sprintf(dishesFilePrefixTpl, dishID)
	publicURL, err := u.s3Client.UploadPicture(ctx, filePath, image)
	if err != nil {
		return "", err
	}
	err = u.dishRepo.SetDishImageURL(ctx, dishID, publicURL)
	if err != nil {
		return "", err
	}
	return publicURL, nil
}

func (u *Usecase) SetIngredientImage(ctx context.Context, ingredientID int64, image *multipart.FileHeader) (string, error) {
	filePath := fmt.Sprintf(ingredientsFilePrefixTpl, ingredientID)
	publicURL, err := u.s3Client.UploadPicture(ctx, filePath, image)
	if err != nil {
		return "", err
	}
	err = u.dishRepo.SetIngredientImageURL(ctx, ingredientID, publicURL)
	if err != nil {
		return "", err
	}
	return publicURL, nil
}

func (u *Usecase) GetCategoryTitleByDishID(ctx context.Context, dishID int64) (string, error) {
	return u.dishRepo.GetCategoryTitleByDishID(ctx, dishID)
}

func (u *Usecase) Create(
	ctx context.Context,
	dish *entity.Dish,
	nutrition *entity.Nutrition,
	sizes []entity.Size,
	ingredients []entity.Ingredient,
) (int64, error) {
	if dish == nil {
		return 0, fmt.Errorf("dish is nil")
	}
	return u.dishRepo.CreateDish(ctx, dish, nutrition, sizes, ingredients)
}

func (u *Usecase) Update(
	ctx context.Context,
	dish *entity.Dish,
	nutrition *entity.Nutrition,
	sizes []entity.Size,
	ingredients []entity.Ingredient,
) (int64, error) {
	if dish == nil {
		return 0, fmt.Errorf("dish is nil")
	}
	if dish.ID == 0 {
		return 0, fmt.Errorf("dish id is 0")
	}
	rating, err := u.dishRepo.GetDishRatingByID(ctx, dish.ID)
	if err != nil {
		return 0, err
	}
	err = u.dishRepo.DeleteDish(ctx, dish.ID)
	if err != nil {
		return 0, err
	}
	dish.Rating = rating.Rating
	dish.ReviewsCount = rating.ReviewsCount
	return u.dishRepo.CreateDish(ctx, dish, nutrition, sizes, ingredients)
}

func (u *Usecase) Delete(ctx context.Context, dishID int64) error {
	return u.dishRepo.DeleteDish(ctx, dishID)
}

func (u *Usecase) GetAllIngredients(ctx context.Context) ([]entity.Ingredient, error) {
	return u.dishRepo.GetAllIngredients(ctx)
}

func (u *Usecase) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	return u.dishRepo.GetAllCategories(ctx)
}
