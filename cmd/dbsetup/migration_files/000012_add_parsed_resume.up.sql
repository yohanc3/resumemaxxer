CREATE TABLE IF NOT EXISTS parsed_resumes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscriber_id UUID NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    parsed_resume JSONB,
    FOREIGN KEY (subscriber_id) REFERENCES subscribers(id)
);
