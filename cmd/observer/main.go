package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"errors"

	"github.com/yohanc3/resumemaxxer/internal/config"
	"github.com/yohanc3/resumemaxxer/internal/observer"
	"github.com/yohanc3/resumemaxxer/internal/storage/db"
)

func main() {
	
	config.LoadConfig()

	obs := &observer.Observer{
		Interval: time.Minute * 5,
		URL:      "https://raw.githubusercontent.com/SimplifyJobs/Summer2026-Internships/refs/heads/dev/.github/scripts/listings.json",
	}
	
	// DB setup
	db, err := db.GetDB()
	if err != nil {
		slog.Error("error when getting db", slog.String("error", err.Error()))
		return
	}

	obs.DB = db
	
	// Initializing observer
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	if err := obs.Watch(ctx); err != nil && errors.Is(err, context.Canceled){
		slog.Error("error when watching with observer", slog.String("error", err.Error()))
	}

}
