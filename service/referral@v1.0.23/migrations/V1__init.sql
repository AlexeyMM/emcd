CREATE SCHEMA IF NOT EXISTS referral;

CREATE TABLE IF NOT EXISTS referral.referrals
(
    user_id        UUID        NOT NULL,
    product        VARCHAR(20) NOT NULL,
    coin           VARCHAR(20) NOT NULL,
    whitelabel_id  UUID,
    fee            DECIMAL     NOT NULL DEFAULT 0,
    whitelabel_fee DECIMAL     NOT NULL DEFAULT 0,
    referral_fee   DECIMAL     NOT NULL DEFAULT 0,
    referral_id    UUID,
    created_at     TIMESTAMP   NOT NULL DEFAULT TIMEZONE('utc'::text, now()),

    PRIMARY KEY (user_id, product, coin)
);

COMMENT ON COLUMN referral.referrals.product IS 'product identifier';
COMMENT ON COLUMN referral.referrals.user_id IS 'system user identifier';
COMMENT ON COLUMN referral.referrals.whitelabel_id IS 'whitelabel identifier(optional)';
COMMENT ON COLUMN referral.referrals.coin IS 'belonging to the coin';
COMMENT ON COLUMN referral.referrals.fee IS 'fee';
COMMENT ON COLUMN referral.referrals.whitelabel_fee IS 'whitelabel fee';
COMMENT ON COLUMN referral.referrals.referral_fee IS 'referral fee';
COMMENT ON COLUMN referral.referrals.created_at IS 'time of creation';

CREATE INDEX IF NOT EXISTS t_referrals_user_id_product_id_coin_index ON referral.referrals USING btree (user_id, product, coin);

CREATE TABLE IF NOT EXISTS referral.default_settings
(
    product      VARCHAR(20) NOT NULL,
    coin         VARCHAR(20) NOT NULL,
    fee          DECIMAL     NOT NULL DEFAULT 0,
    referral_fee DECIMAL     NOT NULL DEFAULT 0,
    created_at   TIMESTAMP   NOT NULL DEFAULT TIMEZONE('utc'::text, now()),

    PRIMARY KEY (product, coin)
);

COMMENT ON COLUMN referral.default_settings.product IS 'product identifier';
COMMENT ON COLUMN referral.default_settings.coin IS 'belonging to the coin';
COMMENT ON COLUMN referral.default_settings.fee IS 'fee';
COMMENT ON COLUMN referral.default_settings.referral_fee IS 'referral fee';
COMMENT ON COLUMN referral.default_settings.created_at IS 'time of creation';

CREATE TABLE IF NOT EXISTS referral.default_whitelabel_settings
(
    product       VARCHAR(20) NOT NULL,
    whitelabel_id UUID        NOT NULL,
    coin          VARCHAR(20) NOT NULL,
    fee           DECIMAL     NOT NULL DEFAULT 0,
    referral_fee  DECIMAL     NOT NULL DEFAULT 0,
    created_at    TIMESTAMP   NOT NULL DEFAULT TIMEZONE('utc'::text, now()),

    PRIMARY KEY (product, whitelabel_id, coin)
);

COMMENT ON COLUMN referral.default_whitelabel_settings.product IS 'product identifier';
COMMENT ON COLUMN referral.default_whitelabel_settings.whitelabel_id IS 'whitelabel identifier';
COMMENT ON COLUMN referral.default_whitelabel_settings.coin IS 'belonging to the coin';
COMMENT ON COLUMN referral.default_whitelabel_settings.fee IS 'fee';
COMMENT ON COLUMN referral.default_whitelabel_settings.referral_fee IS 'referral fee';
COMMENT ON COLUMN referral.default_whitelabel_settings.created_at IS 'time of creation';

CREATE TABLE IF NOT EXISTS referral.referrals_logs
(
    user_id UUID        NOT NULL,
    product VARCHAR(20) NOT NULL,
    coin    VARCHAR(20) NOT NULL,
    history JSONB
);

COMMENT ON COLUMN referral.referrals_logs.product IS 'product identifier';
COMMENT ON COLUMN referral.referrals_logs.user_id IS 'system user identifier';
COMMENT ON COLUMN referral.referrals_logs.coin IS 'belonging to the coin';
COMMENT ON COLUMN referral.referrals_logs.history IS 'history of changes';

CREATE INDEX IF NOT EXISTS t_referrals_user_id_product_id_coin_index ON referral.referrals USING btree (user_id, product, coin);

-- function adds/updates records in table referral.t_referrals_logs using key t_referral_id
CREATE OR REPLACE FUNCTION referral.referrals_logs_upsert_data(user_id2 UUID, product2 VARCHAR(20), coin2 VARCHAR(20),
                                                               new_data json) RETURNS void
    LANGUAGE PLPGSQL AS
$$
DECLARE
    existing_id UUID;
BEGIN
    -- Check if a record with the given id and name exists.
    SELECT user_id
    INTO existing_id
    FROM referral.referrals_logs
    WHERE user_id = user_id2
      AND product = product2
      AND coin = coin2;

    IF FOUND THEN
        -- If a record with the given id and name exists, update it.
        UPDATE referral.referrals_logs
        SET history = history || new_data::jsonb
        WHERE user_id = user_id2
          AND product = product2
          AND coin = coin2;
    ELSE
        -- If no matching record exists, insert a new record.
        INSERT INTO referral.referrals_logs (user_id, product, coin, history)
        VALUES (user_id2, product2, coin2, jsonb_build_array(new_data::jsonb));
    END IF;
END;
$$;

-- when updating to the referral.t_referrals table, the old value is logged in the referral.t_referrals_logs table
CREATE OR REPLACE FUNCTION referral.referrals_logger() RETURNS TRIGGER
    LANGUAGE PLPGSQL AS
$$
DECLARE
    element json;
BEGIN
    -- old value is logged in table referral.t_referrals_logs
    FOR element IN SELECT row_to_json(referrals)
                   FROM referral.referrals
                   WHERE user_id = OLD.user_id
                     AND product = OLD.product
                     AND coin = OLD.coin
        LOOP
            PERFORM referral.referrals_logs_upsert_data(OLD.user_id, OLD.product, OLD.coin, element);
        END LOOP;

    RETURN NEW;
END;
$$;

CREATE TRIGGER tr_referrals_logger
    AFTER UPDATE
    ON referral.referrals
    FOR EACH ROW
EXECUTE PROCEDURE referral.referrals_logger();