CREATE TABLE IF NOT EXISTS resume_generation_queue(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_posting_id UUID NOT NULL,
    subscriber_id UUID NOT NULL,
    company_name TEXT NOT NULL,
    job_posting_url TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'enqueued' CHECK (
        status IN (
            'enqueued',
            'processing_resume',
            'completed_resume',
            'delivering',
            'delivered',
            'failed'
        )
    ),
    retries INTEGER DEFAULT 0 NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    delivered_at TIMESTAMPTZ DEFAULT NULL,
    fulfilled_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (subscriber_id) REFERENCES subscribers(id),
    CONSTRAINT unique_subscriber_jobposting UNIQUE (subscriber_id, job_posting_id)
);
