ALTER TABLE resume_generation_queue
ADD CONSTRAINT resume_gen_job_posting_id_key UNIQUE (job_posting_id);
