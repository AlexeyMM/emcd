alter table invoice_form
    add column external_id varchar(36),
    add column expires_at  timestamptz;

create index invoice_form_external_id_idx on invoice_form (external_id);
