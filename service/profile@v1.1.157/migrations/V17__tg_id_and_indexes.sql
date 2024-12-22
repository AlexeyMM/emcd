alter table notification_settings alter column
    tg_id TYPE bigint;

CREATE UNIQUE INDEX users_username_uindex ON users USING btree (LOWER(username));

CREATE UNIQUE INDEX users_email_uindex ON users USING btree (LOWER(email));
