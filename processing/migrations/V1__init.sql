-- this table is used to keep track of merchants who have access to processing
create table if not exists merchant
(
    id uuid primary key -- emcd user_id
);

create type invoice_status as enum (
    'waiting_for_deposit',
    'payment_confirmation',
    'partially_paid',
    'payment_accepted',
    'finished',
    'expired',
    'cancelled' -- kyc, by merchant himslef or whatever
    );

create table if not exists deposit_address
(
    address     text    not null primary key,
    network_id  text    not null,
    merchant_id uuid    not null references merchant (id),
    available   boolean not null
);

create index if not exists deposit_address_user_id_network_id_ix on deposit_address (merchant_id, network_id) where available;

create table if not exists invoice
(
    id      uuid primary key,
    merchant_id     uuid           not null references merchant (id),
    coin_id         text           not null, -- use new id
    network_id      text           not null,
    deposit_address text           not null references deposit_address (address),
    amount          decimal        not null,
    buyer_fee       decimal        not null,
    merchant_fee    decimal        not null,
    title           text           not null default '',
    description     text           not null default '',
    checkout_url    text           not null default '',
    status          invoice_status not null,
    expires_at      timestamptz    not null,
    external_id     varchar(36)    not null default '',
    buyer_email     varchar(50)    not null

    -- to get already paid amount query table with transactions
);

create table if not exists merchant_tariff
(
    merchant_id uuid primary key references merchant (id),
    upper_fee   decimal not null,
    lower_fee   decimal not null,
    min_pay     decimal not null,
    max_pay     decimal not null
);

create table if not exists invoice_form
(
    id            uuid primary key,
    merchant_id   uuid not null references merchant (id),
    coin_id       text,
    network_id    text,
    amount        decimal,
    title         text,
    description   text,
    buyer_email   varchar(50),
    checkout_url  text not null
);
