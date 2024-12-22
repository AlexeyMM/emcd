create table emcd.wallets_accounting_actions

(
    action_id           uuid                                           not null
        primary key,
    status              varchar   default ''::character varying        not null,
    created_at          timestamp default timezone('utc'::text, now()) not null,
    amount              numeric                                        not null,
    fee                 numeric,
    sender_account_id   integer                                        not null,
    receiver_account_id integer                                        not null,
    coin_id             integer                                        not null,
    type                integer                                        not null,
    receiver_address    varchar(255),
    comment             text,
    token_id            integer,
    user_id             integer,
    send_message        boolean   default false
)
;

CREATE OR REPLACE FUNCTION wallets_notify_event() RETURNS TRIGGER AS
$$
DECLARE
data json;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        data = row_to_json(NEW);
END IF;
    PERFORM pg_notify('wallets-api-events', data::text);
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_wallets_accounting_actions
    AFTER INSERT
    ON emcd.wallets_accounting_actions
    FOR EACH ROW
    EXECUTE PROCEDURE wallets_notify_event()
;

ALTER TABLE emcd.wallets_accounting_actions ALTER COLUMN type DROP NOT NULL
;

ALTER TABLE emcd.wallets_accounting_actions ALTER COLUMN status DROP DEFAULT
;


alter table emcd.transactions
    add column "action_id" uuid references emcd.wallets_accounting_actions;
create unique index transactions_action_id_amount_type_index
    on emcd.transactions (action_id, amount, type);