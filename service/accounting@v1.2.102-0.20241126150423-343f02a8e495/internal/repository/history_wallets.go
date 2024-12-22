package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"

	"code.emcdtech.com/emcd/service/accounting/model"
)

const (
	Internal = "Internal server error"

	tokenUSDTBEP20ID = 1
	tokenUSDCBEP20ID = 2
	tokenUSDTTRC20ID = 3
	tokenUSDCTRC20ID = 4

	txStatusInProgress = 1
	txStatusDone       = 2
	txStatusDeclined   = 3

	bep20 = "bep20"
	trc20 = "trc20"

	coinUsdtID = 10
	coinUsdcID = 11
)

const (
	selectHistoryCountStatement = "COUNT(*)"

	selectHistoryListStmt = `
		o.id,
		o.amount,
		EXTRACT(EPOCH FROM o.created_at)::integer AS created_at,
		o.created_at AS date,
		o.type,
		o.coin_id,
		COALESCE(t.token_id, 0) AS token_id,
		t.hash AS tx_id,
		t.comment,
		tfee.amount AS fee,
		t.receiver_address AS address,
		ac.id AS coinhold_id,
		coalesce(rt.coin_id, ep.coin_id) AS exchange_to_coin_id,
		coalesce(st.amount, t.amount) AS exchange_amount_sent,
		coalesce(rt.amount, ep.amount) AS exchange_amount_receive,
		coalesce(e.rate, trunc(ep.amount / t.amount, 8)) AS exchange_rate,
		coalesce(e.is_success, case when ep.id is not null then true end) AS exchange_is_success,
		CASE
           WHEN t.type in (22, 23, 54, 55)
               THEN (SELECT u.email
                     FROM emcd.users u
                              INNER JOIN users_accounts ua ON u.id = ua.user_id
                     WHERE ua.id = t.sender_account_id)
           END                                   AS referral_email,
		CASE
            WHEN COALESCE(p2p_balance_block.status_id, p2p_unblock.status_id, p2p_commission.status_id) IN (5, 7)
                THEN 2
            WHEN COALESCE(p2p_balance_block.status_id, p2p_unblock.status_id, p2p_commission.status_id) IN (2, 8)
                THEN 3
            ELSE 1
        END AS p2p_status,
        COALESCE(p2p_balance_block.id, p2p_unblock.id, p2p_commission.id, 0) AS p2p_order_id,
        CASE
            WHEN t.type = 1
                THEN COALESCE(wallet_t.network_id, '')
            WHEN t.type = 30
                THEN COALESCE(withdraw_t.network_id, '')
            ELSE ''
            END                                   AS network_id,
        CASE
            WHEN t.type = 1
                THEN COALESCE(wallet_t.coin_str_id, '')
            ELSE ''
            END                                   as coin_str_id`

	getWalletHistoryBaseQuery = `
	SELECT
		%s
	FROM emcd.operations o
	RIGHT JOIN emcd.users_accounts ua ON o.account_id = ua.id
	LEFT JOIN emcd.coins c ON c.id = ua.coin_id
	RIGHT JOIN emcd.transactions t ON o.transaction_id = t.id
	LEFT JOIN emcd.accounts_coinhold ac ON o.account_id = ac.account_id
	LEFT JOIN emcd.coinhold_types ct ON ac.coinhold_type_id = ct.id
	LEFT JOIN emcd.exchanges e ON e.send_transaction_id = o.transaction_id
	LEFT JOIN emcd.transactions rt ON rt.id = e.receive_transaction_id
	LEFT JOIN emcd.transactions st ON st.id = e.send_transaction_id
	LEFT JOIN emcd.transactions er ON er.action_id = t.action_id AND er.type = 58 -- exchange rollback
	LEFT JOIN emcd.transactions ep ON ep.action_id = t.action_id AND ep.type = 60 -- exchange payout
	LEFT JOIN emcd.wallet_transactions wallet_t
					   ON wallet_t.tx_id = t.hash and t.type = 1 AND wallet_t.coin_id = ua.coin_id AND wallet_t.amount = o.amount -- TODO: for hotFix duplicates
	LEFT JOIN emcd.withdraw_transactions withdraw_t ON withdraw_t.tx_id = o.transaction_id
	LEFT JOIN emcd.transactions tfee on withdraw_t.fee_tx_id = tfee.id
	LEFT JOIN emcd.p2p_orders p2p_balance_block ON o.type <> 67 AND o.transaction_id = p2p_balance_block.balance_block_transaction_id
    LEFT JOIN emcd.transactions_blocks tb ON tb.unblock_transaction_id = o.transaction_id
    LEFT JOIN emcd.p2p_orders p2p_unblock ON p2p_unblock.balance_block_transaction_id = tb.block_transaction_id
    LEFT JOIN emcd.p2p_orders p2p_commission ON (o.transaction_id = p2p_commission.commission_block_transaction_id)
	WHERE
	    ua.user_id = $1
	    AND e.is_success IS NOT FALSE
		AND er.id IS NULL -- no rollback transaction
		AND (
			CASE
				WHEN o.type = 74 THEN ua.account_type_id <> 1
				WHEN o.type = 75 THEN ua.account_type_id <> 7
				ELSE true
			END)
		%s %s %s`

	orderStmt = "ORDER BY o.id DESC"
)

var (
	lateStatusDoneTxTypeIDs = map[int]struct{}{
		int(model.IncomeBillTrTypeID):               {},
		int(model.UserPaysPoolComsTrTypeID):         {},
		int(model.UserPaysPoolDonationsTrTypeID):    {},
		int(model.PoolPaysUsersBalanceTrTypeID):     {},
		int(model.WithdrawalWithCommissionTrTypeID): {},
		int(model.WalletMiningTransferTrTypeID):     {},
		int(model.ExchBlockTrTypeID):                {},
	}

	p2pTxTypeIDs = map[int]struct{}{
		int(model.P2PSellTrType):           {},
		int(model.P2PBuyTrType):            {},
		int(model.P2PSellCommissionTrType): {},
	}
)

type WalletsHistory interface {
	GetWalletHistory(ctx context.Context, filter *model.HistoryInput) (*model.HistoryOutput, error)
	GetWalletHistoryTotal(ctx context.Context, param *model.HistoryInput) (int64, error)
}

type walletsHistory struct {
	pool              *pgxpool.Pool
	walletCoinsStrIDs map[int]string
}

func NewWalletsHistory(pool *pgxpool.Pool, walletCoinsStrIDs map[int]string) WalletsHistory {
	return &walletsHistory{
		pool:              pool,
		walletCoinsStrIDs: walletCoinsStrIDs,
	}
}

func (r *walletsHistory) GetWalletHistory(ctx context.Context, filter *model.HistoryInput) (*model.HistoryOutput, error) {
	var sqlParams []any
	sqlParams = append(sqlParams, int(filter.UserID))
	whereStmt, sqlParams, err := r.buildFilterQuery(filter, sqlParams)
	if err != nil {
		return nil, fmt.Errorf("buildFilterQuery: %w", err)
	}
	limitStmt, sqlParams := r.buildLimits(filter, sqlParams)

	q := fmt.Sprintf(getWalletHistoryBaseQuery, selectHistoryListStmt, whereStmt, orderStmt, limitStmt)
	rows, err := r.pool.Query(ctx, q, sqlParams...)
	if err != nil {
		return nil, fmt.Errorf("GetWalletHistory Query: %w: %s", err, Internal)
	}
	defer rows.Close()

	var historyTxs []*model.Wallet
	for rows.Next() {
		item := new(model.Wallet)
		if err = rows.Scan(
			&item.ID,
			&item.Amount,
			&item.CreatedAt,
			&item.Date,
			&item.Type,
			&item.CoinID,
			&item.TokenID,
			&item.TxID,
			&item.Comment,
			&item.Fee,
			&item.Address,
			&item.CoinholdID,
			&item.ExchangeToCoinID,
			&item.ExchangeAmountSent,
			&item.ExchangeAmountReceive,
			&item.ExchangeRate,
			&item.ExchangeIsSuccess,
			&item.ReferralEmail,
			&item.P2PStatus,
			&item.P2POrderID,
			&item.NetworkID,
			&item.CoinStrID,
		); err != nil {
			return nil, fmt.Errorf("GetWalletHistory Scan: %w: %s", err, Internal)
		}

		var coinName string
		coinNameVal, ok := r.walletCoinsStrIDs[item.CoinID]
		if ok {
			coinName = coinNameVal
		}

		switch item.Type {
		case int(model.IncomeBillTrTypeID), int(model.WithdrawalWithCommissionTrTypeID):
			if item.NetworkID == "" {
				if item.TokenID == tokenUSDTBEP20ID || item.TokenID == tokenUSDCBEP20ID {
					item.NetworkID = bep20
				} else if item.TokenID == tokenUSDTTRC20ID || item.TokenID == tokenUSDCTRC20ID {
					item.NetworkID = trc20
				} else {
					item.NetworkID = coinName
				}
			}
			if item.CoinStrID == "" {
				item.CoinStrID = coinName
			}

		case int(model.PoolPaysUsersBalanceTrTypeID), int(model.WalletMiningTransferTrTypeID):
			item.NetworkID, item.CoinStrID = coinName, coinName
		default:
			item.CoinStrID = coinName
		}

		if item.TokenID == 0 {
			switch item.CoinID {
			case coinUsdtID:
				item.TokenID = tokenUSDTBEP20ID
			case coinUsdcID:
				item.TokenID = tokenUSDCBEP20ID
			}
		}
		if item.TxID != nil && *item.TxID == "" && item.Type == int(model.MiningWalletPayoutTrType) {
			item.TxID = nil
		}

		item.Status = getStatus(item)

		if item.Type != int(model.WalletToWalletTrTypeID) && item.Type != int(model.MiningToWalletTrTypeID) {
			item.Comment = nil
		}

		if _, ok := p2pTxTypeIDs[item.Type]; ok {
			item.OrderID = &item.P2POrderID
		}

		historyTxs = append(historyTxs, item)
	}

	return &model.HistoryOutput{
		TotalCount:    0,
		IncomesSum:    nil,
		PayoutsSum:    nil,
		HasNewIncome:  nil,
		HasNewPayouts: nil,
		Incomes:       nil,
		Payouts:       nil,
		Wallets:       historyTxs,
	}, nil
}

func (r *walletsHistory) GetWalletHistoryTotal(ctx context.Context, input *model.HistoryInput) (int64, error) {
	var sqlParams []any
	sqlParams = append(sqlParams, input.UserID)
	whereStatements, params, err := r.buildFilterQuery(input, sqlParams)
	if err != nil {
		return 0, fmt.Errorf("buildFilterQuery: %w", err)
	}

	q := fmt.Sprintf(getWalletHistoryBaseQuery, selectHistoryCountStatement, whereStatements, "", "")
	var total int64
	if err := r.pool.QueryRow(ctx, q, params...).Scan(&total); err != nil {
		return 0, fmt.Errorf("GetWalletHistoryTotal QueryRow: %w: %s", err, Internal)
	}

	return total, nil
}

func (r *walletsHistory) buildFilterQuery(params *model.HistoryInput, sqlParams []any) (string, []any, error) {
	var filterQuery string
	statementCounter := 2

	if len(params.CoinsIDs) > 0 {
		filterQuery += fmt.Sprintf(" AND o.coin_id = ANY($%d)", statementCounter)
		sqlParams = append(sqlParams, pq.Array(params.CoinsIDs))
		statementCounter++
	}
	filterQuery += fmt.Sprintf(" AND o.type = ANY($%d)", statementCounter)
	sqlParams = append(sqlParams, pq.Array(params.TransactionTypesIDs))
	statementCounter++
	if len(params.AccountTypeIDs) > 0 {
		filterQuery += fmt.Sprintf(" AND ua.account_type_id = ANY($%d)", statementCounter)
		sqlParams = append(sqlParams, pq.Array(params.AccountTypeIDs))
		statementCounter++
	}
	if params.From != "" {
		from, err := time.Parse(time.DateOnly, params.From)
		if err != nil {
			return "", nil, fmt.Errorf("params.From: time.Parse: %w", err)
		}
		filterQuery += fmt.Sprintf(" AND date(o.created_at) >= $%d", statementCounter)
		year, month, day := from.Date()
		sqlParams = append(sqlParams, fmt.Sprintf("%d-%d-%d", year, month, day))
		statementCounter++
	}
	if params.To != "" {
		to, err := time.Parse(time.DateOnly, params.To)
		if err != nil {
			return "", nil, fmt.Errorf("params.To: time.Parse: %w", err)
		}
		filterQuery += fmt.Sprintf(" AND date(o.created_at) <= $%d", statementCounter)
		year, month, day := to.Date()
		sqlParams = append(sqlParams, fmt.Sprintf("%d-%d-%d", year, month, day))
		statementCounter++
	}
	if params.CoinholdID > 0 {
		filterQuery += fmt.Sprintf(" AND ac.id = $%d", statementCounter)
		sqlParams = append(sqlParams, params.CoinholdID)
	}

	return filterQuery, sqlParams, nil
}

func (r *walletsHistory) buildLimits(input *model.HistoryInput, params []any) (string, []any) {
	statementCounter := len(params) + 1
	filterQuery := fmt.Sprintf(" LIMIT $%d OFFSET $%d", statementCounter, statementCounter+1)
	params = append(params, input.Limit, input.Offset)

	return filterQuery, params
}

// // setReferralType - is not used
// func (r *walletsHistory) setReferralType(item *model.Wallet) {
//	if item.Type == int(model.PoolPaysUsersReferralsTrTypeID) || item.Type == int(model.PoolPaysBenefitOtherUserTrTypeID) {
//		t := int(enum.MiningAccountTypeID)
//		item.ReferralType = &t
//	} else if item.Type == int(model.CnhldRefBonusTrTypeID) || item.Type == int(model.PromoAuthorBonusTrTypeID) {
//		t := int(enum.CoinholdAccountTypeID)
//		item.ReferralType = &t
//	}
// }

func getStatus(tx *model.Wallet) int {
	if _, ok := p2pTxTypeIDs[tx.Type]; ok {
		return tx.P2PStatus
	}

	if tx.ExchangeIsSuccess != nil && *tx.ExchangeIsSuccess {
		return txStatusDone
	}

	if _, ok := lateStatusDoneTxTypeIDs[tx.Type]; !ok {
		return txStatusDone
	}

	if tx.TxID == nil {
		return txStatusInProgress
	}

	if tx.TxID != nil && *tx.TxID != "" {
		if !strings.Contains(*tx.TxID, "err") {
			return txStatusDone
		}

		return txStatusDeclined
	}

	return txStatusInProgress
}
