package app

import (
	"github.com/labstack/echo/v4"
)

const (
	RouteInvoice = "/invoicebox-webhook"
)

func (a *App) HandleConfirmation(c echo.Context) error {

	paymentStatus, chatId := a.icWh.HandleConfirmation(c.Response(), c.Request())
	a.Printf("paymentStatus: %s\n", paymentStatus)
	a.bm.PayStatusHandler(c.Request().Context(), a.b, paymentStatus, chatId)
	return nil
}
