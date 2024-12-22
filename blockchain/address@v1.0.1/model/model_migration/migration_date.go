package model_migration

import "time"

type MigrationTableName string

const MigrationTableAddresses MigrationTableName = "addresses"
const MigrationTableUsersAccounts MigrationTableName = "users_accounts"
const MigrationTablePersonalAddresses MigrationTableName = "personal_addresses"

type MigrationDate struct {
	TableName string
	LastAt    time.Time
}
