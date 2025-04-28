package cart

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"

	cartentity "domashka-backend/internal/entity/cart"
	dishentity "domashka-backend/internal/entity/dishes"
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

func (r *Repository) Create(ctx context.Context, userID int64) error {
	_, err := r.pg.Pool.Exec(ctx, "INSERT INTO carts (user_id) VALUES ($1)", userID)
	return err
}

func (r *Repository) AddItem(
	ctx context.Context,
	userID int64,
	dish dishentity.Dish,
	sizeID int64,
	addedIngredients []int64,
	removedIngredients []int64,
	notes string,
) (cartItemID int64, err error) {
	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO cart_items (user_id, dish_id, chef_id, dish_size_id, customer_notes, quantity)
		VALUES ($1, $2, $3, $4, $5, 1)
		RETURNING id
	`, userID, dish.ID, dish.ChefID, sizeID, notes).Scan(&cartItemID)

	if err != nil {
		return 0, fmt.Errorf("failed to add item to cart: %w", err)
	}
	// Добавляем добавленные ингредиенты
	for _, ingredientID := range addedIngredients {
		_, err := tx.Exec(ctx, `
			INSERT INTO cart_item_added_ingredients (cart_item_id, ingredient_id)
			VALUES ($1, $2)
		`, cartItemID, ingredientID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert added ingredient %d: %w", ingredientID, err)
		}
	}

	// Добавляем удалённые ингредиенты
	for _, ingredientID := range removedIngredients {
		_, err := tx.Exec(ctx, `
			INSERT INTO cart_item_removed_ingredients (cart_item_id, ingredient_id)
			VALUES ($1, $2)
		`, cartItemID, ingredientID)
		if err != nil {
			return 0, fmt.Errorf("insert failed removed ingredient %d: %w", ingredientID, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	return cartItemID, nil
}

func (r *Repository) RemoveItem(ctx context.Context, cartItemID int64) error {
	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "DELETE FROM cart_items WHERE id = $1", cartItemID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM cart_item_added_ingredients WHERE cart_item_id = $1", cartItemID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM cart_item_removed_ingredients WHERE cart_item_id = $1", cartItemID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error) {
	rows, err := r.pg.Pool.Query(ctx, `
		SELECT
			ci.id, ci.quantity, ci.customer_notes,
			d.id, d.name, d.description, d.chef_id, d.image_url,
			ds.id, ds.dish_id, ds.label, ds.weight_value, ds.weight_unit, ds.price_value, ds.price_currency
		FROM cart_items ci
		JOIN dishes d ON ci.dish_id = d.id
		JOIN dish_sizes ds ON ci.dish_size_id = ds.id
		WHERE ci.user_id = $1
	`, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return []cartentity.CartItem{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []cartentity.CartItem
	cartItemIDs := make([]int64, 0)

	for rows.Next() {
		var item cartentity.CartItem
		var dish dishentity.Dish
		var size dishentity.Size

		err := rows.Scan(
			&item.ID, &item.Quantity, &item.Notes,
			&dish.ID, &dish.Name, &dish.Description, &dish.ChefID, &dish.ImageURL,
			&size.ID, &size.DishID, &size.Label, &size.WeightValue, &size.WeightUnit, &size.PriceValue, &size.PriceCurrency,
		)
		if err != nil {
			return nil, err
		}

		item.Dish = dish
		item.Size = size
		cartItemIDs = append(cartItemIDs, item.ID)
		items = append(items, item)
	}

	// ingredients (added and removed)
	addedMap, err := r.getCartItemIngredientsMap(ctx, "cart_item_added_ingredients", cartItemIDs)
	if err != nil {
		return nil, err
	}
	removedMap, err := r.getCartItemIngredientsMap(ctx, "cart_item_removed_ingredients", cartItemIDs)
	if err != nil {
		return nil, err
	}

	// заполняем ингредиенты
	for i, item := range items {
		items[i].AddedIngredients = addedMap[item.ID]
		items[i].RemovedIngredients = removedMap[item.ID]
	}

	return items, nil
}

func (r *Repository) getCartItemIngredientsMap(ctx context.Context, tableName string, cartItemIDs []int64) (map[int64][]dishentity.Ingredient, error) {
	if len(cartItemIDs) == 0 {
		return map[int64][]dishentity.Ingredient{}, nil
	}

	rows, err := r.pg.Pool.Query(ctx, fmt.Sprintf(`
		SELECT
			ci.cart_item_id,
			i.id, i.name, i.image_url, i.is_allergen, i.category_id
		FROM %s ci
		JOIN ingredients i ON ci.ingredient_id = i.id
		WHERE ci.cart_item_id = ANY($1)
	`, tableName), cartItemIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]dishentity.Ingredient)

	for rows.Next() {
		var cartItemID int64
		var ingr dishentity.Ingredient

		err := rows.Scan(&cartItemID, &ingr.ID, &ingr.Name, &ingr.ImageURL, &ingr.IsAllergen, &ingr.CategoryID)
		if err != nil {
			return nil, err
		}

		result[cartItemID] = append(result[cartItemID], ingr)
	}

	return result, nil
}

func (r *Repository) Clear(ctx context.Context, userID int64) error {
	var carItemIDs []int64
	rows, err := r.pg.Pool.Query(ctx, `SELECT id FROM cart_items WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cartItemID int64
		err := rows.Scan(&cartItemID)
		if err != nil {
			return err
		}
		carItemIDs = append(carItemIDs, cartItemID)
	}
	_, err = r.pg.Pool.Exec(ctx, `DELETE FROM cart_items WHERE user_id=$1`, userID)
	for _, cartItemID := range carItemIDs {
		err = r.RemoveItem(ctx, cartItemID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) IncrementCartItemQuantity(ctx context.Context, cartItemID int64) (int32, error) {
	var newQuantity int32

	err := r.pg.Pool.QueryRow(ctx, `
		UPDATE cart_items
		SET quantity = quantity + 1
		WHERE id = $1
		RETURNING quantity
	`, cartItemID).Scan(&newQuantity)

	if err != nil {
		return 0, fmt.Errorf("failed to increment quantity for cart item %d: %w", cartItemID, err)
	}

	return newQuantity, nil
}

func (r *Repository) DecrementCartItemQuantity(ctx context.Context, cartItemID int64) (int32, error) {
	var quantity int32
	err := r.pg.Pool.QueryRow(ctx, `
		SELECT quantity FROM cart_items where id = $1
	`, cartItemID).Scan(&quantity)
	if err != nil {
		return 0, fmt.Errorf("failed to get quantity for cart item %d: %w", cartItemID, err)
	}
	if quantity <= 1 {
		return 0, r.RemoveItem(ctx, cartItemID)
	}
	var newQuantity int32
	err = r.pg.Pool.QueryRow(ctx, `
		UPDATE cart_items
		SET quantity = quantity - 1
		WHERE id = $1
		RETURNING quantity
	`, cartItemID).Scan(&newQuantity)

	if err != nil {
		return 0, fmt.Errorf("failed to decrement quantity for cart item %d: %w", cartItemID, err)
	}

	return newQuantity, nil
}

func (r *Repository) GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error) {
	rows, err := r.pg.Pool.Query(ctx, `
		SELECT
			ci.id, ci.quantity, ci.customer_notes,
			d.id, d.name, d.description, d.chef_id, d.image_url,
			ds.id, ds.dish_id, ds.label, ds.weight_value, ds.weight_unit, ds.price_value, ds.price_currency
		FROM orders_cart_items oci
		JOIN cart_items ci ON oci.cart_item_id = ci.id
		JOIN dishes d ON ci.dish_id = d.id
		JOIN dish_sizes ds ON ci.dish_size_id = ds.id
		WHERE oci.order_id = $1
	`, orderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return []cartentity.CartItem{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []cartentity.CartItem
	cartItemIDs := make([]int64, 0)

	for rows.Next() {
		var item cartentity.CartItem
		var dish dishentity.Dish
		var size dishentity.Size

		err := rows.Scan(
			&item.ID, &item.Quantity, &item.Notes,
			&dish.ID, &dish.Name, &dish.Description, &dish.ChefID, &dish.ImageURL,
			&size.ID, &size.DishID, &size.Label, &size.WeightValue, &size.WeightUnit, &size.PriceValue, &size.PriceCurrency,
		)
		if err != nil {
			return nil, err
		}

		item.Dish = dish
		item.Size = size
		cartItemIDs = append(cartItemIDs, item.ID)
		items = append(items, item)
	}

	// ingredients (added and removed)
	addedMap, err := r.getCartItemIngredientsMap(ctx, "cart_item_added_ingredients", cartItemIDs)
	if err != nil {
		return nil, err
	}
	removedMap, err := r.getCartItemIngredientsMap(ctx, "cart_item_removed_ingredients", cartItemIDs)
	if err != nil {
		return nil, err
	}

	// заполняем ингредиенты
	for i, item := range items {
		items[i].AddedIngredients = addedMap[item.ID]
		items[i].RemovedIngredients = removedMap[item.ID]
	}

	return items, nil
}
