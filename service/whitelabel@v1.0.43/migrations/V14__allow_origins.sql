
CREATE TABLE allow_origins (
    user_id INTEGER NOT NULL DEFAULT 0,
    origin VARCHAR NOT NULL DEFAULT '',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now())
)