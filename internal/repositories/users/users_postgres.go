package users

import (
	"context"
	"database/sql"
	"domashka-backend/internal/custom_errors"
	chefEntity "domashka-backend/internal/entity/chefs"
	dishesEntity "domashka-backend/internal/entity/dishes"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
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
		if errors.Is(err, pgx.ErrNoRows) {
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
		    updated_at = NOW()
		WHERE id = $19
	`
	query = `UPDATE users SET`
	if user.Alias != "" {
		query += fmt.Sprintf(" alias = '%s',", user.Alias)
	}
	if user.Name != "" {
		query += fmt.Sprintf(" name = '%s',", user.Name)
	}
	if user.FirstName != "" {
		query += fmt.Sprintf(" first_name = '%s',", user.FirstName)
	}
	if user.LastName != nil {
		query += fmt.Sprintf(" last_name = '%s',", *user.LastName)
	}
	if user.Email != nil {
		query += fmt.Sprintf(" email = '%s',", *user.Email)
	}
	if user.NumberPhone != nil {
		query += fmt.Sprintf(" number_phone = '%s',", *user.NumberPhone)
	}
	query += " updated_at = NOW() WHERE id = $1"
	_, err := r.pg.Pool.Exec(ctx, query, id)
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

func (r *Repository) CheckIfUserIsChef(ctx context.Context, userID int64) (*int64, bool, error) {
	query := `SELECT chef_id FROM users_chefs WHERE user_id = $1`
	var chefID int64
	err := r.pg.Pool.QueryRow(ctx, query, userID).Scan(&chefID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &chefID, true, nil
}

// GetFavoritesDishesByUserID возвращает список любимых блюд пользователя.
func (r *Repository) GetFavoritesDishesByUserID(ctx context.Context, userID int64) ([]dishesEntity.Dish, error) {
	const query = `
        SELECT d.id, d.name, d.description, d.chef_id, d.image_url
        FROM user_favorite_dishes ufd
        JOIN dishes d ON ufd.dish_id = d.id
        WHERE ufd.user_id = $1
        ORDER BY ufd.created_at DESC
    `

	rows, err := r.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetFavoritesDishesByUserID query error: %w", err)
	}
	defer rows.Close()

	var dishes []dishesEntity.Dish
	for rows.Next() {
		var d dishesEntity.Dish
		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Description,
			&d.ChefID,
			&d.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("GetFavoritesDishesByUserID scan error: %w", err)
		}
		dishes = append(dishes, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFavoritesDishesByUserID rows error: %w", err)
	}

	return dishes, nil
}

// GetFavoritesChefsByUserID возвращает список любимых шеф-поваров пользователя.
func (r *Repository) GetFavoritesChefsByUserID(ctx context.Context, userID int64) ([]chefEntity.Chef, error) {
	const query = `
        SELECT c.id, c.name, c.image_url, c.small_image_url, c.description
        FROM user_favorite_chefs ufc
        JOIN chefs c ON ufc.chef_id = c.id
        WHERE ufc.user_id = $1
        ORDER BY ufc.created_at DESC
    `

	rows, err := r.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetFavoritesChefsByUserID query error: %w", err)
	}
	defer rows.Close()

	var chefs []chefEntity.Chef
	for rows.Next() {
		var c chefEntity.Chef
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.ImageURL,
			&c.SmallImageURL,
			&c.Description,
		); err != nil {
			return nil, fmt.Errorf("GetFavoritesChefsByUserID scan error: %w", err)
		}
		chefs = append(chefs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFavoritesChefsByUserID rows error: %w", err)
	}

	return chefs, nil
}
