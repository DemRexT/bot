package app

import (
	"context"
	"github.com/go-telegram/bot/models"
	"lotBot/pkg/invoicebox"
	"lotBot/pkg/yougile"
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
	Yougile struct {
		Login    string
		Password string
		Token    string
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
	icWh    *invoicebox.WebhookHandler
	update  *models.Update
	Yougile *yougile.YougileClient
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

	a.bm = botLogic.NewBotManager(a.db, a.Logger, a.cfg.Bot.AdminChatID, a.cfg.InvoiceConfig, yougile.Config(a.cfg.Yougile))

	b, err := bot.New(cfg.Bot.Token)
	if err != nil {
		panic(err)
	}
	a.b = b

	a.icWh = invoicebox.NewWebhookHandler(a.db, a.Logger)

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
