package geo

import (
	"context"
	addressentity "domashka-backend/internal/entity/geo"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type GeoRepository interface {
	AddClientAddress(ctx context.Context, clientID int, address addressentity.Address) error
	AddChefAddress(ctx context.Context, chefID int, address addressentity.Address) error
	GetClientAddresses(ctx context.Context, clientID int) ([]addressentity.Address, error)
	GetChefAddress(ctx context.Context, chefID int) (addressentity.Address, error)
	UpdateClientAddress(ctx context.Context, clientID int, addressID int, address addressentity.Address) error
	UpdateChefAddress(ctx context.Context, chefID int, address addressentity.Address) error
	GetChefsAddrByRange(ctx context.Context, clientAddressID int, radius float64) ([]addressentity.Address, error)
	GetClientsAddrByRange(ctx context.Context, chefID int, radius float64) ([]addressentity.Address, error)
	GetLastUpdatedClientAddress(ctx context.Context, clientID int64) (*addressentity.Address, error)
	GetAddressByID(ctx context.Context, id int64) (*addressentity.Address, error)
	PushClientAddress(ctx context.Context, addressID int64) error
}
