package telegram

import tele "gopkg.in/telebot.v4"

type ContactUseCase interface {
	HandleContact(c tele.Context, contact *tele.Contact) error
}
