ALTER TABLE whitelabel
    DROP COLUMN fee,
    DROP COLUMN whitelabel_fee,
    DROP COLUMN reward;


ALTER TABLE whitelabel_history
    ADD COLUMN coin text;

ALTER TABLE whitelabel_history
    DROP CONSTRAINT whitelabel_history_pkey;

ALTER TABLE whitelabel_history
    ADD PRIMARY KEY (whitelabel_id,coin,created_at);

CREATE TABLE whitelabel_commissions(
    whitelabel_id uuid,
    fee numeric,
    whitelabel_fee numeric,
    reward numeric,
    coin text,
    PRIMARY KEY (whitelabel_id,coin),
    CONSTRAINT fk_wl_id
        FOREIGN KEY (whitelabel_id)
            REFERENCES whitelabel(id)
);
