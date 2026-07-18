CREATE TABLE IF NOT EXISTS resume_generation_queue(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_posting_id VARCHAR(40) UNIQUE NOT NULL,
    job_posting_company_name TEXT NOT NULL,
    job_posting_is_active BOOLEAN NOT NULL,
    retries INTEGER DEFAULT 0 NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    fulfilled_at TIMESTAMPTZ DEFAULT NULL
);
