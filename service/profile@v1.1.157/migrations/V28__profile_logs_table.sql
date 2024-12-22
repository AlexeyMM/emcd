CREATE TABLE IF NOT EXISTS profile_logs (
    id BIGSERIAL PRIMARY KEY,
    originator TEXT NOT NULL,
    change_type TEXT NOT NULL,
    details TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
)
