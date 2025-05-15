package tg

import (
	"domashka-backend/internal/entity/users"
	"github.com/golang/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
)

func TestUseCase_HandleContact(t *testing.T) {
	type fields struct {
		Redis      redisClient
		usersRepo  usersRepo
		jwtUsecase jwtUsecase
	}
	type args struct {
		c       telebot.Context
		contact *telebot.Contact
	}
	tests := []struct {
		name       string
		redis      func(ctrl *gomock.Controller) redisClient
		usersRepo  func(ctrl *gomock.Controller) usersRepo
		jwtUsecase func(ctrl *gomock.Controller) jwtUsecase
		args       args
		wantErr    bool
	}{
		{
			name: "success",
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().Get(gomock.Any()).Return("login", nil)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Delete(gomock.Any()).Return(nil)
				return m
			},
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(&users.User{}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(&users.User{}, nil)
				return m
			},
			jwtUsecase: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				m.EXPECT().GenerateJWT(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
				return m
			},
			args: args{
				contact: &telebot.Contact{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.redis(ctrl), tt.usersRepo(ctrl), tt.jwtUsecase(ctrl))
			if err := u.HandleContact(tt.args.c, tt.args.contact); (err != nil) != tt.wantErr {
				t.Errorf("HandleContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_tgLogin(t *testing.T) {
	type fields struct {
		Redis      redisClient
		usersRepo  usersRepo
		jwtUsecase jwtUsecase
	}
	type args struct {
		contact *telebot.Contact
	}
	tests := []struct {
		name       string
		redis      func(ctrl *gomock.Controller) redisClient
		usersRepo  func(ctrl *gomock.Controller) usersRepo
		jwtUsecase func(ctrl *gomock.Controller) jwtUsecase
		args       args
		wantErr    bool
	}{
		{
			name: "success",
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Delete(gomock.Any()).Return(nil)
				return m
			},
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(&users.User{}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(&users.User{}, nil)
				return m
			},
			jwtUsecase: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				m.EXPECT().GenerateJWT(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
				return m
			},
			args: args{
				contact: &telebot.Contact{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.redis(ctrl), tt.usersRepo(ctrl), tt.jwtUsecase(ctrl))
			if err := u.tgLogin(tt.args.contact); (err != nil) != tt.wantErr {
				t.Errorf("tgLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_tgRegister(t *testing.T) {
	type fields struct {
		Redis      redisClient
		usersRepo  usersRepo
		jwtUsecase jwtUsecase
	}
	type args struct {
		contact *telebot.Contact
	}
	tests := []struct {
		name       string
		redis      func(ctrl *gomock.Controller) redisClient
		usersRepo  func(ctrl *gomock.Controller) usersRepo
		jwtUsecase func(ctrl *gomock.Controller) jwtUsecase
		args       args
		wantErr    bool
	}{
		{
			name: "success",
			redis: func(ctrl *gomock.Controller) redisClient {
				m := NewMockredisClient(ctrl)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Delete(gomock.Any()).Return(nil)

				return m
			},
			usersRepo: func(ctrl *gomock.Controller) usersRepo {
				m := NewMockusersRepo(ctrl)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().GetByPhone(gomock.Any(), gomock.Any()).Return(&users.User{}, nil)
				return m
			},
			jwtUsecase: func(ctrl *gomock.Controller) jwtUsecase {
				m := NewMockjwtUsecase(ctrl)
				m.EXPECT().GenerateJWT(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
				return m
			},
			args: args{
				contact: &telebot.Contact{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := New(tt.redis(ctrl), tt.usersRepo(ctrl), tt.jwtUsecase(ctrl))
			if err := u.tgRegister(tt.args.contact); (err != nil) != tt.wantErr {
				t.Errorf("tgRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
