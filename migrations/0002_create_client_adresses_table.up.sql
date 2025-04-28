-- Включаем расширение PostGIS для работы с географическими данными
CREATE EXTENSION IF NOT EXISTS postgis;
-- Создание таблицы адресов
CREATE TABLE IF NOT EXISTS client_addresses
(
    id           BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,              -- Уникальный идентификатор адреса
    client_id    INT                    NOT NULL, -- Идентификатор клиента (связь с клиентом)
    address_type VARCHAR(50)            NOT NULL, -- Тип адреса (например, домашний, рабочий)
    name         VARCHAR(255)           NOT NULL, -- Название (например, название улицы)
    full_address TEXT                   NOT NULL, -- Полный адрес (например, "ул. Ленина, д. 10, кв. 5")
    comment      TEXT,                            -- Дополнительные комментарии
    created_at   TIMESTAMP DEFAULT NOW(),         -- Дата и время создания
    updated_at   TIMESTAMP DEFAULT NOW(),         -- Дата и время последнего обновления
    geom         GEOGRAPHY(Point, 4326) NOT NULL  -- Географическая точка (широта/долгота)
);

-- Создание GIST индекса для поля geom для эффективного поиска по гео-данным
CREATE INDEX IF NOT EXISTS idx_client_addresses_geom ON client_addresses USING GIST (geom);

-- Создание уникального индекса на основе географической точки и client_id для предотвращения дублирования адресов
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_address_per_client ON client_addresses (client_id, geom);
