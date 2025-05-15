package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	cartentity "domashka-backend/internal/entity/cart"
	dishentity "domashka-backend/internal/entity/dishes"
	"domashka-backend/internal/entity/orders"
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

func (r *Repository) CreateOrder(
	ctx context.Context,
	chefID int64,
	shiftID int64,
	userID int64,
	clientAddressID int64,
	totalCost float32,
	leaveByTheDoor bool,
	callBeforehand bool,
) (int64, error) {
	var orderID int64
	err := r.pg.Pool.QueryRow(ctx, `
		INSERT INTO orders (chef_id, shift_id, total_cost, leave_by_the_door, call_beforehand, client_address_id, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, chefID, shiftID, totalCost, leaveByTheDoor, callBeforehand, clientAddressID, userID).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *Repository) AddCartItemToOrder(ctx context.Context, cartItem *cartentity.CartItem, orderID, userID int64) error {
	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	var orderedItemID int64
	err = r.pg.Pool.QueryRow(ctx, `INSERT INTO ordered_items
    	(user_id, dish_id, chef_id, order_id, dish_size_id, quantity, customer_notes)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`,
		userID, cartItem.Dish.ID, cartItem.Dish.ChefID, orderID, cartItem.Size.ID, cartItem.Quantity, cartItem.Notes).
		Scan(&orderedItemID)
	if err != nil {
		return err
	}
	// Добавляем добавленные ингредиенты
	for _, ingredient := range cartItem.AddedIngredients {
		_, err := tx.Exec(ctx, `
			INSERT INTO ordered_item_added_ingredients (ordered_item_id, ingredient_id)
			VALUES ($1, $2)
		`, orderedItemID, ingredient.ID)
		if err != nil {
			return fmt.Errorf("failed to insert added ingredient %d: %w", ingredient.ID, err)
		}
	}

	// Добавляем удалённые ингредиенты
	for _, ingredient := range cartItem.RemovedIngredients {
		_, err := tx.Exec(ctx, `
			INSERT INTO ordered_item_added_ingredients (ordered_item_id, ingredient_id)
			VALUES ($1, $2)
		`, orderedItemID, ingredient.ID)
		if err != nil {
			return fmt.Errorf("insert failed removed ingredient %d: %w", ingredient.ID, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return err
}

func (r *Repository) GetCartItems(ctx context.Context, userID int64) ([]cartentity.CartItem, error) {
	rows, err := r.pg.Pool.Query(ctx, `
		SELECT
			ci.id, ci.quantity, ci.customer_notes,
			d.id, d.name, d.description, d.chef_id, d.image_url,
			ds.id, ds.dish_id, ds.label, ds.weight_value, ds.weight_unit, ds.price_value, ds.price_currency
		FROM ordered_items ci
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
	orderedItemIDs := make([]int64, 0)

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
		orderedItemIDs = append(orderedItemIDs, item.ID)
		items = append(items, item)
	}

	// ingredients (added and removed)
	addedMap, err := r.getCartItemIngredientsMap(ctx, "ordered_item_added_ingredients", orderedItemIDs)
	if err != nil {
		return nil, err
	}
	removedMap, err := r.getCartItemIngredientsMap(ctx, "ordered_item_removed_ingredients", orderedItemIDs)
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

func (r *Repository) GetCartItemsByOrderID(ctx context.Context, orderID int64) ([]cartentity.CartItem, error) {
	rows, err := r.pg.Pool.Query(ctx, `
		SELECT
			ci.id, ci.quantity, ci.customer_notes,
			d.id, d.name, d.description, d.chef_id, d.image_url,
			ds.id, ds.dish_id, ds.label, ds.weight_value, ds.weight_unit, ds.price_value, ds.price_currency
		FROM ordered_items ci
		JOIN dishes d ON ci.dish_id = d.id
		JOIN dish_sizes ds ON ci.dish_size_id = ds.id
		WHERE ci.order_id = $1
	`, orderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return []cartentity.CartItem{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []cartentity.CartItem
	orderedItemIDs := make([]int64, 0)

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
		orderedItemIDs = append(orderedItemIDs, item.ID)
		items = append(items, item)
	}

	// ingredients (added and removed)
	addedMap, err := r.getCartItemIngredientsMap(ctx, "ordered_item_added_ingredients", orderedItemIDs)
	if err != nil {
		return nil, err
	}
	removedMap, err := r.getCartItemIngredientsMap(ctx, "ordered_item_removed_ingredients", orderedItemIDs)
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

func (r *Repository) getCartItemIngredientsMap(ctx context.Context, tableName string, orderedItemIDs []int64) (map[int64][]dishentity.Ingredient, error) {
	if len(orderedItemIDs) == 0 {
		return map[int64][]dishentity.Ingredient{}, nil
	}

	rows, err := r.pg.Pool.Query(ctx, fmt.Sprintf(`
		SELECT
			ci.ordered_item_id,
			i.id, i.name, i.image_url, i.is_allergen, i.category_id
		FROM %s ci
		JOIN ingredients i ON ci.ingredient_id = i.id
		WHERE ci.ordered_item_id = ANY($1)
	`, tableName), orderedItemIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]dishentity.Ingredient)

	for rows.Next() {
		var orderedItemID int64
		var ingr dishentity.Ingredient

		err := rows.Scan(&orderedItemID, &ingr.ID, &ingr.Name, &ingr.ImageURL, &ingr.IsAllergen, &ingr.CategoryID)
		if err != nil {
			return nil, err
		}

		result[orderedItemID] = append(result[orderedItemID], ingr)
	}

	return result, nil
}

func (r *Repository) GetOrdersByShiftID(ctx context.Context, shiftID int64) ([]orders.Order, error) {
	query := `
        SELECT 
			id, status
        FROM 
            orders
        WHERE shift_id = $1
    `
	rows, err := r.pg.Pool.Query(ctx, query, shiftID)
	if errors.Is(err, pgx.ErrNoRows) {
		return []orders.Order{}, nil
	}
	if err != nil {
		return nil, err
	}
	var o []orders.Order
	for rows.Next() {
		var order orders.Order
		if err := rows.Scan(&order.ID, &order.Status); err != nil {
			return nil, err
		}
		o = append(o, order)
	}
	return o, nil
}

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status int32) error {
	_, err := r.pg.Pool.Exec(ctx, `UPDATE orders SET status = $1, updated_at = now() WHERE id = $2`, status, orderID)
	return err
}

func (r *Repository) GetShiftIDByOrderID(ctx context.Context, orderID int64) (int64, error) {
	query := `
        SELECT 
			shift_id
        FROM 
            orders
        WHERE id = $1
    `
	var shiftID int64
	err := r.pg.Pool.QueryRow(ctx, query, orderID).Scan(&shiftID)
	return shiftID, err
}

func (r *Repository) GetOrderByID(ctx context.Context, orderID int64) (*orders.Order, error) {
	query := `
        SELECT 
			id,
			shift_id,
			chef_id,
			status,
			total_cost,
			leave_by_the_door,
			client_address_id
        FROM 
            orders
        WHERE id = $1
    `
	var order orders.Order
	err := r.pg.Pool.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.ShiftID,
		&order.ChefID,
		&order.Status,
		&order.TotalCost,
		&order.LeaveByTheDoor,
		&order.ClientAddressID,
	)
	return &order, err
}

func (r *Repository) GetOrderedDishesIDsAndChefsIDs(ctx context.Context, userID int64, dishesLimit, chefsLimit int) ([]int64, []int64, error) {
	query := `
        SELECT 
			dish_id, chef_id
        FROM 
            ordered_items
        WHERE user_id = $1
        ORDER BY added_at DESC
    `
	chefsIDs := make([]int64, 0)
	dishesIDs := make([]int64, 0)
	rows, err := r.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []int64{}, []int64{}, nil
		}
		return nil, nil, err
	}
	dishesIDsSet := map[int64]bool{}
	chefsIDsSet := map[int64]bool{}
	for rows.Next() {
		var dishID int64
		var chefID int64
		err := rows.Scan(&dishID, &chefID)
		if err != nil {
			return nil, nil, err
		}
		dishesIDsSet[dishID] = true
		chefsIDsSet[chefID] = true
	}

	for dishID := range dishesIDsSet {
		dishesIDs = append(dishesIDs, dishID)
		if len(dishesIDs) >= dishesLimit {
			break
		}
	}
	for chefID := range chefsIDsSet {
		chefsIDs = append(chefsIDs, chefID)
		if len(chefsIDs) >= chefsLimit {
			break
		}
	}
	return dishesIDs, chefsIDs, nil
}

func (r *Repository) GetStatus(ctx context.Context, orderID int64) (int32, error) {
	query := `
        SELECT 
			status
        FROM 
            orders
        WHERE id = $1
    `
	var status int32
	err := r.pg.Pool.QueryRow(ctx, query, orderID).Scan(&status)
	return status, err
}
func (r *Repository) GetOrdersByUserID(ctx context.Context, userID int64) ([]orders.Order, error) {
	const query = `
        SELECT
            id,
            shift_id,
            chef_id,
            status,
            created_at,
            updated_at,
            total_cost,
            leave_by_the_door,
            client_address_id
        FROM orders
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetOrdersByUserID query: %w", err)
	}
	defer rows.Close()

	var result []orders.Order
	for rows.Next() {
		var o orders.Order
		if err := rows.Scan(
			&o.ID,
			&o.ShiftID,
			&o.ChefID,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.TotalCost,
			&o.LeaveByTheDoor,
			&o.ClientAddressID,
		); err != nil {
			return nil, fmt.Errorf("GetOrdersByUserID scan: %w", err)
		}
		result = append(result, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetOrdersByUserID rows: %w", err)
	}

	return result, nil
}

func (r *Repository) CountDishInOrders(ctx context.Context, dishID int64) (int, error) {
	var count int
	err := r.pg.Pool.QueryRow(ctx, `SELECT count(*) FROM ordered_items WHERE dish_id = $1`, dishID).Scan(&count)
	return count, err
}
