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

func (u *Usecase) OpenShift(ctx context.Context, chefID int64) error {
	return u.repo.OpenShift(ctx, chefID)
}

func (u *Usecase) CloseShift(ctx context.Context, chefID int64) error {
	return u.repo.CloseActiveShiftByChefID(ctx, chefID)
}

func (u *Usecase) GetDailyProfits(ctx context.Context, chefID int64) ([]shifts.DailyProfit, error) {
	return u.repo.GetDailyProfits(ctx, chefID)
}
