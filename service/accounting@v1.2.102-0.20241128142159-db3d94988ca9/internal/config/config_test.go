package config_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/config"
	"github.com/google/uuid"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigSuccess(t *testing.T) {
	defer func() {
		err := os.Unsetenv("POSTGRES_CONNECTION_STRING")
		require.NoError(t, err)
	}()

	dsn := "postgresql://"
	err := os.Setenv("POSTGRES_CONNECTION_STRING", dsn)
	require.NoError(t, err)
	id1 := uuid.New()
	id2 := uuid.New()
	err = os.Setenv("IGNORE_REFERRAL_PAYMENT_BY_USER_ID", id1.String()+","+id2.String())
	require.NoError(t, err)

	cfg, err := config.NewConfig()
	require.NoError(t, err)
	assert.Equal(t, dsn, cfg.PostgresConnectionString)
	assert.Equal(t, ":8080", cfg.Port)
	assert.Equal(t, []uuid.UUID{id1, id2}, cfg.IgnoreReferralPaymentByUserID)
}

func TestNewConfigSuccessEmptyValues(t *testing.T) {
	cfg, err := config.NewConfig()
	require.NoError(t, err)
	assert.Empty(t, cfg.PostgresConnectionString)
	assert.Equal(t, ":8080", cfg.Port)
}
