ALTER TABLE email_settings_password_restoration
    DROP CONSTRAINT email_settings_password_restoration_pkey;

ALTER TABLE email_settings_password_restoration
    ADD PRIMARY KEY (whitelabel_id);