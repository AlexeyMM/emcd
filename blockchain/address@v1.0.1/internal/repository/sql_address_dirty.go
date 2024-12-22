package repository

import (
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/blockchain/address/model"
)

type addressDirtySql struct {
	*model.AddressDirty
}

func newAddressDirtySql(addressDirty *model.AddressDirty) *addressDirtySql {

	return &addressDirtySql{
		AddressDirty: addressDirty,
	}
}

func (s *addressDirtySql) toNamedArgs() *pgx.NamedArgs {

	return &pgx.NamedArgs{
		"address":    s.Address,
		"network":    s.Network.ToString(),
		"is_dirty":   s.IsDirty,
		"updated_at": s.UpdatedAt,
		"created_at": s.CreatedAt,
	}
}
