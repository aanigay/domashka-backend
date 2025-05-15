BEGIN;

-- 1. Создаём sequence, если её нет
CREATE SEQUENCE IF NOT EXISTS reviews_id_seq
    OWNED BY reviews.id;

-- 2. Устанавливаем текущий курсор sequence на максимальное существующее id
SELECT setval('reviews_id_seq', COALESCE(MAX(id), 1)) FROM reviews;

-- 3. Приводим тип колонки к BIGINT на всякий случай
ALTER TABLE reviews
    ALTER COLUMN id TYPE BIGINT
        USING id::BIGINT;

-- 4. Добавляем DEFAULT nextval для автозаполнения
ALTER TABLE reviews
    ALTER COLUMN id SET DEFAULT nextval('reviews_id_seq');

COMMIT;