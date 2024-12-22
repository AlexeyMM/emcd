package repository_migration

import (
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressOldSql struct {
	*model.AddressOld
}

func newAddressOldSql(addressOld *model.AddressOld) *addressOldSql {

	return &addressOldSql{
		AddressOld: addressOld,
	}
}

func (s *addressOldSql) toNamedArgs() *pgx.NamedArgs {

	return &pgx.NamedArgs{
		"id":              s.Id,
		"address":         s.Address,
		"user_uuid":       s.UserUuid,
		"address_type":    s.AddressType.Number(),
		"network":         s.Network.ToString(),
		"user_account_id": s.UserAccountId,
		"coin":            s.Coin,
		"created_at":      s.CreatedAt,
	}
}
