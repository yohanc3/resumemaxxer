package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yohanc3/resumemaxxer/internal/config"
	"github.com/yohanc3/resumemaxxer/internal/storage/db"
)

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	shutdownContext, cancel := context.WithTimeout(ctx, time.Second * 10) 
	defer cancel()

	jobHandler := jobHandler.NewJobHandler(shutdownContext, db)

	if err := jobHandler.Start(); err != nil {
		slog.Error("error when running job handler", slog.String("error", err.Error()))
	}

}
