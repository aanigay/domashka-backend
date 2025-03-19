package notifications

import (
	"context"
	notifEntity "domashka-backend/internal/entity/notifications"
	"log"
)

type NotificationUsecase struct {
	notificationRepo NotificationRepository
	smtpClient       SMTPClient
}

func New(notificationRepo NotificationRepository, smtp SMTPClient) *NotificationUsecase {
	return &NotificationUsecase{
		notificationRepo: notificationRepo,
		smtpClient:       smtp,
	}
}

func (u *NotificationUsecase) CreateNotification(ctx context.Context, n notifEntity.Notification) (int, error) {
	id, err := u.notificationRepo.CreateNotification(ctx, n)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *NotificationUsecase) SendEmailNotification(ctx context.Context, n notifEntity.Notification) error {
	log.Printf("DEBUG: Запуск SendEmailNotification для уведомления id=%d, получатель=%s", n.ID, n.Recipient)

	attempts, err := u.smtpClient.SendEmail(ctx, n.Recipient, n.Subject.String, notifEntity.EmailData{
		Title:    n.Subject.String,
		UserName: n.Recipient,
		Body:     n.Message,
		Footer:   "Footer",
	})
	if err != nil {
		log.Printf("DEBUG: Ошибка отправки email для уведомления id=%d: %v", n.ID, err)
		// При ошибке отправки обновляем статус уведомления на error
		if updateErr := u.notificationRepo.UpdateNotification(ctx, n.ID, notifEntity.Notification{
			Status:       notifEntity.StatusError,
			SendAttempts: attempts,
		}); updateErr != nil {
			log.Printf("DEBUG: Ошибка обновления статуса уведомления на error для id=%d: %v", n.ID, updateErr)
			return updateErr
		}
		return err
	}

	log.Printf("DEBUG: Email успешно отправлен для уведомления id=%d с попытками=%d", n.ID, attempts)
	// Обновляем статус уведомления на sent
	if updateErr := u.notificationRepo.UpdateNotification(ctx, n.ID, notifEntity.Notification{
		Status:       notifEntity.StatusSent,
		SendAttempts: attempts,
	}); updateErr != nil {
		log.Printf("DEBUG: Ошибка обновления статуса уведомления на sent для id=%d: %v", n.ID, updateErr)
		return updateErr
	}

	log.Printf("DEBUG: Уведомление id=%d успешно обновлено после отправки email", n.ID)
	return nil
}

func (u *NotificationUsecase) GetNotifications(ctx context.Context, filters map[string]string, page, limit int) ([]notifEntity.Notification, int, error) {
	notifications, total, err := u.notificationRepo.GetNotifications(ctx, filters, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (u *NotificationUsecase) GetNotificationByID(ctx context.Context, id int) (*notifEntity.Notification, error) {
	n, err := u.notificationRepo.GetNotificationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (u *NotificationUsecase) ResendNotification(ctx context.Context, id int) error {
	n, err := u.notificationRepo.GetNotificationByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.SendEmailNotification(ctx, *n)
	if err != nil {
		return err
	}

	return nil
}
