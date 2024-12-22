alter table invoice
    add column created_at  timestamptz NOT NULL default NOW(),
    add column finished_at timestamptz;