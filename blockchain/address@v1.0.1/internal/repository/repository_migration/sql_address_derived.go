package repository_migration

import (
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressDerivedSql struct {
	*model.AddressDerived
}

func newAddressDerivedSql(addressDerived *model.AddressDerived) *addressDerivedSql {

	return &addressDerivedSql{
		AddressDerived: addressDerived,
	}
}

func (s *addressDerivedSql) toNamedArgs() *pgx.NamedArgs {

	return &pgx.NamedArgs{
		"address_uuid":   s.AddressUuid,
		"network_group":  s.NetworkGroup.NetworkGroupEnum,
		"master_key_id":  s.MasterKeyId,
		"derived_offset": s.DerivedOffset,
	}
}
