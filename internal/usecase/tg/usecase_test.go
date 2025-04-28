package tg

import (
	"gopkg.in/telebot.v4"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		redis      redisClient
		usersRepo  usersRepo
		jwtUsecase jwtUsecase
	}
	tests := []struct {
		name string
		args args
		want *UseCase
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.redis, tt.args.usersRepo, tt.args.jwtUsecase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_HandleContact(t *testing.T) {
	type fields struct {
		Redis      redisClient
		usersRepo  usersRepo
		jwtUsecase jwtUsecase
	}
	type args struct {
		c       telebot.Context
		contact *tele.Contact
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				Redis:      tt.fields.Redis,
				usersRepo:  tt.fields.usersRepo,
				jwtUsecase: tt.fields.jwtUsecase,
			}
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
		contact *tele.Contact
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				Redis:      tt.fields.Redis,
				usersRepo:  tt.fields.usersRepo,
				jwtUsecase: tt.fields.jwtUsecase,
			}
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
		contact *tele.Contact
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				Redis:      tt.fields.Redis,
				usersRepo:  tt.fields.usersRepo,
				jwtUsecase: tt.fields.jwtUsecase,
			}
			if err := u.tgRegister(tt.args.contact); (err != nil) != tt.wantErr {
				t.Errorf("tgRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
