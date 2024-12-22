
CREATE TABLE frontend_configs (
    ref_id INTEGER NOT NULL DEFAULT 0,
    origin VARCHAR NOT NULL DEFAULT '',
    media_id VARCHAR NOT NULL DEFAULT '',
    colors jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now())
)