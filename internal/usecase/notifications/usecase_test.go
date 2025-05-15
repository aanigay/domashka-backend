package notifications

import (
	"context"
	"domashka-backend/internal/entity/notifications"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestNotificationUsecase_CreateNotification(t *testing.T) {
	type args struct {
		ctx context.Context
		n   notifications.Notification
	}
	tests := []struct {
		name              string
		notificationsRepo func(ctrl *gomock.Controller) NotificationRepository
		smtpClient        func(ctrl *gomock.Controller) SMTPClient
		args              args
		want              int
		wantErr           bool
	}{
		{
			name: "success",
			notificationsRepo: func(ctrl *gomock.Controller) NotificationRepository {
				m := NewMockNotificationRepository(ctrl)
				m.EXPECT().CreateNotification(gomock.Any(), gomock.Any()).Return(1, nil)
				return m
			},
			smtpClient: func(ctrl *gomock.Controller) SMTPClient {
				m := NewMockSMTPClient(ctrl)
				return m
			},
			args: args{
				ctx: context.Background(),
				n:   notifications.Notification{},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.notificationsRepo(ctrl), tt.smtpClient(ctrl))
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

func TestNotificationUsecase_GetNotificationByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name              string
		notificationsRepo func(ctrl *gomock.Controller) NotificationRepository
		smtpClient        func(ctrl *gomock.Controller) SMTPClient
		args              args
		want              *notifications.Notification
		wantErr           bool
	}{
		{
			name: "success",
			notificationsRepo: func(ctrl *gomock.Controller) NotificationRepository {
				m := NewMockNotificationRepository(ctrl)
				m.EXPECT().GetNotificationByID(gomock.Any(), gomock.Any()).Return(&notifications.Notification{}, nil)
				return m
			},
			smtpClient: func(ctrl *gomock.Controller) SMTPClient {
				m := NewMockSMTPClient(ctrl)
				return m
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &notifications.Notification{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.notificationsRepo(ctrl), tt.smtpClient(ctrl))
			got, err := u.GetNotificationByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotificationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotificationByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationUsecase_GetNotifications(t *testing.T) {
	type fields struct {
		notificationRepo NotificationRepository
		smtpClient       SMTPClient
	}
	type args struct {
		ctx     context.Context
		filters map[string]string
		page    int
		limit   int
	}
	tests := []struct {
		name              string
		notificationsRepo func(ctrl *gomock.Controller) NotificationRepository
		smtpClient        func(ctrl *gomock.Controller) SMTPClient
		args              args
		want              []notifications.Notification
		want1             int
		wantErr           bool
	}{
		{
			name: "success",
			notificationsRepo: func(ctrl *gomock.Controller) NotificationRepository {
				m := NewMockNotificationRepository(ctrl)
				m.EXPECT().GetNotifications(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]notifications.Notification{}, 0, nil)
				return m
			},
			smtpClient: func(ctrl *gomock.Controller) SMTPClient {
				m := NewMockSMTPClient(ctrl)
				return m
			},
			args: args{
				ctx:     context.Background(),
				filters: map[string]string{},
				page:    1,
				limit:   10,
			},
			want:    []notifications.Notification{},
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.notificationsRepo(ctrl), tt.smtpClient(ctrl))
			got, got1, err := u.GetNotifications(tt.args.ctx, tt.args.filters, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotifications() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetNotifications() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNotificationUsecase_ResendNotification(t *testing.T) {
	type fields struct {
		notificationRepo NotificationRepository
		smtpClient       SMTPClient
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name              string
		notificationsRepo func(ctrl *gomock.Controller) NotificationRepository
		smtpClient        func(ctrl *gomock.Controller) SMTPClient
		args              args
		wantErr           bool
	}{
		{
			name: "success",
			notificationsRepo: func(ctrl *gomock.Controller) NotificationRepository {
				m := NewMockNotificationRepository(ctrl)
				m.EXPECT().GetNotificationByID(gomock.Any(), gomock.Any()).Return(&notifications.Notification{}, nil)
				m.EXPECT().UpdateNotification(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			smtpClient: func(ctrl *gomock.Controller) SMTPClient {
				m := NewMockSMTPClient(ctrl)
				m.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
				return m
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.notificationsRepo(ctrl), tt.smtpClient(ctrl))
			if err := u.ResendNotification(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ResendNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotificationUsecase_SendEmailNotification(t *testing.T) {
	type fields struct {
		notificationRepo NotificationRepository
		smtpClient       SMTPClient
	}
	type args struct {
		ctx context.Context
		n   notifications.Notification
	}
	tests := []struct {
		name              string
		notificationsRepo func(ctrl *gomock.Controller) NotificationRepository
		smtpClient        func(ctrl *gomock.Controller) SMTPClient
		args              args
		wantErr           bool
	}{
		{
			name: "success",
			notificationsRepo: func(ctrl *gomock.Controller) NotificationRepository {
				m := NewMockNotificationRepository(ctrl)
				m.EXPECT().UpdateNotification(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			smtpClient: func(ctrl *gomock.Controller) SMTPClient {
				m := NewMockSMTPClient(ctrl)
				m.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
				return m
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.notificationsRepo(ctrl), tt.smtpClient(ctrl))
			if err := u.SendEmailNotification(tt.args.ctx, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("SendEmailNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
