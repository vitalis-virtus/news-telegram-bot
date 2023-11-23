package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/vitalis-virtus/news-telegram-bot/internal/config"
	"github.com/vitalis-virtus/news-telegram-bot/internal/fetcher"
	"github.com/vitalis-virtus/news-telegram-bot/internal/notifier"
	"github.com/vitalis-virtus/news-telegram-bot/internal/storage"
	"github.com/vitalis-virtus/news-telegram-bot/internal/summary"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("[ERROR] failed create telegram bot: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("[ERROR] failed connect ot DB: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcher        = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		notifier = notifier.New(
			articleStorage,
			summary.NewOpenAISummarizer(config.Get().OpenAIKey, config.Get().OpenAIPrompt),
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().NotificationInterval,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start fetch: %v", err)
				return
			}

			log.Printf("[INFO] fetcher stopped")
		}
	}(ctx)

	// go func(ctx context.Context) {
	if err := notifier.SelectAndSendArticle(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to start notifier: %v", err)
			return
		}

		log.Printf("[INFO] notifier stopped")
	}
	// }(ctx)

}
