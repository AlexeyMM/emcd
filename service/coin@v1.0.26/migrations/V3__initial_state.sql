DROP TABLE s_coin.t_wallet_coin_networks;
DROP TABLE s_coin.t_wallet_coins;
DROP TABLE s_coin.t_mining_coins;
DROP SCHEMA s_coin;

CREATE SCHEMA emcd;

CREATE TABLE emcd.coin
(
    id                   VARCHAR(10)  NOT NULL PRIMARY KEY,
    title                VARCHAR(255) NOT NULL DEFAULT '',
    description          VARCHAR      NOT NULL DEFAULT '',
    media_id             UUID         NOT NULL,
    decimals             NUMERIC      NOT NULL DEFAULT 0,
    withdrawal_fee       NUMERIC      NOT NULL DEFAULT 0,
    withdrawal_min_limit NUMERIC      NOT NULL DEFAULT 0,
    is_active            BOOLEAN      NOT NULL DEFAULT true,
    created_at           TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at           TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    is_mining            BOOLEAN      NOT NULL DEFAULT false,
    mining_fee           NUMERIC      NOT NULL DEFAULT 0,
    mining_merge         VARCHAR(10),
    sort_priority        INTEGER      NOT NULL
);

CREATE TABLE emcd.network
(
    id         VARCHAR(10) NOT NULL PRIMARY KEY,
    is_active  BOOLEAN     NOT NULL DEFAULT true,
    created_at TIMESTAMP   NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at TIMESTAMP   NOT NULL DEFAULT TIMEZONE('utc'::text, now())
);

CREATE TABLE emcd.coin_network
(
    coin_id     VARCHAR(10)  NOT NULL references emcd.coin,
    network_id  VARCHAR(10)  NOT NULL references emcd.network,
    title       VARCHAR(255) NOT NULL DEFAULT '',
    description VARCHAR      NOT NULL DEFAULT '',
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    created_at  TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    updated_at  TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, now()),
    PRIMARY KEY (coin_id, network_id)
);
