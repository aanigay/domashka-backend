package cart

import (
	"context"
	"database/sql"
	"domashka-backend/internal/custom_errors"
	"errors"

	entities "domashka-backend/internal/entity/carts"
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

func (r *Repository) CreateCart(ctx context.Context, cart *entities.Cart) error {
	row := r.pg.Pool.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", cart.UserID)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return err
	}
	_, err := r.pg.Pool.Exec(ctx, "INSERT INTO carts (user_id) VALUES ($1)", userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetCartByUserID(ctx context.Context, userID string) (*entities.Cart, error) {
	row := r.pg.Pool.QueryRow(ctx, "SELECT user_id, created_at, updated_at FROM carts WHERE user_id = $1", userID)
	var cart entities.Cart
	if err := row.Scan(&cart.UserID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

func (r *Repository) GetCartItemsByUserID(ctx context.Context, userID string) ([]entities.CartItem, error) {
	rows, err := r.pg.Pool.Query(ctx, "SELECT * FROM cart_items WHERE user_id = $1", userID)
	cartItems := make([]entities.CartItem, 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cartItems, nil
		}
	}
	for rows.Next() {
		var cartItem entities.CartItem
		if err := rows.Scan(&cartItem.ID, &cartItem.UserID, &cartItem.DishID, &cartItem.ChefID, &cartItem.AdditionalIngredientsIDs, &cartItem.RemovedIngredientsIDs, &cartItem.AddedAt, &cartItem.CustomerNotes); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, cartItem)
	}
	return cartItems, nil
}

func (r *Repository) ClearCartByUserID(ctx context.Context, userID string) error {
	_, err := r.pg.Pool.Exec(ctx, "DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddCartItem(ctx context.Context, cartItem *entities.CartItem) (*string, error) {
	var cartItemID string
	err := r.pg.Pool.QueryRow(ctx, "INSERT INTO cart_items (user_id, dish_id, chef_id, additional_ingredients_ids, removed_ingredients_ids, added_at, customer_notes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		cartItem.UserID, cartItem.DishID, cartItem.ChefID, cartItem.AdditionalIngredientsIDs, cartItem.RemovedIngredientsIDs, cartItem.AddedAt, cartItem.CustomerNotes).Scan(&cartItemID)
	if err != nil {
		return nil, err
	}
	return &cartItemID, nil
}

func (r *Repository) UpdateCartItem(ctx context.Context, cartItem *entities.CartItem) (*string, error) {
	var cartItemID string
	err := r.pg.Pool.QueryRow(ctx, "UPDATE cart_items SET dish_id = $1, chef_id = $2, additional_ingredients_ids = $3, removed_ingredients_ids = $4, added_at = $5, customer_notes = $6 WHERE id = $7 RETURNING id",
		cartItem.DishID, cartItem.ChefID, cartItem.AdditionalIngredientsIDs, cartItem.RemovedIngredientsIDs, cartItem.AddedAt, cartItem.CustomerNotes, cartItem.ID).
		Scan(&cartItemID)
	if err != nil {
		return nil, err
	}
	return &cartItemID, nil
}

func (r *Repository) DeleteCartItem(ctx context.Context, cartItemID string) error {
	_, err := r.pg.Pool.Exec(ctx, "DELETE FROM cart_items WHERE id = $1", cartItemID)
	if err != nil {
		return err

	}
	return nil
}

func (r *Repository) DeleteCart(ctx context.Context, userID string) error {
	_, err := r.pg.Pool.Exec(ctx, "DELETE FROM carts WHERE user_id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}
