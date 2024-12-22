CREATE TABLE email_settings_password_restoration (
    whitelabel_id uuid,
    provider text,
    sender text,
    title text,
    body text,
    login text,
    password text,
    domain text,
    api_key text,
    PRIMARY KEY (whitelabel_id, provider)
);