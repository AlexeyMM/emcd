DROP TABLE IF EXISTS emcd.coin_network;
DROP TABLE IF EXISTS emcd.network;
DROP TABLE IF EXISTS emcd.coin;

DROP SCHEMA IF EXISTS emcd;

CREATE TABLE networks
(
	id          VARCHAR(10)  NOT NULL PRIMARY KEY,
	is_active   BOOLEAN      NOT NULL,
	title       VARCHAR(255) NOT NULL DEFAULT '',
	description VARCHAR      NOT NULL DEFAULT '',
	created_at  TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, NOW()),
	updated_at  TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, NOW())
);

CREATE TABLE coins
(
	id                               VARCHAR(10)  NOT NULL PRIMARY KEY,
	is_active                        BOOLEAN      NOT NULL,
	title                            VARCHAR(255) NOT NULL DEFAULT '',
	description                      VARCHAR      NOT NULL DEFAULT '',
	sort_priority                    INTEGER      NOT NULL,
	media_url                        VARCHAR      NOT NULL DEFAULT '',
	is_withdrawals_disabled          BOOLEAN      NOT NULL,
	withdrawals_disabled_description TEXT,
	created_at                       TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, NOW()),
	updated_at                       TIMESTAMP    NOT NULL DEFAULT TIMEZONE('utc'::text, NOW())
);

CREATE TABLE coins_networks
(
	coin_id                          VARCHAR(10)      NOT NULL REFERENCES coins,
	network_id                       VARCHAR(10)      NOT NULL REFERENCES networks,
	is_active                        BOOLEAN          NOT NULL,
	title                            VARCHAR(255)     NOT NULL DEFAULT '',
	description                      VARCHAR          NOT NULL DEFAULT '',
	contract_address                 VARCHAR(64),
	decimals                         INTEGER          NOT NULL,
	is_wallet                        BOOLEAN          NOT NULL,
	withdrawal_fee                   DOUBLE PRECISION NOT NULL,
	withdrawal_min_limit             DOUBLE PRECISION NOT NULL,
	withdrawal_fee_ttl_seconds       INTEGER          NOT NULL,
	is_mining                        BOOLEAN          NOT NULL,
	is_free_withdraw                 BOOLEAN          NOT NULL,
	mining_fee                       DOUBLE PRECISION,
	is_withdrawals_disabled          BOOLEAN          NOT NULL,
	withdrawals_disabled_description TEXT,
	created_at                       TIMESTAMP        NOT NULL DEFAULT TIMEZONE('utc'::text, NOW()),
	updated_at                       TIMESTAMP        NOT NULL DEFAULT TIMEZONE('utc'::text, NOW()),
	PRIMARY KEY (coin_id, network_id)
);
