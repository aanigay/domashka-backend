ALTER TABLE IF EXISTS chefs
    DROP COLUMN IF EXISTS is_block,
    DROP COLUMN IF EXISTS  is_archive,
    DROP COLUMN IF EXISTS  is_self_employed,
    DROP COLUMN IF EXISTS  role_id;