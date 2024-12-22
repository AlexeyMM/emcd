CREATE TABLE notification_settings (
                                       user_id uuid REFERENCES users (id),
                                       is_tg_notifications_on boolean,
                                       tg_id int,
                                       is_email_notifications_on boolean,
                                       is_push_notifications_on boolean,
                                       PRIMARY KEY (user_id)
);

ALTER TABLE users ADD COLUMN language text;

ALTER TABLE push_tokens ALTER COLUMN device_id TYPE text;