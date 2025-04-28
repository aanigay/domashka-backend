package orders

import (
	"context"
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
	clientAddressID int64,
	totalCost float32,
	leaveByTheDoor bool,
	callBeforehand bool,
) (int64, error) {
	var orderID int64
	err := r.pg.Pool.QueryRow(ctx, `
		INSERT INTO orders (chef_id, shift_id, total_cost, leave_by_the_door, call_beforehand, client_address_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, chefID, shiftID, totalCost, leaveByTheDoor, callBeforehand, clientAddressID).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *Repository) AddCartItemToOrder(ctx context.Context, cartItemID int64, orderID int64) error {
	_, err := r.pg.Pool.Exec(ctx, "INSERT INTO orders_cart_items (cart_item_id, order_id) VALUES ($1, $2)", cartItemID, orderID)
	return err
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
