CREATE TABLE IF NOT EXISTS resume_queue_jobs_dlq (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resume_queue_job_id UUID NOT NULL,
    job_posting_id UUID NOT NULL
);
