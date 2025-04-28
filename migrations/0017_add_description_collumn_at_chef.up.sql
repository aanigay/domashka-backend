ALTER TABLE chefs
    ADD COLUMN description varchar(255) NULL DEFAULT NULL;

COMMENT ON COLUMN chefs.description IS 'BIO шефа';
