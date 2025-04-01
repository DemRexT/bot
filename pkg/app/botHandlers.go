package app

import (
	botLib "github.com/go-telegram/bot"
	"log"
	"lotBot/pkg/lotBot/bot"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(botLib.HandlerTypeMessageText, "/start", botLib.MatchTypeExact, bot.StartHandler)
}

func (a *App) CallbackQueryDataHandler() {
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "role_", botLib.MatchTypePrefix, bot.CallbackHandler)
	log.Println("Bot handlers registered successfully")
}
