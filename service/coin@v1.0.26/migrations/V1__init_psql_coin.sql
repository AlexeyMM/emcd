CREATE SCHEMA IF NOT EXISTS s_coin;

CREATE TABLE IF NOT EXISTS s_coin.t_wallet_coins (
    coin                      VARCHAR(255) NOT NULL DEFAULT '',
    title                     VARCHAR(255) NOT NULL DEFAULT '',
    description               VARCHAR NOT NULL DEFAULT '',
    media_id                  VARCHAR NOT NULL DEFAULT '',
    decimals                  NUMERIC NOT NULL DEFAULT 0,
    withdrawal_fee            NUMERIC NOT NULL DEFAULT 0,
    withdrawal_min_limit      NUMERIC NOT NULL DEFAULT 0,
    is_active                 BOOLEAN NOT NULL DEFAULT true,
    created_at                TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at                TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now())
);

CREATE UNIQUE INDEX IF NOT EXISTS t_wallet_coins_coin_uindex ON s_coin.t_wallet_coins (coin);
COMMENT ON COLUMN s_coin.t_wallet_coins.coin IS 'coin code, which must be unique';
COMMENT ON COLUMN s_coin.t_wallet_coins.title IS 'coin name';
COMMENT ON COLUMN s_coin.t_wallet_coins.description IS 'additional field for detailed description';
COMMENT ON COLUMN s_coin.t_wallet_coins.media_id IS 'image url(logo)';
COMMENT ON COLUMN s_coin.t_wallet_coins.decimals IS 'a number of symbols after comma';
COMMENT ON COLUMN s_coin.t_wallet_coins.withdrawal_fee IS 'withdrawal fee';
COMMENT ON COLUMN s_coin.t_wallet_coins.withdrawal_fee IS 'withdrawal minimum limit';


CREATE TABLE IF NOT EXISTS s_coin.t_wallet_coin_networks (
    coin                      VARCHAR(255) NOT NULL REFERENCES s_coin.t_wallet_coins(coin) ON UPDATE CASCADE ON DELETE RESTRICT,
    code                      VARCHAR(255) NOT NULL DEFAULT '',
    title                     VARCHAR(255) NOT NULL DEFAULT '',
    description               VARCHAR NOT NULL DEFAULT '',
    is_active                 BOOLEAN NOT NULL DEFAULT true,
    created_at                TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at                TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now())
);

CREATE INDEX IF NOT EXISTS t_wallet_coin_networks_coin_index ON s_coin.t_wallet_coin_networks (coin);
CREATE UNIQUE INDEX IF NOT EXISTS t_wallet_coin_networks_code_uindex ON s_coin.t_wallet_coin_networks (code);
COMMENT ON COLUMN s_coin.t_wallet_coin_networks.coin IS 'coin code';
COMMENT ON COLUMN s_coin.t_wallet_coin_networks.code IS 'network code, which must be unique';
COMMENT ON COLUMN s_coin.t_wallet_coin_networks.title IS 'network name';
COMMENT ON COLUMN s_coin.t_wallet_coin_networks.description IS 'additional field for detailed description';

CREATE TABLE IF NOT EXISTS s_coin.t_mining_coins (
    coin                      VARCHAR(255) NOT NULL DEFAULT '',
    title                     VARCHAR(255) NOT NULL DEFAULT '',
    description               VARCHAR NOT NULL DEFAULT '',
    fee                       NUMERIC NOT NULL DEFAULT 0,
    marge_coin                VARCHAR(255) NOT NULL DEFAULT '',
    is_active                 BOOLEAN NOT NULL DEFAULT true,
    created_at                TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at                TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc'::text, now())
);


CREATE UNIQUE INDEX IF NOT EXISTS t_mining_coins_coin_uindex ON s_coin.t_mining_coins (coin);
COMMENT ON COLUMN s_coin.t_mining_coins.coin IS 'coin code, which must be unique';
COMMENT ON COLUMN s_coin.t_mining_coins.title IS 'coin name';
COMMENT ON COLUMN s_coin.t_mining_coins.description IS 'additional field for detailed description';
COMMENT ON COLUMN s_coin.t_mining_coins.fee IS 'default mining fee (different users have different fees)';
COMMENT ON COLUMN s_coin.t_mining_coins.marge_coin IS 'marge mining coin';