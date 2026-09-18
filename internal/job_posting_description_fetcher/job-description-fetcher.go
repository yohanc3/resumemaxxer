package jobpostingdescriptionfetcher

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type JobDescriptionFetcher struct {
	db      *sql.DB
	sem     chan struct{}
	wg      sync.WaitGroup
	fetcher *Scraper
}

type MissingJobDescriptionURL struct {
	url            string
	job_posting_id string
}

func NewJobDescriptionFetcher(db *sql.DB, SEMAPHORE_LENGTH int) *JobDescriptionFetcher {
	return &JobDescriptionFetcher{db: db, sem: make(chan struct{}, SEMAPHORE_LENGTH), fetcher: NewScraper()}
}

func (j *JobDescriptionFetcher) Start(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			if err := j.processMissingJobDescriptions(ctx); err != nil {
				slog.ErrorContext(ctx, "error when processing missing job descriptions", slog.String("error", err.Error()))
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (j *JobDescriptionFetcher) processMissingJobDescriptions(ctx context.Context) error {

	missingJobDescriptions, err := j.getMissingJobDescriptionsURLs(ctx)
	if err != nil {
		return err
	}

	for _, m := range missingJobDescriptions {
		go func() {
			j.fetchJobDescription(ctx, m)
		}()
	}

	return nil
}

func (j *JobDescriptionFetcher) getMissingJobDescriptionsURLs(ctx context.Context) ([]*MissingJobDescriptionURL, error) {

	tx, err := j.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error when beginning tx before fetching job descriptions")
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
			WITH descriptions_jobs AS (
				SELECT j.url, q.id AS queue_id, q.job_posting_id
				FROM job_description_fetch_queue AS q
				JOIN job_postings AS j ON j.id = q.job_posting_id
				WHERE q.status = 'enqueued'
				ORDER BY q.created_at
				LIMIT $1
				FOR UPDATE OF q SKIP LOCKED
			)
			UPDATE job_description_fetch_queue AS q
			SET status = 'processing'
			FROM descriptions_jobs AS d
			WHERE q.id = d.queue_id
			RETURNING d.url, d.job_posting_id
		`, cap(j.sem))

	if err != nil {
		return nil, fmt.Errorf("error when fetching job description queue jobs. %w", err)
	}

	var missingJobDescriptions []*MissingJobDescriptionURL

	for rows.Next() {

		missingJobDescription := &MissingJobDescriptionURL{}

		if err := rows.Scan(&missingJobDescription.url, &missingJobDescription.job_posting_id); err != nil {
			return nil, fmt.Errorf("error when parsing row: %w", err)
		}

		missingJobDescriptions = append(missingJobDescriptions, missingJobDescription)

	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error caught after scanning all rows: %w", err)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close queue rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit queue jobs: %w", err)
	}

	return missingJobDescriptions, nil
}

func (j *JobDescriptionFetcher) fetchJobDescription(ctx context.Context, jobDescriptionURL *MissingJobDescriptionURL) error {
	slog.Log(ctx, slog.LevelInfo, "fetching job description", slog.String("job data: ", fmt.Sprintf("%+v", jobDescriptionURL)))

	bodyText, err := j.fetcher.GetPageText(ctx, jobDescriptionURL.url)
	if err != nil {
		return fmt.Errorf("error when fetching page text. error: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("url body: \n %v", bodyText))

	return nil
}
