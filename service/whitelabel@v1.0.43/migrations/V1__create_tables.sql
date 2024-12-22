CREATE TABLE whitelabel (
    id uuid,
    name text,
    description text,
    fee numeric,
    whitelabel_fee numeric,
    reward numeric,
    logo text,
    PRIMARY KEY (id)
);

CREATE TABLE whitelabel_history (
    whitelabel_id uuid,
    fee numeric,
    whitelabel_fee numeric,
    reward numeric,
    created_at timestamp,
    PRIMARY KEY (whitelabel_id, created_at)
);