package vault_test

import (
	"context"
	"testing"
	"time"

	apivault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/sdk/config"
	"code.emcdtech.com/emcd/sdk/config/vault"
)

func TestVaultProvider(t *testing.T) {
	testCases := []struct {
		name   string
		source string
		want   string
		err    bool
	}{
		{
			name:   "unregistered scheme",
			source: "secret://some-secret-string",
			want:   "secret://some-secret-string",
			err:    false,
		},
		{
			name:   "secret-vault scheme",
			source: "secret-vault://secret/key1/value1",
			want:   "value_1",
			err:    false,
		},
		{
			name:   "no scheme",
			source: "123",
			want:   "123",
			err:    false,
		},
		{
			name:   "empty value",
			source: "",
			want:   "",
			err:    false,
		},
		{
			name:   "secret-vault scheme (not found mount)",
			source: "secret-vault://secret_not_found/key1/value1",
			want:   "secret-vault://secret_not_found/key1/value1",
			err:    true,
		},
		{
			name:   "secret-vault scheme (not found key)",
			source: "secret-vault://secret/not_found_key/value1",
			want:   "secret-vault://secret/not_found_key/value1",
			err:    true,
		},
		{
			name:   "secret-vault scheme (not found value)",
			source: "secret-vault://secret/key1/value_not_found",
			want:   "secret-vault://secret/key1/value_not_found",
			err:    true,
		},
		{
			name:   "secret-vault scheme (incorrect path)",
			source: "secret-vault://secret/key1/value_not_found/12",
			want:   "secret-vault://secret/key1/value_not_found/12",
			err:    true,
		},
	}

	cfgVault := apivault.DefaultConfig()
	require.NoError(t, cfgVault.ReadEnvironment())
	client, err := apivault.NewClient(cfgVault)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = client.KVv2("secret").Put(ctx, "key1", map[string]interface{}{"value1": "value_1"})
	require.NoError(t, err)
	vaultProvider, err := vault.NewProvider()
	require.NoError(t, err)

	onSet := config.ValuesProviderToOnSetFn(vaultProvider)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r, err := onSet(testCase.source)
			if testCase.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, testCase.want, r)
		})
	}
}
