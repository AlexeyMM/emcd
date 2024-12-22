CREATE SCHEMA IF NOT EXISTS emcd;
CREATE SCHEMA IF NOT EXISTS histories;

CREATE TABLE emcd.users (
    email CHARACTER VARYING,
    username CHARACTER VARYING NOT NULL,
    password CHARACTER VARYING NOT NULL,
    nopay BOOLEAN DEFAULT FALSE,
    is_email_notifications_on BOOLEAN DEFAULT TRUE,
    is_tg_notifications_on BOOLEAN DEFAULT FALSE,
    tg_id BIGINT,
    auth_secret CHARACTER VARYING,
    api_key CHARACTER VARYING,
    LANGUAGE CHARACTER VARYING DEFAULT 'en'::CHARACTER VARYING,
    timezone CHARACTER VARYING,
    is_donatation_on BOOLEAN DEFAULT FALSE,
    is_coinhold_enabled BOOLEAN DEFAULT TRUE,
    is_employee BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    id INTEGER NOT NULL,
    ref_id INTEGER DEFAULT 1 NOT NULL,
    parent_id INTEGER,
    master_fee NUMERIC,
    master_id INTEGER,
    phone CHARACTER VARYING,
    is_phone_verified BOOLEAN DEFAULT FALSE,
    is_phone_2fa_enabled BOOLEAN DEFAULT FALSE NOT NULL,
    twa_code CHARACTER VARYING,
    is_hedge_enabled BOOLEAN DEFAULT FALSE NOT NULL,
    company_name CHARACTER VARYING(255),
    is_autopay_disabled BOOLEAN DEFAULT FALSE,
    kyc_status CHARACTER VARYING,
    wb_link_id CHARACTER VARYING,
    pass_updated_at TIMESTAMP WITHOUT TIME ZONE,
    primary_currency CHARACTER VARYING DEFAULT 'usd'::CHARACTER VARYING NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    free_withdraw BOOLEAN DEFAULT TRUE NOT NULL,
    def_coin_id INTEGER DEFAULT 0 NOT NULL,
    new_id UUID,
    kyc_idenfy_status SMALLINT DEFAULT 0 NOT NULL,
    suspended TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE emcd.vip_users (
    id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    title CHARACTER VARYING(64),
    description CHARACTER VARYING(1024),
    is_divided_address BOOLEAN NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    is_blocked_ref_tier BOOLEAN DEFAULT FALSE NOT NULL,
    spec_ref_fee NUMERIC
);

CREATE TABLE histories.segment_userids (
    user_id INTEGER NOT NULL,
    segment_id INTEGER DEFAULT 1 NOT NULL
);
