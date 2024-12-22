package model_migration

import (
	"time"

	"github.com/google/uuid"
)

// SELECT
// ua.user_id_new as user_uuid,
// ua.coin_new as coin_code,
// ua.minpay as min_payout,
// COALESCE(ua.address, '') AS ua_address,
// COALESCE(a.address, '') AS a_address,
// aa.address AS aa_address,
// ua.created_at as created_at
// FROM
// emcd.autopay_addresses aa
// LEFT JOIN
// emcd.users_accounts ua ON aa.user_account_id = ua.id
// LEFT JOIN
// emcd.addresses a ON a.user_account_id = ua.id and token_id is null
// where aa.address <> '' and aa.address is not null
// and ua.account_type_id in (1,2)
// and ((ua.address <> '' and ua.address is not null and ua.address <> aa.address) or (a.address <> '' and a.address is not null and a.address <> aa.address))
// and ua.created_at > '2001-01-01 00:00:00.000' limit 100;

type AddressPersonalMigration struct {
	UserUuid  uuid.UUID
	CoinCode  string
	MinPayout float64
	UaAddress string
	AAddress  string
	AaAddress string
	CreatedAt time.Time
}

type AddressPersonalMigrations []*AddressPersonalMigration
