package main

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	cfg, err := LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	b, err := bot.New(cfg.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, myStartHand)

	ctx := context.Background()
	log.Println("Бот стартовал")
	b.Start(ctx)
}

func myStartHand(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello world",
	})
}
