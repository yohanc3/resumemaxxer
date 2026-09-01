ALTER TABLE resume_generation_queue 
ADD COLUMN status TEXT NOT NULL DEFAULT 'enqueued'
CHECK (status IN ('enqueued', 'processing_resume', 'completed_resume', 'delivering', 'delivered', 'failed'));
