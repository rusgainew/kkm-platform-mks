-- Event Log for Idempotency
-- Tracks processed events to prevent duplicate processing

CREATE TABLE IF NOT EXISTS event_log (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL UNIQUE,
    event_type VARCHAR(100) NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for fast event_id lookup
CREATE INDEX IF NOT EXISTS idx_event_log_event_id ON event_log(event_id);
CREATE INDEX IF NOT EXISTS idx_event_log_event_type ON event_log(event_type);
CREATE INDEX IF NOT EXISTS idx_event_log_processed_at ON event_log(processed_at DESC);

-- Comment
COMMENT ON TABLE event_log IS 'Tracks processed events for idempotency in CQRS read models';
COMMENT ON COLUMN event_log.event_id IS 'Unique event identifier from RabbitMQ message';
COMMENT ON COLUMN event_log.event_type IS 'Event type for categorization';
