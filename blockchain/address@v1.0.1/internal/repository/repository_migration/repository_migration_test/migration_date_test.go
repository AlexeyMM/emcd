package repository_migration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/blockchain/address/internal/repository/repository_migration"
	"code.emcdtech.com/emcd/blockchain/address/model/model_migration"
)

func TestMigrationDate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer truncateTables(ctx, t, dbPool, "migration_date")

	repo := repository_migration.NewMigrationRepository(dbPool, dbPool)

	t.Run("put and get migration date addresses success", func(t *testing.T) {
		lastAt := time.Now().UTC().Truncate(time.Millisecond)

		err := repo.UpdateMigrationLastAt(ctx, model_migration.MigrationTableAddresses, lastAt)
		require.NoError(t, err)

		lastAtGet, err := repo.GetMigrationLastAt(ctx, model_migration.MigrationTableAddresses)
		require.NoError(t, err)
		require.NotNil(t, lastAtGet)

		require.Equal(t, *lastAtGet, lastAt)

	})

	t.Run("put and get migration date users accounts success", func(t *testing.T) {
		lastAt := time.Now().Add(10 * time.Second).UTC().Truncate(time.Millisecond)

		err := repo.UpdateMigrationLastAt(ctx, model_migration.MigrationTableUsersAccounts, lastAt)
		require.NoError(t, err)

		lastAtGet, err := repo.GetMigrationLastAt(ctx, model_migration.MigrationTableUsersAccounts)
		require.NoError(t, err)
		require.NotNil(t, lastAtGet)

		require.Equal(t, *lastAtGet, lastAt)

	})
}
