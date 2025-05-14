package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"lotBot/pkg/db"
	"lotBot/pkg/embedlog"
	"lotBot/pkg/invoicebox"
	"strconv"
	"strings"
)

type BotManager struct {
	embedlog.Logger
	adminChatID int
	ic          *invoicebox.InvoiceClient
	repo        db.LotbotRepo
}

func NewBotManager(logger embedlog.Logger, adminChatID int, cfg invoicebox.Config) *BotManager {
	return &BotManager{
		Logger:      logger,
		adminChatID: adminChatID,
		ic:          invoicebox.NewInvoiceClient(logger, cfg),
	}
}

const (
	PatternStart             = "start"
	PatternRole              = "role_"
	PatternRegister          = "register_"
	PatternSubmitModeration  = "submit_for_moderation_"
	PatternAction            = "action_"
	PatternViewTask          = "view_tasks"
	PatternReady             = "ready_"
	PatternCall              = "call"
	PatternNot               = "not_"
	PatternCreateTask        = "create_task"
	PatternLater             = "later"
	PatternTaskCheckResponse = "check_response"
	UrlRegisterStudent       = "https://docs.google.com/forms/d/e/1FAIpQLSemsbNWCx2ewY25WlvQP_baBef6RUs1jF0w1p4obb99ieXFAw/viewform?usp=pp_url&entry.1082496981="
	UrlRegisterBusiness      = "https://docs.google.com/forms/d/e/1FAIpQLSdz5iYc9UB6M3wOOrGGl-4jTywltlkl7AZgqXrNKIBqrY87mA/viewform?usp=pp_url&entry.213949143="
	UrlCreateTask            = "https://docs.google.com/forms/d/e/1FAIpQLScQgB6T74K87rZHi8a9qi-l565V3rrO5sKUlHe9LStZiRM3YA/viewform?usp=pp_url&entry.995903952="
	UrlTelegrammChat         = "https://web.telegram.org/a/#"
)

func (bm BotManager) PrivateOnly(handler bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		if update.Message != nil && update.Message.Chat.Type != "private" {
			return
		}

		handler(ctx, b, update)
	}
}

func (bm BotManager) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func (bm BotManager) PayHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	redirectURL, err := bm.ic.AskApi()
	if err != nil {
		bm.Errorf("Ошибка при вызове InvoiceBox API: %v", err)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Произошла ошибка при создании счёта. Попробуйте позже.",
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Счёт успешно создан! Перейдите по ссылке для оплаты:\n%s", redirectURL),
	})
}
func (bm BotManager) CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("Error answering callback: %v", err)
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
		bm.Errorf("Error sending response: %v", err)
	}
}

func (bm BotManager) Register(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("Error answering callback: %v", err)
		return
	}

	var kb *models.InlineKeyboardMarkup
	switch update.CallbackQuery.Data {
	case PatternRegister + "Teen":

		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "Пройти регистрацию",
						URL:          UrlRegisterStudent + strconv.FormatInt(update.CallbackQuery.From.ID, 10),
						CallbackData: PatternSubmitModeration + update.CallbackQuery.Data,
					},
				},
			},
		}

	case PatternRegister + "Business":

		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "Пройти регистрацию",
						URL:          UrlRegisterBusiness + strconv.FormatInt(update.CallbackQuery.From.ID, 10),
						CallbackData: PatternSubmitModeration + update.CallbackQuery.Data,
					},
				},
			},
		}

	default:
		bm.Errorf("Сломалась регистрация")
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        "Пожалуйста, заполните данные о себе в этой форме",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Error sending response: %v", err)
	}
}

func (bm BotManager) ModerationStudent(ctx context.Context, b *bot.Bot, update *models.Update) {

	if b == nil {
		log.Println("Ошибка: бот не инициализирован (nil)")
		return
	}

	if update == nil || update.CallbackQuery == nil {
		log.Println("Ошибка: некорректный update объект")
		return
	}

	var data StudentData
	if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &data); err != nil {
		bm.Errorf("Ошибка парсинга JSON: %v\nДанные: %s", err, update.CallbackQuery.Data)

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: bm.adminChatID,
			Text:   "Ошибка обработки данных студента",
		}); err != nil {
			bm.Errorf("Ошибка отправки сообщения: %v", err)
		}

		return
	}

	userID := data.Tgid

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Принять",
					CallbackData: PatternAction + "accept_" + userID + "_" + "Teen",
				},
			},
			{
				{
					Text:         "Отклонить",
					CallbackData: PatternAction + "reject_" + userID + "_" + "Teen",
				},
			},
		},
	}
	response := fmt.Sprintf(ResponseStudentModeration,
		data.Name, data.Birthday, data.City, data.Skill, data.Email)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      bm.adminChatID,
		Text:        response,
		ParseMode:   "Markdown",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        "Твоя заявка принята!\nПозвоним тебе для подтверждения и в течение часа подтвердим твою регистрацию в сервисе",
		ParseMode:   "Markdown",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

}

func (bm BotManager) ModerationBusines(ctx context.Context, b *bot.Bot, update *models.Update) {
	if b == nil {
		log.Println("Ошибка: бот не инициализирован (nil)")
		return
	}

	if update == nil || update.CallbackQuery == nil {
		log.Println("Ошибка: некорректный update объект")
		return
	}

	var data BusinesData
	if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &data); err != nil {
		bm.Errorf("Ошибка парсинга JSON: %v\nДанные: %s", err, update.CallbackQuery.Data)

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: bm.adminChatID,
			Text:   "Ошибка обработки данных студента",
		}); err != nil {
			bm.Errorf("Ошибка отправки сообщения: %v", err)
		}

		return
	}

	userID := data.Tgid

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Принять",
					CallbackData: PatternAction + "accept_" + userID + "_" + "Business",
				},
			},
			{
				{
					Text:         "Отклонить",
					CallbackData: PatternAction + "reject_" + userID + "_" + "Business",
				},
			},
		},
	}

	response := fmt.Sprintf(ResponceBusinessModeration,
		data.CompanyName, data.INN, data.FieldOfActivity, data.ContactPersonFullName, data.ContactPersonPhone)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      bm.adminChatID,
		Text:        response,
		ParseMode:   "Markdown",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        "Ваша заявка отправлена на модерацию\nВ течение часа вернемся с результатом",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}
}

func (bm BotManager) ModerationResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("Error answering callback: %v", err)
		return
	}
	parts := strings.Split(update.CallbackQuery.Data, "_")

	if len(parts) < 4 {
		bm.Errorf("не удалось отобразить карточку пользователя, len(parts) < 4\n")
		return
	}

	actionID, err := strconv.Atoi(parts[2])
	if err != nil {
		bm.Errorf("Проблема с ID: %v", err)
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
							CallbackData: PatternCreateTask,
						},
						{
							Text:         "Позже",
							CallbackData: PatternLater,
						},
					},
				},
			}
		case "Teen":
			response = "Твои данные подтверждены!\nГотовимся отправить тебе первое задание!"
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
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      responceAdmin,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

}

func (bm BotManager) ViewTasks(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
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
		bm.Errorf("%v", err)
	}

}

func (bm BotManager) StudentReadiness(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
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
		bm.Errorf("%v", err)
	}
}

func (bm BotManager) Call(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
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
		Text:        "Отлично!\nДавай назначим созвон для выяснения деталей, затем ты сможешь приступить",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("%v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: bm.adminChatID,
		Text:   "Запрос на новый созвон от пользователя!",
	})
	if err != nil {
		bm.Errorf("%v", err)
	}

}

func (bm BotManager) NotReady(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
	}

	parts := strings.Split(update.CallbackQuery.Data, "_")

	if len(parts) < 2 {
		bm.Errorf("не удалось отобразить карточку пользователя, len(parts) < 2\n")
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
		bm.Errorf("%v", err)
	}
}

func (bm BotManager) CreateTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "Создать лот",
					URL:  UrlCreateTask + strconv.FormatInt(update.CallbackQuery.From.ID, 10),
				},
			},
		},
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        "Данные о лоте отправлены модераторам",
		ReplyMarkup: kb,
	})

}

func (bm BotManager) ModerationTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	if b == nil {
		log.Println("Ошибка: бот не инициализирован (nil)")
		return
	}

	if update == nil || update.CallbackQuery == nil {
		log.Println("Ошибка: некорректный update объект")
		return
	}

	var data TaskData
	if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &data); err != nil {
		bm.Errorf("Ошибка парсинга JSON: %v\nДанные: %s", err, update.CallbackQuery.Data)

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: bm.adminChatID,
			Text:   "Ошибка обработки данных студента",
		}); err != nil {
			bm.Errorf("Ошибка отправки сообщения: %v", err)
		}

		return
	}

	response := fmt.Sprintf(ResponceTaskModeration,
		data.NameTask, data.Direction, data.Description, data.Deadline, data.SlotCall)
	var kb *models.InlineKeyboardMarkup
	if data.Link != "" {
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "Файлы к лоту",
						URL:  data.Link,
					},
					{
						Text: "создатель",
						URL:  UrlTelegrammChat + data.TgId,
					},
				},
			},
		}
	}

	params := &bot.SendMessageParams{
		ChatID:    bm.adminChatID,
		Text:      response,
		ParseMode: "Markdown",
	}

	if kb != nil {
		params.ReplyMarkup = kb
	}

	_, err := b.SendMessage(ctx, params)
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}
}

func (bm BotManager) Later(ctx context.Context, b *bot.Bot, update *models.Update) {

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
	}

	bm.Errorf("%v", update.CallbackQuery.Message.Message.Chat.ID)

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      "Ок!\nКак будете готовы, выберите в меню пункт \"Разместить задание\"",
	})

	newCmd := models.BotCommand{
		Command:     "place_task",
		Description: "Создать задание",
	}

	_, err = b.SetChatMenuButton(ctx, &bot.SetChatMenuButtonParams{
		ChatID:     update.CallbackQuery.Message.Message.Chat.ID,
		MenuButton: models.MenuButtonCommands{Type: "commands"},
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}

	commands, err := b.GetMyCommands(ctx, &bot.GetMyCommandsParams{})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}

	commands = append(commands, newCmd)

	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}
}
func (bm BotManager) VerificationTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
	}

	bm.Errorf("%v", update.CallbackQuery.Message.Message.Chat.ID)

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      "Ок!\nКак будет готово, выберите в меню пункт \"Готово для проверки\"",
	})

	newCmd := models.BotCommand{
		Command:     "ready_verification",
		Description: "Готово для проверки",
	}

	_, err = b.SetChatMenuButton(ctx, &bot.SetChatMenuButtonParams{
		ChatID:     update.CallbackQuery.Message.Message.Chat.ID,
		MenuButton: models.MenuButtonCommands{Type: "commands"},
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}

	commands, err := b.GetMyCommands(ctx, &bot.GetMyCommandsParams{})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}

	commands = append(commands, newCmd)

	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}
}

func (bm BotManager) VerificationRequest(ctx context.Context, b *bot.Bot, update *models.Update) {

	kbAdmin := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Готово",
					CallbackData: PatternTaskCheckResponse + "_completed_" + strconv.FormatInt(update.Message.From.ID, 10),
				},
				{
					Text:         "Отправить на доработку",
					CallbackData: PatternTaskCheckResponse + "_revision_" + strconv.FormatInt(update.Message.From.ID, 10),
				},
			},
		},
	}

	kbBusiness := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Назначить созвон для проверки",
					CallbackData: PatternCall,
				},
			},
		},
	}
	nameTask := "Название задания"
	businessID := int64(1098511932)

	response := fmt.Sprintf(RequestTaskVerification,
		nameTask, businessID, update.Message.From.ID)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      bm.adminChatID,
		Text:        response,
		ParseMode:   "Markdown",
		ReplyMarkup: kbAdmin,
	})

	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      businessID,
		Text:        response,
		ParseMode:   "Markdown",
		ReplyMarkup: kbBusiness,
	})

	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}
}

func (bm BotManager) ResponseVerificationTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
	}

	parts := strings.Split(update.CallbackQuery.Data, "_")

	if len(parts) < 4 {
		bm.Errorf("не удалось отобразить карточку пользователя, len(parts) < 4\n")
		return
	}
	var response string
	switch parts[2] {
	case "completed":
		response = "Задание проверено - все ок, но нужно кое-что доработать!"
	case "revision":
		response = "Принято!\nЗаказчик принял твою работу! Ожидай оплаты)"
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: parts[3],
		Text:   response,
	})

	if err != nil {
		bm.Errorf("%v", err)
	}

}
