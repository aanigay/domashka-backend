package telegram

import (
	tele "gopkg.in/telebot.v4"
)

func NewBot(bot *tele.Bot, contactUseCase ContactUseCase) {
	tgHandler := NewContactHandler(contactUseCase)

	bot.Handle("/start", tgHandler.Start)
	bot.Handle(tele.OnContact, tgHandler.HandleContact)
}
