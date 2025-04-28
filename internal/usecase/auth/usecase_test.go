package auth

import (
	"context"
	"domashka-backend/internal/entity/auth"
	userentity "domashka-backend/internal/entity/users"
	"testing"
	"time"
)

type mockUsersRepo struct{}
type mockRedisClient struct{}
type mockJwtUsecase struct{}
type mockSMSClient struct{}

func TestUseCase_Auth(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		ctx context.Context
		req auth.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				usersRepo: mockUsersRepo{},
				redis:     mockRedisClient{},
				sms:       mockSMSClient{},
				jwt:       mockJwtUsecase{},
			},
			args: args{
				ctx: context.Background(),
				req: auth.Request{
					Phone: "1234567890",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			if err := u.Auth(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_AuthViaTg(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				usersRepo: mockUsersRepo{},
				redis:     mockRedisClient{},
				sms:       mockSMSClient{},
				jwt:       mockJwtUsecase{},
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "1234567890",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			if err := u.AuthViaTg(tt.args.ctx, tt.args.phoneNumber); (err != nil) != tt.wantErr {
				t.Errorf("AuthViaTg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_AuthViaTgStatus(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		in0         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				usersRepo: mockUsersRepo{},
				redis:     mockRedisClient{},
				sms:       mockSMSClient{},
				jwt:       mockJwtUsecase{},
			},
			args: args{
				in0:         context.Background(),
				phoneNumber: "1234567890",
			},
			want:    "success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			got, err := u.AuthViaTgStatus(tt.args.in0, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthViaTgStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthViaTgStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Verify(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		ctx   context.Context
		phone string
		otp   string
		role  string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantUserID int64
		wantToken  string
		wantErr    bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			gotUserID, gotToken, err := u.Verify(tt.args.ctx, tt.args.phone, tt.args.otp, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("Verify() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
			if gotToken != tt.wantToken {
				t.Errorf("Verify() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestUseCase_login(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		ctx  context.Context
		user *userentity.User
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
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			if err := u.login(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_register(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		ctx   context.Context
		phone string
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
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			if err := u.register(tt.args.ctx, tt.args.phone); (err != nil) != tt.wantErr {
				t.Errorf("register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_sendOTP(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		ctx  context.Context
		user *userentity.User
		otp  string
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
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			if err := u.sendOTP(tt.args.ctx, tt.args.user, tt.args.otp); (err != nil) != tt.wantErr {
				t.Errorf("sendOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_validateOTP(t *testing.T) {
	type fields struct {
		usersRepo usersRepo
		redis     redisClient
		sms       SMSClient
		jwt       jwtUsecase
	}
	type args struct {
		phone string
		otp   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				usersRepo: tt.fields.usersRepo,
				redis:     tt.fields.redis,
				sms:       tt.fields.sms,
				jwt:       tt.fields.jwt,
			}
			got, err := u.validateOTP(tt.args.phone, tt.args.otp)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateOTP() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateOTP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateOTP(); got != tt.want {
				t.Errorf("generateOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSMSDelay(t *testing.T) {
	type args struct {
		attempts int
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSMSDelay(tt.args.attempts); got != tt.want {
				t.Errorf("getSMSDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}
