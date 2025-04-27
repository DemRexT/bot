package app

import (
	botLib "github.com/go-telegram/bot"
	"lotBot/pkg/lotBot/bot"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(botLib.HandlerTypeMessageText, "/start", botLib.MatchTypeExact, a.bm.StartHandler)
	a.b.RegisterHandler(botLib.HandlerTypeMessageText, "/pay", botLib.MatchTypeExact, a.bm.PayHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternStart, botLib.MatchTypePrefix, a.bm.StartHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternRole, botLib.MatchTypePrefix, a.bm.CallbackHandler)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternRegister, botLib.MatchTypePrefix, a.bm.Register)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternAction, botLib.MatchTypePrefix, a.bm.ModerationResponse)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternViewTask, botLib.MatchTypePrefix, a.bm.ViewTasks)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternReady, botLib.MatchTypePrefix, a.bm.StudentReadiness)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternCall, botLib.MatchTypePrefix, a.bm.Call)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternNot, botLib.MatchTypePrefix, a.bm.NotReady)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternCreateTask, botLib.MatchTypePrefix, a.bm.CreateTask)
	a.b.RegisterHandler(botLib.HandlerTypeCallbackQueryData, bot.PatternLater, botLib.MatchTypePrefix, a.bm.Later)
}
