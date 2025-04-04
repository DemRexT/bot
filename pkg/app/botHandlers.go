package app

import (
	botLib "github.com/go-telegram/bot"
	"lotBot/pkg/lotBot/bot"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(botLib.HandlerTypeMessageText, "/start", botLib.MatchTypeExact, bot.StartHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "role_", botLib.MatchTypePrefix, bot.CallbackHandler)
}
