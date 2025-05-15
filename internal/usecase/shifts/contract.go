package shifts

import (
	"context"

	"domashka-backend/internal/entity/shifts"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type ShiftsRepo interface {
	GetActiveShiftByChefID(ctx context.Context, chefID int64) (*shifts.Shift, error)
	OpenShift(ctx context.Context, chefID int64) error
	CloseActiveShiftByChefID(ctx context.Context, chefID int64) error
	GetDailyProfits(ctx context.Context, id int64) ([]shifts.DailyProfit, error)
}
