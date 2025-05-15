-- 1. Создание таблицы
CREATE TABLE IF NOT EXISTS reviews (
                         id BIGSERIAL,
                         chef_id BIGINT NOT NULL,
                         user_id BIGINT NOT NULL,
                         stars SMALLINT NOT NULL CHECK (stars BETWEEN 1 AND 5),
                         comment VARCHAR(255) NOT NULL,
                         is_verified BOOLEAN NOT NULL DEFAULT FALSE,
                         include_in_rating BOOLEAN NOT NULL DEFAULT FALSE,
                         is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Комментарии к столбцам
COMMENT ON COLUMN reviews.id              IS 'Уникальный идентификатор отзыва';
COMMENT ON COLUMN reviews.chef_id         IS 'Идентификатор шеф-повара, которого оценивают';
COMMENT ON COLUMN reviews.user_id         IS 'Идентификатор пользователя, оставившего отзыв';
COMMENT ON COLUMN reviews.stars           IS 'Оценка в виде количества звёзд (1–5)';
COMMENT ON COLUMN reviews.comment         IS 'Текстовый комментарий отзыва';
COMMENT ON COLUMN reviews.is_verified     IS 'Флаг верификации отзыва';
COMMENT ON COLUMN reviews.include_in_rating IS 'Флаг включения отзыва в расчёт рейтинга шеф-повара';
COMMENT ON COLUMN reviews.is_deleted      IS 'Флаг мягкого удаления отзыва';
COMMENT ON COLUMN reviews.created_at      IS 'Время создания отзыва';
COMMENT ON COLUMN reviews.updated_at      IS 'Время последнего обновления отзыва';

-- 3. Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_reviews_chef_include
    ON reviews (chef_id)
    WHERE include_in_rating;

-- для быстрого получения всех отзывов пользователя
CREATE INDEX IF NOT EXISTS idx_reviews_user
    ON reviews (user_id);