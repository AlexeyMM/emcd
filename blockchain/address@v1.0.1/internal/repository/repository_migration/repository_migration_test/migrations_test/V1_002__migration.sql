-- migrate:up

CREATE TABLE migration_date
(
    table_name TEXT      NOT NULL,
    last_at    TIMESTAMP NOT NULL, -- last migration date
    CONSTRAINT migration_date_uniq UNIQUE (table_name)
);

insert into migration_date (table_name, last_at)
values ('addresses', '2000-01-01 00:00:00.000');
insert into migration_date (table_name, last_at)
values ('users_accounts', '2000-01-01 00:00:00.000');

-- migrate:down

-- drop table if exists migration_date cascade;
