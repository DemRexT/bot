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

type BotManager struct {
	adminChatID int
}

func NewBotManager(adminChatID int) *BotManager {
	return &BotManager{adminChatID: adminChatID}
}

const (
	PatternStart            = "start"
	PatternRole             = "role_"
	PatternRegister         = "register_"
	PatternSubmitModeration = "submit_for_moderation_"
	PatternAction           = "action_"
	PatternViewTask         = "view_tasks"
	PatternReady            = "ready_"
	PatternCall             = "call"
	PatternNot              = "not_"
)

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Предприниматель", CallbackData: PatternRole + "1"},
				{Text: "Исполнитель", CallbackData: PatternRole + "2"},
			},
		},
	}

	if update.CallbackQuery != nil && update.CallbackQuery.Data == PatternStart {
		_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})
		if err != nil {
			return
		}

		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
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
	case PatternRole + "1":
		response = "Приветствие:\n✨ Добро пожаловать на EAZZY — сервис подросткового аутсорсинга!\n\n✔️ " +
			"Возьмем ответственность за выполнение задачи на себя как полноценный бизнес-партнёр\n✔️ " +
			"Подберем проверенных исполнителей, обучаем их и сопровождаем.\n✔️ " +
			"Проконтролируем качество и отдадим результат, соответствующий ожиданиям\n\n" +
			"Для начала давайте познакомимся\n🚀 Погнали!\n"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Зарегистрироваться", CallbackData: PatternRegister + "Business"},
				},
			},
		}
	case PatternRole + "2":
		response = "Приветствие:\n✨ Добро пожаловать на EAZZY — сервис подросткового аутсорсинга!\n\n✔️ " +
			"Поможем тебе сформулировать и описать твои умения и превратить их в доход\n✔️ " +
			"Предоставим безопасные и честные рабочие возможности\n✔️ " +
			"Дадим старт твоей карьере, поддержим и поможем в процессе\n\n\nДля начала давай знакомиться\n🚀 Погнали!"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Зарегистрироваться", CallbackData: PatternRegister + "Teen"},
				},
			},
		}
	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
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
	case PatternRegister + "Business":

		response = "Регистрация заказчика"

	case PatternRegister + "Teen":
		response = "Регистрация исполнителя"

	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Отправить на модерацию", CallbackData: PatternSubmitModeration + update.CallbackQuery.Data},
			},
		},
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func (bm BotManager) Moderation(ctx context.Context, b *bot.Bot, update *models.Update) {

	userID := update.CallbackQuery.Message.Message.Chat.ID
	parts := strings.Split(update.CallbackQuery.Data, "_")

	if len(parts) < 5 {
		log.Printf("не удалось отобразить карточку пользователя, len(parts) < 5\n")
		return
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Принять",
					CallbackData: PatternAction + "accept_" + strconv.FormatInt(userID, 10) + "_" + parts[4],
				},
			},
			{
				{
					Text:         "Отклонить",
					CallbackData: PatternAction + "reject_" + strconv.FormatInt(userID, 10) + "_" + parts[4],
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
		ChatID:      bm.adminChatID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      "Заявка на модерацию отправлена",
	})
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

	if len(parts) < 4 {
		log.Printf("не удалось отобразить карточку пользователя, len(parts) < 4\n")
		return
	}

	actionID, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Проблема с ID: %v", err)
		return
	}

	var kb *models.InlineKeyboardMarkup
	var response string
	var responceAdmin string
	switch parts[1] {
	case "reject":

		response = "Заявка не прошла модерацию(\n\n" +
			"Были введены некорректные или недостоверные данные.\n" +
			"Пожалуйста, вернись к первому шагу и проверь, не допущена ли ошибка"

		responceAdmin = "Пользователь отклонен"

		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "Вернутся назад",
						CallbackData: PatternStart,
					},
				},
			},
		}
	case "accept":

		responceAdmin = "Пользователь подтвержден"
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
							CallbackData: PatternViewTask,
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

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      responceAdmin,
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
				{Text: "Готов", CallbackData: PatternReady + "yes"},
				{Text: "Не готов", CallbackData: PatternReady + "not"},
			},
		},
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
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
	case PatternReady + "yes":
		response = "Отлично!\nДавай назначим созвон с заказчиком длявыяснения деталей,затем ты сможеш приступить"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Окей", CallbackData: PatternCall},
				},
			},
		}
	case PatternReady + "not":
		response = "Подскажи, почему именно:"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "еще занят с предыдущей задачей", CallbackData: PatternNot + "busy"},
				},
				{
					{Text: "задача мне не интересна", CallbackData: PatternNot + "interesting"},
				},
				{
					{Text: "не понял задание и/или не уверен, что справлюсь", CallbackData: PatternNot + "understand"},
				},
			},
		}
	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("%v", err)
	}
}

func (bm BotManager) Call(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        "Отлично! Давай назначим созвон",
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("%v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: bm.adminChatID,
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

	if len(parts) < 2 {
		log.Printf("не удалось отобразить карточку пользователя, len(parts) < 2\n")
		return
	}
	var kb *models.InlineKeyboardMarkup
	var response string
	switch parts[1] {
	case "busy":
		response = "Хочешь взять это задание следующим после текущего?"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Да", CallbackData: PatternCall},
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
				},
				{
					{Text: "Нет, отправляйте другие, не зашло именно это", CallbackData: "_"},
				},
			},
		}
	case "understand":
		response = "Хочешь задать вопросы и получить \\более подробное пояснение?"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "да, пожалуйста  (связь в личке или созвон)", CallbackData: PatternCall},
				},
				{
					{Text: "нет, спасибо", CallbackData: "following_tasks"},
				},
			},
		}
	default:
		response = "Неизвестная команда: " + update.CallbackQuery.Data
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("%v", err)
		log.Printf("%v", update.CallbackQuery.Data)
	}
}
