package jobdelegator

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	resumebuilder "github.com/yohanc3/resumemaxxer/internal/resume-builder"
)

var SEMAPHORE_LENGTH = 3

type jobDelegator struct {
	db            *sql.DB
	interval      time.Duration
	resumeBuilder *resumebuilder.ResumeBuilder
	wg            sync.WaitGroup
	sem           chan struct{}
}

func NewJobDelegator(db *sql.DB, interval time.Duration, resumeBuilder *resumebuilder.ResumeBuilder) *jobDelegator {
	return &jobDelegator{
		db:            db,
		interval:      interval,
		resumeBuilder: resumeBuilder,
		sem:           make(chan struct{}, SEMAPHORE_LENGTH),
	}
}

func (j *jobDelegator) Start(ctx context.Context) error {

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := j.DelegateJobs(ctx); err != nil {
				slog.ErrorContext(ctx, "error when delegating jobs", slog.String("error", err.Error()))
			}
		}
	}

}

func (j *jobDelegator) DelegateJobs(ctx context.Context) error {

	tx, err := j.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "error when starting transaction for delegating jobs")
		return fmt.Errorf("error when starting transaction for delegating jobs. error: %w", err)
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		WITH selected_jobs AS (
			SELECT id	
			FROM resume_generation_queue
			WHERE status = 'enqueued'
			ORDER BY created_at
			LIMIT 3
			FOR UPDATE SKIP LOCKED
		)
		UPDATE resume_generation_queue AS queue 
		SET status = 'processing'	
		FROM selected_jobs
		WHERE queue.id = selected_jobs.id
		RETURNING queue.id, queue.job_posting_id, queue.user_id, queue.company_name, 
				  queue.retries, queue.status, queue.created_at, queue.fulfilled_at;
	`)

	if err != nil {
		return fmt.Errorf("error when getting jobs from resume_queue_jobs table. %w", err)
	}

	var queueJobs []*resumebuilder.QueueJob

	for rows.Next() {

		q := &resumebuilder.QueueJob{}

		err := rows.Scan(&q.ID, &q.JobPostingID, &q.UserID, &q.CompanyName,
			&q.Retries, &q.Status, &q.CreatedAt, &q.FulfilledAt)

		if err != nil {
			slog.Log(ctx, slog.LevelError, fmt.Sprintf("error when scanning job from resume_queue_jobs table. %v", err.Error()))
			continue
		}

		queueJobs = append(queueJobs, q)

	}

	if err = rows.Close(); err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("errors when scanning rows. %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, q := range queueJobs {

		// Trigger resume creation step
		j.sem <- struct{}{}
		j.wg.Add(1)

		go func(q *resumebuilder.QueueJob) {
			defer j.wg.Done()
			defer func() { <-j.sem }()

			err := j.resumeBuilder.CreateResume(ctx, q)
			if err != nil {
				j.OnJobError(ctx, q)
			}

			// resumeBuilder.CreateResume should return the metadata so we can
			// push the completed job to a new queue. for now, it returns an error
		}(q)

	}

	j.wg.Wait()


	slog.Log(ctx, slog.LevelInfo, fmt.Sprintf("successfully enqueued %v jobs: %+v", len(queueJobs), queueJobs))
	return nil

}

func (j *jobDelegator) OnJobError(ctx context.Context, q *resumebuilder.QueueJob) {
	// if job errored out >3 times, insert into dead letter queue and alert me
	// else, mark it as not_processed
	// if anything fails, just log the queue object and the error itself
	if q.Retries >= 3 {
		slog.Log(ctx, slog.LevelInfo, fmt.Sprintf("queue job retried 3+ times (%v). inserting into dead letter queue: %+v", q.Retries, q))

		tx, err := j.db.BeginTx(ctx, nil)
		if err != nil {
			slog.ErrorContext(ctx, "error when starting transaction", slog.String("error", err.Error()))
			return
		}

		defer tx.Rollback()

		_, err = tx.ExecContext(ctx, `
			INSERT INTO resume_queue_jobs_dlq(
				resume_queue_job_id, job_posting_id, failed_at
			)
			VALUES($1, $2, NOW());
		`,
			q.ID, q.JobPostingID,
		)
		if err != nil {
			slog.ErrorContext(ctx, "error when inserting into resume_queue_jobs_dlq.", 
				slog.String("queue job: ", fmt.Sprintf("%+v", q)),
				slog.String("error", err.Error()),
			)
			return
		}

		if err = tx.Commit(); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error when committing transaction %v", err.Error()), slog.String("queue job", fmt.Sprintf("%+v", q)))
			return
		}

		if err != nil {
			slog.Log(ctx, slog.LevelError, fmt.Sprintf("error when inserting queue job into dlq: %+v - error: %v", q, err.Error()))
			return
		}
	} else {
		slog.Log(ctx, slog.LevelInfo, fmt.Sprintf("queue job retried %v times. retrying job... - q: %+v", q.Retries, q))
		go func() {
			// increment queue job retry counter and recursively call j.DelegateJobs(ctx, q[with retries += 1])
			_, err := j.db.ExecContext(ctx, `
				UPDATE resume_generation_queue
				SET retries = retries + 1, status = 'enqueued'
				WHERE id = $1
			`, q.ID)
			if err != nil {
				slog.ErrorContext(ctx, "error when re-queueing errored job", slog.String("err", err.Error()))
			}
		}()
	}
}
