ALTER TABLE profile
    ADD COLUMN created_by uuid;

ALTER TABLE profile_history
    ADD COLUMN created_by uuid;