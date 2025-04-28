ALTER TABLE IF EXISTS users
    DROP COLUMN IF EXISTS last_sms_request,
    DROP COLUMN IF EXISTS sms_attempts,
    DROP COLUMN IF EXISTS is_spam;
