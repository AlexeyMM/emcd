CREATE SCHEMA swap;

CREATE TABLE swap.users
(
    id    UUID PRIMARY KEY,
    email varchar NOT NULL
);

CREATE TABLE swap.swaps
(
    id                  UUID PRIMARY KEY,
    user_id             UUID                                      NOT NULL,
    account_from        INT                                       NOT NULL,
    coin_from           VARCHAR(10)                               NOT NULL,
    address_from        VARCHAR                                   NOT NULL,
    network_from        VARCHAR(50)                               NOT NULL,
    tag_from            VARCHAR(50) DEFAULT ''                    NOT NULL,
    coin_to             VARCHAR(10)                               NOT NULL,
    address_to          VARCHAR                                   NOT NULL,
    network_to          VARCHAR(50)                               NOT NULL,
    tag_to              VARCHAR(50)                               NOT NULL,
    amount_from         DECIMAL                                   NOT NULL,
    amount_to           DECIMAL                                   NOT NULL,
    status              SMALLINT                                  NOT NULL,
    start_time          TIMESTAMP                                 NOT NULL,
    end_time            TIMESTAMP   DEFAULT '1970-01-01 00:00:00' NOT NULL,
    ready_for_execution BOOLEAN     DEFAULT FALSE                 NOT NULL
);

CREATE TABLE swap.accounts
(
    id       INT PRIMARY KEY,
    is_valid BOOLEAN NOT NULL
);

CREATE TABLE swap.secrets
(
    account_id INT     NOT NULL REFERENCES swap.accounts (id),
    api_key    VARCHAR NOT NULL,
    api_secret VARCHAR NOT NULL
);

CREATE TABLE swap.deposits
(
    tx_id      VARCHAR PRIMARY KEY,
    swap_id    UUID        NOT NULL REFERENCES swap.swaps (id),
    coin       varchar(10) NOT NULL,
    amount     DECIMAL     NOT NULL,
    fee        DECIMAL     NOT NULL,
    status     SMALLINT    NOT NULL,
    updated_at TIMESTAMP   NOT NULL
);

CREATE TABLE swap.orders
(
    id          UUID PRIMARY KEY,
    swap_id     UUID        NOT NULL REFERENCES swap.swaps (id),
    account_id  INT         NOT NULL REFERENCES swap.accounts (id),
    category    VARCHAR(10) NOT NULL,
    symbol      VARCHAR(20) NOT NULL,
    direction   SMALLINT    NOT NULL,
    amount_from DECIMAL     NOT NULL,
    amount_to   DECIMAL     NOT NULL,
    status      SMALLINT    NOT NULL,
    is_first    BOOLEAN     NOT NULL
);

CREATE TABLE swap.internal_transfers
(
    id                UUID PRIMARY KEY,
    coin              VARCHAR(10) NOT NULL,
    amount            DECIMAL     NOT NULL,
    from_account_id   INT         NOT NULL,
    to_account_id     INT         NOT NULL,
    from_account_type VARCHAR(10) NOT NULL,
    to_account_type   VARCHAR(10) NOT NULL,
    status            SMALLINT    NOT NULL,
    updated_at        TIMESTAMP   NOT NULL
);

CREATE TABLE swap.withdraws
(
    id                    INT DEFAULT 0 NOT NULL,
    internal_id           UUID          NOT NULL PRIMARY KEY, -- нужен, потому что id (биржи) мы получаем после вставки в эту таблицу
    swap_id               UUID          NOT NULL REFERENCES swap.swaps (id),
    hash_id               VARCHAR       NOT NULL,
    coin                  VARCHAR(10)   NOT NULL,
    network               VARCHAR(50)   NOT NULL,
    address               VARCHAR       NOT NULL,
    tag                   VARCHAR(50)   NOT NULL,
    amount                DECIMAL       NOT NULL,
    include_fee_in_amount BOOLEAN       NOT NULL,
    status                SMALLINT      NOT NULL,
    explorer_link         VARCHAR       NOT NULL
);

CREATE INDEX idx_withdraws_swap_id ON swap.withdraws (swap_id);
