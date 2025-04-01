package bot

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Предприниматель", CallbackData: "role_1"},
				{Text: "Исполнитель", CallbackData: "role_2"},
			},
		},
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Выберите роль",
		ReplyMarkup: kb,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		return
	}
}

func CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Всегда сначала отвечаем на callback
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("Error answering callback: %v", err)
		return
	}

	var response string
	switch update.CallbackQuery.Data {
	case "role_1":
		response = "Приветствие:\n✨ Добро пожаловать на EAZZY — сервис подросткового аутсорсинга!\n\n✔️ " +
			"Возьмем ответственность за выполнение задачи на себя как полноценный бизнес-партнёр\n✔️ " +
			"Подберем проверенных исполнителей, обучаем их и сопровождаем.\n✔️ " +
			"Проконтролируем качество и отдадим результат, соответствующий ожиданиям\n\n" +
			"Для начала давайте познакомимся\n🚀 Погнали!\n"
	case "role_2":
		response = "Приветствие:\n✨ Добро пожаловать на EAZZY — сервис подросткового аутсорсинга!\n\n✔️ " +
			"Поможем тебе сформулировать и описать твои умения и превратить их в доход\n✔️ " +
			"Предоставим безопасные и честные рабочие возможности\n✔️ " +
			"Дадим старт твоей карьере, поддержим и поможем в процессе\n\n\nДля начала давай знакомиться\n🚀 Погнали!"
	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   response,
	})
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
