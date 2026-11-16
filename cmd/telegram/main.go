package main

import (
	"context"
	"crypto/rand"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yohanc3/resumemaxxer/internal/config"
	"github.com/yohanc3/resumemaxxer/internal/llm/deepseek"
	resumehandler "github.com/yohanc3/resumemaxxer/internal/resume-handler"
	"github.com/yohanc3/resumemaxxer/internal/storage/db"
	objectstorage "github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
	"github.com/yohanc3/resumemaxxer/internal/telegram"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		slog.Error("Cannot load env variables. Exiting observer.", slog.String("error", err.Error()))
		return
	}

	db, err := db.GetDB()
	if err != nil {
		slog.Error("Cannot get db connection", slog.String("error", err.Error()))
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	r2, err := objectstorage.NewR2(context.Background(), config.Cfg.R2BucketName, config.Cfg.R2AccessKeyID, config.Cfg.R2SecretAccessKey, config.Cfg.R2AccountID)
	if err != nil {
		slog.Error("error when creating new r2 client", slog.String("error", err.Error()))
		return
	}

	llmHandler, err := deepseek.NewDeepseekClient(config.Cfg.DeepSeekAPIURL, config.Cfg.DeepSeekAPIKey)
	if err != nil {
		slog.Error("Error when initializing new deepseek client", slog.String("error", err.Error()))
		return
	}

	resHandler := resumehandler.NewResumeHandler(db, r2, llmHandler)

	bot, err := tgbotapi.NewBotAPI(config.Cfg.TelegramApiKey)
	if err != nil {
		slog.Error("error when starting up telegram bot", slog.String("error", err.Error()))
		return
	}
	bot.Debug = true
	telegramHandler := telegram.NewTelegramHandler(db, r2, bot, resHandler)

	path := "/" + rand.Text()

	wh, err := tgbotapi.NewWebhook(config.Cfg.TelegramWebhookURL + path)
	if err != nil {
		slog.Error("error when initializing telegram bot webhook", slog.String("error", err.Error()))
		return
	}

	_, err = bot.Request(wh)
	if err != nil {
		slog.Error("error when attributing webhook to bot", slog.String("error", err.Error()))
		return
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		slog.Error("error when requesting webhook info", slog.String("error", err.Error()))
		return
	}

	slog.Info("webhook information successfully obtained", slog.String("webhook url", info.URL))

	updates := bot.ListenForWebhook(path)
	go func() {
		slog.Info("listening on port :8443")
		err = http.ListenAndServe("0.0.0.0:8443", nil)
		if err != nil {
			slog.Error("error when serving on port :8443", slog.String("error", err.Error()))
			return
		}
	}()

	if err := telegramHandler.Run(ctx, updates); err != nil {
		slog.ErrorContext(ctx, "error when running telegram handler", slog.String("error", err.Error()))
		return
	}
}
