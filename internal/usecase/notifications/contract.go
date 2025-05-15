package notifications

import (
	"context"
	"domashka-backend/internal/entity/notifications"
)

//go:generate mockgen -source=contract.go -destination contract_mocks_test.go -package $GOPACKAGE

type NotificationRepository interface {
	CreateNotification(ctx context.Context, n notifications.Notification) (int, error)
	GetNotifications(ctx context.Context, filters map[string]string, page, limit int) ([]notifications.Notification, int, error)
	GetNotificationByID(ctx context.Context, id int) (*notifications.Notification, error)
	UpdateNotification(ctx context.Context, id int, n notifications.Notification) error
}

type SMTPClient interface {
	SendEmail(ctx context.Context, to, subject string, data notifications.EmailData) (int, error)
}
