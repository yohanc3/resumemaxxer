package jobdelegator

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type jobDelegator struct {
	db                *sql.DB
	workflowsEntryURL string
	interval          time.Duration
}

type queueJob struct {
	id             string
	job_posting_id string
	user_id        string
	company_name   string
	retries        int
	status         string
	created_at     int64
	fulfilled_at   int64
}

func NewJobDelegator(db *sql.DB, workflowsEntryURL string) *jobDelegator {
	return &jobDelegator{
		db:                db,
		workflowsEntryURL: workflowsEntryURL,
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

	var queueJobs []queueJob

	for rows.Next() {

		var q queueJob
		err := rows.Scan(&q.id, &q.job_posting_id, &q.user_id, &q.company_name,
			&q.retries, &q.status, &q.created_at, &q.fulfilled_at)

		if err != nil {
			return fmt.Errorf("error when scanning job from resume_queue_jobs table. %w", err.Error())
		}

		queueJobs = append(queueJobs, q)

	}

}
