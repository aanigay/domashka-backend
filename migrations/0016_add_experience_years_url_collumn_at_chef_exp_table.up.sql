ALTER TABLE chefs_experience
    ADD COLUMN experience_years int NULL DEFAULT NULL;

COMMENT ON COLUMN chefs_experience.experience_years IS 'Количество лет опыта';
