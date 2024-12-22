package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	sdkConfig "code.emcdtech.com/emcd/sdk/config"
)

func TestConfig(t *testing.T) {
	type testCfg struct {
		DefaultDonations DefaultDonations `env:"DEFAULT_DONATIONS_JSON,required"`
	}
	id1, id2 := uuid.New(), uuid.New()
	d1, d2 := decimal.NewFromFloat(1.5534), decimal.NewFromFloat(2.0)
	os.Setenv("DEFAULT_DONATIONS_JSON", fmt.Sprintf(`{"%v":%s, "%v":%s}`, id1, d1.String(), id2, d2.String()))
	cfg, err := sdkConfig.New[testCfg]()
	require.NoError(t, err)
	require.Equal(t, cfg.DefaultDonations[id1], d1)
	require.Equal(t, cfg.DefaultDonations[id2], d2)
}
