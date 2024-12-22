create type transaction_confirmation_status as enum ('confirming', 'confirmed');

create table invoice_transaction
(
    hash                text                            not null primary key,
    invoice_id          uuid                            not null references invoice (id),
    amount              decimal                         not null,
    received_address    text                            not null references deposit_address (address),
    confirmation_status transaction_confirmation_status not null,
    created_at          timestamptz                     not null default now()
);

create index invoice_transaction_invoice_id_idx on invoice_transaction (invoice_id, created_at);