CREATE TABLE whitelabel_coins (
    wl_id       UUID            NOT NULL,
    coin_id     VARCHAR(10)     NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at  TIMESTAMP       NOT NULL DEFAULT TIMEZONE('utc'::text, now()),

    PRIMARY KEY (wl_id, coin_id)
);

-- adding some test data
INSERT INTO whitelabel_coins
    (wl_id, coin_id)
VALUES
    ('c3addbb7-5a09-4e70-8497-14c4d7011c26', 'btc');

INSERT INTO whitelabel_coins
    (wl_id, coin_id)
VALUES
    ('c3addbb7-5a09-4e70-8497-14c4d7011c26', 'kas');

INSERT INTO whitelabel_coins
    (wl_id, coin_id)
VALUES
    ('c3addbb7-5a09-4e70-8497-14c4d7011c44', 'dash')