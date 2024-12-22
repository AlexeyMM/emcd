CREATE TABLE push_tokens (
    user_id uuid,
    device_id int,
    token text,
    PRIMARY KEY (user_id, device_id)
);