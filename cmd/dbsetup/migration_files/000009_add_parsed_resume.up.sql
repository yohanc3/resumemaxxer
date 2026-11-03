CREATE TABLE IF NOT EXISTS parsed_resumes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    parsed_resume JSONB NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
