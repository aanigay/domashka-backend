package dishes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	dishEntity "domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/utils/pointers"
	"domashka-backend/pkg/postgres"
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
	row := r.pg.Pool.QueryRow(ctx, "SELECT id, chef_id, name, description, image_url FROM dishes WHERE id = $1", dishID)
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
	if err = r.pg.Pool.QueryRow(ctx, `SELECT category_id FROM dishes_dish_categories WHERE dish_id = $1`, dishID).Scan(&dish.CategoryID); err != nil {
		fmt.Println(err)
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
func (r *Repository) GetDishesByChefID(ctx context.Context, chefID int64, limit *int) ([]dishEntity.Dish, error) {
	var dishes []dishEntity.Dish
	query := "SELECT id, chef_id, name, description, image_url, is_deleted FROM dishes WHERE chef_id = $1 AND is_deleted = false"
	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}
	rows, err := r.pg.Pool.Query(ctx, query, chefID)
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
			&dish.IsDeleted,
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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
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
        WHERE d.is_deleted = false
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

func (r *Repository) SetDishImageURL(ctx context.Context, dishID int64, imageURL string) error {
	if _, err := r.pg.Pool.Exec(ctx,
		"UPDATE dishes SET image_url=$1 WHERE id=$2",
		imageURL, dishID,
	); err != nil {
		return fmt.Errorf("db update: %w", err)
	}

	return nil
}

func (r *Repository) SetIngredientImageURL(ctx context.Context, ingredientID int64, imageURL string) error {
	if _, err := r.pg.Pool.Exec(ctx,
		"UPDATE ingredients SET image_url=$1 WHERE id=$2",
		imageURL, ingredientID,
	); err != nil {
		return fmt.Errorf("db update: %w", err)
	}

	return nil
}

func (r *Repository) GetCategoryTitleByDishID(ctx context.Context, dishID int64) (string, error) {
	var categoryID int64
	err := r.pg.Pool.QueryRow(ctx, "SELECT category_id FROM dishes_dish_categories WHERE dish_id = $1", dishID).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "Другое", nil
		}
		return "", err
	}
	var categoryTitle string
	err = r.pg.Pool.QueryRow(ctx, "SELECT title FROM dish_categories WHERE id = $1", categoryID).Scan(&categoryTitle)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "Другое", nil
		}
		return "", err
	}
	return categoryTitle, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]dishEntity.Dish, error) {
	var dishes []dishEntity.Dish
	rows, err := r.pg.Pool.Query(ctx, "SELECT id, chef_id, name, description, image_url FROM dishes")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []dishEntity.Dish{}, nil
		}
		return nil, err
	}
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

func (r *Repository) CreateDish(ctx context.Context, dish *dishEntity.Dish, nutrition *dishEntity.Nutrition, sizes []dishEntity.Size, ingredients []dishEntity.Ingredient) (int64, error) {
	if dish == nil {
		return 0, fmt.Errorf("dish is nil")
	}
	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	var dishID int64
	err = tx.QueryRow(ctx, "INSERT INTO dishes (chef_id, name, description, image_url) VALUES ($1, $2, $3, '') RETURNING id", dish.ChefID, dish.Name, dish.Description).Scan(&dishID)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(ctx, "INSERT INTO dishes_dish_categories (dish_id, category_id) VALUES ($1, $2)", dishID, dish.CategoryID)
	if err != nil {
		return 0, err
	}
	if nutrition != nil {
		_, err = tx.Exec(ctx, "INSERT INTO nutritions (dish_id, calories, fat, carbohydrates, protein) VALUES ($1, $2, $3, $4, $5)", dishID, nutrition.Calories, nutrition.Fat, nutrition.Carbohydrates, nutrition.Protein)
	}
	for _, size := range sizes {
		_, err = tx.Exec(ctx, "INSERT INTO dish_sizes (dish_id, label, weight_value, weight_unit, price_value, price_currency) VALUES ($1, $2, $3, $4, $5, $6)", dishID, size.Label, size.WeightValue, size.WeightUnit, size.PriceValue, size.PriceCurrency)
		if err != nil {
			return 0, err
		}
	}
	for _, ingredient := range ingredients {
		_, err = tx.Exec(ctx, "INSERT INTO dish_ingredients (dish_id, ingredient_id, is_removable) VALUES ($1, $2, $3)", dishID, ingredient.ID, ingredient.IsRemovable)
		if err != nil {
			return 0, err
		}
	}
	_, err = tx.Exec(ctx, "INSERT INTO dish_ratings (dish_id, rating, reviews_count) VALUES ($1, $2, $3)", dishID, pointers.From(dish.Rating), pointers.From(dish.Rating))
	if err != nil {
		return 0, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return dishID, nil
}

func (r *Repository) DeleteDish(ctx context.Context, dishID int64) error {
	_, err := r.pg.Pool.Exec(ctx, "UPDATE dishes SET is_deleted = true WHERE id = $1", dishID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetRating(ctx context.Context, dishID int64, rating float32, reviewsCount int64) error {
	_, err := r.pg.Pool.Exec(ctx, "UPDATE dish_ratings SET rating = $1, reviews_count = $2 WHERE dish_id = $3", rating, reviewsCount, dishID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllIngredients(ctx context.Context) ([]dishEntity.Ingredient, error) {
	query := `SELECT id, name, image_url FROM ingredients`
	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	ingredients := make([]dishEntity.Ingredient, 0)
	for rows.Next() {
		var ingredient dishEntity.Ingredient
		err := rows.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.ImageURL,
		)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]dishEntity.Category, error) {
	query := `SELECT id, title FROM dish_categories`
	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	categories := make([]dishEntity.Category, 0)
	for rows.Next() {
		var category dishEntity.Category
		err := rows.Scan(
			&category.ID,
			&category.Title,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
