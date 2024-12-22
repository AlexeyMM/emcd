-- migrate:up

create schema emcd;

create sequence users_id_seq
    as integer;

create table emcd.users
(
    email                     varchar,
    username                  varchar                                             not null,
    password                  varchar                                             not null,
    nopay                     boolean   default false,
    is_email_notifications_on boolean   default true,
    is_tg_notifications_on    boolean   default false,
    tg_id                     bigint,
    auth_secret               varchar,
    api_key                   varchar,
    language                  varchar,
    timezone                  varchar,
    is_donatation_on          boolean   default false,
    is_coinhold_enabled       boolean   default true,
    is_employee               boolean   default false,
    created_at                timestamp default now()                             not null,
    updated_at                timestamp default now()                             not null,
    id                        integer   default nextval('users_id_seq'::regclass) not null
        constraint users_pk primary key,
    ref_id                    integer   default 1                                 not null,
    parent_id                 integer,
    master_fee                numeric,
    master_id                 integer,
    phone                     varchar,
    is_phone_verified         boolean   default false,
    is_phone_2fa_enabled      boolean   default false                             not null,
    twa_code                  varchar,
    is_hedge_enabled          boolean   default false                             not null,
    company_name              varchar(255),
    is_autopay_disabled       boolean   default false,
    kyc_status                varchar,
    wb_link_id                varchar,
    pass_updated_at           timestamp,
    primary_currency          varchar   default 'usd'::character varying          not null,
    is_active                 boolean   default true                              not null,
    free_withdraw             boolean   default true                              not null,
    def_coin_id               integer   default 0                                 not null,
    new_id                    uuid
);

create table emcd.users_accounts
(
    id              serial
        constraint users_accounts_pk
            primary key,
    user_id         integer not null
        constraint users_accounts_users_id_fk
            references emcd.users
            on update cascade on delete cascade,
    coin_id         integer not null,
    account_type_id integer not null,
    minpay          numeric not null,
    address         varchar,
    changed_at      timestamp,
    img1            numeric,
    img2            numeric,
    is_active       boolean,
    created_at      timestamp default now(),
    updated_at      timestamp default now(),
    fee             numeric   default 0.015,
    user_id_new     uuid    null,
    coin_new        varchar null
);

CREATE TABLE emcd.addresses
(
    id              serial4           NOT NULL,
    user_account_id int4              NULL,
    coin_id         int4              NOT NULL,
    token_id        int4              NULL,
    address         varchar(128)      NOT NULL,
    created_at      timestamp         NOT NULL,
    deleted_at      timestamp         NULL,
    is_tracked      bool DEFAULT true NOT NULL,
    network_id      varchar(10)       NULL,
    coin_str_id     varchar(10)       NULL,
    public_key_id   int4              NULL,
    address_offset  int4              NULL,
    CONSTRAINT addresses_pk PRIMARY KEY (id)
);

CREATE TABLE emcd.autopay_addresses
(
    id              serial4                 NOT NULL,
    user_account_id int4                    NOT NULL,
    address         varchar                 NOT NULL,
    "percent"       int4      DEFAULT 0     NOT NULL,
    "label"         varchar                 NULL,
    created_at      timestamp DEFAULT now() NOT NULL,
    updated_at      timestamp DEFAULT now() NOT NULL,
    CONSTRAINT autopay_addresses_pkey PRIMARY KEY (id)
);


-- migrate:down

-- drop table if exists migration_date cascade;
