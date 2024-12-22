create schema emcd

    create table emcd.account_types
    (
        id         serial
            constraint account_types_pk
                primary key,
        name       varchar not null,
        comment    text,
        created_at timestamp default now()
    );


create unique index account_types_id_uindex
    on emcd.account_types (id);

create unique index account_types_name_uindex
    on emcd.account_types (name);

create table emcd.coinhold_default_settings
(
    id         serial
        constraint coinhold_default_settings_pk
            primary key,
    asset      varchar,
    name       varchar,
    value      varchar,
    created_at date default CURRENT_DATE,
    updated_at date default CURRENT_DATE
);



create unique index coinhold_default_settings_id_uindex
    on emcd.coinhold_default_settings (id);


create table emcd.coins
(
    id                               serial
        constraint coins_pk
            primary key,
    name                             varchar                                        not null,
    description                      varchar,
    code                             varchar                                        not null,
    rate                             numeric,
    created_at                       timestamp        default now(),
    updated_at                       timestamp        default now(),
    is_mining                        boolean          default false                 not null,
    is_wallet                        boolean          default false                 not null,
    is_deposit                       boolean          default false                 not null,
    is_hedge                         boolean          default false                 not null,
    is_p2p                           boolean          default false                 not null,
    is_withdrawals_disabled          boolean          default false,
    withdrawals_disabled_description text,
    is_free_withdraw                 boolean          default false                 not null,
    withdraw_min_limit               double precision default 0                     not null,
    withdraw_fee                     double precision default 0                     not null,
    title                            varchar(32)      default ''::character varying not null,
    picture_url                      varchar(256)     default ''::character varying not null,
    tokens_network_title             varchar(16)      default ''::character varying not null
);


create table emcd.coin_rates
(
    id         serial
        constraint coin_rates_pk
            primary key,
    coin_id    integer not null
        constraint coin_rates_to_coins_fk
            references emcd.coins
            on update cascade on delete restrict,
    rate       double precision,
    date       varchar,
    updated_at timestamp default now()
);



create unique index coins_code_uindex
    on emcd.coins (code);

create unique index coins_id_uindex
    on emcd.coins (id);

create unique index coins_name_uindex
    on emcd.coins (name);

create index coins_is_mining_is_wallet_is_deposit_is_hedge_index
    on emcd.coins (is_mining, is_wallet, is_deposit, is_hedge);



create table emcd.daily_pool_profits
(
    id           serial
        constraint daily_pool_profits_pk
            primary key,
    coin_id      integer           not null,
    amount       numeric,
    status       integer           not null,
    profit       numeric,
    hashrate     numeric,
    created_at   date    default CURRENT_DATE,
    updated_at   date    default CURRENT_DATE,
    hashrate_ext numeric default 0 not null
);



create unique index daily_pool_profits_id_uindex
    on emcd.daily_pool_profits (id);



create table emcd.notification_groups
(
    id          serial
        constraint notification_groups_pk
            primary key,
    name        varchar                                 not null,
    description varchar   default ''::character varying not null,
    editable    boolean   default false                 not null,
    created_at  timestamp default CURRENT_TIMESTAMP     not null,
    updated_at  timestamp default CURRENT_TIMESTAMP,
    segment_key varchar   default ''::character varying not null
);



create table emcd.default_notification_settings
(
    id            serial
        constraint default_notification_settings_pk
            primary key,
    send_email    boolean   default false,
    send_telegram boolean   default false,
    group_id      integer                             not null
        references emcd.notification_groups
            on delete cascade,
    send_push     boolean   default false,
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    updated_at    timestamp default CURRENT_TIMESTAMP
);



create table emcd.notification_types
(
    id          serial
        constraint notification_types_pk
            primary key,
    name        varchar                                 not null,
    segment_key varchar                                 not null,
    group_id    integer                                 not null
        references emcd.notification_groups
            on delete cascade,
    description varchar   default ''::character varying not null,
    created_at  timestamp default CURRENT_TIMESTAMP     not null,
    updated_at  timestamp default CURRENT_TIMESTAMP
);



create table emcd.promocodes
(
    id                          serial
        constraint promocodes_pk
            primary key,
    code                        varchar          not null,
    valid_days_amount           integer   default 0,
    has_no_limit                boolean   default false,
    bonus_fee                   double precision not null,
    referral_enabled            boolean   default false,
    is_summable                 boolean   default false,
    is_active                   boolean   default true,
    created_at                  timestamp default now(),
    coin_id                     varchar          not null,
    ref_id                      integer   default 1,
    bonus_fee_days_amount       integer   default 0,
    bonus_fee_has_no_limit      boolean   default false,
    partner                     varchar,
    is_only_for_registration    boolean   default false,
    is_only_for_private_cabinet boolean   default false,
    is_disposable               boolean   default false
);



create unique index promocodes_code_uindex
    on emcd.promocodes (code);

create index promocodes_id_index
    on emcd.promocodes (id);

create table emcd.referral_rewards
(
    id           serial
        constraint referral_rewards_pk
            primary key,
    tier         integer not null,
    min_hashrate numeric,
    max_hashrate numeric,
    unit         bigint,
    referral_fee numeric,
    commission   numeric,
    coin_id      integer not null
        constraint referral_rewards_coins_id_fk
            references emcd.coins
            on update cascade on delete cascade,
    created_at   timestamp default now(),
    updated_at   timestamp default now()
);



create index referral_rewards_coin_id_index
    on emcd.referral_rewards (coin_id);

create unique index referral_rewards_id_uindex
    on emcd.referral_rewards (id);



create table emcd.registration_sources
(
    id                             serial,
    user_id                        integer,
    token                          varchar,
    used                           boolean   default false,
    ip                             varchar,
    yid                            varchar,
    gid                            varchar,
    utm                            varchar,
    created_at                     timestamp default now() not null,
    updated_at                     timestamp default now(),
    phone                          varchar,
    code                           varchar,
    is_send_notifications_accepted boolean
);



create unique index registration_sources_id_uindex
    on emcd.registration_sources (id);



create table emcd.transaction_types
(
    id              serial
        constraint transaction_types_pk
            primary key,
    description     text                    not null,
    created_at      timestamp default now() not null,
    account_type_id integer                 not null
        constraint transaction_types___account_type_fk
            references emcd.account_types
);



create unique index transaction_types_id_uindex
    on emcd.transaction_types (id);

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
        constraint users_pk
            primary key,
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


create table emcd.subscribers
(
    id         serial
        constraint subscribers_pk
            primary key,
    user_id    integer
        constraint subscribers_users_id_fk
            references emcd.users
            on update cascade on delete cascade,
    email      varchar not null,
    mining     boolean   default false,
    mined      boolean   default false,
    lang       varchar,
    created_at timestamp default now(),
    updated_at timestamp default now()
);



create table emcd.user_logs
(
    id                serial
        constraint user_logs_pk
            primary key,
    user_id           integer                 not null
        constraint user_logs_users_id_fk
            references emcd.users
            on update cascade on delete cascade,
    change_type       varchar,
    is_segment_sended boolean   default false,
    used              boolean   default false,
    ip                varchar,
    token             varchar,
    old_value         varchar,
    value             varchar,
    created_at        timestamp default now(),
    updated_at        timestamp default now(),
    active            boolean   default false not null
);



create index user_logs_is_segment_sended_index
    on emcd.user_logs (is_segment_sended);

create index user_logs_user_id_index
    on emcd.user_logs (user_id);



create table emcd.user_notification_types
(
    id            serial
        constraint user_notification_types_pk
            primary key,
    user_id       integer not null
        constraint user_notification_types_users_id_fk
            references emcd.users
            on delete cascade,
    type_id       integer not null
        constraint user_notification_types_notification_types_id_fk
            references emcd.notification_types
            on delete cascade,
    send_email    boolean default false,
    send_telegram boolean default false
);



create index hedge_orders_user_id_index
    on emcd.users (id);

create unique index users_email_uindex
    on emcd.users (email);

create unique index users_username_uindex
    on emcd.users (username);

create index users_pass_updated_at_index
    on emcd.users (pass_updated_at);

create index users_ref_id_index
    on emcd.users (ref_id);



create table emcd.users_accounts
(
    id              serial
        constraint users_accounts_pk
            primary key,
    user_id         integer not null
        constraint users_accounts_users_id_fk
            references emcd.users
            on update cascade on delete cascade,
    coin_id         integer not null
        constraint users_accounts_coins_id_fk
            references emcd.coins
            on update cascade on delete cascade,
    account_type_id integer not null
        constraint users_accounts_account_types_id_fk
            references emcd.account_types
            on update cascade on delete cascade,
    minpay          numeric not null,
    address         varchar,
    changed_at      timestamp,
    img1            numeric,
    img2            numeric,
    is_active       boolean,
    created_at      timestamp default now(),
    updated_at      timestamp default now(),
    fee             numeric   default 0.015,
    user_id_new     uuid null,
    coin_new        varchar null
);



create table emcd.accounts_balances
(
    account_id integer                 not null
        constraint accounts_balances_users_accounts_id_fk
            references emcd.users_accounts
            on delete cascade,
    balance    numeric   default 0     not null,
    updated_at timestamp default now() not null
);



create unique index accounts_balances_account_id_uindex
    on emcd.accounts_balances (account_id);

create index accounts_balances_balance_index
    on emcd.accounts_balances (balance);



create table emcd.accounts_referral
(
    id               serial
        constraint accounts_referral_pk
            primary key,
    account_id       integer not null
        constraint accounts_referral_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    tier             integer not null,
    referral_fee     numeric not null,
    active_referrals integer,
    coin_id          integer not null
        constraint accounts_referral_coins_id_fk
            references emcd.coins
            on update cascade on delete cascade,
    created_at       timestamp default now(),
    updated_at       timestamp default now()
);



create index accounts_referral_account_id_index
    on emcd.accounts_referral (account_id);

create unique index accounts_referral_id_uindex
    on emcd.accounts_referral (id);



create table emcd.coinhold_events
(
    id         serial
        constraint coinhold_events_pk
            primary key,
    user_id    integer not null
        constraint coinhold_events_users_id_fk
            references emcd.users
            on update cascade on delete cascade,
    account_id integer not null
        constraint coinhold_events_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    name       varchar not null,
    params     varchar,
    created_at timestamp default now(),
    updated_at timestamp default now()
);


create index coinhold_events_account_id_index
    on emcd.coinhold_events (account_id);

create unique index coinhold_events_id_uindex
    on emcd.coinhold_events (id);

create index coinhold_events_user_id_index
    on emcd.coinhold_events (user_id);



create table emcd.hedge_orders
(
    id                      serial
        constraint hedge_orders_pk
            primary key,
    user_id                 integer                    not null
        constraint hedge_orders_to_users_fk
            references emcd.users
            on update cascade on delete restrict,
    status                  varchar,
    json_request_data       varchar,
    json_response_data      varchar,
    cost                    numeric          default 0 not null,
    account_id              integer
        constraint hedge_orders_to_users_accounts_fk
            references emcd.users_accounts
            on update cascade on delete restrict,
    created_at              timestamp        default now(),
    block_transaction_id    integer          default 0 not null,
    unblock_transaction_id  integer,
    pay_transaction_id      integer,
    fee_transaction_id      integer,
    coin_id                 integer          default 0 not null,
    deal_id                 varchar,
    pnl_percentage          double precision default 0 not null,
    payout_transaction_id   integer,
    structure_product_price numeric          default 0 not null
);



create index hedge_orders_transactions_order_id_index
    on emcd.hedge_orders (id);

create table emcd.referral_statistics_log
(
    id                serial
        constraint referral_statistics_log_pk
            primary key,
    user_id           integer not null
        constraint referral_statistics_log_users_id_fk
            references emcd.users
            on update cascade on delete cascade,
    account_id        integer not null
        constraint referral_statistics_log_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    referral_hashrate numeric,
    tier              integer not null,
    reward            numeric,
    active_referrals  integer,
    month_earnings    numeric,
    report_date       date,
    coin_id           integer not null
        constraint referral_statistics_log_coins_id_fk
            references emcd.coins
            on update cascade on delete cascade,
    created_at        timestamp default now(),
    updated_at        timestamp default now()
);



create unique index referral_statistics_log_id_uindex
    on emcd.referral_statistics_log (id);



create index hedge_orders_account_id_index
    on emcd.users_accounts (id);

create unique index users_accounts_id_uindex
    on emcd.users_accounts (id);

create index users_accounts_user_id_coin_id_account_type_id_index
    on emcd.users_accounts (user_id, coin_id, account_type_id);

create index users_accounts_user_id_new_coin_new_account_type_id_index
    on emcd.users_accounts (user_id_new, coin_new, account_type_id);

create index whitebird_orders_account_id_index
    on emcd.users_accounts (id);



create table emcd.users_promocodes
(
    id           serial
        constraint table_name_pk
            primary key,
    promocode_id integer not null
        constraint table_name_promocodes_id_fk
            references emcd.promocodes
            on delete cascade,
    user_id      integer not null
        constraint table_name_users_id_fk
            references emcd.users
            on delete cascade,
    created_at   timestamp default now(),
    expires_at   timestamp
);



create index table_name_user_id_index
    on emcd.users_promocodes (user_id);

create index users_promocodes_promocode_id_index
    on emcd.users_promocodes (promocode_id);

create table emcd.whitebird_orders
(
    id                   serial
        constraint whitebird_orders_pk
            primary key,
    request_id           varchar,
    status               varchar,
    from_currency        varchar             not null,
    to_currency          varchar             not null,
    amount               double precision    not null,
    user_id              integer             not null,
    created_at           timestamp default now(),
    updated_at           timestamp default now(),
    amount_received      numeric,
    block_transaction_id integer,
    ext_service_id       integer   default 1 not null,
    currency_rate        numeric,
    order_code           varchar,
    is_visible           boolean,
    two_fa_type          varchar,
    card_id              uuid,
    account_id           integer
        constraint whitebird_orders_to_users_accounts_fk
            references emcd.users_accounts
            on update cascade on delete restrict,
    garantex_order_id    integer,
    comment              varchar,
    hashed_pan           varchar,
    market_status        varchar
);



create table emcd.whitebird_logs
(
    id         serial
        constraint whitebird_logs_pk
            primary key,
    type       varchar,
    old_value  varchar,
    new_value  varchar,
    user_id    integer not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    order_id   integer
        references emcd.whitebird_orders
            on update cascade on delete restrict
);



create index whitebird_logs_order_id_idx
    on emcd.whitebird_logs (order_id);

create index hashed_pan_whitebird_orders_idx
    on emcd.whitebird_orders using hash (hashed_pan)
    where (hashed_pan IS NOT NULL);

create table emcd.coinhold_interest_groups
(
    id   serial
        constraint coinhold_interest_groups_pk
            primary key,
    name varchar(64) not null
);


create table emcd.coinhold_types
(
    id                 serial
        constraint coinhold_types_pk
            primary key,
    created_at         timestamp                     not null,
    updated_at         timestamp                     not null,
    deleted_at         timestamp,
    is_active          boolean                       not null,
    name               varchar(255)                  not null,
    coin_id            integer
        constraint coinhold_types_coins_id_fk
            references emcd.coins,
    initial_amount     numeric                       not null,
    period_days        smallint                      not null,
    period_months      smallint                      not null,
    is_auto_renewal    boolean                       not null,
    is_capitalized     boolean                       not null,
    is_fixed           boolean                       not null,
    min_mining_percent double precision              not null,
    interest_type_id   smallint                      not null,
    limit_amount       double precision default 0    not null,
    interest_group_id  integer                       not null
        constraint coinhold_types_coinhold_interest_groups_id_fk
            references emcd.coinhold_interest_groups,
    is_reducible       boolean          default true not null
);


create table emcd.accounts_coinhold
(
    id                     serial
        constraint accounts_coinhold_pk
            primary key,
    account_id             integer                 not null
        constraint accounts_coinhold_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    deposit_mining_percent numeric,
    created_at             timestamp default now(),
    updated_at             timestamp default now(),
    is_capitalized         boolean   default false not null,
    end_time               timestamp,
    coinhold_type_id       integer
        constraint coinhold_types_id_fk
            references emcd.coinhold_types,
    deleted_at             timestamp,
    is_auto_renewal        boolean   default false not null,
    promocode              varchar(32)
);



create index accounts_coinhold_account_id_index
    on emcd.accounts_coinhold (account_id);

create unique index accounts_coinhold_id_uindex
    on emcd.accounts_coinhold (id);

create index accounts_coinhold_promocode_index
    on emcd.accounts_coinhold (promocode);



create table emcd.coinhold_interest_tiers
(
    id                serial
        constraint coinhold_interest_tiers_pk
            primary key,
    interest_group_id integer          not null
        constraint coinhold_interest_tiers_coinhold_interest_groups_id_fk
            references emcd.coinhold_interest_groups,
    level             smallint         not null,
    min_balance       double precision not null,
    max_balance       double precision,
    percent           double precision not null
);



create index coinhold_interest_tiers_interest_group_id_index
    on emcd.coinhold_interest_tiers (interest_group_id);

create table emcd.accounts_pool
(
    id                   serial
        constraint accounts_pool_pk
            primary key,
    emcd_address_autopay boolean   default false,
    account_id           integer not null
        constraint accounts_pool_to_users_accounts_fk
            references emcd.users_accounts
            on update cascade on delete restrict,
    created_at           timestamp default now(),
    updated_at           timestamp default now(),
    autopay_percent      integer   default 100
        constraint autopay_pecent_val
            check ((autopay_percent <= 100) AND (autopay_percent >= 0))
);


create table emcd.service_users
(
    user_id     integer     not null
        constraint service_users_users_id_fk
            references emcd.users
            on delete cascade,
    username    varchar(64) not null,
    email       varchar(64) not null,
    description varchar
);



create unique index service_users_email_uindex
    on emcd.service_users (email);

create unique index service_users_user_id_uindex
    on emcd.service_users (user_id);

create unique index service_users_username_uindex
    on emcd.service_users (username);

create table emcd.tokens
(
    id                               serial
        constraint tokens_pk
            primary key,
    coin_id                          integer
        constraint tokens_coins_id_fk
            references emcd.coins
        constraint tokens_to_coin_fk
            references emcd.coins,
    code                             varchar      not null,
    description                      varchar      not null,
    contract_address                 varchar(128) not null,
    network_fee                      numeric      not null,
    main_coin_id                     integer      not null
        constraint tokens_coins_id_fk_2
            references emcd.coins,
    decimals                         smallint     not null,
    is_withdrawals_disabled          boolean default false,
    withdrawals_disabled_description text
);



create table emcd.transactions
(
    id                  serial
        constraint transactions_pk
            primary key,
    type                integer                 not null,
    sender_account_id   integer                 not null
        constraint transactions_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    receiver_account_id integer
        constraint transactions_users_id_fk_2
            references emcd.users_accounts
            on update cascade on delete cascade,
    coin_id             integer                 not null
        constraint transactions_coins_id_fk
            references emcd.coins
            on update cascade on delete cascade,
    amount              numeric,
    comment             text,
    hashrate            bigint,
    created_at          timestamp default now() not null,
    hash                varchar,
    fee                 numeric,
    gas_price           numeric,
    from_referral_id    integer
        constraint transactions_to_ref_user_fk
            references emcd.users
            on update cascade on delete restrict,
    is_viewed           boolean   default false,
    receiver_address    varchar(255),
    token_id            integer
        constraint transactions_token_id_fk
            references emcd.tokens
);



create table emcd.early_payment_events
(
    id             serial
        constraint early_payment_events_pk
            primary key,
    coin_id        integer not null
        constraint early_payment_events_coins_id_fk
            references emcd.coins
            on update cascade on delete cascade,
    account_id     integer not null
        constraint early_payment_events_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    event_name     varchar not null,
    event_data     varchar,
    status         integer not null,
    transaction_id integer
        constraint early_payment_events_transactions_id_fk
            references emcd.transactions
            on update cascade on delete cascade,
    amount         numeric not null,
    created_at     timestamp default now(),
    updated_at     timestamp default now()
);



create unique index early_payment_events_id_uindex
    on emcd.early_payment_events (id);



create table emcd.operations
(
    id             serial
        constraint operations_pk
            primary key,
    type           integer not null,
    transaction_id integer not null
        constraint operations_transactions_id_fk
            references emcd.transactions
            on update cascade on delete cascade,
    account_id     integer not null
        constraint operations_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    coin_id        integer not null
        constraint operations_coins_id_fk
            references emcd.coins
            on update cascade on delete cascade,
    amount         numeric,
    created_at     timestamp default now(),
    hash           varchar,
    token_id       integer
        constraint operations_token_id_fk
            references emcd.tokens
);



create index operations_account_id_index
    on emcd.operations (account_id);

create index operations_created_at_index
    on emcd.operations (created_at);

create unique index operations_id_uindex
    on emcd.operations (id);

create index operations_type_index
    on emcd.operations (type);

create index operations_token_id_index
    on emcd.operations (token_id desc);



create index hedge_orders_transactions_transaction_id_index
    on emcd.transactions (id);

create index transactions_coin_id_index
    on emcd.transactions (coin_id);

create index transactions_created_at_index
    on emcd.transactions (created_at);

create unique index transactions_id_uindex
    on emcd.transactions (id);

create index transactions_receiver_account_id__index
    on emcd.transactions (receiver_account_id);

create index transactions_sender_account_id_index
    on emcd.transactions (sender_account_id);

create index transactions_type_index
    on emcd.transactions (type);

create index transactions_token_id_index
    on emcd.transactions (token_id desc);

create table emcd.transactions_blocks
(
    id                     serial
        constraint transactions_blocks_pk
            primary key,
    block_transaction_id   integer   not null
        constraint transactions_blocks_block_transaction_id_fk
            references emcd.transactions,
    unblock_transaction_id integer
        constraint transactions_blocks_unblock_transaction_id_fk
            references emcd.transactions,
    unblock_to_account_id  integer   not null
        constraint users_accounts_unblock_to_account_id_fk
            references emcd.users_accounts,
    blocked_till           timestamp not null
);



create index transactions_blocks_block_till_index
    on emcd.transactions_blocks (blocked_till);

create index transactions_blocks_unblock_transaction_id_index
    on emcd.transactions_blocks (unblock_transaction_id);

create index tokens_coin_id_index
    on emcd.tokens (coin_id);

create table emcd.addresses
(
    id              serial
        constraint addresses_pk
            primary key,
    user_account_id integer
        constraint tokens_coins_id_fk
            references emcd.users_accounts,
    coin_id         integer      not null
        constraint addresses_to_coin_fk
            references emcd.coins,
    token_id        integer
        constraint addresses_to_token_fk
            references emcd.tokens,
    address         varchar(128) not null,
    created_at      timestamp    not null,
    deleted_at      timestamp
);



create table emcd.wallet_transactions
(
    id                  serial
        constraint wallet_transactions_pk
            primary key,
    coin_id             integer      not null,
    tx_id               varchar(255) not null,
    amount              numeric,
    confirmations       integer   default 0,
    is_confirmed        boolean   default false,
    created_at          timestamp default now(),
    receiver_account_id integer
        constraint wallet_transactions_users_accounts_id_fk
            references emcd.users_accounts,
    updated_at          timestamp,
    confirmed_at        timestamp,
    deleted_at          timestamp,
    receiver_address_id integer
        constraint wallet_transactions_receiver_address_id_fk
            references emcd.addresses,
    token_id            integer
        constraint wallet_transactions_token_id_fk
            references emcd.tokens
);



create unique index wallet_transactions_coin_id_tx_id_receiver_account_id_uindex
    on emcd.wallet_transactions (coin_id, tx_id, receiver_account_id);

create index wallet_transactions_receiver_account_id_index
    on emcd.wallet_transactions (receiver_account_id);

create index wallet_transactions_tx_id_index
    on emcd.wallet_transactions (tx_id);

create index wallet_transactions_receiver_address_id_index
    on emcd.wallet_transactions (receiver_address_id desc);

create index wallet_transactions_token_id_index
    on emcd.wallet_transactions (token_id desc);

create index addresses_coin_id_index
    on emcd.addresses (coin_id);

create index addresses_token_id_index
    on emcd.addresses (token_id);

create table emcd.subscription_form
(
    id         serial
        primary key,
    user_id    integer not null
        constraint service_users_users_id_fk
            references emcd.users
            on delete cascade,
    created_at timestamp default now()
);



create unique index subscription_form_user_id_uindex
    on emcd.subscription_form (user_id);

create table emcd.fiat_commissions
(
    id                 serial
        constraint fiat_commissions_pk
            primary key,
    emcd_commission    numeric,
    partner_commission numeric,
    service_id         integer not null,
    comment            varchar,
    created_at         timestamp default now(),
    updated_at         timestamp default now()
);



create table emcd.cryptocurrency_price
(
    id                  serial
        primary key,
    code                varchar                             not null
        unique,
    price_usd           numeric   default 0,
    price_btc           numeric   default 0,
    changed_percent_usd numeric   default 0,
    updated_at          timestamp default CURRENT_TIMESTAMP not null,
    price_rub           numeric   default 0,
    price_cny           numeric   default 0,
    price_irr           numeric   default 0,
    price_eur           numeric   default 0,
    price_kzt           numeric   default 0
);



create table emcd.user_notification_settings
(
    id            serial
        primary key,
    user_id       integer                             not null
        references emcd.users
            on delete cascade,
    group_id      integer                             not null
        references emcd.notification_groups
            on delete cascade,
    send_email    boolean   default false,
    send_telegram boolean   default false,
    send_push     boolean   default false,
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    updated_at    timestamp default CURRENT_TIMESTAMP
);



create table emcd.daily_pool_discounts
(
    id         serial
        constraint daily_pool_discounts_pk
            primary key,
    coin_id    integer not null,
    percent    numeric,
    amount     numeric,
    created_at date default CURRENT_DATE
);



create table emcd.withdraw_transactions
(
    tx_id     integer not null
        primary key,
    fee_tx_id integer not null
);



create index withdraw_transactions_fee_tx_id_index
    on emcd.withdraw_transactions (fee_tx_id);

create table emcd.worker_groups
(
    id         serial
        constraint worker_groups_pk
            primary key,
    group_name varchar(64) not null,
    user_id    integer     not null
        constraint worker_groups_users_id_fk
            references emcd.users
            on delete cascade,
    created_at timestamp default now(),
    updated_at timestamp default now()
);


create unique index worker_groups_id_uindex
    on emcd.worker_groups (id);

create unique index worker_groups_user_id_group_name_unique
    on emcd.worker_groups (user_id, group_name);

create table emcd.worker_groups_items
(
    id              serial
        constraint worker_groups_items_pk
            primary key,
    worker_name     varchar(64) not null,
    worker_group_id integer     not null
        constraint worker_groups_id_fk
            references emcd.worker_groups
            on delete cascade,
    created_at      timestamp default now(),
    updated_at      timestamp default now()
);



create unique index worker_groups_items_id_uindex
    on emcd.worker_groups_items (id);

create unique index worker_groups_id_worker_name_uindex
    on emcd.worker_groups_items (worker_group_id, worker_name);

create table emcd.autopay_addresses
(
    id              serial
        primary key,
    user_account_id integer                 not null
        constraint fk_user_account_id
            references emcd.users_accounts,
    address         varchar                 not null,
    percent         integer   default 0     not null,
    label           varchar,
    created_at      timestamp default now() not null,
    updated_at      timestamp default now() not null
);



create index autopay_addresses_user_account_id_index
    on emcd.autopay_addresses (user_account_id);

create table emcd.vip_users
(
    id                  serial
        constraint vip_users_id_primary_key
            primary key,
    user_id             integer                 not null
        constraint foreign_key_vip_users_userid
            references emcd.users,
    title               varchar(64),
    description         varchar(1024),
    is_divided_address  boolean                 not null,
    created_at          timestamp default now() not null,
    updated_at          timestamp default now() not null,
    is_blocked_ref_tier boolean   default false not null,
    spec_ref_fee        numeric
);

comment on column emcd.vip_users.is_divided_address is 'divided mining payout address options';


create unique index vip_users_user_id_index
    on emcd.vip_users (user_id);

create table emcd.fiat_promocodes
(
    id             serial
        primary key,
    promocode      varchar(32)             not null,
    percent        double precision        not null,
    is_disposable  boolean   default false not null,
    per_user_limit integer,
    created_at     timestamp default now() not null,
    deleted_at     timestamp,
    expire_at      timestamp
);



create table emcd.card_orders
(
    id                     serial
        primary key,
    user_id                integer                                 not null
        references emcd.users
            on update cascade on delete restrict,
    account_id             integer                                 not null
        references emcd.users_accounts
            on update cascade on delete restrict,
    pair                   varchar   default ''::character varying not null,
    status                 varchar   default ''::character varying,
    amount                 numeric                                 not null,
    amount_fiat            numeric                                 not null,
    order_code             varchar   default ''::character varying,
    request_id             varchar   default ''::character varying,
    block_transaction_id   integer,
    ext_service_id         integer   default 1                     not null,
    rate                   numeric                                 not null,
    is_visible             boolean,
    two_fa_type            varchar,
    card_id                uuid,
    hashed_pan             varchar,
    market_order_id        integer,
    market_status          varchar   default ''::character varying,
    market_amount_received numeric   default 0                     not null,
    comment                varchar   default ''::character varying,
    created_at             timestamp default now(),
    updated_at             timestamp default now(),
    promocode_id           integer
        constraint card_orders_fiat_promocodes_id_fk
            references emcd.fiat_promocodes
            on update restrict on delete restrict
);



create index card_orders_promocode_id_index
    on emcd.card_orders (promocode_id);

create table emcd.card_order_logs
(
    id         serial
        primary key,
    order_id   integer                                 not null
        references emcd.card_orders
            on update cascade on delete restrict,
    type       varchar   default ''::character varying not null,
    old_status varchar   default ''::character varying not null,
    new_status varchar   default ''::character varying not null,
    comment    varchar   default ''::character varying,
    created_at timestamp default now()
);



create index card_order_logs_order_id_index
    on emcd.card_order_logs (order_id);

create table emcd.fiat_currencies
(
    id         serial
        primary key,
    code       varchar(255) default ''::character varying,
    name       varchar(255) default ''::character varying,
    sign       varchar(255) default ''::character varying,
    is_p2p     boolean      default true              not null,
    created_at timestamp    default CURRENT_TIMESTAMP not null
);


create table emcd.p2p_offer_statuses
(
    id          serial
        primary key,
    code        varchar(255) default ''::character varying,
    description varchar      default ''::character varying not null,
    created_at  timestamp    default CURRENT_TIMESTAMP     not null
);



create table emcd.p2p_statuses
(
    id          serial
        primary key,
    code        varchar(255) default ''::character varying,
    description varchar      default ''::character varying not null,
    ttl         integer      default 0                     not null,
    created_at  timestamp    default CURRENT_TIMESTAMP     not null
);



create table emcd.p2p_users
(
    id               serial
        primary key,
    user_id          integer                not null
        references emcd.users
            on update cascade on delete restrict,
    rating           numeric   default 0    not null,
    success_trades   integer   default 0    not null,
    unsuccess_trades integer   default 0    not null,
    is_active        boolean   default true not null,
    created_at       timestamp default now(),
    updated_at       timestamp default now()
);



create unique index p2p_users_user_id_uindex
    on emcd.p2p_users (user_id);

create table emcd.p2p_commissions
(
    id         serial
        primary key,
    commission numeric   default 0 not null,
    direction  varchar   default ''::character varying,
    comment    varchar   default ''::character varying,
    created_at timestamp default now(),
    updated_at timestamp default now()
);



create table emcd.p2p_status_reasons
(
    id          serial
        primary key,
    code        varchar(255) default ''::character varying,
    description varchar      default ''::character varying not null,
    created_at  timestamp    default CURRENT_TIMESTAMP     not null
);



create table emcd.exchange_promocodes
(
    id             serial
        primary key,
    promocode      varchar(32)             not null,
    percent        double precision        not null,
    is_disposable  boolean   default false not null,
    per_user_limit integer,
    created_at     timestamp default now() not null,
    deleted_at     timestamp,
    expire_at      timestamp
);



create table emcd.exchanges
(
    id                     serial
        constraint exchanges_pk
            primary key,
    pair                   varchar(32) not null,
    send_transaction_id    integer     not null
        constraint exchanges_send_transactions_id_fk
            references emcd.transactions,
    rate                   numeric     not null,
    receive_transaction_id integer
        constraint exchanges_receive_transactions_id_fk
            references emcd.transactions,
    exchange_order_id      varchar(128),
    exchange_amount        numeric,
    is_success             boolean,
    comment                varchar(256),
    promocode_id           integer
        constraint exchanges_promocodes_id_fk
            references emcd.exchange_promocodes
            on update restrict on delete restrict
);



create index exchanges_receive_transaction_id_index
    on emcd.exchanges (receive_transaction_id);

create index exchanges_send_transaction_id_index
    on emcd.exchanges (send_transaction_id);

create index exchanges_exchange_order_id_index
    on emcd.exchanges (exchange_order_id);

create index exchanges_promocode_id_index
    on emcd.exchanges (promocode_id);

create table emcd.users_mining_forecast
(
    account_id integer                 not null
        constraint account_id_key
            unique
        constraint users_mining_forecast_users_accounts_id_fk
            references emcd.users_accounts
            on update cascade on delete cascade,
    img1       numeric   default 0     not null,
    img2       numeric   default 0     not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);



create table emcd.rates
(
    id                  serial
        primary key,
    crypto_currency_id  integer                             not null
        references emcd.coins
            on update cascade on delete restrict,
    fiat_currency_id    integer                             not null
        references emcd.fiat_currencies
            on update cascade on delete restrict,
    crypto_to_fiat_rate numeric   default 0                 not null,
    updated_at          timestamp default CURRENT_TIMESTAMP not null
);



create table emcd.p2p_pay_method_types
(
    id          serial
        primary key,
    code        varchar(255) default ''::character varying,
    description varchar      default ''::character varying not null,
    created_at  timestamp    default CURRENT_TIMESTAMP     not null
);



create table emcd.p2p_pay_methods
(
    id               serial
        primary key,
    name             varchar(255) default ''::character varying,
    description      varchar      default ''::character varying not null,
    is_active        boolean      default true                  not null,
    created_at       timestamp    default CURRENT_TIMESTAMP     not null,
    fiat_currency_id integer      default 1                     not null
        references emcd.fiat_currencies
            on update cascade on delete restrict,
    key              varchar      default ''::character varying not null,
    type             integer      default 1                     not null
        references emcd.p2p_pay_method_types
            on update cascade on delete restrict
);



create table emcd.p2p_offers
(
    id              serial
        primary key,
    owner_id        integer                                 not null
        references emcd.p2p_users
            on update cascade on delete restrict,
    status_id       integer                                 not null
        references emcd.p2p_offer_statuses
            on update cascade on delete restrict,
    crypto_currency integer                                 not null
        references emcd.coins
            on update cascade on delete restrict,
    fiat_currency   integer                                 not null
        references emcd.fiat_currencies
            on update cascade on delete restrict,
    pay_method_id   integer                                 not null
        references emcd.p2p_pay_methods
            on update cascade on delete restrict,
    direction       varchar   default ''::character varying not null,
    rate            numeric   default 0                     not null,
    sell_min_amount numeric   default 0                     not null,
    sell_max_amount numeric   default 0                     not null,
    counter_details bytea                                   not null,
    terms           varchar   default ''::character varying not null,
    created_at      timestamp default now(),
    updated_at      timestamp default now(),
    deleted_at      timestamp,
    percent         numeric   default 0                     not null,
    auto_update     numeric   default 0                     not null
);


create table emcd.p2p_orders
(
    id                              serial
        primary key,
    offer_id                        integer                                 not null
        references emcd.p2p_offers
            on update cascade on delete restrict,
    client_account_id               integer                                 not null
        references emcd.users_accounts
            on update cascade on delete restrict,
    owner_account_id                integer                                 not null
        references emcd.users_accounts
            on update cascade on delete restrict,
    crypto_currency                 integer                                 not null
        references emcd.coins
            on update cascade on delete restrict,
    fiat_currency                   integer                                 not null
        references emcd.fiat_currencies
            on update cascade on delete restrict,
    status_id                       integer                                 not null
        references emcd.p2p_statuses
            on update cascade on delete restrict,
    pay_method_id                   integer                                 not null
        references emcd.p2p_pay_methods
            on update cascade on delete restrict,
    order_code                      varchar   default ''::character varying not null,
    rate                            numeric   default 0                     not null,
    amount_crypto                   numeric   default 0                     not null,
    amount_fiat                     numeric   default 0                     not null,
    counter_details                 bytea                                   not null,
    balance_block_transaction_id    integer,
    commission_block_transaction_id integer,
    chat_id                         varchar   default ''::character varying not null,
    seller_dispute_chat_id          varchar   default ''::character varying not null,
    buyer_dispute_chat_id           varchar   default ''::character varying not null,
    comment                         varchar   default ''::character varying,
    created_at                      timestamp default now(),
    updated_at                      timestamp default now()
);



create index p2p_orders_status_id_index
    on emcd.p2p_orders (status_id);

create table emcd.p2p_order_logs
(
    id            serial
        primary key,
    order_id      integer not null
        references emcd.p2p_orders
            on update cascade on delete restrict,
    old_status_id integer not null
        references emcd.p2p_statuses
            on update cascade on delete restrict,
    new_status_id integer not null
        references emcd.p2p_statuses
            on update cascade on delete restrict,
    reason_id     integer not null
        references emcd.p2p_status_reasons
            on update cascade on delete restrict,
    comment       varchar   default ''::character varying,
    created_at    timestamp default now()
);



create index p2p_order_logs_order_id_index
    on emcd.p2p_order_logs (order_id);

