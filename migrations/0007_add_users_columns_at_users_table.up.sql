-- Добавление колонки username, если она отсутствует
ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS username VARCHAR(255);

-- Добавление колонки name, если она отсутствует
ALTER TABLE IF EXISTS  users
    ADD COLUMN IF NOT EXISTS name VARCHAR(255);

-- Добавление колонки chat_id, если она отсутствует
ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS chat_id VARCHAR(255);
