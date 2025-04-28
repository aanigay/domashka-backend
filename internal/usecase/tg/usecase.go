package tg

import (
	"context"
	"domashka-backend/internal/custom_errors"
	usersentity "domashka-backend/internal/entity/users"
	tele "gopkg.in/telebot.v4"
	"strconv"
	"time"
)

type UseCase struct {
	Redis      redisClient
	usersRepo  usersRepo
	jwtUsecase jwtUsecase
}

func New(redis redisClient, usersRepo usersRepo, jwtUsecase jwtUsecase) *UseCase {
	return &UseCase{
		Redis:      redis,
		usersRepo:  usersRepo,
		jwtUsecase: jwtUsecase,
	}
}

func (u *UseCase) HandleContact(c tele.Context, contact *tele.Contact) error {
	method, err := u.Redis.Get("tg_auth:" + contact.PhoneNumber)
	if err != nil {
		return err
	}
	if method == "" {
		return custom_errors.ErrPhoneNumberMismatch
	}

	// Номер совпадает, вызываем нужный метод
	switch method {
	case "register":
		return u.tgRegister(contact)
	case "login":
		return u.tgLogin(contact)
	default:
		return custom_errors.ErrPhoneNumberMismatch
	}
}

func (u *UseCase) tgLogin(contact *tele.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := u.usersRepo.GetByPhone(ctx, contact.PhoneNumber)
	if err != nil {
		return err
	}

	user.Name = contact.FirstName
	if contact.LastName != "" {
		user.Name += " " + contact.LastName
	}
	user.ChatID = strconv.FormatInt(contact.UserID, 10)

	if err := u.usersRepo.Update(ctx, user.ID, *user); err != nil {
		return err
	}

	user, err = u.usersRepo.GetByPhone(ctx, contact.PhoneNumber)
	if err != nil {
		return err
	}
	jwt, err := u.jwtUsecase.GenerateJWT(user.ID, "user")
	if err != nil {
		return err
	}

	if err := u.Redis.Set("token:"+contact.PhoneNumber, jwt, 24*time.Hour); err != nil {
		return err
	}

	return u.Redis.Delete("tg_auth:" + contact.PhoneNumber)
}

func (u *UseCase) tgRegister(contact *tele.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &usersentity.User{
		Name:        contact.FirstName,
		ChatID:      strconv.FormatInt(contact.UserID, 10),
		NumberPhone: &contact.PhoneNumber,
	}
	if contact.LastName != "" {
		user.Name += " " + contact.LastName
	}

	if err := u.usersRepo.Create(ctx, user); err != nil {
		return err
	}
	user, err := u.usersRepo.GetByPhone(ctx, contact.PhoneNumber)
	if err != nil {
		return err
	}

	jwt, err := u.jwtUsecase.GenerateJWT(user.ID, "user")
	if err != nil {
		return err
	}

	if err := u.Redis.Set("token:"+contact.PhoneNumber, jwt, 24*time.Hour); err != nil {
		return err
	}

	return u.Redis.Delete("tg_auth:" + contact.PhoneNumber)
}
