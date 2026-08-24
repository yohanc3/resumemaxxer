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

var SEMAPHORE_LENGTH = 3

func main(){

	err := config.LoadConfig()
	if err != nil {
		panic("Cannot load env variables. Exiting job-handler...")
	}

	db, err := db.GetDB()
	if err != nil {
		slog.Error("error when getting db", slog.String("error", err.Error()))
		return
	}

	storage := objectstorage.NewR2(config.Cfg.DBName, config.Cfg.DBPassword)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	shutdownContext, cancel := context.WithTimeout(ctx, time.Second * 10) 
	defer cancel()
	
	c := make(chan struct{}, SEMAPHORE_LENGTH)
	resumebuilder := resumebuilder.NewResumeBuilder(c, db, storage)

	jobDelegator := jobdelegator.NewJobDelegator(db, time.Second * 10, &resumebuilder)

	if err := jobDelegator.Start(shutdownContext); err != nil {
		slog.Error("error when running job handler", slog.String("error", err.Error()))
	}

}
