package auth

import (
	"context"
	"crypto/rand"
	"domashka-backend/internal/entity/auth"
	userentity "domashka-backend/internal/entity/users"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"

	"domashka-backend/internal/custom_errors"
)

const (
	MaxAttempts  = 11
	InitialDelay = 1 * time.Minute
	SecondDelay  = 2 * time.Minute
	MaxDelay     = 5 * time.Minute
)

type UseCase struct {
	usersRepo usersRepo
	redis     redisClient
	sms       SMSClient
	jwt       jwtUsecase
}

func New(repo usersRepo, redis redisClient, jwt jwtUsecase, sms SMSClient) *UseCase {
	return &UseCase{
		usersRepo: repo,
		redis:     redis,
		jwt:       jwt,
		sms:       sms,
	}
}

func (u *UseCase) sendOTP(ctx context.Context, user *userentity.User, otp string) error {
	if user.IsSpam == 1 {
		return custom_errors.ErrUserIsSpam
	}

	if user.LastSMSRequest != nil {
		delay := getSMSDelay(user.SMSAttempts)
		if time.Since(*user.LastSMSRequest) < delay {
			return errors.New("пожалуйста, подождите перед повторной отправкой")
		}
	}

	if user.NumberPhone == nil || *user.NumberPhone == "" {
		return errors.New("номер телефона отсутствует")
	}

	userID := user.ID

	if user.SMSAttempts >= MaxAttempts {
		user.IsSpam = 1

		if err := u.usersRepo.Update(ctx, userID, *user); err != nil {
			return err
		}
		return custom_errors.ErrUserIsSpam
	}

	if err := u.sms.Send(*user.NumberPhone, otp); err != nil {
		return err
	}

	now := time.Now()
	user.LastSMSRequest = &now
	user.SMSAttempts++
	if err := u.usersRepo.Update(ctx, userID, *user); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) Auth(ctx context.Context, req auth.Request) error {
	log.Printf("DEBUG: Начало аутентификации для номера: %s", req.Phone)
	user, err := u.usersRepo.GetByPhone(ctx, req.Phone)

	if user == nil {
		log.Printf("DEBUG: Пользователь с номером %s не найден (user == nil), запускаем регистрацию", req.Phone)
		return u.register(ctx, req.Phone)
	}
	if err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
			log.Printf("DEBUG: Пользователь с номером %s не найден, запускаем регистрацию", req.Phone)
			return u.register(ctx, req.Phone)
		}
		log.Printf("DEBUG: Пользователь: %s", user)
		log.Printf("DEBUG: Ошибка получения пользователя по номеру %s: %v", req.Phone, err)
		return err
	}
	log.Printf("DEBUG: Пользователь с номером %s найден, запускаем логин", req.Phone)
	return u.login(ctx, user)
}

func (u *UseCase) AuthViaTg(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("tg_auth:%s", phoneNumber)

	user, err := u.usersRepo.GetByPhone(ctx, phoneNumber)
	if err != nil && !errors.Is(err, custom_errors.ErrUserNotFound) {
		return err
	}
	if user == nil {
		err = u.redis.Set(key, "register", 5*time.Minute)
		return err
	}

	err = u.redis.Set(key, "login", 5*time.Minute)
	return err
}

func (u *UseCase) AuthViaTgStatus(_ context.Context, phoneNumber string) (string, error) {
	key := fmt.Sprintf("tg_auth:%s", phoneNumber)
	isExpired, err := u.redis.IsExpired(key)
	if err != nil {
		return "", err
	}
	if isExpired {
		return "", custom_errors.ErrExpiredTTL
	}

	token, err := u.redis.Get("token:" + phoneNumber)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", custom_errors.ErrConfirmationNotReceived
		}
		return "", err
	}

	if token == "" {
		return "", custom_errors.ErrConfirmationNotReceived
	}

	return token, nil
}

func (u *UseCase) register(ctx context.Context, phone string) error {
	user, err := u.usersRepo.CreateWithPhone(ctx, phone)
	if err != nil {
		return err
	}

	otp := generateOTP()

	if err := u.redis.Set(phone, otp, 5*time.Minute); err != nil {
		return err
	}

	if err := u.sendOTP(ctx, user, otp); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) Verify(ctx context.Context, phone string, otp string, role string) (userID int64, token string, err error) {
	isValid, err := u.validateOTP(phone, otp)
	if err != nil || !isValid {
		return 0, "", fmt.Errorf("invalid or expired OTP")
	}
	user, err := u.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return 0, "", err
	}
	token, err = u.jwt.GenerateJWT(user.ID, role)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

func (u *UseCase) login(ctx context.Context, user *userentity.User) error {
	log.Printf("DEBUG: Начало логина для пользователя ID: %d", user.ID)
	otp := generateOTP()
	log.Printf("DEBUG: Сгенерирован OTP для логина пользователя ID: %d: %s", user.ID, otp)
	if err := u.redis.Set(*user.NumberPhone, otp, 5*time.Minute); err != nil {
		log.Printf("DEBUG: Ошибка установки OTP в Redis для пользователя ID: %d: %v", user.ID, err)
		return err
	}

	if err := u.sendOTP(ctx, user, otp); err != nil {
		log.Printf("DEBUG: Ошибка отправки OTP при логине для пользователя ID: %d: %v", user.ID, err)
		return err
	}
	log.Printf("DEBUG: Логин успешно завершен для пользователя ID: %d", user.ID)
	return nil
}

func (u *UseCase) validateOTP(phone string, otp string) (bool, error) {
	// todo: Убрать костыль перед выгрузкой в прод
	if otp == "0123" {
		return true, nil
	}
	storedOTP, err := u.redis.Get(phone)
	if err != nil {
		return false, err
	}
	return storedOTP == otp, nil
}

func generateOTP() string {
	otp := make([]byte, 4)
	_, err := rand.Read(otp)
	if err != nil {
		log.Fatal("failed to generate OTP")
	}
	for i := range otp {
		otp[i] = '0' + otp[i]%10
	}
	return string(otp)
}

func getSMSDelay(attempts int) time.Duration {
	switch attempts {
	case 1:
		return InitialDelay
	case 2:
		return SecondDelay
	default:
		return MaxDelay
	}
}
