package dishes

import (
	"context"
	entity "domashka-backend/internal/entity/dishes"
	"domashka-backend/pkg/postgres"
)

type Repository struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]entity.DishCategory, error) {
	categories := make([]entity.DishCategory, 0)
	rows, err := r.pg.Pool.Query(ctx, "SELECT * FROM dishes_categories")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var category entity.DishCategory
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *Repository) CreateDish(ctx context.Context, dish *entity.Dish) (*string, error) {
	var id string
	err := r.pg.Pool.QueryRow(ctx, "INSERT INTO dishes (chef_id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id", dish.ChefID, dish.Name, dish.Description, dish.Price, dish.Stock).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// TODO: implement
//func (r *Repository) GetDishByID(ctx context.Context, id string) (*entity.Dish, error)            {
//	var dish entity.Dish
//	err := r.pg.Pool.QueryRow(ctx, "SELECT * FROM dishes WHERE id = $1", id).Scan(&dish.ID, &dish.ChefID, &dish.Name, &dish.Description, &dish.Price, &dish.Stock, Cre)
//}
//func (r *Repository) GetDishesByChefID(ctx context.Context, chefID string) ([]entity.Dish, error) {}
//func (r *Repository) GetDishesByCategoryID(ctx context.Context, categoryID string) ([]entity.Dish, error) {
//}
//func (r *Repository) UpdateDish(ctx context.Context, dish *entity.Dish) error {}
//func (r *Repository) RemoveDish(ctx context.Context, dishID string) error     {}
