ALTER TABLE users
    ADD COLUMN IF NOT EXISTS secret_key text default '' NOT NULL;