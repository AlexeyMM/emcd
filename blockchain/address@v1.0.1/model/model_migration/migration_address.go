package model_migration

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// CREATE TABLE emcd.addresses (
// id serial4 NOT NULL,
// user_account_id int4 NULL,
// coin_id int4 NOT NULL,
// token_id int4 NULL,
// address varchar(128) NOT NULL,
// network_id VARCHAR(10) NULL,
// created_at timestamp NOT NULL,
// deleted_at timestamp NULL,
// address_offset int4 NULL,
// );

type AddressMigration struct {
	Id            int32
	UserAccountId int32 // select count(*) from emcd.addresses where user_account_id is null == 0
	CoinId        int32
	TokenId       sql.NullInt32
	Address       string
	NetworkId     sql.NullString
	CreatedAt     time.Time
	DeletedAt     sql.NullTime
	AddressOffset sql.NullInt32
	UserUuid      uuid.UUID
}

type AddressMigrations []*AddressMigration

// CREATE TABLE emcd.autopay_addresses
// (
// id              serial4                 NOT NULL,
// user_account_id int4                    NOT NULL,
// address         varchar                 NOT NULL,
// "percent"       int4      DEFAULT 0     NOT NULL,
// "label"         varchar                 NULL,
// created_at      timestamp DEFAULT now() NOT NULL,
// updated_at      timestamp DEFAULT now() NOT NULL,
// CONSTRAINT autopay_addresses_pkey PRIMARY KEY (id)
// );

type AutopayAddressMigration struct {
	Id            int32
	UserAccountId int32
	Address       string
	Percent       int32
	Label         sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
