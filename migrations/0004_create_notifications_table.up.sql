CREATE TABLE IF NOT EXISTS notifications
(
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id       BIGINT,
    channel       VARCHAR(50) NOT NULL,
    scenario      VARCHAR(50) NOT NULL,
    subject       VARCHAR(255),
    message       TEXT,
    recipient     VARCHAR(255),
    status        VARCHAR(50) NOT NULL CHECK (status IN ('created', 'sent', 'error', 'retry')),
    error_message TEXT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    send_attempts INT       DEFAULT 0,
    metadata      JSONB
);

-- Indexes for optimized queries
CREATE INDEX idx_notifications_user_id ON notifications (user_id);
CREATE INDEX idx_notifications_created_at ON notifications (created_at);
CREATE INDEX idx_notifications_status ON notifications (status);
CREATE INDEX idx_notifications_channel ON notifications (channel);
