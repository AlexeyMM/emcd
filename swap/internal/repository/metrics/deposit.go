package metrics

import (
	"context"
	"strconv"
	"time"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type Deposit struct {
	deposit   repository.Deposit
	histogram *prometheus.HistogramVec
}

func NewDeposit(deposit repository.Deposit, histogram *prometheus.HistogramVec) *Deposit {
	return &Deposit{
		deposit:   deposit,
		histogram: histogram,
	}
}

func (d *Deposit) Add(ctx context.Context, deposit *model.Deposit) error {
	start := time.Now()
	err := d.deposit.Add(ctx, deposit)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues("deposit.Add", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (d *Deposit) Find(ctx context.Context, filter *model.DepositFilter) (model.Deposits, error) {
	start := time.Now()
	deps, err := d.deposit.Find(ctx, filter)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues("deposit.Find", strconv.FormatBool(err == nil)).Observe(duration)
	return deps, err
}

func (d *Deposit) FindOne(ctx context.Context, filter *model.DepositFilter) (*model.Deposit, error) {
	start := time.Now()
	deps, err := d.deposit.FindOne(ctx, filter)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues("deposit.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return deps, err
}

func (d *Deposit) WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) error {
	start := time.Now()
	err := d.deposit.WithinTransaction(ctx, txFn)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues("deposit.WithinTransaction", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (d *Deposit) WithinTransactionWithOptions(ctx context.Context, txFn func(ctx context.Context) error, opts pgx.TxOptions) error {
	start := time.Now()
	err := d.deposit.WithinTransactionWithOptions(ctx, txFn, opts)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues("deposit.WithinTransactionWithOptions", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (d *Deposit) Runner(ctx context.Context) transactor.PgxQueryRunner {
	return d.deposit.Runner(ctx)
}
