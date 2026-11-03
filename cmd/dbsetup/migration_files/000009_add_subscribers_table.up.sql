CREATE TABLE IF NOT EXISTS subscribers (
    id UUID PRIMARY KEY,
    chat_id VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(), 
);
