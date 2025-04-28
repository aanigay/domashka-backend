package shifts

import (
	"context"
	"domashka-backend/internal/entity/shifts"

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

func (r *Repository) GetActiveShiftIDByChefID(ctx context.Context, chefID int64) (int64, error) {
	var shiftID int64
	err := r.pg.Pool.QueryRow(ctx, "SELECT id FROM shifts WHERE chef_id = $1 AND is_active = true", chefID).Scan(&shiftID)
	if err != nil {
		return 0, err
	}
	return shiftID, nil
}

func (r *Repository) GetActiveShiftByChefID(ctx context.Context, chefID int64) (*shifts.Shift, error) {
	var shift shifts.Shift
	err := r.pg.Pool.QueryRow(ctx, "SELECT id, chef_id, is_active, created_at, closed_at, total_profit FROM shifts WHERE chef_id = $1 AND is_active = true", chefID).Scan(
		&shift.ID,
		&shift.ChefID,
		&shift.IsActive,
		&shift.CreatedAt,
		&shift.ClosedAt,
		&shift.TotalProfit,
	)
	if err != nil {
		return nil, err
	}
	return &shift, nil
}
