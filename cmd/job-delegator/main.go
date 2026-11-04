package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yohanc3/resumemaxxer/internal/config"
	jobdelegator "github.com/yohanc3/resumemaxxer/internal/job-delegator"
	"github.com/yohanc3/resumemaxxer/internal/llm/deepseek"
	resumehandler "github.com/yohanc3/resumemaxxer/internal/resume-handler"
	"github.com/yohanc3/resumemaxxer/internal/storage/db"
	objectstorage "github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
)

var SEMAPHORE_LENGTH = 3

func main() {

	err := config.LoadConfig()
	if err != nil {
		slog.Error("Cannot load env variables. Exiting observer.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db, err := db.GetDB()
	if err != nil {
		slog.Error("error when getting db", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	storage, err := objectstorage.NewR2(ctx, config.Cfg.R2BucketName, config.Cfg.R2AccessKeyID, config.Cfg.R2SecretAccessKey, config.Cfg.R2AccountID)
	if err != nil {
		slog.Error("error when getting r2 objectstorage object", slog.String("error", err.Error()))
		return
	}

	deepseekClient, err := deepseek.NewDeepseekClient(config.Cfg.DeepSeekAPIURL, config.Cfg.DeepSeekAPIKey) 
	if err != nil {
		slog.ErrorContext(ctx, "error when creating deepseek client", slog.String("error", err.Error()))
	}

	resumebuilder := resumehandler.NewResumeHandler(db, storage, deepseekClient)

	jobDelegator := jobdelegator.NewJobDelegator(db, time.Second * time.Duration(config.Cfg.DelegatorIntervalSeconds), &resumebuilder, SEMAPHORE_LENGTH)

	if err := jobDelegator.Start(ctx); err != nil {
		slog.Error("error when running job handler", slog.String("error", err.Error()))
	}

}
