package users

import (
	"context"
	"database/sql"
	"domashka-backend/internal/custom_errors"
	"errors"
	"log"

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

// Create вставляет нового пользователя с новыми колонками: username, name и chat_id.
func (r *Repository) Create(ctx context.Context, user *usersEntity.User) error {
	_, err := r.pg.Pool.Exec(ctx, `
		INSERT INTO users (
			username, 
			name, 
			alias, 
			first_name, 
			second_name, 
			last_name, 
			email, 
			number_phone, 
			status, 
			external_type, 
			telegram_name, 
			chat_id,
			notification_flag, 
			role, 
			birthday
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`,
		user.Username,         // $1
		user.Name,             // $2
		user.Alias,            // $3
		user.FirstName,        // $4
		user.SecondName,       // $5
		user.LastName,         // $6
		user.Email,            // $7
		user.NumberPhone,      // $8
		user.Status,           // $9
		user.ExternalType,     // $10
		user.TelegramName,     // $11
		user.ChatID,           // $12
		user.NotificationFlag, // $13
		user.Role,             // $14
		user.Birthday,         // $15
	)
	if err != nil {
		return err
	}
	return nil
}

// GetByID возвращает пользователя по идентификатору.
func (r *Repository) GetByID(ctx context.Context, id int64) (*usersEntity.User, error) {
	var user usersEntity.User
	row := r.pg.Pool.QueryRow(ctx, `
		SELECT 
			id, 
			username, 
			name, 
			alias, 
			first_name, 
			second_name, 
			last_name, 
			email, 
			number_phone, 
			status, 
			external_type, 
			telegram_name, 
			external_id, 
			chat_id,
			notification_flag, 
			role, 
			birthday, 
			created_at, 
			updated_at 
		FROM users 
		WHERE id = $1
	`, id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Alias,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.NumberPhone,
		&user.Status,
		&user.ExternalType,
		&user.TelegramName,
		&user.ExternalID,
		&user.ChatID,
		&user.NotificationFlag,
		&user.Role,
		&user.Birthday,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update обновляет данные пользователя с учетом новых полей: name и chat_id.
func (r *Repository) Update(ctx context.Context, id int64, user usersEntity.User) error {
	query := `
		UPDATE users
		SET alias = $1,
		    name = $2,
		    first_name = $3,
		    second_name = $4,
		    last_name = $5,
		    email = $6,
		    number_phone = $7,
		    is_spam = $8,
		    sms_attempts = $9,
		    last_sms_request = $10,
		    status = $11,
		    external_type = $12,
		    telegram_name = $13,
		    external_id = $14,
		    chat_id = $15,
		    notification_flag = $16,
		    role = $17,
		    birthday = $18,
		    updated_at = NOW()
		WHERE id = $19
	`
	_, err := r.pg.Pool.Exec(ctx, query,
		user.Alias,            // $1
		user.Name,             // $2
		user.FirstName,        // $3
		user.SecondName,       // $4
		user.LastName,         // $5
		user.Email,            // $6
		user.NumberPhone,      // $7
		user.IsSpam,           // $8
		user.SMSAttempts,      // $9
		user.LastSMSRequest,   // $10
		user.Status,           // $11
		user.ExternalType,     // $12
		user.TelegramName,     // $13
		user.ExternalID,       // $14
		user.ChatID,           // $15
		user.NotificationFlag, // $16
		user.Role,             // $17
		user.Birthday,         // $18
		id,                    // $19
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.pg.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

// GetByPhone возвращает пользователя по номеру телефона.
func (r *Repository) GetByPhone(ctx context.Context, phone string) (*usersEntity.User, error) {
	var user usersEntity.User
	row := r.pg.Pool.QueryRow(ctx, `
		SELECT 
			id, 
			username, 
			name, 
			alias, 
			first_name, 
			second_name, 
			last_name, 
			email, 
			number_phone, 
			status, 
			external_type, 
			telegram_name, 
			external_id, 
			chat_id,
			notification_flag, 
			role, 
			birthday, 
			created_at, 
			updated_at 
		FROM users 
		WHERE number_phone = $1
	`, phone)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Alias,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.NumberPhone,
		&user.Status,
		&user.ExternalType,
		&user.TelegramName,
		&user.ExternalID,
		&user.ChatID,
		&user.NotificationFlag,
		&user.Role,
		&user.Birthday,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateWithPhone(ctx context.Context, phone string) (*usersEntity.User, error) {
	log.Printf("DEBUG: Начало CreateWithPhone для номера: %s", phone)
	var user usersEntity.User

	query := `
		INSERT INTO users (
			username, 
			name, 
			alias, 
			first_name, 
			second_name, 
			last_name, 
			email, 
			number_phone, 
			status, 
			external_type, 
			telegram_name, 
			chat_id,
			notification_flag, 
			role, 
			birthday
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING 
			id, 
			username, 
			name, 
			alias, 
			first_name, 
			second_name, 
			last_name, 
			email, 
			number_phone, 
			status, 
			external_type, 
			telegram_name, 
			external_id, 
			COALESCE(chat_id, '') as chat_id,
			notification_flag, 
			role, 
			birthday, 
			created_at, 
			updated_at
	`
	err := r.pg.Pool.QueryRow(ctx, query,
		phone,    // username
		"",       // name
		phone,    // alias
		"",       // first_name
		nil,      // second_name
		nil,      // last_name
		nil,      // email
		phone,    // number_phone
		0,        // status
		0,        // external_type
		nil,      // telegram_name
		nil,      // chat_id
		1,        // notification_flag
		"client", // role
		nil,      // birthday
	).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Alias,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.NumberPhone,
		&user.Status,
		&user.ExternalType,
		&user.TelegramName,
		&user.ExternalID,
		&user.ChatID,
		&user.NotificationFlag,
		&user.Role,
		&user.Birthday,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("DEBUG: Ошибка в CreateWithPhone для номера %s: %v", phone, err)
		return nil, err
	}
	log.Printf("DEBUG: Пользователь успешно создан для номера: %s, user: %+v", phone, user)
	return &user, nil
}
