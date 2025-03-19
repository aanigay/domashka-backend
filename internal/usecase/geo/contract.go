package geo

import (
	"context"
	addressentity "domashka-backend/internal/entity/geo"
)

type GeoRepository interface {
	AddClientAddress(ctx context.Context, clientID int, address addressentity.Address) error
	AddChefAddress(ctx context.Context, chefID int, address addressentity.Address) error
	GetClientAddresses(ctx context.Context, clientID int) ([]addressentity.Address, error)
	GetChefAddress(ctx context.Context, chefID int) (addressentity.Address, error)
	UpdateClientAddress(ctx context.Context, clientID int, addressID int, address addressentity.Address) error
	UpdateChefAddress(ctx context.Context, chefID int, address addressentity.Address) error
	GetChefsAddrByRange(ctx context.Context, clientAddressID int, radius float64) ([]addressentity.Address, error)
	GetClientsAddrByRange(ctx context.Context, chefID int, radius float64) ([]addressentity.Address, error)
}
