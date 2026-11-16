package telegram

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
	"uuid"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	resumehandler "github.com/yohanc3/resumemaxxer/internal/resume-handler"
	objectstorage "github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
)

type TelegramHandler struct {
	db            *sql.DB
	objectStorage objectstorage.ObjectStorage
	bot           *tgbotapi.BotAPI
	resumeHandler *resumehandler.ResumeHandler
	wg            sync.WaitGroup
	sem           chan struct{}
}

type SignupQueueJob struct {
	subscriberID        string
	chatID              int64
	telegramDocumentURL string
	status              string
	objectKey           string
	filename            string
}

func NewTelegramHandler(db *sql.DB, objectStorage objectstorage.ObjectStorage, bot *tgbotapi.BotAPI, resumeHandler *resumehandler.ResumeHandler) *TelegramHandler {
	return &TelegramHandler{db: db, objectStorage: objectStorage, bot: bot, resumeHandler: resumeHandler}
}

func (t *TelegramHandler) Run(ctx context.Context, updatesChannel tgbotapi.UpdatesChannel) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	enqueueDone := make(chan error, 1)
	go func() {
		enqueueDone <- t.enqueueSignups(ctx, updatesChannel)
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-enqueueDone:
			return err

		case <-ticker.C:
			if err := t.handleSignups(ctx); err != nil {
				slog.ErrorContext(ctx, "error processing signups", slog.String("error", err.Error()))
			}
		}

	}
}

func (t *TelegramHandler) enqueueSignups(ctx context.Context, updates tgbotapi.UpdatesChannel) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case update, ok := <-updates:
			if !ok {
				return nil
			}

			// Skip if message doesnt contain a document
			if update.Message == nil || update.Message.Document == nil {
				slog.ErrorContext(ctx, "error when sending message: 'message type not supported': %w", slog.String("message", fmt.Sprintf("%+v", update)))
				continue
			}

			chatID := update.Message.Chat.ID
			fileURL, err := t.bot.GetFileDirectURL(update.Message.Document.FileID)
			if err != nil {
				return fmt.Errorf("error when getting direct file url when enqueing update: %w", err)
			}

			subscriberID := uuid.New()

			tx, err := t.db.BeginTx(ctx, nil)
			if err != nil {
				return err
			}
			defer tx.Rollback()

			_, err = tx.ExecContext(ctx, "INSERT INTO subscribers (id, chat_id) VALUES ($1, $2)", subscriberID, chatID)
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, "INSERT INTO subscription_queue (subscriber_id, telegram_document_url) VALUES ($1, $2)", subscriberID, fileURL)
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, "INSERT INTO parsed_resumes (subscriber_id, parsed_resume) VALUES ($1, NULL)", subscriberID)
			if err != nil {
				return err
			}

			if err = tx.Commit(); err != nil {
				return err
			}

		}

		return nil
	}
}

func (t *TelegramHandler) handleSignups(ctx context.Context) error {

	rows, err := t.db.QueryContext(ctx, `
		SELECT q.subscriber_id, q.status, q.telegram_document_url, s.chat_id, COALESCE(r.object_key, ''), COALESCE(r.filename, '')
		FROM subscription_queue as q 
		LEFT JOIN resumes AS r ON r.subscriber_id = q.subscriber_id
		JOIN subscribers AS s ON s.id = q.subscriber_id
		WHERE q.status != 'finished'`,
	)
	if err != nil {
		return fmt.Errorf("error when selecting signup jobs: %w", err)
	}

	defer rows.Close()

	var signupQueueJobs []*SignupQueueJob

	for rows.Next() {
		signupQueueJob := &SignupQueueJob{}

		err := rows.Scan(&signupQueueJob.subscriberID, &signupQueueJob.status, &signupQueueJob.telegramDocumentURL, &signupQueueJob.chatID, &signupQueueJob.objectKey, &signupQueueJob.filename)

		if err != nil {
			return fmt.Errorf("error when scanning row: %w", err)
		}

		signupQueueJobs = append(signupQueueJobs, signupQueueJob)
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	for _, signupQueueJob := range signupQueueJobs {
		slog.InfoContext(ctx, "handling update", slog.String("document url", signupQueueJob.telegramDocumentURL))
		err := t.handleSignup(ctx, signupQueueJob)
		if err != nil {
			slog.ErrorContext(ctx, "error when handling signup", slog.String("error", err.Error()))
		}
	}

	return nil
}

func (t *TelegramHandler) handleSignup(ctx context.Context, s *SignupQueueJob) error {

	for {
		var err error

		switch s.status {

		// Possible statuses are: archive, parse_resume_text, notify_user,
		//
		// fyi: each function is responsible for its own cleanup in case of a failure
		case "archive":
			err = t.archiveResume(ctx, s)

		case "parse_resume_text":
			err = t.parseResume(ctx, s)

		case "notify_user":
			err = t.notifyUser(ctx, s)

		case "finished":
			slog.InfoContext(ctx, "successfully signed up a new user", slog.String("subscriberID", s.subscriberID))
			return nil

		default:
			return fmt.Errorf("unknown step: %s", s.status)
		}

		if err != nil {
			return fmt.Errorf("error when completing signup: %w", err)
		}

		// Persists outputs and the next step atomically.
		// If this fails, we don't execute the next step.
		if err := t.checkpoint(ctx, s); err != nil {
			return err
		}

	}

}

func (t *TelegramHandler) archiveResume(ctx context.Context, s *SignupQueueJob) error {

	if s.filename == "" && s.objectKey == "" {
		s.filename, s.objectKey = objectstorage.CreateFileMetadata(s.subscriberID)
	}

	_, err := t.db.ExecContext(ctx, "INSERT INTO resumes (subscriber_id, object_key, filename, status) VALUES ($1, $2, $3, 'pending')", s.subscriberID, s.objectKey, s.filename)
	if err != nil {
		return err
	}

	resp, err := http.Get(s.telegramDocumentURL)
	if err != nil {
		return fmt.Errorf("error when downloading user document: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return errors.New("document download status isn't OK")
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error when reading file bytes: %w", err)
	}

	_, err = t.objectStorage.Put(ctx, bytes, s.objectKey, "application/pdf")
	if err != nil {
		return fmt.Errorf("error when storing resume")
	}

	_, err = t.db.ExecContext(ctx, "UPDATE resumes SET status = 'archived' WHERE subscriber_id = $1", s.subscriberID)
	if err != nil {
		return err
	}

	_, err = t.db.ExecContext(ctx, "UPDATE subscription_queue SET status = 'parse_resume_text' WHERE subscriber_id = $1", s.subscriberID)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "successfully archived user's resume", slog.String("subscriber_id", s.subscriberID), slog.String("object key", s.objectKey))

	return nil

}

func (t *TelegramHandler) parseResume(ctx context.Context, s *SignupQueueJob) error {

	resp, err := http.Get(s.telegramDocumentURL)
	if err != nil {
		slog.Error("error when downloading user document", slog.String("error", err.Error()))
		return fmt.Errorf("error when downloading user document") 
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return errors.New("document download status isn't OK")
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error when reading file bytes: %w", err)
	}

	parsedResume, err := t.resumeHandler.ParseResume(ctx, bytes)
	if err != nil {
		return fmt.Errorf("parsing resume went wrong: %w", err)
	}

	parsedResumeJSONB, err := json.Marshal(parsedResume)
	if err != nil {
		return err
	}

	_, err = t.db.ExecContext(ctx, "UPDATE parsed_resumes SET parsed_resume = $1 WHERE subscriber_id = $2", parsedResumeJSONB, s.subscriberID)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "successfully parsed user's resume", slog.String("subscriber_id", s.subscriberID))

	return nil
}

func (t *TelegramHandler) notifyUser(ctx context.Context, s *SignupQueueJob) error {

	_, err := t.bot.Send(tgbotapi.NewMessage(s.chatID, "You've been successfully signed up for new internship postings!"))
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "successfully sent subscription notification to user", slog.String("subscriber_id", s.subscriberID))
	return nil
}

func (t *TelegramHandler) checkpoint(ctx context.Context, s *SignupQueueJob) error {

	switch s.status {

	case "archive":
		s.status = "parse_resume_text"

	case "parse_resume_text":
		s.status = "notify_user"

	case "notify_user":
		s.status = "finished"
	}

	_, err := t.db.ExecContext(ctx, "UPDATE subscription_queue SET status = $1 WHERE subscriber_id = $2", s.status, s.subscriberID)
	if err != nil {
		return err
	}

	return nil
}
