package pg

import (
	"code.emcdtech.com/b2b/swap/model"
	"github.com/Masterminds/squirrel"
)

type internalTransferPartialSql struct {
	*model.InternalTransferPartial
}

func newInternalTransferPartialSql(filter *model.InternalTransferPartial) *internalTransferPartialSql {
	return &internalTransferPartialSql{
		InternalTransferPartial: filter,
	}
}

func (it *internalTransferPartialSql) applyToSql(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
	if it.Status != nil {
		query = query.
			Set("status", it.Status)
	}
	return query
}
