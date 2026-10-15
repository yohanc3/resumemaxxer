package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/yohanc3/resumemaxxer/internal/telegram"
)

func main() {
	if err := telegram.SendText(context.Background(), os.Getenv("TELEGRAM_BOT_TOKEN"),
		os.Getenv("TELEGRAM_CHAT_ID"), "Hello from ResumeMaxxer!"); err != nil {
		slog.Error("could not send Telegram message", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Telegram message sent successfully")
}
