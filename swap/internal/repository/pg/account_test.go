package pg

import (
	"context"
	"testing"

	"code.emcdtech.com/b2b/swap/model"
	"github.com/stretchr/testify/require"
)

func TestAccount_AddGetSubAccount(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	acc := model.Account{
		ID: 1,
		Keys: &model.Secrets{
			ApiKey:    "12345",
			ApiSecret: "67890",
		},
		IsValid: true,
	}

	rep := NewAccount(db)

	err := rep.Add(ctx, &acc)
	require.NoError(t, err)

	account, err := rep.FindOne(ctx, &model.AccountFilter{
		ID: &acc.ID,
	})
	require.NoError(t, err)
	require.Equal(t, &acc, account)
}

func TestAccount_AddGetSubAccounts(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	acc := model.Account{
		ID: 1,
		Keys: &model.Secrets{
			ApiKey:    "12345",
			ApiSecret: "67890",
		},
		IsValid: true,
	}
	acc2 := model.Account{
		ID: 2,
		Keys: &model.Secrets{
			ApiKey:    "987654",
			ApiSecret: "99999",
		},
		IsValid: true,
	}

	rep := NewAccount(db)

	err := rep.Add(ctx, &acc)
	require.NoError(t, err)
	err = rep.Add(ctx, &acc2)
	require.NoError(t, err)

	accounts, err := rep.Find(ctx, &model.AccountFilter{})
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	var count int
	for _, account := range accounts {
		switch account.ID {
		case acc.ID:
			require.Equal(t, &acc, account)
			count++
		case acc2.ID:
			require.Equal(t, &acc2, account)
			count++
		}
	}
	require.Equal(t, count, 2)
}
