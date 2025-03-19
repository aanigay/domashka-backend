package notifications

import (
	"context"
	"database/sql"
	"domashka-backend/internal/entity/notifications"
	"domashka-backend/pkg/postgres"
	"errors"
	"fmt"
	"log"
)

type Repository struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{
		pg: pg,
	}
}

func (r *Repository) CreateNotification(ctx context.Context, n notifications.Notification) (int, error) {
	var ID int
	err := r.pg.Pool.QueryRow(ctx, `
		INSERT INTO notifications (user_id, channel, scenario, subject, message, recipient, status, created_at, updated_at, send_attempts, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, 'created', NOW(), NOW(), 0, $7)
		RETURNING id`,
		n.UserID, n.Channel, n.Scenario, n.Subject, n.Message, n.Recipient, n.Metadata).
		Scan(&ID)

	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (r *Repository) GetNotifications(ctx context.Context, filters map[string]string, page, limit int) ([]notifications.Notification, int, error) {
	var conditions []string
	var params []interface{}
	query := "SELECT id, user_id, channel, scenario, subject, message, recipient, status, created_at, updated_at, send_attempts, metadata FROM notifications"

	i := 1
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, i))
		params = append(params, value)
		i++
	}

	// WHERE если есть фильтры
	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for j := 1; j < len(conditions); j++ {
			query += " AND " + conditions[j]
		}
	}

	// Пагинация
	offset := (page - 1) * limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", i, i+1)
	params = append(params, limit, offset)

	rows, err := r.pg.Pool.Query(ctx, query, params...)
	if err != nil {
		log.Printf("Ошибка запроса списка уведомлений: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	// Заполняем результат
	var allNotifications []notifications.Notification
	for rows.Next() {
		var n notifications.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Channel, &n.Scenario, &n.Subject, &n.Message, &n.Recipient, &n.Status, &n.CreatedAt, &n.UpdatedAt, &n.SendAttempts, &n.Metadata)
		if err != nil {
			log.Printf("Ошибка чтения строки: %v", err)
			continue
		}
		allNotifications = append(allNotifications, n)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM notifications"
	if len(conditions) > 0 {
		countQuery += " WHERE " + conditions[0]
		for j := 1; j < len(conditions); j++ {
			countQuery += " AND " + conditions[j]
		}
	}
	err = r.pg.Pool.QueryRow(ctx, countQuery, params[:len(params)-2]...).Scan(&total)
	if err != nil {
		log.Printf("Ошибка при подсчете записей: %v", err)
	}

	return allNotifications, total, nil
}

func (r *Repository) GetNotificationByID(ctx context.Context, id int) (*notifications.Notification, error) {
	var n notifications.Notification
	err := r.pg.Pool.QueryRow(ctx, `
		SELECT id, user_id, channel, scenario, subject, message, recipient, status, created_at, updated_at, send_attempts, metadata
		FROM notifications WHERE id = $1`, id).
		Scan(&n.ID, &n.UserID, &n.Channel, &n.Scenario, &n.Subject, &n.Message, &n.Recipient, &n.Status, &n.CreatedAt, &n.UpdatedAt, &n.SendAttempts, &n.Metadata)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Printf("Ошибка получения уведомления ID=%d: %v", id, err)
		return nil, err
	}

	return &n, nil
}

func (r *Repository) UpdateNotification(ctx context.Context, id int, n notifications.Notification) error {
	log.Printf("DEBUG: Обновление уведомления id=%d, новый статус=%s, send_attempts=%d", id, n.Status, n.SendAttempts)
	_, err := r.pg.Pool.Exec(ctx,
		`UPDATE notifications
		SET user_id = $1, channel = $2, scenario = $3, subject = $4, message = $5, 
		    recipient = $6, status = $7, updated_at = NOW(), 
		    send_attempts = send_attempts + $8, metadata = $9
		WHERE id = $10`,
		n.UserID, n.Channel, n.Scenario, n.Subject, n.Message, n.Recipient, n.Status, n.SendAttempts, n.Metadata, id)

	if err != nil {
		log.Printf("ERROR: Ошибка обновления уведомления ID=%d: %v", id, err)
		return err
	}

	log.Printf("DEBUG: Уведомление id=%d успешно обновлено", id)
	return nil
}
