-- Удаление уникального индекса на основе географической точки и chef_id
DROP INDEX IF EXISTS idx_unique_address_per_chef;

-- Удаление GIST индекса для поля geom
DROP INDEX IF EXISTS idx_chef_addresses_geom;

-- Удаление таблицы chef_addresses
DROP TABLE IF EXISTS chef_addresses;

-- Удаление расширения PostGIS (если оно не используется другими объектами)
DROP EXTENSION IF EXISTS postgis;
