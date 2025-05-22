package app

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
)

const (
	RouteInvoice = "/invoicebox-webhook"
)

func (a *App) HandleConfirmation(c echo.Context) error {
	a.icWh.HandleConfirmation(c.Response(), c.Request())
	fmt.Printf("paymentStatus: %s\n", a.icWh.PaymentStatus)
	a.bm.PayStatusHandler(context.Background(), a.b, a.icWh.PaymentStatus, a.icWh.TgChatID)
	return nil
}
