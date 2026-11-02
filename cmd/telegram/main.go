package main

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yohanc3/resumemaxxer/internal/config"
)


func main(){

	err := config.LoadConfig()

	if err != nil {
		slog.Error("Cannot load env variables. Exiting observer.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(config.Cfg.TelegramApiKey)
	if err != nil {
		slog.Error("error when starting up telegram bot", slog.String("error", err.Error()))
		os.Exit(1)
	}

	bot.Debug = true
	
	path := "/" + rand.Text()

	wh, err := tgbotapi.NewWebhook(config.Cfg.TelegramWebhookURL + path)
	if err != nil {
		slog.Error("error when initializing telegram bot webhook", slog.String("error", err.Error()))
		os.Exit(1)
	}

	_, err = bot.Request(wh)
	if err != nil {
		slog.Error("error when attributing webhook to bot", slog.String("error", err.Error()))
		os.Exit(1)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		slog.Error("error when requesting webhook info", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("webhook information successfully obtained", slog.String("webhook url", info.URL))

	updates := bot.ListenForWebhook(path) 
	go http.ListenAndServe("0.0.0.0:8443", nil)

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "gotcha big dawg")

		bot.Send(msg)
	}
}
