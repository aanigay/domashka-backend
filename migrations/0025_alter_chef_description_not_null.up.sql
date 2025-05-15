UPDATE chefs
SET description = ''
WHERE description IS NULL;

ALTER TABLE IF EXISTS chefs
ALTER COLUMN description SET NOT NULL,
ALTER COLUMN description SET DEFAULT '';