ALTER TABLE IF EXISTS chefs
ALTER COLUMN description DROP NOT NULL,
ALTER COLUMN description SET DEFAULT NULL;