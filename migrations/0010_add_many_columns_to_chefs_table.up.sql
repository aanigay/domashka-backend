ALTER TABLE IF EXISTS chefs
    ADD COLUMN IF NOT EXISTS role_id int,
    ADD COLUMN IF NOT EXISTS is_self_employed bool,
    ADD COLUMN IF NOT EXISTS is_archive bool,
    ADD COLUMN IF NOT EXISTS is_block bool;

COMMENT ON COLUMN chefs.role_id IS 'Идентификатор роли шефа (связь с таблицей roles)';
COMMENT ON COLUMN chefs.is_self_employed IS 'Флаг, указывающий на то, что шеф-повар работает как самозанятый';
COMMENT ON COLUMN chefs.is_archive IS 'Флаг, указывающий на активность профиля шеф-повара';
COMMENT ON COLUMN chefs.is_block IS 'Флаг, определяющий, заблокирован профиль шеф-повара';

