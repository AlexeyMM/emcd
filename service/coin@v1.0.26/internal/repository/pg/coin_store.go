package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"code.emcdtech.com/emcd/service/coin/internal/repository"
	"code.emcdtech.com/emcd/service/coin/model"
)

const (
	selectGetCoinFromLegacyIDSQL = `
SELECT
	id,
	is_active,
	title,
	description,
	sort_priority_mining,
	sort_priority_wallet,
	media_url,
	is_withdrawals_disabled,
	legacy_coin_id,
	mining_reward_type
FROM coins
WHERE legacy_coin_id = $1
`

	selectGetCoinSQL = `
SELECT
	id,
	is_active,
	title,
	description,
	sort_priority_mining,
	sort_priority_wallet,
	media_url,
	is_withdrawals_disabled,
	legacy_coin_id,
	mining_reward_type
FROM coins
WHERE id = $1
`

	selectCoins = `
SELECT
	id,
	is_active,
	title,
	description,
	sort_priority_mining,
	sort_priority_wallet,
	media_url,
	is_withdrawals_disabled,
	legacy_coin_id,
	mining_reward_type
FROM coins
ORDER BY sort_priority_wallet, sort_priority_mining
LIMIT $1 
OFFSET $2
`

	selectCoinsCountSQL = `
SELECT COUNT(*) FROM coins
`

	selectCoinsNetworksSQL = `
SELECT
	cn.coin_id,
	cn.network_id,
	cn.is_active,
	cn.title,
	cn.description,
	COALESCE(cn.contract_address, ''),
	cn.decimals,
	cn.minpay_mining,
	cn.withdrawal_fee,
	cn.withdrawal_min_limit,
	cn.withdrawal_fee_ttl_seconds,
	cn.is_mining,
	cn.is_wallet,
	cn.is_free_withdraw,
	cn.is_withdrawals_disabled,
	COALESCE(cn.hash_divisor_power_of_ten, 0),
	n.explorer_url,
	cn.priority 
FROM coins_networks cn
INNER JOIN networks n ON n.id = cn.network_id
`

	selectCoinNetworkSQL = `
SELECT
	cn.coin_id,
	cn.network_id,
	cn.is_active,
	cn.title,
	cn.description,
	COALESCE(cn.contract_address, ''),
	cn.decimals,
	cn.withdrawal_fee,
	cn.withdrawal_min_limit,
	cn.withdrawal_fee_ttl_seconds,
	cn.is_mining,
	cn.is_wallet,
	cn.is_free_withdraw,
	cn.is_withdrawals_disabled,
	COALESCE(cn.hash_divisor_power_of_ten, 0),
	n.explorer_url
FROM coins_networks cn
INNER JOIN networks n ON n.id = cn.network_id
WHERE cn.coin_id = $1
  AND cn.network_id = $2
`
)

type CoinStore struct {
	pool *pgxpool.Pool
}

func NewCoinStore(pool *pgxpool.Pool) *CoinStore {
	return &CoinStore{pool: pool}
}

func (s *CoinStore) GetCoinFromLegacyID(ctx context.Context, legacyCoinID int32) (*model.Coin, error) {
	rows, err := s.pool.Query(ctx, selectGetCoinFromLegacyIDSQL, legacyCoinID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("exec selectGetCoinFromLegacyIDSQL: %w", err)
	}
	defer rows.Close()
	var coin model.Coin
	for rows.Next() {
		err := rows.Scan(
			&coin.ID,
			&coin.IsActive,
			&coin.Title,
			&coin.Description,
			&coin.SortPriorityMining,
			&coin.SortPriorityWallet,
			&coin.MediaURL,
			&coin.IsWithdrawalsDisabled,
			&coin.LegacyCoinID,
			&coin.MiningRewardType,
		)
		if err != nil {
			return nil, fmt.Errorf("scan dataset selectGetCoinFromLegacyIDSQL: %w", err)
		}
		return &coin, nil //nolint:staticcheck
	}
	return nil, repository.ErrNotFound
}

func (s *CoinStore) GetCoin(ctx context.Context, coinID string) (*model.Coin, error) {
	rows, err := s.pool.Query(ctx, selectGetCoinSQL, coinID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("exec selectGetCoinSQL: %w", err)
	}
	defer rows.Close()
	var coin model.Coin
	for rows.Next() {
		err := rows.Scan(
			&coin.ID,
			&coin.IsActive,
			&coin.Title,
			&coin.Description,
			&coin.SortPriorityMining,
			&coin.SortPriorityWallet,
			&coin.MediaURL,
			&coin.IsWithdrawalsDisabled,
			&coin.LegacyCoinID,
			&coin.MiningRewardType,
		)
		if err != nil {
			return nil, fmt.Errorf("scan dataset selectGetCoinSQL: %w", err)
		}
		return &coin, nil //nolint:staticcheck
	}
	return nil, repository.ErrNotFound
}

func (s *CoinStore) GetCoins(ctx context.Context, limit, offset int32) ([]*model.Coin, int32, error) {
	var count int32
	if err := s.pool.QueryRow(ctx, selectCoinsCountSQL).Scan(&count); err != nil {
		return nil, 0, errors.WithStack(err)
	}

	rows, err := s.pool.Query(ctx, selectCoins, limit, offset)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	defer rows.Close()

	coins := make([]*model.Coin, 0)
	for rows.Next() {
		var item model.Coin
		if err = rows.Scan(
			&item.ID,
			&item.IsActive,
			&item.Title,
			&item.Description,
			&item.SortPriorityMining,
			&item.SortPriorityWallet,
			&item.MediaURL,
			&item.IsWithdrawalsDisabled,
			&item.LegacyCoinID,
			&item.MiningRewardType,
		); err != nil {
			return nil, 0, errors.WithStack(err)
		}
		coins = append(coins, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return coins, count, nil
}

func (s *CoinStore) GetCoinsNetworks(ctx context.Context) ([]*model.CoinNetwork, error) {
	rows, err := s.pool.Query(ctx, selectCoinsNetworksSQL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	coinsNetworks := make([]*model.CoinNetwork, 0)
	for rows.Next() {
		var item model.CoinNetwork
		if err = rows.Scan(
			&item.CoinID,
			&item.NetworkID,
			&item.IsActive,
			&item.Title,
			&item.Description,
			&item.ContractAddress,
			&item.Decimals,
			&item.MinPayMining,
			&item.WithdrawalFee,
			&item.WithdrawalMinLimit,
			&item.WithdrawalFeeTtlSeconds,
			&item.IsMining,
			&item.IsWallet,
			&item.IsFreeWithdraw,
			&item.IsWithdrawalsDisabled,
			&item.HashDivisorPowerOfTen,
			&item.ExplorerUrl,
			&item.Priority,
		); err != nil {
			return nil, errors.WithStack(err)
		}
		coinsNetworks = append(coinsNetworks, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return coinsNetworks, nil
}

func (s *CoinStore) GetCoinNetwork(ctx context.Context, coinID, networkID string) (*model.CoinNetwork, error) {
	var coinNetwork model.CoinNetwork
	if err := s.pool.QueryRow(ctx, selectCoinNetworkSQL, coinID, networkID).Scan(
		&coinNetwork.CoinID,
		&coinNetwork.NetworkID,
		&coinNetwork.IsActive,
		&coinNetwork.Title,
		&coinNetwork.Description,
		&coinNetwork.ContractAddress,
		&coinNetwork.Decimals,
		&coinNetwork.WithdrawalFee,
		&coinNetwork.WithdrawalMinLimit,
		&coinNetwork.WithdrawalFeeTtlSeconds,
		&coinNetwork.IsMining,
		&coinNetwork.IsWallet,
		&coinNetwork.IsFreeWithdraw,
		&coinNetwork.IsWithdrawalsDisabled,
		&coinNetwork.HashDivisorPowerOfTen,
		&coinNetwork.ExplorerUrl,
	); err != nil {
		return nil, errors.WithStack(err)
	}

	return &coinNetwork, nil
}
