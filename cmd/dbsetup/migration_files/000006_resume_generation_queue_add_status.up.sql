ALTER TABLE resume_generation_queue 
ADD COLUMN status TEXT NOT NULL DEFAULT 'queued'
CHECK (status IN ('queued', 'processing', 'completed', 'failed'));
