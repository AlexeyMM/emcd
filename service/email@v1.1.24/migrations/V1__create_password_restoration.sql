CREATE TABLE smtp_settings (
    username text,
    password text,
    server_address text,
    server_port int,
    from_addr text,
    whitelabel_id uuid,
    PRIMARY KEY (whitelabel_id)
);