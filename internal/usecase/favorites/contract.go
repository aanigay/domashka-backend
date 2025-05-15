package favorites

import (
	"context"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type favRepo interface {
	AddFavoriteChef(ctx context.Context, chefID, userID int64) error
	RemoveFavoriteChef(ctx context.Context, chefID, userID int64) error
	AddFavoriteDish(ctx context.Context, dishID, userID int64) error
	RemoveFavoriteDish(ctx context.Context, dishID, userID int64) error
}
