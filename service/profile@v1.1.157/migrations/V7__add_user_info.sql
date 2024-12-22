ALTER TABLE users
    ADD COLUMN phone text,
    ADD COLUMN email text,
    ADD COLUMN password bytea,
    ADD COLUMN created_at timestamp,
    ADD COLUMN whitelabel_id uuid,
    ADD COLUMN api_key text;

ALTER TABLE users
    DROP COLUMN vip,
    DROP COLUMN segment_id;


ALTER TABLE users
    ADD CONSTRAINT users_unq UNIQUE(email, whitelabel_id);

CREATE TABLE user_settings (
    user_id uuid,
    vip bool,
    segment_id int,
    PRIMARY KEY (user_id)
);

ALTER TABLE user_settings
    ADD CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users(id);


CREATE TABLE roles (
    user_id uuid,
    whitelabel_id uuid,
    role text,
    PRIMARY KEY (user_id,whitelabel_id)
);

ALTER TABLE roles
    ADD CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users(id);
