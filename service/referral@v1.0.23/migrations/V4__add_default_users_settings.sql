CREATE TABLE IF NOT EXISTS referral.default_settings_referrals
(
    referral_id  UUID        NOT NULL,
    product      VARCHAR(20) NOT NULL,
    coin         VARCHAR(20) NOT NULL,
    fee          DECIMAL     NOT NULL DEFAULT 0,
    referral_fee DECIMAL     NOT NULL DEFAULT 0,
    created_at   TIMESTAMP   NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at   TIMESTAMP   NOT NULL DEFAULT TIMEZONE('utc'::text, now()),

    PRIMARY KEY (referral_id, product, coin)
);