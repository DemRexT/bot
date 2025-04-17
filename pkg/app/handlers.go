package app

import (
	"fmt"
	"github.com/go-telegram/bot/models"
	"io"
	"net/http"
	_ "net/http/pprof"

	"lotBot/pkg/rpc"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vmkteam/rpcgen/v2"
	zm "github.com/vmkteam/zenrpc-middleware"
	"github.com/vmkteam/zenrpc/v2"
)

// runHTTPServer is a function that starts http listener using labstack/echo.
func (a *App) runHTTPServer(host string, port int) error {
	listenAddress := fmt.Sprintf("%s:%d", host, port)
	a.Printf("starting http listener at http://%s\n", listenAddress)

	return a.echo.Start(listenAddress)
}

func (a *App) registerHandlers() {
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{"Authorization", "Authorization2", "Origin", "X-Requested-With", "Content-Type", "Accept", "Platform", "Version"},
	}))

	// sentry middleware
	a.echo.Use(sentryecho.New(sentryecho.Options{
		Repanic:         true,
		WaitForDelivery: true,
	}))

	a.echo.Use(zm.EchoIPContext(), zm.EchoSentryHubContext())
}

// registerDebugHandlers adds /debug/pprof handlers into a.echo instance.
func (a *App) registerDebugHandlers() {
	dbg := a.echo.Group("/debug")

	// add pprof integration
	dbg.Any("/pprof/*", func(c echo.Context) error {
		if h, p := http.DefaultServeMux.Handler(c.Request()); p != "" {
			h.ServeHTTP(c.Response(), c.Request())
			return nil
		}
		return echo.NewHTTPError(http.StatusNotFound)
	})

	a.echo.GET("/status", func(c echo.Context) error {
		// test postgresql connection
		_, err := a.db.Exec(`SELECT 1`)
		if err != nil {
			return c.String(http.StatusInternalServerError, "DB error")
		}
		return c.String(http.StatusOK, "OK")
	})
}

func (a *App) registerAPIHandlers() {
	srv := rpc.New(a.db, a.Logger, a.cfg.Server.IsDevel)
	gen := rpcgen.FromSMD(srv.SMD())

	a.echo.Any("/formresulstudent", a.handleFormResultStudent)
	a.echo.Any("/formresultbusines", a.handleFormResultBusines)
	a.echo.Any("/formresultlot", a.handleFormResultLot)

	a.echo.Any("/v1/rpc/", zm.EchoHandler(zm.XRequestID(srv)))
	a.echo.Any("/v1/rpc/doc/", echo.WrapHandler(http.HandlerFunc(zenrpc.SMDBoxHandler)))
	a.echo.Any("/v1/rpc/openrpc.json", echo.WrapHandler(http.HandlerFunc(rpcgen.Handler(gen.OpenRPC("apisrv", "http://localhost:8075/v1/rpc")))))
	a.echo.Any("/v1/rpc/api.ts", echo.WrapHandler(http.HandlerFunc(rpcgen.Handler(gen.TSClient(nil)))))
}

func (a *App) handleFormResultStudent(c echo.Context) error {

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
	a.bm.ModerationStudent(c.Request().Context(), a.b, update)

	return c.JSON(http.StatusOK, map[string]string{"status": "Данные переданы на модерацию"})
}

func (a *App) handleFormResultBusines(c echo.Context) error {
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
	a.bm.ModerationBusines(c.Request().Context(), a.b, update)

	return c.JSON(http.StatusOK, map[string]string{"status": "Данные переданы на модерацию"})
}

func (a *App) handleFormResultLot(c echo.Context) error {
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
	a.bm.ModerationTask(c.Request().Context(), a.b, update)

	return c.JSON(http.StatusOK, map[string]string{"status": "Данные переданы на модерацию"})
}
