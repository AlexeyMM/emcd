ALTER TABLE IF EXISTS users
    ADD COLUMN pool_type varchar(20) DEFAULT 'emcd';
