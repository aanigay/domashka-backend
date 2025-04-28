package shifts

import (
	"context"
	
	"domashka-backend/internal/entity/shifts"
)

type ShiftsRepo interface {
	GetActiveShiftByChefID(ctx context.Context, chefID int64) (*shifts.Shift, error)
}
