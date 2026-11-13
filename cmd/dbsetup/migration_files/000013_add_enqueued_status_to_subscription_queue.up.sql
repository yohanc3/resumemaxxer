ALTER TABLE subscription_queue
DROP CONSTRAINT subscription_queue_status_check;

ALTER TABLE subscription_queue
ADD CONSTRAINT subscription_queue_status_check
CHECK (status IN (
    'enqueued',
    'archive',
    'parse_resume_text',
    'notify_user',
    'finished'
));
