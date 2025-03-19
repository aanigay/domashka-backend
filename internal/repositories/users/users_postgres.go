package users

import (
	"context"
	"database/sql"

	"domashka-backend/internal/custom_errors"
	usersEntity "domashka-backend/internal/entity/users"
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

func (r *Repository) Create(ctx context.Context, user *usersEntity.User) (*string, error) {
	var id string
	err := r.pg.Pool.QueryRow(ctx, `
    INSERT INTO users (username, alias, first_name, second_name, last_name, email, number_phone, status, external_type, telegram_name, notification_flag, role, birthday)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    RETURNING id`,
		user.Username, user.Alias, user.FirstName, user.SecondName, user.LastName, user.Email, user.NumberPhone,
		user.Status, user.ExternalType, user.TelegramName, user.NotificationFlag, user.Role, user.Birthday).
		Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*usersEntity.User, error) {
	var user usersEntity.User
	row := r.pg.Pool.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(
		&user.ID, &user.Username, &user.Alias, &user.FirstName, &user.SecondName,
		&user.LastName, &user.Email, &user.NumberPhone, &user.Status, &user.ExternalType,
		&user.TelegramName, &user.ExternalID, &user.NotificationFlag, &user.Role, &user.Birthday,
		&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Update(ctx context.Context, id string, user usersEntity.User) error {
	_, err := r.pg.Pool.Exec(ctx, `
		UPDATE users
		SET alias = $1, first_name = $2, second_name = $3, last_name = $4, email = $5, 
		    number_phone = $6, status = $7, external_type = $8, role = $9
		WHERE id = $10`,
		user.Alias, user.FirstName, user.SecondName, user.LastName, user.Email, user.NumberPhone,
		user.Status, user.ExternalType, user.Role, id)

	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.pg.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*usersEntity.User, error) {
	var user usersEntity.User
	row := r.pg.Pool.QueryRow(ctx, "SELECT * FROM users WHERE number_phone = $1", phone)
	err := row.Scan(
		&user.ID, &user.Username, &user.Alias, &user.FirstName, &user.SecondName,
		&user.LastName, &user.Email, &user.NumberPhone, &user.Status, &user.ExternalType,
		&user.TelegramName, &user.ExternalID, &user.NotificationFlag, &user.Role, &user.Birthday,
		&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateWithPhone(ctx context.Context, phone string) error {
	_, err := r.pg.Pool.Exec(ctx, `
		INSERT INTO users (
			username, 
			alias, 
			first_name, 
			second_name, 
			last_name, 
			email, 
			number_phone, 
			status, 
			external_type, 
			telegram_name, 
			notification_flag, 
			role, 
			birthday
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`,
		phone,
		phone,
		"",
		nil,
		nil,
		nil,
		phone,
		0,
		0,
		nil,
		1,
		"client",
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
