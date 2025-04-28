package telegram

import (
	"domashka-backend/internal/custom_errors"
	"errors"
	tele "gopkg.in/telebot.v4"
)

type ContactHandler struct {
	contactUsecase ContactUseCase
}

func NewContactHandler(tg ContactUseCase) *ContactHandler {
	return &ContactHandler{contactUsecase: tg}
}

func (h *ContactHandler) Start(c tele.Context) error {
	replyKeyboard := &tele.ReplyMarkup{}

	btnContact := replyKeyboard.Contact("Отправить контакт")
	replyKeyboard.Reply(replyKeyboard.Row(btnContact))
	replyKeyboard.OneTimeKeyboard = true

	return c.Reply(
		"Чтобы подтвердить авторизацию, отправьте, пожалуйста, свой контакт, нажав на кнопку ниже.",
		replyKeyboard,
	)
}

func (h *ContactHandler) HandleContact(c tele.Context) error {
	contact := c.Message().Contact
	if contact == nil {
		return c.Reply("Пожалуйста, отправьте контакт, используя кнопку ниже.")
	}

	err := h.contactUsecase.HandleContact(c, contact)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPhoneNumberMismatch) {
			replyKeyboard := &tele.ReplyMarkup{}
			btnRetry := replyKeyboard.Contact("Отправить контакт ещё раз")
			replyKeyboard.Reply(replyKeyboard.Row(btnRetry))
			replyKeyboard.OneTimeKeyboard = true

			return c.Reply("Номер телефона не совпадает с введённым. Попробуйте снова.", replyKeyboard)
		}

		return c.Reply("Произошла ошибка. Попробуйте снова позже.")
	}

	return c.Reply("Авторизация успешна!")
}
