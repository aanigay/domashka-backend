package shifts

import (
	"context"
	
	"domashka-backend/internal/entity/shifts"
)

type Usecase struct {
	repo ShiftsRepo
}

func New(repo ShiftsRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetActiveShiftByChefID(ctx context.Context, chefID int64) (*shifts.Shift, error) {
	return u.repo.GetActiveShiftByChefID(ctx, chefID)
}
