ALTER TABLE chefs
    ADD COLUMN legal_info varchar(255) NULL DEFAULT NULL;

COMMENT ON COLUMN chefs.legal_info IS 'Юредискиая информация о шефе';
