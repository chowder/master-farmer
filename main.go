package main

import (
	"context"
	"github.com/chowder/master-farmer/internal/pkg/app"
	"github.com/chowder/master-farmer/internal/pkg/handlers"
	"github.com/go-telegram/bot"
	"log"
	"os"
	"os/signal"
)

const (
	TelegramApiTokenEnvVar = "TELEGRAM_TOKEN"
	AppDsnEnvVar           = "MASTER_FARMER_DSN"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.ShowCategory),
		bot.WithCallbackQueryDataHandler(handlers.ExpandCategoryPrefix, bot.MatchTypePrefix, handlers.ExpandCategoryHandler),
		bot.WithCallbackQueryDataHandler(handlers.ReturnPrefix, bot.MatchTypePrefix, handlers.ShowCategory),
		bot.WithCallbackQueryDataHandler(handlers.EndConversationPrefix, bot.MatchTypePrefix, handlers.EndConversationHandler),
	}

	token := os.Getenv(TelegramApiTokenEnvVar)
	if token == "" {
		log.Fatalf("Environment variable '%s' was not provided", TelegramApiTokenEnvVar)
	}

	tgBot, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv(AppDsnEnvVar)
	if dsn == "" {
		dsn = "app.db"
	}

	a := app.New(ctx, tgBot, dsn)

	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/offset", bot.MatchTypePrefix, a.SetOffsetHandler)
	tgBot.RegisterHandler(bot.HandlerTypeCallbackQueryData, handlers.CreateTimerPrefix, bot.MatchTypePrefix, a.CreateTimerHandler)
	tgBot.RegisterHandler(bot.HandlerTypeCallbackQueryData, handlers.CancelTimerPrefix, bot.MatchTypePrefix, a.CancelTimerHandler)

	a.Start()

	tgBot.Start(ctx)
}
