package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"uuid"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yohanc3/resumemaxxer/internal/config"
	"github.com/yohanc3/resumemaxxer/internal/llm/deepseek"
	resumehandler "github.com/yohanc3/resumemaxxer/internal/resume-handler"
	"github.com/yohanc3/resumemaxxer/internal/storage/db"
	objectstorage "github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		slog.Error("Cannot load env variables. Exiting observer.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db, err := db.GetDB()
	if err != nil {
		slog.Error("Cannot get db connection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	r2, err := objectstorage.NewR2(context.Background(), config.Cfg.R2BucketName, config.Cfg.R2AccessKeyID, config.Cfg.R2SecretAccessKey, config.Cfg.R2AccountID)
	if err != nil {
		slog.Error("error when creating new r2 client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	llmHandler, err := deepseek.NewDeepseekClient(config.Cfg.DeepSeekAPIURL, config.Cfg.DeepSeekAPIKey)
	if err != nil {
		slog.Error("Error when initializing new deepseek client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	resHandler := resumehandler.NewResumeHandler(db, r2, llmHandler)

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
	go func() {
		slog.Info("listening on port :8443")
		err = http.ListenAndServe("0.0.0.0:8443", nil)
		if err != nil {
			slog.Error("error when serving on port :8443", slog.String("error", err.Error()))
		}
	}()

	for update := range updates {
		j, err := json.Marshal(update)
		if err != nil {
			slog.Error("error when marshaling update struct", slog.String("error", err.Error()))
			continue
		}
		slog.Info("message received", slog.String("message", fmt.Sprintf("message: %v", string(j))))

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message != nil && update.Message.Document == nil {

			switch update.Message.Command() {
			case "start":
				msg.Text = "Upload your resume to begin🚨"
			}

			if msg.Text != "" {
				bot.Send(msg)
			}

		} else if update.Message != nil && update.Message.Document != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hold tight while I process your resume..."))
			doc := update.Message.Document

			if doc.MimeType != "application/pdf" {
				msg.Text = "Please upload your resume in PDF format"
				continue
			}

			url, err := bot.GetFileDirectURL(doc.FileID)
			if err != nil {
				slog.Error("error when getting file url", slog.String("error", err.Error()))
				continue
			}

			resp, err := http.Get(url)
			if err != nil {
				slog.Error("error when downloading user document", slog.String("error", err.Error()))
				continue
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				slog.Error("document download status isn't OK", slog.Int("status code", resp.StatusCode))
				continue
			}
			defer resp.Body.Close()

			bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("error when reading file bytes", slog.String("error", err.Error()))
				continue
			}

			subscriberID := uuid.New().String()

			parsedResume, err := resHandler.ParseResume(context.Background(), bytes)
			if err != nil {
				slog.Error("parsing resume went wrong", slog.String("error", err.Error()))
				continue
			}

			fileName, objectKey := objectstorage.CreateFileMetadata(subscriberID, parsedResume.Header.FirstName, parsedResume.Header.LastName)

			eTag, versionID, err := r2.Put(context.Background(), bytes, objectKey, "application/pdf")
			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				slog.Error("document download status isn't OK", slog.Int("status code", resp.StatusCode))
				continue
			}

			slog.Info("file uploaded successfully", slog.String("etag", eTag), slog.String("version_id", versionID))

			parsedResumeJSON, err := json.MarshalIndent(parsedResume, "", "  ")
			if err != nil {
				slog.Error("error when marshaling parsed resume struct", slog.String("error", err.Error()))
				continue
			}

			slog.Info("parsed resume", slog.String("parsed", string(parsedResumeJSON)))

			msg.Text = fmt.Sprintf("uploaded resume successfully\n objectKey: %s\n filename: %s\n eTag: %s\n", objectKey, fileName, eTag)
			if _, err := bot.Send(msg); err != nil {
				slog.Error("error when sending bot message", slog.String("error", err.Error()))
				continue
			}

		}

	}
}
