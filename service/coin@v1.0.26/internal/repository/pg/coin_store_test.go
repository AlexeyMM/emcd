package pg_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/coin/internal/repository/pg"
	"code.emcdtech.com/emcd/service/coin/model"
)

func TestContract_Create(t *testing.T) {
	ctx := context.Background()
	repo := pg.NewCoinStore(db)

	t.Run("Full", func(t *testing.T) {
		in := &model.Coin{
			ID:                    "etc_2",
			Title:                 "title",
			Description:           "desc",
			IsActive:              true,
			SortPriorityMining:    1,
			SortPriorityWallet:    1,
			MediaURL:              "media",
			IsWithdrawalsDisabled: true,
			MiningRewardType:      "mining_reward_type",
		}

		const query = `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
`
		_, err := db.Exec(ctx, query, in.ID, in.Title, in.Description, in.MediaURL, in.SortPriorityMining,
			in.SortPriorityWallet, in.IsWithdrawalsDisabled, in.LegacyCoinID, in.MiningRewardType, in.IsActive)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins")
		}()

		coins, size, err := repo.GetCoins(ctx, 100, 0)
		require.NoError(t, err)
		require.Equal(t, size, int32(14)) // together with migrations
		require.Len(t, coins, 14)

		found := false
		for _, coin := range coins {
			if coin.ID == "etc_2" {
				found = true
				require.True(t, reflect.DeepEqual(in, coin))
			}
		}
		require.True(t, found)
	})

	t.Run("CoinFromLegacyID", func(t *testing.T) {
		in := &model.Coin{
			ID:                    "etc_2_3",
			Title:                 "title",
			Description:           "desc",
			IsActive:              true,
			SortPriorityMining:    1,
			SortPriorityWallet:    1,
			MediaURL:              "media",
			LegacyCoinID:          100,
			IsWithdrawalsDisabled: true,
			MiningRewardType:      "mining_reward_type",
		}

		query := `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
		_, err := db.Exec(ctx, query, in.ID, in.Title, in.Description, in.MediaURL, in.SortPriorityMining,
			in.SortPriorityWallet, in.IsWithdrawalsDisabled, in.LegacyCoinID, in.MiningRewardType, in.IsActive)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins")
		}()

		coinID, err := repo.GetCoinFromLegacyID(ctx, in.LegacyCoinID)
		require.NoError(t, err)
		require.EqualValues(t, coinID, in)
	})

	t.Run("GetCoin", func(t *testing.T) {
		in := &model.Coin{
			ID:                    "etc_2_4",
			Title:                 "title",
			Description:           "desc",
			IsActive:              true,
			SortPriorityMining:    1,
			SortPriorityWallet:    1,
			MediaURL:              "media",
			LegacyCoinID:          101,
			IsWithdrawalsDisabled: true,
			MiningRewardType:      "mining_reward_type",
		}

		query := `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
		_, err := db.Exec(ctx, query, in.ID, in.Title, in.Description, in.MediaURL, in.SortPriorityMining,
			in.SortPriorityWallet, in.IsWithdrawalsDisabled, in.LegacyCoinID, in.MiningRewardType, in.IsActive)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins")
		}()

		coinID, err := repo.GetCoinFromLegacyID(ctx, in.LegacyCoinID)
		require.NoError(t, err)
		require.EqualValues(t, coinID, in)
	})
}
