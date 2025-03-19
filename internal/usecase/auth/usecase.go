package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"domashka-backend/internal/custom_errors"
)

type UseCase struct {
	usersRepo usersRepo
	redis     redisClient
	jwt       jwtUsecase
}

func New(repo usersRepo, redis redisClient, jwt jwtUsecase) *UseCase {
	return &UseCase{
		usersRepo: repo,
		redis:     redis,
		jwt:       jwt,
	}
}

func (u *UseCase) SendOTP(phone string, otp string) error {
	return nil // TODO: implement
}

func (u *UseCase) Register(ctx context.Context, phone string) error {
	user, err := u.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
			user = nil
		}
	}

	if user != nil {
		return custom_errors.ErrUserExists
	}

	if err = u.usersRepo.CreateWithPhone(ctx, phone); err != nil {
		return err
	}

	OTP := generateOTP()

	if err = u.redis.Set(phone, OTP, 2*time.Minute); err != nil {
		return err
	}

	if err = u.SendOTP(phone, OTP); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) Verify(_ context.Context, phone string, otp string) (string, error) {
	isValid, err := u.ValidateOTP(phone, otp)
	if err != nil || !isValid {
		return "", fmt.Errorf("invalid or expired OTP")
	}
	token, err := u.jwt.GenerateJWT()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UseCase) Login(ctx context.Context, phone string) error {
	if _, err := u.usersRepo.GetByPhone(ctx, phone); err != nil {
		return err
	}

	OTP := generateOTP()

	// Store OTP
	if err := u.redis.Set(phone, OTP, 5*time.Minute); err != nil {
		return err
	}

	if err := u.SendOTP(phone, OTP); err != nil {
		return err
	}
	return nil
}

func (u *UseCase) ValidateOTP(phone string, otp string) (bool, error) {
	storedOTP, err := u.redis.Get(phone)
	if err != nil {
		return false, err
	}
	return storedOTP == otp, nil
}

func generateOTP() string {
	otp := make([]byte, 6)
	_, err := rand.Read(otp)
	if err != nil {
		log.Fatal("Failed to generate OTP")
	}
	for i := range otp {
		otp[i] = '0' + otp[i]%10
	}
	return string(otp)
}
