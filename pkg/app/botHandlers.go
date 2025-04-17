package app

import (
	botLib "github.com/go-telegram/bot"
	"lotBot/pkg/lotBot/bot"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(botLib.HandlerTypeMessageText, "/start", botLib.MatchTypeExact, bot.StartHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternStart, botLib.MatchTypePrefix, bot.StartHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternRole, botLib.MatchTypePrefix, bot.CallbackHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternRegister, botLib.MatchTypePrefix, bot.Register)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternAction, botLib.MatchTypePrefix, bot.ModerationResponse)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternViewTask, botLib.MatchTypePrefix, bot.ViewTasks)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternReady, botLib.MatchTypePrefix, bot.StudentReadiness)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternCall, botLib.MatchTypePrefix, a.bm.Call)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternNot, botLib.MatchTypePrefix, bot.NotReady)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternCreateTask, botLib.MatchTypePrefix, bot.CreateTask)
}
