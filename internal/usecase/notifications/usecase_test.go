package notifications

import (
	"context"
	"domashka-backend/internal/entity/notifications"
	"testing"
)

func TestNotificationUsecase_CreateNotification(t *testing.T) {
	type fields struct {
		notificationRepo NotificationRepository
		smtpClient       SMTPClient
	}
	type args struct {
		ctx context.Context
		n   notifications.Notification
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &NotificationUsecase{
				notificationRepo: tt.fields.notificationRepo,
				smtpClient:       tt.fields.smtpClient,
			}
			got, err := u.CreateNotification(tt.args.ctx, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateNotification() got = %v, want %v", got, tt.want)
			}
		})
	}
}
