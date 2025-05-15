package postgres

import (
	"context"
	"domashka-backend/pkg/postgres"
	"errors"
)

type Repository struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) AddFavoriteChef(ctx context.Context, userID, chefID int64) error {
	_, err := r.pg.Pool.Exec(ctx,
		`INSERT INTO user_favorite_chefs (user_id, chef_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, chefID,
	)
	return err
}

func (r *Repository) RemoveFavoriteChef(ctx context.Context, userID, chefID int64) error {
	cmd, err := r.pg.Pool.Exec(ctx,
		`DELETE FROM user_favorite_chefs WHERE user_id = $1 AND chef_id = $2`,
		userID, chefID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("favorite chef not found")
	}
	return nil
}

func (r *Repository) AddFavoriteDish(ctx context.Context, userID, dishID int64) error {
	_, err := r.pg.Pool.Exec(ctx,
		`INSERT INTO user_favorite_dishes (user_id, dish_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, dishID,
	)
	return err
}

func (r *Repository) RemoveFavoriteDish(ctx context.Context, userID, dishID int64) error {
	cmd, err := r.pg.Pool.Exec(ctx,
		`DELETE FROM user_favorite_dishes WHERE user_id = $1 AND dish_id = $2`,
		userID, dishID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("favorite dish not found")
	}
	return nil
}
