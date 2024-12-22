alter table invoice
    add column
        required_payment decimal null;

update invoice set required_payment = amount + amount * buyer_fee;

alter table invoice
    alter column
        required_payment set not null;