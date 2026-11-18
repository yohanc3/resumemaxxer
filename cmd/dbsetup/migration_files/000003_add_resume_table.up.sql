CREATE TABLE IF NOT EXISTS resumes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscriber_id UUID NOT NULL,
    object_key VARCHAR(100) NOT NULL,
    filename VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(20) CHECK (status IN ('pending', 'archived')),
    FOREIGN KEY (subscriber_id) REFERENCES subscribers(id)
);
