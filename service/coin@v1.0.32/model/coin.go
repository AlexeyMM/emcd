package model

type (
	Coin struct {
		ID                    string         `json:"id"                      db:"id"`
		Title                 string         `json:"title"                   db:"title"`
		Description           string         `json:"description"             db:"description"`
		IsActive              bool           `json:"is_active"               db:"is_active"`
		SortPriorityMining    int32          `json:"sort_priority_mining"    db:"sort_priority_mining"`
		SortPriorityWallet    int32          `json:"sort_priority_wallet"    db:"sort_priority_wallet"`
		MediaURL              string         `json:"media_url"               db:"media_url"`
		IsWithdrawalsDisabled bool           `json:"is_withdrawals_disabled" db:"is_withdrawals_disabled"`
		LegacyCoinID          int32          `json:"legacy_coin_id"          db:"legacy_coin_id"`
		MiningRewardType      string         `json:"mining_reward_type"      db:"mining_reward_type"`
		Networks              []*CoinNetwork `json:"networks"                db:"-"`
		SortPrioritySwap      int32          `json:"sort_priority_swap" db:"sort_priority_swap"`
	}

	CoinNetwork struct {
		CoinID                  string  `json:"coin_id"                    db:"coin_id"`
		NetworkID               string  `json:"network_id"                 db:"network_id"`
		Title                   string  `json:"title"                      db:"title"`
		IsActive                bool    `json:"is_active"                  db:"is_active"`
		Description             string  `json:"description"                db:"description"`
		ContractAddress         string  `json:"contract_address"           db:"contract_address"`
		Decimals                int32   `json:"decimals"                   db:"decimals"`
		MinPayMining            float64 `json:"minpay_mining"              db:"minpay_mining"`
		WithdrawalFee           float64 `json:"withdrawal_fee"             db:"withdrawal_fee"`
		WithdrawalMinLimit      float64 `json:"withdrawal_min_limit"       db:"withdrawal_min_limit"`
		WithdrawalFeeTtlSeconds int32   `json:"withdrawal_fee_ttl_seconds" db:"withdrawal_fee_ttl_seconds"`
		IsMining                bool    `json:"is_mining"                  db:"is_mining"`
		IsWallet                bool    `json:"is_wallet"                  db:"is_wallet"`
		IsFreeWithdraw          bool    `json:"is_free_withdraw"           db:"is_free_withdraw"`
		IsWithdrawalsDisabled   bool    `json:"is_withdrawals_disabled"    db:"is_withdrawals_disabled"`
		HashDivisorPowerOfTen   int32   `json:"hash_divisor_power_of_ten"  db:"hash_divisor_power_of_ten"`
		ExplorerUrl             string  `json:"explorer_url"               db:"explorer_url"`
		Priority                int32   `json:"priority"                   db:"priority"`
	}
)
