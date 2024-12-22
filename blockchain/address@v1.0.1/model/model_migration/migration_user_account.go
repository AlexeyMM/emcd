package model_migration

import (
	"database/sql"

	accountingEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	"github.com/google/uuid"
)

type UserAccountMigration struct {
	Id            int32
	UserId        int32
	CoinId        int32
	AccountTypeId accountingEnum.AccountTypeIdWrapper
	Address       sql.NullString
	IsActive      sql.NullBool
	CreatedAt     sql.NullTime
	UserUuid      uuid.UUID
}

type UserAccountMigrations []*UserAccountMigration
