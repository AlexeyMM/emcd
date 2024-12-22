CREATE TABLE sent_email_messages (
    id      uuid default gen_random_uuid() primary key,
    email text,
    type text,
    created_at timestamp
);