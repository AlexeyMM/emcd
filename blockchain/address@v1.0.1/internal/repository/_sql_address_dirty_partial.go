package repository

// import (
// 	"github.com/Masterminds/squirrel"
//
// 	"code.emcdtech.com/emcd/blockchain/address/model"
// )
//
// type addressDirtyPartialSql struct {
// 	*model.AddressDirtyPartial
// }
//
// func newAddressDirtyPartialSql(partial *model.AddressDirtyPartial) *addressDirtyPartialSql {
//
// 	return &addressDirtyPartialSql{
// 		AddressDirtyPartial: partial,
// 	}
// }
//
// func (partial *addressDirtyPartialSql) applyToQuery(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
// 	if partial.IsDirty != nil {
// 		query = query.Set("is_dirty", *partial.IsDirty)
//
// 	}
//
// 	return query
// }
