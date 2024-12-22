package pg_test

import (
	"context"
	"reflect"
	"testing"
	"time"

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
			SortPrioritySwap:      10,
		}

		const query = `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active, sort_priority_swap)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
`
		_, err := db.Exec(ctx, query, in.ID, in.Title, in.Description, in.MediaURL, in.SortPriorityMining,
			in.SortPriorityWallet, in.IsWithdrawalsDisabled, in.LegacyCoinID, in.MiningRewardType, in.IsActive, in.SortPrioritySwap)
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
			SortPrioritySwap:      10,
		}

		query := `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active, sort_priority_swap)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`
		_, err := db.Exec(ctx, query, in.ID, in.Title, in.Description, in.MediaURL, in.SortPriorityMining,
			in.SortPriorityWallet, in.IsWithdrawalsDisabled, in.LegacyCoinID, in.MiningRewardType, in.IsActive, in.SortPrioritySwap)
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
			SortPrioritySwap:      10,
		}

		query := `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active, sort_priority_swap)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`
		_, err := db.Exec(ctx, query, in.ID, in.Title, in.Description, in.MediaURL, in.SortPriorityMining,
			in.SortPriorityWallet, in.IsWithdrawalsDisabled, in.LegacyCoinID, in.MiningRewardType, in.IsActive, in.SortPrioritySwap)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins")
		}()

		coinID, err := repo.GetCoinFromLegacyID(ctx, in.LegacyCoinID)
		require.NoError(t, err)
		require.EqualValues(t, coinID, in)
	})

	t.Run("Get_Coin", func(t *testing.T) {
		inc := &model.Coin{
			ID:                    "etc_3",
			Title:                 "title",
			Description:           "desc",
			IsActive:              true,
			SortPriorityMining:    1,
			SortPriorityWallet:    1,
			MediaURL:              "media",
			IsWithdrawalsDisabled: true,
			MiningRewardType:      "mining_reward_type",
			SortPrioritySwap:      10,
		}

		const query = `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active, sort_priority_swap)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
`
		_, err := db.Exec(ctx, query, inc.ID, inc.Title, inc.Description, inc.MediaURL, inc.SortPriorityMining,
			inc.SortPriorityWallet, inc.IsWithdrawalsDisabled, inc.LegacyCoinID, inc.MiningRewardType, inc.IsActive, inc.SortPrioritySwap)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins")
		}()

		coins, err := repo.GetCoin(ctx, inc.ID)
		require.NoError(t, err)
		require.Equal(t, inc, coins) // together with migrations
	})

	t.Run("GetCoinNetworks", func(t *testing.T) {
		inc := &model.Coin{
			ID:                    "DOGE",
			Title:                 "title",
			Description:           "desc",
			IsActive:              true,
			SortPriorityMining:    1,
			SortPriorityWallet:    1,
			MediaURL:              "media",
			IsWithdrawalsDisabled: true,
			MiningRewardType:      "mining_reward_type",
		}

		const queryCoin = `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
`
		_, err := db.Exec(ctx, queryCoin, inc.ID, inc.Title, inc.Description, inc.MediaURL, inc.SortPriorityMining,
			inc.SortPriorityWallet, inc.IsWithdrawalsDisabled, inc.LegacyCoinID, inc.MiningRewardType, inc.IsActive)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins")
		}()
		id := "mainnet"
		isActive := true
		title := "Main Network"
		description := "Main network description"
		createdAt := time.Now()
		updatedAt := time.Now()
		explorer_url := "test"

		queryNetwork := `INSERT INTO networks (id, is_active, title, description, explorer_url, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

		// Execute the SQL statement
		_, err = db.Exec(ctx, queryNetwork, id, isActive, title, description, explorer_url, createdAt, updatedAt)
		require.NoError(t, err)
		defer func() {
			_, _ = db.Exec(ctx, "delete from networks")
		}()

		in := &model.CoinNetwork{
			CoinID:                  "DOGE",
			NetworkID:               "mainnet",
			IsActive:                true,
			Title:                   "Bitcoin",
			Description:             "Bitcoin mainnet",
			ContractAddress:         "test",
			Decimals:                8,
			MinPayMining:            0,
			WithdrawalFee:           0.001,
			WithdrawalMinLimit:      0.01,
			WithdrawalFeeTtlSeconds: 3600,
			IsMining:                false,
			IsWallet:                true,
			IsFreeWithdraw:          true,
			IsWithdrawalsDisabled:   false,
			HashDivisorPowerOfTen:   0,
			ExplorerUrl:             "test",
			Priority:                0,
		}

		query := `INSERT INTO coins_networks (coin_id, network_id, is_active, title, description, contract_address, 
		decimals, is_wallet, withdrawal_fee, withdrawal_min_limit, withdrawal_fee_ttl_seconds, is_mining, 
		is_free_withdraw, is_withdrawals_disabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

		_, err = db.Exec(ctx, query, in.CoinID,
			in.NetworkID,
			in.IsActive,
			in.Title,
			in.Description,
			in.ContractAddress,
			in.Decimals,
			in.IsWallet,
			in.WithdrawalFee,
			in.WithdrawalMinLimit,
			in.WithdrawalFeeTtlSeconds,
			in.IsMining,
			in.IsFreeWithdraw,
			in.IsWithdrawalsDisabled)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins_networks")
		}()
		coinNet, err := repo.GetCoinNetwork(ctx, in.CoinID, in.NetworkID)
		require.NoError(t, err)
		require.Equal(t, in, coinNet) // together with migrations
	})

	t.Run("GetCoinsNetworks", func(t *testing.T) {
		inc := &model.Coin{
			ID:                    "DOGE",
			Title:                 "title",
			Description:           "desc",
			IsActive:              true,
			SortPriorityMining:    1,
			SortPriorityWallet:    1,
			MediaURL:              "media",
			IsWithdrawalsDisabled: true,
			MiningRewardType:      "mining_reward_type",
		}

		const queryCoin = `
		INSERT INTO 
		    coins (id, title, description, media_url, sort_priority_mining, 
		           sort_priority_wallet, is_withdrawals_disabled, legacy_coin_id, mining_reward_type, is_active)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
`
		_, err := db.Exec(ctx, queryCoin, inc.ID, inc.Title, inc.Description, inc.MediaURL, inc.SortPriorityMining,
			inc.SortPriorityWallet, inc.IsWithdrawalsDisabled, inc.LegacyCoinID, inc.MiningRewardType, inc.IsActive)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins")
		}()
		id := "mainnet"
		isActive := true
		title := "Main Network"
		description := "Main network description"
		createdAt := time.Now()
		updatedAt := time.Now()
		explorer_url := "test"

		queryNetwork := `INSERT INTO networks (id, is_active, title, description, explorer_url, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

		// Execute the SQL statement
		_, err = db.Exec(ctx, queryNetwork, id, isActive, title, description, explorer_url, createdAt, updatedAt)
		require.NoError(t, err)
		defer func() {
			_, _ = db.Exec(ctx, "delete from networks")
		}()

		in := &model.CoinNetwork{
			CoinID:                  "DOGE",
			NetworkID:               "mainnet",
			IsActive:                true,
			Title:                   "Bitcoin",
			Description:             "Bitcoin mainnet",
			ContractAddress:         "test",
			Decimals:                8,
			MinPayMining:            0,
			WithdrawalFee:           0.001,
			WithdrawalMinLimit:      0.01,
			WithdrawalFeeTtlSeconds: 3600,
			IsMining:                false,
			IsWallet:                true,
			IsFreeWithdraw:          true,
			IsWithdrawalsDisabled:   false,
			HashDivisorPowerOfTen:   0,
			ExplorerUrl:             "test",
			Priority:                0,
		}

		query := `INSERT INTO coins_networks (coin_id, network_id, is_active, title, description, contract_address, 
		decimals, is_wallet, withdrawal_fee, withdrawal_min_limit, withdrawal_fee_ttl_seconds, is_mining, 
		is_free_withdraw, is_withdrawals_disabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

		_, err = db.Exec(ctx, query, in.CoinID,
			in.NetworkID,
			in.IsActive,
			in.Title,
			in.Description,
			in.ContractAddress,
			in.Decimals,
			in.IsWallet,
			in.WithdrawalFee,
			in.WithdrawalMinLimit,
			in.WithdrawalFeeTtlSeconds,
			in.IsMining,
			in.IsFreeWithdraw,
			in.IsWithdrawalsDisabled)
		require.NoError(t, err)

		defer func() {
			_, _ = db.Exec(ctx, "delete from coins_networks")
		}()
		coinNet, err := repo.GetCoinsNetworks(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(coinNet)) // together with migrations
	})
}
