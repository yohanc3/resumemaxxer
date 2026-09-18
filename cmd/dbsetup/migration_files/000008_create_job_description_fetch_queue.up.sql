CREATE TABLE IF NOT EXISTS job_description_fetch_queue (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
  job_posting_id UUID NOT NULL UNIQUE,
  status TEXT NOT NULL CHECK (status IN ('enqueued', 'processing', 'processed')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  fulfilled_at TIMESTAMPTZ DEFAULT NULL,
  FOREIGN KEY (job_posting_id) REFERENCES job_postings (id)
)
