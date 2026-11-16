CREATE TABLE IF NOT EXISTS subscription_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscriber_id UUID NOT NULL,
    status VARCHAR(25) DEFAULT 'archive' CHECK (
        status IN ('archive', 'parse_resume_text', 'notify_user', 'finished')
    ),
    telegram_document_url TEXT NOT NULL,
    FOREIGN KEY (subscriber_id) REFERENCES subscribers(id)
);
