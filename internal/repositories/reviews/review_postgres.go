// package review

package review

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	reviewEntity "domashka-backend/internal/entity/reviews"
	"domashka-backend/pkg/postgres"
)

type Repository struct {
	pg *postgres.Postgres
}

var ErrReviewNotFound = errors.New("review not found")

// New создает новый репозиторий для отзывов.
func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg: pg}
}

// GetByID возвращает отзыв по его ID.
func (r *Repository) GetByID(ctx context.Context, id int64) (*reviewEntity.Review, error) {
	const q = `
        SELECT
            id,
            chef_id,
            user_id,
            stars,
            comment,
            is_verified,
            include_in_rating,
            is_deleted,
            created_at,
            updated_at,
            order_id
        FROM reviews
        WHERE id = $1
    `
	row := r.pg.Pool.QueryRow(ctx, q, id)

	var rv reviewEntity.Review
	err := row.Scan(
		&rv.ID,
		&rv.ChefID,
		&rv.UserID,
		&rv.Stars,
		&rv.Comment,
		&rv.IsVerified,
		&rv.IncludeInRating,
		&rv.IsDeleted,
		&rv.CreatedAt,
		&rv.UpdatedAt,
		&rv.OrderID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return &rv, nil
}

// ListByChef возвращает список отзывов по ID шеф-повара.
// Если includeOnly=true, возвращаются только отзывы с include_in_rating = true.
// Если limit указан, к запросу добавляется LIMIT.
func (r *Repository) ListByChef(ctx context.Context, chefID int64, includeOnly bool, limit *int) ([]reviewEntity.Review, error) {
	base := `
        SELECT
            id,
            chef_id,
            user_id,
            stars,
            comment,
            is_verified,
            include_in_rating,
            is_deleted,
            created_at,
            updated_at,
            order_id
        FROM reviews
        WHERE chef_id = $1
    `
	args := []interface{}{chefID}
	if includeOnly {
		base += " AND include_in_rating = true"
	}
	base += " ORDER BY updated_at DESC"
	if limit != nil {
		base += fmt.Sprintf(" LIMIT %d", *limit)
	}

	rows, err := r.pg.Pool.Query(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []reviewEntity.Review
	for rows.Next() {
		var rv reviewEntity.Review
		if err := rows.Scan(
			&rv.ID,
			&rv.ChefID,
			&rv.UserID,
			&rv.Stars,
			&rv.Comment,
			&rv.IsVerified,
			&rv.IncludeInRating,
			&rv.IsDeleted,
			&rv.CreatedAt,
			&rv.UpdatedAt,
			&rv.OrderID,
		); err != nil {
			return nil, err
		}
		result = append(result, rv)
	}
	return result, nil
}

// Create вставляет новый отзыв и возвращает его ID и timestamps.
func (r *Repository) Create(ctx context.Context, rv *reviewEntity.Review) error {
	const q = `
        INSERT INTO reviews
            (chef_id, user_id, stars, comment, is_verified, include_in_rating, is_deleted, order_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `
	err := r.pg.Pool.QueryRow(ctx, q,
		rv.ChefID,
		rv.UserID,
		rv.Stars,
		rv.Comment,
		rv.IsVerified,
		rv.IncludeInRating,
		rv.IsDeleted,
		rv.OrderID,
	).Scan(&rv.ID, &rv.CreatedAt, &rv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create review: %w", err)
	}
	return nil
}

// Update обновляет все изменяемые поля отзыва (кроме created_at).
func (r *Repository) Update(ctx context.Context, rv *reviewEntity.Review) error {
	const q = `
        UPDATE reviews SET
            stars            = $1,
            comment          = $2,
            is_verified      = $3,
            include_in_rating= $4,
            is_deleted       = $5,
            updated_at       = NOW()
        WHERE id = $6
    `
	cmdTag, err := r.pg.Pool.Exec(ctx, q,
		rv.Stars,
		rv.Comment,
		rv.IsVerified,
		rv.IncludeInRating,
		rv.IsDeleted,
		rv.ID,
	)
	if err != nil {
		return fmt.Errorf("update review: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrReviewNotFound
	}
	return nil
}

// SoftDelete помечает отзыв как удаленный (is_deleted = true).
func (r *Repository) SoftDelete(ctx context.Context, id int64) error {
	const q = `
        UPDATE reviews
        SET is_deleted = true,
            updated_at = NOW()
        WHERE id = $1
    `
	cmdTag, err := r.pg.Pool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("soft delete review: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrReviewNotFound
	}
	return nil
}

// GetReviewByOrderAndUserID возвращает отзыв по связке order_id и user_id.
// Если отзыв не найден — возвращает ErrReviewNotFound.
func (r *Repository) GetReviewByOrderAndUserID(ctx context.Context, orderID, userID int64) (*reviewEntity.Review, error) {
	const q = `
			SELECT
				r.id,
				r.chef_id,
				r.user_id,
				r.stars,
				r.comment,
				r.is_verified,
				r.include_in_rating,
				r.is_deleted,
				r.created_at,
				r.updated_at,
				r.order_id
			FROM reviews r
			JOIN orders o
			  ON r.order_id = o.id
			WHERE r.order_id = $1
			  AND r.user_id = $2
			ORDER BY r.updated_at
    		LIMIT 1
`

	row := r.pg.Pool.QueryRow(ctx, q, orderID, userID)

	var rv reviewEntity.Review
	err := row.Scan(
		&rv.ID,
		&rv.ChefID,
		&rv.UserID,
		&rv.Stars,
		&rv.Comment,
		&rv.IsVerified,
		&rv.IncludeInRating,
		&rv.IsDeleted,
		&rv.CreatedAt,
		&rv.UpdatedAt,
		&rv.OrderID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return &rv, nil
}
