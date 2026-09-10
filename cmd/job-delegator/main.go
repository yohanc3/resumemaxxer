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
	resumebuilder "github.com/yohanc3/resumemaxxer/internal/resume-builder"
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
		slog.Error("error when getting db", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	storage := objectstorage.NewR2(config.Cfg.R2BucketName, config.Cfg.R2AccessKeyID, config.Cfg.R2SecretAccessKey, config.Cfg.R2AccountID)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	resumebuilder := resumebuilder.NewResumeBuilder(db, storage)

	jobDelegator := jobdelegator.NewJobDelegator(db, time.Second * time.Duration(config.Cfg.DelegatorIntervalSeconds), &resumebuilder)

	if err := jobDelegator.Start(ctx); err != nil {
		slog.Error("error when running job handler", slog.String("error", err.Error()))
	}

}
