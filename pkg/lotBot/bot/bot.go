package bot

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"strconv"
	"strings"
)

const AdminChatID = -4732218051

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Предприниматель", CallbackData: "role_1"},
				{Text: "Исполнитель", CallbackData: "role_2"},
			},
		},
	}

	if update.CallbackQuery != nil && update.CallbackQuery.Data == "start" {
		_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})
		if err != nil {
			return
		}

		// Вызываем ту же логику, что и для /start
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "Выберите роль",
			ReplyMarkup: kb,
		})
		if err != nil {
			return
		}
	} else {
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

}

func CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("Error answering callback: %v", err)
		return
	}

	var response string
	var kb *models.InlineKeyboardMarkup
	switch update.CallbackQuery.Data {
	case "role_1":
		response = "Приветствие:\n✨ Добро пожаловать на EAZZY — сервис подросткового аутсорсинга!\n\n✔️ " +
			"Возьмем ответственность за выполнение задачи на себя как полноценный бизнес-партнёр\n✔️ " +
			"Подберем проверенных исполнителей, обучаем их и сопровождаем.\n✔️ " +
			"Проконтролируем качество и отдадим результат, соответствующий ожиданиям\n\n" +
			"Для начала давайте познакомимся\n🚀 Погнали!\n"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Зарегистрироваться", CallbackData: "register_Business"},
				},
			},
		}
	case "role_2":
		response = "Приветствие:\n✨ Добро пожаловать на EAZZY — сервис подросткового аутсорсинга!\n\n✔️ " +
			"Поможем тебе сформулировать и описать твои умения и превратить их в доход\n✔️ " +
			"Предоставим безопасные и честные рабочие возможности\n✔️ " +
			"Дадим старт твоей карьере, поддержим и поможем в процессе\n\n\nДля начала давай знакомиться\n🚀 Погнали!"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Зарегистрироваться", CallbackData: "register_Teen"},
				},
			},
		}
	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func Register(ctx context.Context, b *bot.Bot, update *models.Update) {
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
	case "register_Business":

		response = "Регистрация заказчика"

	case "register_Teen":
		response = "Регистрация исполнителя"

	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Отправить на модерацию", CallbackData: "submit_for_moderation_" + update.CallbackQuery.Data},
			},
		},
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func Moderation(ctx context.Context, b *bot.Bot, update *models.Update) {

	userID := update.CallbackQuery.Message.Message.Chat.ID
	parts := strings.Split(update.CallbackQuery.Data, "_")

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Принять",
					CallbackData: "action_accept_" + strconv.FormatInt(userID, 10) + "_" + parts[4],
				},
			},
			{
				{
					Text:         "Отклонить",
					CallbackData: "action_reject_" + strconv.FormatInt(userID, 10) + "_" + parts[4],
				},
			},
		},
	}
	var response string
	switch parts[4] {
	case "Business":

		response = "Модерация заказчика"

	case "Teen":
		response = "Модерация исполнителя"

	default:

		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      AdminChatID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func ModerationResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("Error answering callback: %v", err)
		return
	}
	parts := strings.Split(update.CallbackQuery.Data, "_")

	actionID, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Проблема с ID: %v", err)
		return
	}

	//"action_reject_" + strconv.FormatInt(userID, 10) + "_" + parts[4]
	var kb *models.InlineKeyboardMarkup
	var response string
	switch parts[1] {
	case "reject":

		response = "Заявка не прошла модерацию(\n\n" +
			"Были введены некорректные или недостоверные данные.\n" +
			"Пожалуйста, вернись к первому шагу и проверь, не допущена ли ошибка"

		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "Вернутся назад",
						CallbackData: "start",
					},
				},
			},
		}
	case "accept":
		switch parts[3] {
		case "Business":
			response = "Модерация пройдена!\n\nХотите разместить первое задание?"
			kb = &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{
							Text:         "Да",
							CallbackData: "create_task",
						},
						{
							Text:         "Позже",
							CallbackData: "later",
						},
					},
				},
			}
		case "Teen":
			response = "Твои данные подтверждены!\n\nПосмотрим, есть ли у нас для тебя первое задание?"
			kb = &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{
							Text:         "Посмотрим",
							CallbackData: "view_tasks",
						},
					},
				},
			}
		}

	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      actionID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}

}

func ViewTasks(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("%v", err)
		return
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Готов", CallbackData: "ready_yes"},
				{Text: "Не готов", CallbackData: "ready_not"},
			},
		},
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text: "У нас есть для тебя задание!\nПожалуйста. ознакомься с заданием.\n" +
			"Срок для изучения задания - до ЧЧ.ММ ДД.ММ\n" +
			"Пришлем напоминалку поле этого срока и уточним готовность.\n" +
			"И помни: мы не выполним задание за тебя,\nно обязательно поможем и подскажем,\n" +
			"если будет трудно или непонятно!",
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("%v", err)
	}

}

func StudentReadiness(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("%v", err)
	}

	var response string
	var kb *models.InlineKeyboardMarkup
	switch update.CallbackQuery.Data {
	case "ready_yes":
		response = "Отлично!\nДавай назначим созвон с заказчиком длявыяснения деталей,затем ты сможеш приступить"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Окей", CallbackData: "call"},
				},
			},
		}
	case "ready_not":
		response = "Подскажи, почему именно:"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "еще занят с предыдущей задачей", CallbackData: "not_busy"},
				},
				{
					{Text: "задача мне не интересна", CallbackData: "not_interesting"},
				},
				{
					{Text: "не понял задание и/или не уверен, что справлюсь", CallbackData: "not_understand"},
				},
			},
		}
	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("%v", err)
	}
}

func Call(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("%v", err)
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Готово", CallbackData: "_"},
			},
		},
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        "Отлично! Давай назначим созвон",
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("%v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: AdminChatID,
		Text:   "Запрос на новый созвон от пользователя!",
	})
	if err != nil {
		log.Printf("%v", err)
	}

}

func NotReady(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("%v", err)
	}

	parts := strings.Split(update.CallbackQuery.Data, "_")
	var kb *models.InlineKeyboardMarkup
	var response string
	switch parts[1] {
	case "busy":
		response = "Хочешь взять это задание следующим после текущего?"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Да", CallbackData: "call_you"},
					{Text: "Нет", CallbackData: "following_tasks"},
				},
			},
		}
	case "interesting":
		response = "Больше не отправлять тебе задачи из этого трека?"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Да, не моё направление", CallbackData: "_"},
					{Text: "Нет, отправляйте другие, не зашло именно это", CallbackData: "_"},
				},
			},
		}
	case "understand":
		response = "Хочешь задать вопросы и получить \\более подробное пояснение?"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "да, пожалуйста  (связь в личке или созвон)", CallbackData: "call_you"},
					{Text: "нет, спасибо", CallbackData: "following_tasks"},
				},
			},
		}
	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("%v", err)
		log.Printf("%v", update.CallbackQuery.Data)
	}
}
