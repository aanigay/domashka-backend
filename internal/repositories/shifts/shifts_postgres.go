package shifts

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	"domashka-backend/internal/entity/shifts"
	"domashka-backend/internal/utils/types"
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
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *Repository) AddToTotalProfit(ctx context.Context, shiftID int64, profit float64) error {
	_, err := r.pg.Pool.Exec(ctx, "UPDATE shifts SET total_profit = total_profit + $1 WHERE id = $2", profit, shiftID)
	return err
}

func (r *Repository) CloseActiveShiftByChefID(ctx context.Context, chefID int64) error {
	_, err := r.pg.Pool.Exec(ctx, "UPDATE shifts SET is_active = false, closed_at = now() WHERE chef_id = $1 AND is_active = true", chefID)
	return err
}

func (r *Repository) OpenShift(ctx context.Context, chefID int64) error {
	_, err := r.pg.Pool.Exec(ctx, "INSERT INTO shifts (chef_id, is_active) VALUES ($1, true)", chefID)
	return err
}

func (r *Repository) GetDailyProfits(ctx context.Context, chefID int64) ([]shifts.DailyProfit, error) {
	rows, err := r.pg.Pool.Query(ctx, `
		SELECT
			DATE(closed_at) AS day,
			SUM(total_profit) AS daily_profit
		FROM
			shifts
		WHERE
			chef_id = $1
		  AND is_active = false
		GROUP BY
			DATE(closed_at)
		ORDER BY
			DATE(closed_at);`, chefID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []shifts.DailyProfit{}, nil
		}
		return nil, err
	}
	dailyProfits := make([]shifts.DailyProfit, 0)
	for rows.Next() {
		var dailyProfit shifts.DailyProfit
		var dateTime time.Time
		if err := rows.Scan(&dateTime, &dailyProfit.Profit); err != nil {
			return nil, err
		}
		dailyProfit.Date = dateTime.Format("02.01")
		dailyProfit.Month = types.MonthEngToRus(dateTime.Format("January"))
		dailyProfits = append(dailyProfits, dailyProfit)
	}
	return dailyProfits, nil
}
