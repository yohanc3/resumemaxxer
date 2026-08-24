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
	}
}

func (j *jobDelegator) Start(ctx context.Context) error {

	ticker := time.NewTicker(j.interval)

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			j.DelegateJobs(ctx)
		}
	}

}

func (j *jobDelegator) DelegateJobs(ctx context.Context) error {

	rows, err := j.db.QueryContext(ctx, `SELECT * FROM resume_queue_jobs;`)
	if err != nil {
		return fmt.Errorf("error when getting jobs from resume_queue_jobs table. %w", err.Error())
	}

	for rows.Next() {

		var q *resumebuilder.QueueJob

		err := rows.Scan(&q.ID, &q.JobPostingID, &q.UserID, &q.CompanyName,
			&q.Retries, &q.Status, &q.CreatedAt, &q.FulfilledAt)

		if err != nil {
			slog.Log(ctx, slog.LevelError, fmt.Sprintf("error when scanning job from resume_queue_jobs table. %w", err.Error()))
		}

		go func() {

			j.wg.Add(1)
			defer j.wg.Done()

			j.sem <- struct{}{}
			defer func() { <-j.sem }()

			_, err := j.db.ExecContext(ctx, `
			UPDATE resume_queue_jobs
			SET status = 'pending'
			WHERE resume_queue_jobs.id = $1
			`, q.ID)

			if err != nil {
				slog.Log(ctx, slog.LevelError, fmt.Sprintf("error when updating queue job to pending. %w", err.Error()))
			}

			err = j.resumeBuilder.CreateResume(ctx, q)
			if err != nil {
				j.OnJobError(ctx, q)
			}

		}()

	}

	j.wg.Wait()

	return nil

}

func (j *jobDelegator) OnJobError(ctx context.Context, q *resumebuilder.QueueJob) {
	// if job errored out >3 times, insert into dead letter queue and alert me
	// else, mark it as not_processed
	// if anything fails, just log the queue object and the error itself
	if q.Retries >= 3 {
		slog.Log(ctx, slog.LevelDebug, fmt.Sprintf("queue job retried 3+ times (%v). inserting into dead letter queue: %+v", q.Retries, q))

		_, err := j.db.ExecContext(ctx, `
			INSERT INTO resume_queue_jobs_dlq(
				resume_queue_job_id, job_posting_id, user_id, retries, failed_at
			)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			`,
			q.ID, q.JobPostingID, q.UserID, q.Retries, time.Now().Unix(),
		)

		if err != nil {
			slog.Log(ctx, slog.LevelError, fmt.Sprintf("error when inserting queue job into dlq: %+v.error: %w", q, err.Error()))
			return
		}
	} else {
		slog.Log(ctx, slog.LevelDebug, fmt.Sprintf("queue job retried %v times. retrying job... - q: %+v", q.Retries, q))
		go func(){
			// increment queue job retry counter and recursively call j.DelegateJobs(ctx, q[with retries += 1])
		}()
	}
}
