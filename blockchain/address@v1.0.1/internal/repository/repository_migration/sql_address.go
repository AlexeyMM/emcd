package repository_migration

import (
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressSql struct {
	*model.Address
}

func newAddressSql(address *model.Address) *addressSql {

	return &addressSql{
		Address: address,
	}
}

func (s *addressSql) toNamedArgs() *pgx.NamedArgs {

	return &pgx.NamedArgs{
		"id":              s.Id,
		"address":         s.Address.Address,
		"user_uuid":       s.UserUuid,
		"processing_uuid": s.ProcessingUuid,
		"address_type":    s.AddressType.Number(),
		"network_group":   s.NetworkGroup.ToString(),
		"created_at":      s.CreatedAt,
	}
}
