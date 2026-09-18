package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yohanc3/resumemaxxer/internal/config"
	jobpostingdescriptionfetcher "github.com/yohanc3/resumemaxxer/internal/job_posting_description_fetcher"
	"github.com/yohanc3/resumemaxxer/internal/storage/db"
)

var SEMAPHORE_LENGTH = 3

func main(){

	err := config.LoadConfig()
	if err != nil {
		slog.Error("error when loading .env config", slog.String("error: ", err.Error()))
		return
	}

	db, err := db.GetDB()
	if err != nil {
		slog.Error("error when establishing db connection", slog.String("error: ", err.Error()))
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)	
	defer stop()

	jobPostingDescriptionFetcher := jobpostingdescriptionfetcher.NewJobDescriptionFetcher(db, SEMAPHORE_LENGTH)

	if err := jobPostingDescriptionFetcher.Start(ctx); err != nil {
		slog.ErrorContext(ctx, "error when job posting description fetcher was running", slog.String("error", err.Error()))
	}

}
