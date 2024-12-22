package repository_migration

import (
	"context"
	"time"

	"code.emcdtech.com/emcd/blockchain/address/model"
	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

type MigrationRepository interface {
	GetMigrationLastAt(ctx context.Context, tableName model_migration.MigrationTableName) (*time.Time, error)
	UpdateMigrationLastAt(ctx context.Context, tableName model_migration.MigrationTableName, lastAt time.Time) error

	GetAddressMigrations(ctx context.Context, lastAtGt time.Time, limit uint64) (*uint64, model_migration.AddressMigrations, error)
	GetUserAccountMigrations(ctx context.Context, lastAtGt time.Time, limit uint64) (*uint64, model_migration.UserAccountMigrations, error)
	GetAddressPersonalMigrations(ctx context.Context, lastAtGt time.Time, limit uint64) (*uint64, model_migration.AddressPersonalMigrations, error)

	AddOldAddressDirect(ctx context.Context, address *model.AddressOld) error
	AddNewAddressDirect(ctx context.Context, address *model.Address) error
	AddNewDerivedAddressDirect(ctx context.Context, address *model.AddressDerived) error
	AddPersonalAddressDirect(ctx context.Context, address *model.AddressPersonal) error

	WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) error
}
