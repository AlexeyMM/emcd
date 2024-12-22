-- migrate:up

CREATE TABLE address_dirty
(
    address         TEXT      NOT NULL, -- address
    network         TEXT      NOT NULL, -- network
    is_dirty        BOOLEAN   NOT NULL, -- is_dirty
    updated_at      TIMESTAMP NOT NULL, -- updated date
    created_at      TIMESTAMP NOT NULL, -- created date

    CONSTRAINT address_dirty_address_network_uniq UNIQUE (address, network)

);

-- migrate:down

-- drop table if exists address_dirty cascade;
