package auth

import (
	"context"
	"domashka-backend/internal/utils/pointers"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"domashka-backend/internal/entity/auth"
	userentity "domashka-backend/internal/entity/users"
)

func TestUseCase_Auth(t *testing.T) {
	tests := []struct {
		name      string
		usersRepo func(ctrl *gomock.Controller) usersRepo
		redis     func(ctrl *gomock.Controller) redisClient
		sms       func(ctrl *gomock.Controller) SMSClient
		jwt       func(ctrl *gomock.Controller) jwtUsecase
		req       auth.Request
		wantErr   bool
	}{
		{
			name: "user nil",
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)
				m.EXPECT().CreateWithPhone(gomock.Any(), gomock.Any()).Return(&userentity.User{
					NumberPhone: pointers.To("81231234567"),
				}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			sms: func(ctrl *gomock.Controller) SMSClient {
				m := NewMockSMSClient(ctrl)
				m.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			jwt: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.usersRepo(ctrl), tt.redis(ctrl), tt.jwt(ctrl), tt.sms(ctrl))
			if err := u.Auth(context.Background(), tt.req); (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_AuthViaTg(t *testing.T) {
	tests := []struct {
		name      string
		usersRepo func(ctrl *gomock.Controller) usersRepo
		redis     func(ctrl *gomock.Controller) redisClient
		sms       func(ctrl *gomock.Controller) SMSClient
		jwt       func(ctrl *gomock.Controller) jwtUsecase
		in        string
		wantErr   bool
	}{
		{
			name: "success",
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(&userentity.User{
					NumberPhone: pointers.To("81231234567"),
				}, nil)
				return m
			},
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			sms: func(ctrl *gomock.Controller) SMSClient {
				m := NewMockSMSClient(ctrl)
				return m
			},
			jwt: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.usersRepo(ctrl), tt.redis(ctrl), tt.jwt(ctrl), tt.sms(ctrl))
			if err := u.AuthViaTg(context.Background(), tt.in); (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_AuthViaTgStatus(t *testing.T) {
	tests := []struct {
		name      string
		usersRepo func(ctrl *gomock.Controller) usersRepo
		redis     func(ctrl *gomock.Controller) redisClient
		sms       func(ctrl *gomock.Controller) SMSClient
		jwt       func(ctrl *gomock.Controller) jwtUsecase
		in        string
		want      string
		wantErr   bool
	}{
		{
			name: "success",
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				return m
			},
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().IsExpired(gomock.Any()).Return(false, nil)
				m.EXPECT().Get(gomock.Any()).Return("token", nil)
				return m
			},
			sms: func(ctrl *gomock.Controller) SMSClient {
				m := NewMockSMSClient(ctrl)
				return m
			},
			jwt: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				return m
			},
			want: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.usersRepo(ctrl), tt.redis(ctrl), tt.jwt(ctrl), tt.sms(ctrl))
			if out, err := u.AuthViaTgStatus(context.Background(), tt.in); (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.Equal(t, tt.want, out)
			}
		})
	}
}

func TestUseCase_Verify(t *testing.T) {
	tests := []struct {
		name       string
		usersRepo  func(ctrl *gomock.Controller) usersRepo
		redis      func(ctrl *gomock.Controller) redisClient
		sms        func(ctrl *gomock.Controller) SMSClient
		jwt        func(ctrl *gomock.Controller) jwtUsecase
		phone      string
		otp        string
		role       string
		wantUserID int64
		wantChefID *int64
		wantToken  string
		wantErr    bool
	}{
		{
			name: "success",
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().CheckIfUserIsChef(gomock.Any(), int64(1)).Return(nil, false, nil)
				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(&userentity.User{
					ID: int64(1),
				}, nil)
				return m
			},
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				return m
			},
			sms: func(ctrl *gomock.Controller) SMSClient {
				m := NewMockSMSClient(ctrl)
				return m
			},
			jwt: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				m.EXPECT().GenerateJWT(int64(1), nil, "admin").Return("token", nil)
				return m
			},
			phone:      "81231234567",
			otp:        "0123",
			role:       "admin",
			wantUserID: 1,
			wantToken:  "token",
			wantErr:    false,
			wantChefID: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.usersRepo(ctrl), tt.redis(ctrl), tt.jwt(ctrl), tt.sms(ctrl))
			gotUserID, gotChefID, gotToken, err := u.Verify(context.Background(), tt.phone, tt.otp, tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("Verify() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
			if !reflect.DeepEqual(gotChefID, tt.wantChefID) {
				t.Errorf("Verify() gotChefID = %v, want %v", gotChefID, tt.wantChefID)
			}
			if gotToken != tt.wantToken {
				t.Errorf("Verify() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestUseCase_login(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *userentity.User
	}
	tests := []struct {
		name      string
		usersRepo func(ctrl *gomock.Controller) usersRepo
		redis     func(ctrl *gomock.Controller) redisClient
		sms       func(ctrl *gomock.Controller) SMSClient
		jwt       func(ctrl *gomock.Controller) jwtUsecase
		args      args
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				user: &userentity.User{
					NumberPhone: pointers.To("81231235656"),
				},
			},
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			sms: func(ctrl *gomock.Controller) SMSClient {
				m := NewMockSMSClient(ctrl)
				m.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			jwt: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				return m
			},
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.usersRepo(ctrl), tt.redis(ctrl), tt.jwt(ctrl), tt.sms(ctrl))
			if err := u.login(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_register(t *testing.T) {
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name      string
		usersRepo func(ctrl *gomock.Controller) usersRepo
		redis     func(ctrl *gomock.Controller) redisClient
		sms       func(ctrl *gomock.Controller) SMSClient
		jwt       func(ctrl *gomock.Controller) jwtUsecase
		args      args
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				phone: "81231235656",
			},
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			sms: func(ctrl *gomock.Controller) SMSClient {
				m := NewMockSMSClient(ctrl)
				m.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
			jwt: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				return m
			},
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().CreateWithPhone(gomock.Any(), gomock.Any()).Return(&userentity.User{
					ID:               0,
					Username:         "",
					Alias:            "",
					FirstName:        "",
					SecondName:       nil,
					LastName:         nil,
					Email:            nil,
					NumberPhone:      pointers.To("81231235656"),
					IsSpam:           0,
					SMSAttempts:      0,
					LastSMSRequest:   nil,
					Status:           0,
					ExternalType:     0,
					TelegramName:     nil,
					ExternalID:       nil,
					NotificationFlag: 0,
					Role:             "",
					Birthday:         nil,
					Name:             "",
					ChatID:           "",
					CreatedAt:        time.Time{},
					UpdatedAt:        time.Time{},
				}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.usersRepo(ctrl), tt.redis(ctrl), tt.jwt(ctrl), tt.sms(ctrl))
			if err := u.register(tt.args.ctx, tt.args.phone); (err != nil) != tt.wantErr {
				t.Errorf("register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
