package app

import (
	"context"
	"fmt"
	"log"
	"lotBot/pkg/invoicebox"
	"net/http"
	"time"

	"lotBot/pkg/db"
	"lotBot/pkg/embedlog"
	botLogic "lotBot/pkg/lotBot/bot"

	"github.com/go-pg/pg/v10"
	"github.com/go-telegram/bot"
	"github.com/labstack/echo/v4"
)

// Config describes .toml file structure
type Config struct {
	Database *pg.Options
	Server   struct {
		Host      string
		Port      int
		IsDevel   bool
		EnableVFS bool
	}
	Bot struct {
		Token       string
		AdminChatID int
	}
	InvoiceConfig invoicebox.Config
}

type App struct {
	embedlog.Logger
	appName string
	cfg     Config
	db      db.DB
	dbc     *pg.DB
	echo    *echo.Echo
	b       *bot.Bot
	bm      *botLogic.BotManager
	bot     *bot.Bot
	ic      *invoicebox.InvoiceClient
}

func New(appName string, verbose bool, cfg Config, db db.DB, dbc *pg.DB) *App {
	a := &App{
		appName: appName,
		cfg:     cfg,
		db:      db,
		dbc:     dbc,
		echo:    echo.New(),
	}
	a.SetStdLoggers(verbose)
	a.echo.HideBanner = true
	a.echo.HidePort = true
	a.echo.IPExtractor = echo.ExtractIPFromRealIPHeader()

	a.bm = botLogic.NewBotManager(a.db, a.Logger, a.cfg.Bot.AdminChatID, a.cfg.InvoiceConfig)

	b, err := bot.New(cfg.Bot.Token)
	if err != nil {
		panic(err)
	}
	a.b = b

	a.ic = invoicebox.NewInvoiceClient(a.Logger, a.cfg.InvoiceConfig)

	return a
}

// Run is a function that runs application.
func (a *App) Run() error {
	a.registerMetrics()
	a.registerHandlers()
	a.registerBotHandlers()
	a.registerDebugHandlers()
	a.registerAPIHandlers()
	go a.b.Start(context.Background())

	_, err := a.ic.AskApi()
	if err != nil {
		return err
	}

	go func() {
		invoiceboxHandler := invoicebox.NewWebhookHandler(a.db, a.Logger)

		http.HandleFunc("/invoicebox-webhook", invoiceboxHandler.HandleConfirmation)
		fmt.Println("Webhook port 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	return a.runHTTPServer(a.cfg.Server.Host, a.cfg.Server.Port)

}

// Shutdown is a function that gracefully stops HTTP server.
func (a *App) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := a.echo.Shutdown(ctx); err != nil {
		a.Errorf("shutting down server err=%q", err)
	}
}
