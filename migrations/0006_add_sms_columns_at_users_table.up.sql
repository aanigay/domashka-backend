ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS is_spam INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sms_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_sms_request TIMESTAMP NULL;
