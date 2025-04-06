package app

import (
	botLib "github.com/go-telegram/bot"
	"lotBot/pkg/lotBot/bot"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(botLib.HandlerTypeMessageText, "/start", botLib.MatchTypeExact, bot.StartHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "start", botLib.MatchTypePrefix, bot.StartHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "role_", botLib.MatchTypePrefix, bot.CallbackHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "register_", botLib.MatchTypePrefix, bot.Register)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "submit_", botLib.MatchTypePrefix, bot.Moderation)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "action_", botLib.MatchTypePrefix, bot.ModerationResponse)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "view_tasks", botLib.MatchTypePrefix, bot.ViewTasks)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "ready_", botLib.MatchTypePrefix, bot.StudentReadiness)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "call", botLib.MatchTypePrefix, bot.Call)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, "not_", botLib.MatchTypePrefix, bot.NotReady)
}
