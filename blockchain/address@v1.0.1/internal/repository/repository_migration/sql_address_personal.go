package repository_migration

import (
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressPersonalSql struct {
	*model.AddressPersonal
}

func newAddressPersonalSql(addressPersonal *model.AddressPersonal) *addressPersonalSql {

	return &addressPersonalSql{
		AddressPersonal: addressPersonal,
	}
}

func (s *addressPersonalSql) toNamedArgs() *pgx.NamedArgs {

	return &pgx.NamedArgs{
		"id":         s.Id,
		"address":    s.Address,
		"user_uuid":  s.UserUuid,
		"network":    s.Network.ToString(),
		"min_payout": s.MinPayout,
		"deleted_at": s.DeletedAt,
		"updated_at": s.UpdatedAt,
		"created_at": s.CreatedAt,
	}
}
