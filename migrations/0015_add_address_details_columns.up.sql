ALTER TABLE IF EXISTS client_addresses
    ADD COLUMN IF NOT EXISTS flat_number TEXT,
    ADD COLUMN IF NOT EXISTS floor_number TEXT,
    ADD COLUMN IF NOT EXISTS entrance_number TEXT,
    ADD COLUMN IF NOT EXISTS intercom_number TEXT,
    ADD COLUMN IF NOT EXISTS courier_comment TEXT;