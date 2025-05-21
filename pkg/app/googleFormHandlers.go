package app

import (
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

const (
	RouteSubmitStudentForm  = "/formstudent"
	RouteSubmitBusinessForm = "/formbusiness"
	RouteSubmitLotForm      = "/formlot"
)

func (a *App) handleFormResult(c echo.Context) error {
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

	switch c.Path() {
	case RouteSubmitStudentForm:
		a.bm.ModerationStudent(c.Request().Context(), a.b, update)
	case RouteSubmitBusinessForm:
		a.bm.ModerationBusines(c.Request().Context(), a.b, update)
	case RouteSubmitLotForm:
		a.bm.ModerationTask(c.Request().Context(), a.b, update)
	default:
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Неизвестный путь",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Данные переданы на модерацию"})
}
