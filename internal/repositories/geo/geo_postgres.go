package geo

import (
	"context"
	"fmt"

	addressentity "domashka-backend/internal/entity/geo"
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

func (r *Repository) AddClientAddress(ctx context.Context, clientID int, address addressentity.Address) error {
	_, err := r.pg.Pool.Exec(ctx,
		`INSERT INTO client_addresses 
		 (client_id, address_type, name, full_address, comment, geom) 
		 VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7),4326)::geography)`,
		clientID,
		address.AddressType,
		address.Name,
		address.Address,
		address.Comment,
		address.Longitude,
		address.Latitude,
	)
	return err
}

func (r *Repository) AddChefAddress(ctx context.Context, chefID int, address addressentity.Address) error {
	_, err := r.pg.Pool.Exec(ctx,
		`INSERT INTO chef_addresses 
		 (chef_id, full_address, comment, geom) 
		 VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5),4326)::geography)`,
		chefID,
		address.Address,
		address.Comment,
		address.Longitude,
		address.Latitude,
	)
	return err
}

func (r *Repository) GetClientAddresses(ctx context.Context, clientID int) ([]addressentity.Address, error) {
	rows, err := r.pg.Pool.Query(ctx,
		`SELECT ST_Y(geom) AS latitude, ST_X(geom) AS longitude, address_type, name, full_address, comment 
		 FROM client_addresses 
		 WHERE client_id = $1`,
		clientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []addressentity.Address
	for rows.Next() {
		var address addressentity.Address
		if err := rows.Scan(
			&address.Latitude,
			&address.Longitude,
			&address.AddressType,
			&address.Name,
			&address.Address,
			&address.Comment,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (r *Repository) GetChefAddress(ctx context.Context, chefID int) (addressentity.Address, error) {
	var address addressentity.Address
	err := r.pg.Pool.QueryRow(ctx,
		`SELECT ST_Y(geom) AS latitude, ST_X(geom) AS longitude, full_address, comment 
		 FROM chef_addresses 
		 WHERE chef_id = $1 LIMIT 1`,
		chefID,
	).Scan(
		&address.Latitude,
		&address.Longitude,
		&address.Address,
		&address.Comment,
	)
	if err != nil {
		return addressentity.Address{}, err
	}
	return address, nil
}

func (r *Repository) UpdateClientAddress(ctx context.Context, clientID int, addressID int, address addressentity.Address) error {
	_, err := r.pg.Pool.Exec(ctx,
		`UPDATE client_addresses 
		 SET address_type = $1, name = $2, full_address = $3, comment = $4,
		     geom = ST_SetSRID(ST_MakePoint($5, $6),4326)::geography,
		     updated_at = NOW()
		 WHERE client_id = $7 AND id = $8`,
		address.AddressType,
		address.Name,
		address.Address,
		address.Comment,
		address.Longitude, // (longitude, latitude)
		address.Latitude,
		clientID,
		addressID,
	)
	return err
}

func (r *Repository) UpdateChefAddress(ctx context.Context, chefID int, address addressentity.Address) error {
	_, err := r.pg.Pool.Exec(ctx,
		`UPDATE chef_addresses 
		 SET full_address = $1, comment = $2,
		     geom = ST_SetSRID(ST_MakePoint($3, $4),4326)::geography,
		     updated_at = NOW()
		 WHERE chef_id = $5`,
		address.Address,
		address.Comment,
		address.Longitude,
		address.Latitude,
		chefID,
	)
	return err
}

func (r *Repository) GetChefsAddrByRange(ctx context.Context, clientAddressID int, radius float64) ([]addressentity.Address, error) {
	var clientAddr struct {
		Latitude  float64
		Longitude float64
	}

	err := r.pg.Pool.QueryRow(ctx,
		`SELECT ST_Y(geom) AS latitude, ST_X(geom) AS longitude 
		 FROM client_addresses 
		 WHERE id = $1`,
		clientAddressID,
	).Scan(&clientAddr.Latitude, &clientAddr.Longitude)
	if err != nil {
		return nil, fmt.Errorf("could not find client geo: %w", err)
	}

	rangeMeters := radius * 1000

	rows, err := r.pg.Pool.Query(ctx,
		`SELECT ST_Y(geom) AS latitude, ST_X(geom) AS longitude, full_address, comment
		 FROM chef_addresses
		 WHERE ST_DWithin(
		     geom,
		     ST_SetSRID(ST_MakePoint($1, $2),4326)::geography,
		     $3
		 )`,
		clientAddr.Longitude, // (longitude, latitude)
		clientAddr.Latitude,
		rangeMeters,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get chef addresses within range: %w", err)
	}
	defer rows.Close()

	var addresses []addressentity.Address
	for rows.Next() {
		var address addressentity.Address
		if err := rows.Scan(
			&address.Latitude,
			&address.Longitude,
			&address.Address,
			&address.Comment,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (r *Repository) GetClientsAddrByRange(ctx context.Context, chefID int, radius float64) ([]addressentity.Address, error) {
	var chefAddr struct {
		Latitude  float64
		Longitude float64
	}
	err := r.pg.Pool.QueryRow(ctx,
		`SELECT ST_Y(geom) AS latitude, ST_X(geom) AS longitude 
		 FROM chef_addresses 
		 WHERE chef_id = $1`,
		chefID,
	).Scan(&chefAddr.Latitude, &chefAddr.Longitude)
	if err != nil {
		return nil, fmt.Errorf("could not find chef geo: %w", err)
	}

	rangeMeters := radius * 1000

	rows, err := r.pg.Pool.Query(ctx,
		`SELECT ST_Y(geom) AS latitude, ST_X(geom) AS longitude, address_type, name, full_address, comment
		 FROM client_addresses
		 WHERE ST_DWithin(
		     geom,
		     ST_SetSRID(ST_MakePoint($1, $2),4326)::geography,
		     $3
		 )`,
		chefAddr.Longitude, // (longitude, latitude)
		chefAddr.Latitude,
		rangeMeters,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get client addresses within range: %w", err)
	}
	defer rows.Close()

	var addresses []addressentity.Address
	for rows.Next() {
		var address addressentity.Address
		if err := rows.Scan(
			&address.Latitude,
			&address.Longitude,
			&address.AddressType,
			&address.Name,
			&address.Address,
			&address.Comment,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}
