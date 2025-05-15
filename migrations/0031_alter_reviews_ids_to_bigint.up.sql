BEGIN;

-- 1. Убираем DEFAULT (связь с sequence)
ALTER TABLE reviews
    ALTER COLUMN id DROP DEFAULT;

-- 2. Явно приводим тип к BIGINT
ALTER TABLE reviews
    ALTER COLUMN id TYPE BIGINT
        USING id::BIGINT;

-- 3. Опционально: удаляем неиспользуемую sequence
DROP SEQUENCE IF EXISTS reviews_id_seq;

COMMIT;