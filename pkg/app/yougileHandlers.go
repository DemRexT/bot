package app

import (
	"encoding/json"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"io"
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
	// Временная структура для проверки поля "event"
	var payload struct {
		Event string `json:"event"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Ошибка разбора JSON",
		})
	}

	switch payload.Event {
	case "task-updated":
		a.bm.ViewTasks(c.Request().Context(), a.b, update)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Данные переданы на модерацию"})
}
