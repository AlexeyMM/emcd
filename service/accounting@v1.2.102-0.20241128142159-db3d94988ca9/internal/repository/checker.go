package repository

import (
	"code.emcdtech.com/emcd/service/accounting/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"time"
)

type CheckerRepository interface {
	GetIncomesHashrateByDate(ctx context.Context, date time.Time) ([]*model.HashrateByDate, error)
	GetCoinsOperationsSum(ctx context.Context) ([]*model.SumCheckData, error)
	GetTransactionOperationsIntegrity(ctx context.Context) ([]*model.TransactionOperationsIntegrityData, error)
	GetCheckTransactionCoins(ctx context.Context) ([]int64, error)
	GetCheckOperationsCoins(ctx context.Context) ([]int64, error)
	GetCheckFreezePayoutsBlocks(ctx context.Context) ([]*model.CheckFreezePayoutsBlocksData, error)
}

type checker struct {
	pool *pgxpool.Pool
}

func NewChecker(pool *pgxpool.Pool) CheckerRepository {
	return &checker{
		pool: pool,
	}
}

func (r *checker) GetIncomesHashrateByDate(ctx context.Context, date time.Time) ([]*model.HashrateByDate, error) {
	query := `select sum(t.hashrate) as hashrate, t.coin_id as coincode
--, date(t.created_at) as date 
from emcd.transactions t
where t.type = 5 and date(t.created_at) = date($1)
GROUP BY t.coin_id, date(t.created_at)`
	rows, err := r.pool.Query(ctx, query, date)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.HashrateByDate, 0)

	for rows.Next() {
		var (
			p        model.HashrateByDate
			coincode sql.NullInt64
			hashrate string
		)

		err = rows.Scan(
			&hashrate,
			&coincode,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan hashrate: %w", err)
		}

		p.CoinId = coincode.Int64
		p.Hashrate = hashrate

		result = append(result, &p)
	}

	return result, nil
}

func (r *checker) GetCoinsOperationsSum(ctx context.Context) ([]*model.SumCheckData, error) {

	query := `select round(sum(o.amount),8) as sum, o.coin_id as coincode from emcd.operations o GROUP BY o.coin_id`

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.SumCheckData, 0)

	for rows.Next() {
		var (
			p        model.SumCheckData
			coincode sql.NullInt64
			sum      decimal.Decimal
		)

		err = rows.Scan(
			&sum,
			&coincode,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		p.CoinId = coincode.Int64
		p.Sum = sum

		result = append(result, &p)
	}

	return result, nil

}

func (r *checker) GetTransactionOperationsIntegrity(ctx context.Context) ([]*model.TransactionOperationsIntegrityData, error) {
	query := `
		select count, 
		       tr_id as trId, 
		       op1_id as op1Id, 
		       op2_id as op2Id, 
		       op_pair_chk as opPairCheck, 
		       tr_neg_chk as trNegCheck, 
		       op_sum_chk as opSumChk, 
		       diff_chk as diffChk, 
		       tr_date_chk as trDateChk, 
		       coin_chk as coinChk, 
		       acc_chk as accChk
		from
  --Группируем транзакции и операции по связям и проверяем контрольные данные
 (
     select count(*) as count,
            t.id as tr_id,
            max(o.id) as op2_id, -- id операции 2, если она есть, иначе будут равны
            min(o.id) as op1_id, -- id операции 1
            (max(o.id) != min(o.id)) as op_pair_chk, -- что операции 2 (id не совпадают)
            (t.amount >= 0) as tr_neg_chk, -- проверка транзакции на отрицательное значение
            (sum(o.amount) = 0) as op_sum_chk, -- сумма двух операций, == 0
            ((t.amount - max(o.amount)) = 0) as diff_chk, -- разница суммы транзакции и суммы положительной операции
            ((extract(epoch from (date(t.created_at)))*2) = (extract(epoch from date(max(o.created_at))) + extract(epoch from date(min(o.created_at))))) as tr_date_chk, -- дата транзакции\операции на совпадение
            ((sum(o.coin_id)/2) = t.coin_id) as coin_chk, -- контроль монеты операций\транзакции
            ((t.sender_account_id + t.receiver_account_id) = sum(o.account_id)) as acc_chk -- контрольная сумма аккаунтов
           from emcd.operations o
                    FULL OUTER JOIN emcd.transactions t ON o.transaction_id = t.id
           GROUP BY t.id
 ) chk

where
   tr_id IS NUll -- нет транзакции
   OR op1_id IS NULL -- нет операций
   OR op_pair_chk != true -- операции не парные
   OR tr_neg_chk != true -- отрицательная транзакция
   OR op_sum_chk != true -- сумма операций не равна 0
   OR diff_chk != true -- сумма операции не сходится с транзакцией
   OR tr_date_chk != true -- дата (день) операций\транзакции не сходятся
   OR coin_chk != true -- монета операций\транзакции не сходятся
   OR acc_chk != true -- аккаунты в транзакции\операциях не сходятся;
	`
	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.TransactionOperationsIntegrityData, 0)

	for rows.Next() {
		var (
			p     model.TransactionOperationsIntegrityData
			count sql.NullInt64
			trId  sql.NullInt64
			op2Id sql.NullInt64
			op1Id sql.NullInt64
		)

		err := rows.Scan(&count, &trId, &op2Id, &op1Id,
			&p.OpPairCheck,
			&p.TrNegChk,
			&p.OpSumChk,
			&p.DiffChk,
			&p.TrDateChk,
			&p.CoinChk,
			&p.AccChk)

		if err != nil {
			return nil, fmt.Errorf("row.Scan: %w", err)
		}

		p.Count = count.Int64
		p.TrId = trId.Int64
		p.Op2Id = op2Id.Int64
		p.Op1Id = op1Id.Int64

		result = append(result, &p)
	}

	return result, nil
}

func (r *checker) GetCheckTransactionCoins(ctx context.Context) ([]int64, error) {
	query := `
SELECT t.id as trId--, 
       --t.coin_id as t_coin, 
       --uar.coin_id as uar_coin, 
       --uas.coin_id as uas_coin
FROM emcd.transactions t
         left join emcd.users_accounts uas ON t.sender_account_id = uas.id
         left join emcd.users_accounts uar ON t.receiver_account_id = uar.id
WHERE uar.coin_id != t.coin_id OR uas.coin_id != t.coin_id
`
	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	result := make([]int64, 0)

	for rows.Next() {
		var (
			trId int64
		)

		err = rows.Scan(
			&trId,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		result = append(result, trId)
	}

	return result, nil

}

func (r *checker) GetCheckOperationsCoins(ctx context.Context) ([]int64, error) {
	query := `
SELECT o.id as opId--, 
       --o.coin_id as o_coin, 
       -- ua.coin_id as ua_coin
FROM emcd.operations o
        left join emcd.users_accounts ua ON o.account_id = ua.id
WHERE ua.coin_id != o.coin_id
`
	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	result := make([]int64, 0)

	for rows.Next() {
		var (
			opId int64
		)

		err = rows.Scan(
			&opId,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		result = append(result, opId)
	}

	return result, nil

}

func (r *checker) GetCheckFreezePayoutsBlocks(ctx context.Context) ([]*model.CheckFreezePayoutsBlocksData, error) {

	query := `SELECT t.id, t.type as type, u.id as user_id, t.created_at as created_at FROM emcd.transactions_blocks bt
LEFT JOIN emcd.transactions t ON bt.block_transaction_id = t.id
LEFT JOIN emcd.users_accounts ua ON t.sender_account_id = ua.id
LEFT JOIN emcd.users u ON ua.user_id = u.id
WHERE bt.unblock_transaction_id is null 
  and t.receiver_address is not null 
  and t.type IN ($1,$2)
  and u.nopay = false 
  and t.created_at <= (now()  - INTERVAL '1 day');`

	rows, err := r.pool.Query(ctx, query, model.WalletMiningTransferTrTypeID, model.PoolPaysUsersBalanceTrTypeID)

	if err != nil {
		return nil, fmt.Errorf("can't select: %w", err)
	}

	defer rows.Close()

	result := make([]*model.CheckFreezePayoutsBlocksData, 0)

	for rows.Next() {
		var (
			p model.CheckFreezePayoutsBlocksData
		)

		err = rows.Scan(
			&p.TrId,
			&p.Type,
			&p.UserId,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan: %w", err)
		}

		result = append(result, &p)
	}

	return result, nil
}
