-- Добавление столбцов created_at и updated_at в таблицу chefs

ALTER TABLE IF EXISTS chefs
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Добавление комментариев к новым столбцам
COMMENT ON COLUMN chefs.created_at IS 'Дата и время создания записи';
COMMENT ON COLUMN chefs.updated_at IS 'Дата и время последнего обновления записи';
