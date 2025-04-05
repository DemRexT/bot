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
