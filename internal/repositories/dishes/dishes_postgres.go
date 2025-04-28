package dishes

import (
	"context"
	dishEntity "domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/utils/pointers"
	"domashka-backend/pkg/postgres"
	"errors"
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{
		pg: pg,
	}
}

func (r *Repository) GetDishByID(ctx context.Context, dishID int64) (*dishEntity.Dish, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT * FROM dishes WHERE id = $1", dishID)

	var dish dishEntity.Dish
	err := row.Scan(
		&dish.ID,
		&dish.ChefID,
		&dish.Name,
		&dish.Description,
		&dish.ImageURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dishEntity.ErrDishNotFound
		}
		return nil, err
	}
	return &dish, nil
}
func (r *Repository) GetDishRatingByID(ctx context.Context, dishID int64) (*dishEntity.Dish, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT * FROM dish_ratings WHERE dish_id = $1", dishID)
	var rating DishRating
	err := row.Scan(
		&rating.DishID,
		&rating.Rating,
		&rating.ReviewsCount,
	)
	if err != nil {
		return nil, err
	}
	dish := dishEntity.Dish{Rating: pointers.To(rating.Rating), ReviewsCount: pointers.To(rating.ReviewsCount)}
	return &dish, nil
}
func (r *Repository) GetDishesByChefID(ctx context.Context, chefID int64) ([]dishEntity.Dish, error) {
	var dishes []dishEntity.Dish
	rows, err := r.pg.Pool.Query(ctx, "SELECT * FROM dishes WHERE chef_id = $1", chefID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var dish dishEntity.Dish
		err := rows.Scan(
			&dish.ID,
			&dish.ChefID,
			&dish.Name,
			&dish.Description,
			&dish.ImageURL,
		)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}
	return dishes, nil
}

func (r *Repository) GetNutritionByDishID(ctx context.Context, dishID int64) (*dishEntity.Nutrition, error) {
	var nutrition dishEntity.Nutrition
	err := r.pg.Pool.QueryRow(ctx, "SELECT * FROM nutritions WHERE dish_id = $1", dishID).Scan(
		&nutrition.DishID,
		&nutrition.Calories,
		&nutrition.Fat,
		&nutrition.Carbohydrates,
		&nutrition.Protein,
	)
	if err != nil {
		return nil, err
	}
	return &nutrition, nil
}

func (r *Repository) GetDishSizesByDishID(ctx context.Context, dishID int64) ([]dishEntity.Size, error) {
	var sizes []dishEntity.Size
	rows, err := r.pg.Pool.Query(ctx, "SELECT * FROM dish_sizes WHERE dish_id = $1", dishID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var size dishEntity.Size
		err := rows.Scan(
			&size.ID,
			&size.DishID,
			&size.Label,
			&size.WeightValue,
			&size.WeightUnit,
			&size.PriceValue,
			&size.PriceCurrency,
		)
		if err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}
	return sizes, nil
}

func (r *Repository) GetIngredientsByDishID(ctx context.Context, dishID int64) ([]dishEntity.Ingredient, error) {
	var ingredients []dishEntity.Ingredient
	rows, err := r.pg.Pool.Query(ctx, "SELECT * FROM dish_ingredients WHERE dish_id = $1", dishID)
	if err != nil {
		return nil, err
	}
	var dishIngredients []DishIngredient
	defer rows.Close()
	for rows.Next() {
		var dishIngredient DishIngredient
		err := rows.Scan(
			&dishIngredient.DishID,
			&dishIngredient.IngredientID,
			&dishIngredient.IsRemovable,
		)
		if err != nil {
			return nil, err
		}
		dishIngredients = append(dishIngredients, dishIngredient)
	}
	for _, dishIngredient := range dishIngredients {
		var ingredient dishEntity.Ingredient
		err := r.pg.Pool.QueryRow(ctx, "SELECT * FROM ingredients WHERE id = $1", dishIngredient.IngredientID).Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.ImageURL,
			&ingredient.CategoryID,
			&ingredient.IsAllergen,
		)
		if err != nil {
			return nil, err
		}
		ingredient.IsRemovable = dishIngredient.IsRemovable
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func (r *Repository) GetTopDishes(ctx context.Context, limit int) ([]dishEntity.Dish, error) {
	var dishes []dishEntity.Dish
	query := `
        SELECT 
            d.id,
            d.name,
            d.description,
            d.image_url,
            d.chef_id,
            dr.rating,
            dr.reviews_count
        FROM 
            public.dishes d
        JOIN 
            public.dish_ratings dr ON d.id = dr.dish_id
        ORDER BY 
            dr.rating DESC
        LIMIT $1;
    `
	rows, err := r.pg.Pool.Query(ctx, query, limit)
	if errors.Is(err, pgx.ErrNoRows) {
		return []dishEntity.Dish{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var dish dishEntity.Dish
		err := rows.Scan(
			&dish.ID,
			&dish.Name,
			&dish.Description,
			&dish.ImageURL,
			&dish.ChefID,
			&dish.Rating,
			&dish.ReviewsCount,
		)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}
	return dishes, nil
}
