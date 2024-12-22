-- migrate:up

CREATE TABLE address_old
(
    id              UUID      NOT NULL, -- uuid identifier
    address         TEXT      NOT NULL, -- address
    user_uuid       UUID      NOT NULL, -- user (for tables unification)
    address_type    int4      NOT NULL, -- type of address generation
    network         TEXT      NOT NULL, -- network
    user_account_id int4      NOT NULL, -- wallets user account id
    coin            TEXT      NOT NULL, -- coin
    created_at      TIMESTAMP NOT NULL, -- created date

    CONSTRAINT address_old_network_user_uuid_network_coin_uniq UNIQUE (user_uuid, network, coin),
    CONSTRAINT address_old_address_uniq UNIQUE (address),

    PRIMARY KEY (id)
);

CREATE INDEX address_old_network_user_uuid_network_coin_idx ON address_old (user_uuid, network, coin);
CREATE INDEX address_old_address_idx ON address_old (address);

CREATE TABLE address
(
    id              UUID      NOT NULL, -- uuid identifier
    address         TEXT      NOT NULL, -- address
    user_uuid       UUID      NOT NULL, -- user
    processing_uuid UUID      NOT NULL, -- processing uuid
    address_type    int4      NOT NULL, -- type of address generation
    network_group   TEXT      NOT NULL, -- network group local identifier
    created_at      TIMESTAMP NOT NULL, -- created date

    CONSTRAINT address_user_uuid_address_type_network_group_uniq UNIQUE (user_uuid, processing_uuid, address_type, network_group),
    CONSTRAINT address_address_uniq UNIQUE (address),

    PRIMARY KEY (id)
);

CREATE INDEX address_user_uuid_address_type_network_group_idx ON address USING btree (user_uuid, processing_uuid, address_type, network_group);
CREATE INDEX address_address_address_idx ON address USING btree (address);

CREATE TABLE address_derived
(
    address_uuid   uuid NOT NULL, -- address uuid
    network_group  TEXT NOT NULL, -- network group local identifier
    master_key_id  int4 NOT NULL, -- key id == 0 (reserved column)
    derived_offset int4 NOT NULL, -- offset of derived address

    CONSTRAINT address_derived_helper_key_offset_network_uniq UNIQUE (network_group, master_key_id, derived_offset),
    CONSTRAINT address_derived_uuid_fk FOREIGN KEY (address_uuid) REFERENCES address (id) ON DELETE CASCADE
);

CREATE INDEX address_derived_helper_key_offset_network_idx ON address_derived USING btree (network_group, master_key_id, derived_offset);

CREATE TABLE address_personal
(
    id         UUID             NOT NULL, -- uuid identifier
    address    TEXT             NOT NULL, -- address (possible not uniq)
    user_uuid  UUID             NOT NULL, -- user
    network    TEXT             NOT NULL, -- network group local identifier
    min_payout DOUBLE PRECISION NOT NULL, -- custom minimum payout ge coin.minpay_mining
    deleted_at TIMESTAMP        NULL,     -- deleted date
    updated_at TIMESTAMP        NOT NULL, -- updated date
    created_at TIMESTAMP        NOT NULL, -- created date

    CONSTRAINT address_personal_user_uuid_network_group_uniq UNIQUE (user_uuid, network),

    PRIMARY KEY (id)
);


-- migrate:down

-- drop table if exists address cascade;
-- drop table if exists address_old  cascade;
-- drop table if exists address_derived  cascade;
