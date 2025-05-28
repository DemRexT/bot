package app

import (
	"encoding/json"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"io"
	"lotBot/pkg/lotBot/bot"
	"net/http"
)

const (
	RouteYougile = "/yougile"
)

func (a *App) handleYougileResult(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Ошибка чтения тела запроса",
		})
	}

	update := &models.Update{
		CallbackQuery: &models.CallbackQuery{
			Data: string(body),
		},
	}

	var task struct {
		Event   string `json:"event"`
		Payload struct {
			ColumnId string `json:"columnId"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(body, &task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Ошибка разбора JSON",
		})
	}

	switch task.Event {
	case "task-updated":
		a.bm.ViewTasks(c.Request().Context(), a.b, update)
	case "task-moved":
		if task.Payload.ColumnId == bot.ColumnInProgress {
			a.bm.Printf("Сработало")
			a.bm.VerificationTask(c.Request().Context(), a.b, update)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Данные переданы на модерацию"})
}
