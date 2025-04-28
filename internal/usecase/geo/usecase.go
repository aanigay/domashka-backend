package geo

import (
	"context"
	"domashka-backend/internal/custom_errors"
	addressentity "domashka-backend/internal/entity/geo"
	"domashka-backend/internal/utils/validation"
)

type UseCase struct {
	repo GeoRepository
}

func New(repo GeoRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) PushClientAddress(ctx context.Context, addressID int64) error {
	return u.repo.PushClientAddress(ctx, addressID)
}
func (u *UseCase) GetAddressByID(ctx context.Context, id int64) (*addressentity.Address, error) {
	return u.repo.GetAddressByID(ctx, id)
}

func (s *UseCase) AddClientAddress(ctx context.Context, clientID int, address addressentity.Address) error {
	return s.repo.AddClientAddress(ctx, clientID, address)
}

func (s *UseCase) AddChefAddress(ctx context.Context, chefID int, address addressentity.Address) error {
	// Валидация адреса
	if !validation.IsAddressInRussia(address) {
		return custom_errors.ErrAddressNotInRussia
	}

	return s.repo.AddChefAddress(ctx, chefID, address)
}

func (s *UseCase) UpdateClientAddress(ctx context.Context, clientID int, addressID int, address addressentity.Address) error {
	//if !validation.IsAddressInRussia(address) {
	//	return custom_errors.ErrAddressNotInRussia
	//}

	return s.repo.UpdateClientAddress(ctx, clientID, addressID, address)

}

func (s *UseCase) GetChefAddress(ctx context.Context, chefID int) (addressentity.Address, error) {
	return s.repo.GetChefAddress(ctx, chefID)

}

func (s *UseCase) GetClientAddresses(ctx context.Context, clientID int) ([]addressentity.Address, error) {
	return s.repo.GetClientAddresses(ctx, clientID)
}

func (s *UseCase) UpdateChefAddress(ctx context.Context, chefID int, address addressentity.Address) error {
	if !validation.IsAddressInRussia(address) {
		return custom_errors.ErrAddressNotInRussia
	}

	return s.repo.UpdateChefAddress(ctx, chefID, address)
}

func (s *UseCase) FindChefsNearAddress(ctx context.Context, clientAddressID int, radius float64) ([]addressentity.Address, error) {
	return s.repo.GetChefsAddrByRange(ctx, clientAddressID, radius)
}

func (s *UseCase) FindClientsNearAddress(ctx context.Context, chefID int, radius float64) ([]addressentity.Address, error) {
	return s.repo.GetClientsAddrByRange(ctx, chefID, radius)
}

func (s *UseCase) GetLastUpdatedClientAddress(ctx context.Context, clientID int64) (*addressentity.Address, error) {
	return s.repo.GetLastUpdatedClientAddress(ctx, clientID)
}
