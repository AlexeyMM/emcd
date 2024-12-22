ALTER TABLE whitelabel
    ADD COLUMN is_two_fa_enabled        BOOLEAN DEFAULT FALSE,
    ADD COLUMN is_captcha_enabled       BOOLEAN DEFAULT FALSE,
    ADD COLUMN is_email_confirm_enabled BOOLEAN DEFAULT FALSE;