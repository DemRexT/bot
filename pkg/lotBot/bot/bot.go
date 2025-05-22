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
	"lotBot/pkg/yougile"
	"strconv"
	"strings"
	"time"
)

type BotManager struct {
	embedlog.Logger
	adminChatID int
	ic          *invoicebox.InvoiceClient
	repo        db.LotbotRepo
	icWh        *invoicebox.WebhookHandler
	Yougile     *yougile.YougileClient
}

func NewBotManager(DB db.DB, logger embedlog.Logger, adminChatID int, cfg invoicebox.Config, ygCfg yougile.Config) *BotManager {
	return &BotManager{
		Logger:      logger,
		adminChatID: adminChatID,
		ic:          invoicebox.NewInvoiceClient(logger, cfg),
		repo:        db.NewLotbotRepo(DB),
		Yougile:     yougile.NewYougileClient(logger, ygCfg),
	}
}

const (
	PatternStart                      = "start"
	PatternRole                       = "role_"
	PatternRegister                   = "register_"
	PatternSubmitModeration           = "submit_for_moderation_"
	PatternAction                     = "action_"
	PatternViewTask                   = "view_tasks"
	PatternReady                      = "ready_"
	PatternCall                       = "call"
	PatternNot                        = "not_"
	PatternCreateTask                 = "create_task"
	PatternLater                      = "later"
	PatternTaskCheckResponse          = "check_response"
	PatternVerificationToTheRequester = "verification_requester"
	UrlRegisterStudent                = "https://docs.google.com/forms/d/e/1FAIpQLSemsbNWCx2ewY25WlvQP_baBef6RUs1jF0w1p4obb99ieXFAw/viewform?usp=pp_url&entry.1082496981="
	UrlRegisterBusiness               = "https://docs.google.com/forms/d/e/1FAIpQLSdz5iYc9UB6M3wOOrGGl-4jTywltlkl7AZgqXrNKIBqrY87mA/viewform?usp=pp_url&entry.213949143="
	UrlCreateTask                     = "https://docs.google.com/forms/d/e/1FAIpQLScQgB6T74K87rZHi8a9qi-l565V3rrO5sKUlHe9LStZiRM3YA/viewform?usp=pp_url&entry.995903952="
	UrlTelegrammChat                  = "https://web.telegram.org/a/#"
	UrlTask                           = "https://ru.yougile.com/team/005706c078bc/#"
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
	ChatID := update.Message.Chat.ID
	redirectURL, err := bm.ic.AskApi(ChatID)
	if err != nil {
		bm.Errorf("Ошибка при вызове InvoiceBox API: %v", err)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: ChatID,
			Text:   "Произошла ошибка при создании счёта. Попробуйте позже.",
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: ChatID,
		Text:   fmt.Sprintf("Счёт успешно создан! Перейдите по ссылке для оплаты:\n%s", redirectURL),
	})
}

func (bm BotManager) PayStatusHandler(ctx context.Context, b *bot.Bot, paymentStatus string, TgChatID int64) {
	ChatID := TgChatID
	//Временно:
	SurveyURL := "https://workspace.google.com/intl/ru/products/forms/"
	//SurveyURL, err := survey.handler
	fmt.Printf("TGID (handler): %d\n", ChatID)

	if paymentStatus == "success" {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: ChatID,
			Text:   fmt.Sprintf("Оплату приняли, спасибо за сотрудничество!\nПожалуйста, оцените работу сервиса: \n%s", SurveyURL),
		})
	} else {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: ChatID,
			Text:   "Оплата не прошла. Попробуйте снова или обратитесь в поддержку.",
		})
	}
}

/*func (bm BotManager) SurveyHandler(ctx context.Context, b *bot.Bot, TgChatID int64){

}
*/

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
		response = HiCompany
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Зарегистрироваться", CallbackData: PatternRegister + "Business"},
				},
			},
		}
	case PatternRole + "2":
		response = HiStudent
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
		bm.Printf("Ошибка: бот не инициализирован (nil)")
		return
	}

	if update == nil || update.CallbackQuery == nil {
		bm.Printf("Ошибка: некорректный update объект")
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

	tgid, err := strconv.ParseInt(data.Tgid, 10, 64)
	if err != nil {
		bm.Errorf("Ошибка парсинга TgID: %v", err)
		return
	}

	joinedSkill := strings.Join(data.Skill, ", ")

	student := &db.Student{
		TgID:     tgid,
		Name:     data.Name,
		Birthday: data.Birthday,
		City:     data.City,
		Scope:    joinedSkill,
		Email:    data.Email,
		StatusID: 2,
	}

	_, err = bm.repo.AddStudent(ctx, student)
	if err != nil {
		bm.Errorf("Не удалось записать в бд: %v", err)
	}

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

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      bm.adminChatID,
		Text:        response,
		ParseMode:   "Markdown",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      "Твоя заявка принята!\nПозвоним тебе для подтверждения и в течение часа подтвердим твою регистрацию в сервисе",
		ParseMode: "Markdown",
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

}

func (bm BotManager) ModerationBusines(ctx context.Context, b *bot.Bot, update *models.Update) {
	if b == nil {
		bm.Printf("Ошибка: бот не инициализирован (nil)")
		return
	}

	if update == nil || update.CallbackQuery == nil {
		bm.Printf("Ошибка: некорректный update объект")
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

	tgid, err := strconv.ParseInt(data.Tgid, 10, 64)
	if err != nil {
		bm.Errorf("Ошибка парсинга TgID: %v", err)
		return
	}

	company := &db.Company{
		Name:     data.CompanyName,
		TgID:     tgid,
		Inn:      data.INN,
		Scope:    data.FieldOfActivity,
		Phone:    data.ContactPersonPhone,
		StatusID: 2,
	}

	_, err = bm.repo.AddCompany(ctx, company)
	if err != nil {
		bm.Errorf("Не удалось записать в бд: %v", err)
	}

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

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      bm.adminChatID,
		Text:        response,
		ParseMode:   "Markdown",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   "Ваша заявка отправлена на модерацию\nВ течение часюа вернемся с результатом",
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

	tgID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		log.Printf("ошибка парсинга TgID: %v", err)
		return
	}

	var kb *models.InlineKeyboardMarkup
	var response string
	var responceAdmin string
	switch parts[3] {
	case "Business":

		search := &db.CompanySearch{
			TgID: &tgID,
		}
		pager := db.Pager{Page: 1, PageSize: 1}

		companies, err := bm.repo.CompaniesByFilters(ctx, search, pager)
		if err != nil {
			bm.Printf("ошибка поиска студента: %v", err)
			return
		}
		if len(companies) == 0 {
			bm.Printf("Студент не найден")
			return
		}

		company := companies[0]
		switch parts[1] {
		case "reject":
			responceAdmin = "Пользователь отклонен"

			company.StatusID = 3

			ok, err := bm.repo.UpdateCompany(ctx, &company, db.WithColumns("statusId"))
			if err != nil {
				bm.Printf("ошибка обновления: %v", err)
				return
			}
			if ok {
				bm.Printf("Статус Компанит успешно обновлён")
			} else {
				bm.Printf("Обновление не затронуло ни одной строки")
			}

			response = NoModeration

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

			company.StatusID = 1

			ok, err := bm.repo.UpdateCompany(ctx, &company, db.WithColumns("statusId"))
			if err != nil {
				bm.Printf("ошибка обновления: %v", err)
				return
			}
			if ok {
				bm.Printf("Статус Компанит успешно обновлён")
			} else {
				bm.Printf("Обновление не затронуло ни одной строки")
			}

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
		}
	case "Teen":
		responceAdmin = "Пользователь подтвержден"

		search := &db.StudentSearch{
			TgID: &tgID,
		}
		pager := db.Pager{Page: 1, PageSize: 1}

		students, err := bm.repo.StudentsByFilters(ctx, search, pager)
		if err != nil {
			bm.Printf("ошибка поиска студента: %v", err)
			return
		}
		if len(students) == 0 {
			bm.Printf("Студент не найден")
			return
		}

		student := students[0]
		switch parts[1] {
		case "reject":
			responceAdmin = "Пользователь отклонен"

			student.StatusID = 3

			ok, err := bm.repo.UpdateStudent(ctx, &student, db.WithColumns("statusId"))
			if err != nil {
				bm.Printf("ошибка обновления: %v", err)
				return
			}
			if ok {
				bm.Printf("Статус студента успешно обновлён")
			} else {
				bm.Printf("Обновление не затронуло ни одной строки")
			}
			response = NoModeration

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
			student.StatusID = 1

			ok, err := bm.repo.UpdateStudent(ctx, &student, db.WithColumns("statusId"))
			if err != nil {
				bm.Printf("ошибка обновления: %v", err)
				return
			}
			if ok {
				bm.Printf("Статус студента успешно обновлён")
			} else {
				bm.Printf("Обновление не затронуло ни одной строки")
			}
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

	var data TaskPurpose
	if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &data); err != nil {
		bm.Errorf("Ошибка парсинга JSON: %v\nДанные: %s", err, update.CallbackQuery.Data)
		return
	}

	response, err := bm.Yougile.GetUserByID(data.Payload.Assigned[0])
	if err != nil {
		bm.Errorf("%v", err)
	}

	tasks, err := bm.Yougile.GetTaskByID(data.Payload.Id)
	if err != nil {
		bm.Errorf("%v", err)
	}

	yougileID := &data.Payload.Id

	search := &db.TaskSearch{
		YougileID: yougileID,
	}
	pager := db.Pager{Page: 1, PageSize: 1}

	tasksDB, err := bm.repo.TasksByFilters(ctx, search, pager)
	if err != nil {
		bm.Errorf("ошибка при поиске задачи по YougileID: %v", err)
		return
	}
	if len(tasksDB) == 0 {
		bm.Errorf("Задача с YougileID=[%s] не найдена", *yougileID)
		return
	}

	taskdb := tasksDB[0]

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Готов", CallbackData: PatternReady + "yes_" + strconv.Itoa(taskdb.ID)},
				{Text: "Не готов", CallbackData: PatternReady + "not_" + strconv.Itoa(taskdb.ID)},
			},
		},
	}

	var task ResponceTask
	if err := json.Unmarshal(tasks, &task); err != nil {
		bm.Errorf("Ошибка при разборе JSON: %v", err)
		return
	}

	var user ResponceUser
	if err := json.Unmarshal(response, &user); err != nil {
		bm.Errorf("Ошибка при разборе JSON: %v", err)
		return
	}

	searchStudent := &db.StudentSearch{
		Email: &user.Email,
	}
	pagerStudent := db.Pager{Page: 1, PageSize: 1}

	students, err := bm.repo.StudentsByFilters(ctx, searchStudent, pagerStudent)
	if err != nil {
		bm.Errorf("ошибка поиска студента по email: %v", err)
		return
	}
	kbAdmin := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "Ссылка на задачу",
					URL:  UrlTask + task.IdTaskProject,
				},
			},
		},
	}

	if len(students) == 0 {
		bm.Errorf("Студент с таким email не найден")
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      bm.adminChatID,
			Text:        "Пользователь с таким email не найден",
			ReplyMarkup: kbAdmin,
		})
		if err != nil {
			bm.Errorf("%v", err)
		}
		return
	}

	student := students[0]

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      student.TgID,
		Text:        NewTask,
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

	parts := strings.Split(update.CallbackQuery.Data, "_")

	if len(parts) < 3 {
		bm.Errorf("не удалось отобразить карточку пользователя, len(parts) < 3\n")
		return
	}

	var response string
	var kb *models.InlineKeyboardMarkup
	switch parts[1] {
	case "yes":
		response = "Отлично!\nДавай назначим созвон с заказчиком длявыяснения деталей,затем ты сможеш приступить"
		kb = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Окей", CallbackData: PatternCall},
				},
			},
		}
		tgid := update.CallbackQuery.From.ID

		search := &db.StudentSearch{
			TgID: &tgid,
		}
		pager := db.Pager{Page: 1, PageSize: 1}

		students, err := bm.repo.StudentsByFilters(ctx, search, pager)
		if err != nil {
			bm.Errorf("ошибка поиска студента: %v", err)
			return
		}
		if len(students) == 0 {
			bm.Printf("Студент не найден")
			return
		}

		student := students[0]

		taskID, err := strconv.Atoi(parts[2])
		if err != nil {
			bm.Errorf("Ошибка парсинга taskID: %v", err)
			return
		}

		task, err := bm.repo.TaskByID(ctx, taskID)
		if err != nil {
			return
		}

		task.StudentID = &student.ID

		ok, err := bm.repo.UpdateTask(ctx, task, db.WithColumns("studentId"))
		if err != nil {
			bm.Errorf("ошибка обновления студента: %v", err)
			return
		}
		if ok {
			bm.Printf("Статус студента успешно обновлён")
		} else {
			bm.Printf("Обновление не затронуло ни одной строки")
		}

		student.StatusID = 2

		ok, err = bm.repo.UpdateStudent(ctx, &student, db.WithColumns("statusId"))
		if err != nil {
			bm.Errorf("ошибка обновления студента: %v", err)
			return
		}
		if ok {
			bm.Printf("Статус студента успешно обновлён")
		} else {
			bm.Printf("Обновление не затронуло ни одной строки")
		}

	case "not":
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
	var userID int64
	var chatID int64

	if update.Message != nil {
		userID = update.Message.From.ID
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.From.ID
		chatID = update.CallbackQuery.Message.Message.Chat.ID

		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})
	} else {
		bm.Errorf("CreateTask: ни Message, ни CallbackQuery не найдены")
		return
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "Создать лот",
					URL:  UrlCreateTask + strconv.FormatInt(userID, 10),
				},
			},
		},
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Пожалуйста, заполните форму по ссылке",
		ReplyMarkup: kb,
	})
	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
	}
}

func (bm BotManager) ModerationTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	if b == nil {
		bm.Printf("Ошибка: бот не инициализирован (nil)")
		return
	}

	if update == nil || update.CallbackQuery == nil {
		bm.Printf("Ошибка: некорректный update объект")
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

	payload := yougile.TaskPayload{
		Title:       data.NameTask,
		ColumnID:    "5bfcc202-886a-4457-b037-15f8d5604558",
		Description: data.Description,
		Archived:    false,
		Completed:   false,
	}

	taskID, err := bm.Yougile.CreateTask(payload)
	if err != nil {
		bm.Errorf("Ошибка создания задачи: %v", err)
	}
	bm.Printf("Создана задача с ID: %s\n", taskID)

	tgid, err := strconv.ParseInt(data.TgId, 10, 64)
	if err != nil {
		bm.Errorf("Ошибка парсинга TgID: %v", err)
		return
	}

	parsedDeadline, err := time.Parse("02.01.2006", data.Deadline) // формат должен соответствовать строке
	if err != nil {
		bm.Errorf("Ошибка парсинга даты: %v", err)
		return
	}

	search := &db.CompanySearch{
		TgID: &tgid,
	}
	pager := db.Pager{Page: 1, PageSize: 1}

	companies, err := bm.repo.CompaniesByFilters(ctx, search, pager)
	if err != nil || len(companies) == 0 {
		bm.Errorf("Ошибка при поиске компании по TgID=%d: %v", tgid, err)
		return
	}

	company := companies[0]

	budget, err := strconv.ParseFloat(data.Budget, 64)
	if err != nil {
		bm.Errorf("Ошибка при парсинге даты: %v", err)
		return
	}

	task := &db.Task{
		CompanyID:   company.ID,
		Scope:       data.Direction,
		Description: data.Description,
		Link:        data.Link,
		Deadline:    parsedDeadline,
		ContactSlot: data.SlotCall,
		StatusID:    1,
		StudentID:   nil,
		Budget:      budget,
		YougileID:   &taskID,
	}

	_, err = bm.repo.AddTask(ctx, task)
	if err != nil {
		bm.Errorf("Не удалось записать в бд: %v", err)
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

	_, err = b.SendMessage(ctx, params)
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

	businessID := int64(1098511932)

	kbAdmin := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Отправить на проверку заказчику",
					CallbackData: PatternVerificationToTheRequester + "_" + strconv.FormatInt(update.Message.From.ID, 10) + "_" + strconv.FormatInt(businessID, 10),
				},
			},
		},
	}

	nameTask := "Название задания"

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
}

func (bm BotManager) VerificationToTheRequester(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		bm.Errorf("%v", err)
	}

	parts := strings.Split(update.CallbackQuery.Data, "_")

	if len(parts) < 4 {
		bm.Errorf("не удалось отобразить карточку пользователя, len(parts) < 3\n")
		return
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

	kbAdmin := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Готово",
					CallbackData: PatternTaskCheckResponse + "_completed_" + parts[2],
				},
				{
					Text:         "Отправить на доработку",
					CallbackData: PatternTaskCheckResponse + "_revision_" + parts[2],
				},
			},
		},
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      parts[3],
		Text:        update.CallbackQuery.Message.Message.Text,
		ParseMode:   "Markdown",
		ReplyMarkup: kbBusiness,
	})

	if err != nil {
		bm.Errorf("Ошибка отправки сообщения: %v", err)
		return
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      bm.adminChatID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        update.CallbackQuery.Message.Message.Text,
		ParseMode:   "Markdown",
		ReplyMarkup: kbAdmin,
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
	bm.Printf("%v", update.CallbackQuery.Data)
	var response string
	switch parts[2] {
	case "completed":
		response = "Принято!\nЗаказчик принял твою работу! Ожидай оплаты)"
	case "revision":
		response = "Задание проверено - все ок, но нужно кое-что доработать!"
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: parts[3],
		Text:   response,
	})

	if err != nil {
		bm.Errorf("%v", err)
	}

}
