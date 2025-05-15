package favorites

import (
	"context"
)

type Usecase struct {
	favRepo favRepo
}

func New(favRepo favRepo) *Usecase {
	return &Usecase{
		favRepo: favRepo,
	}
}

func (u *Usecase) AddFavoriteChef(ctx context.Context, userID, chefID int64) error {
	return u.favRepo.AddFavoriteChef(ctx, userID, chefID)
}

func (u *Usecase) RemoveFavoriteChef(ctx context.Context, userID, chefID int64) error {
	return u.favRepo.RemoveFavoriteChef(ctx, userID, chefID)
}

func (u *Usecase) AddFavoriteDish(ctx context.Context, userID, dishID int64) error {
	return u.favRepo.AddFavoriteDish(ctx, userID, dishID)
}

func (u *Usecase) RemoveFavoriteDish(ctx context.Context, userID, dishID int64) error {
	return u.favRepo.RemoveFavoriteDish(ctx, userID, dishID)
}
