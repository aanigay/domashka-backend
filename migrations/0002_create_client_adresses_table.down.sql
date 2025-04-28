-- Удаление индексов
DROP INDEX IF EXISTS idx_unique_address_per_client;
DROP INDEX IF EXISTS idx_client_addresses_geom;

-- Удаление таблиц
DROP TABLE IF EXISTS client_addresses;
