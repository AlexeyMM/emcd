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

type Withdraw struct {
	withdraw  repository.Withdraw
	histogram *prometheus.HistogramVec
}

func NewWithdraw(withdraw repository.Withdraw, histogram *prometheus.HistogramVec) *Withdraw {
	return &Withdraw{
		withdraw:  withdraw,
		histogram: histogram,
	}
}

func (w *Withdraw) Add(ctx context.Context, withdraw *model.Withdraw) error {
	start := time.Now()
	err := w.withdraw.Add(ctx, withdraw)
	duration := time.Since(start).Seconds()

	w.histogram.WithLabelValues("withdraw.Add", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (w *Withdraw) Find(ctx context.Context, filter *model.WithdrawFilter) (model.Withdraws, error) {
	start := time.Now()
	wt, err := w.withdraw.Find(ctx, filter)
	duration := time.Since(start).Seconds()

	w.histogram.WithLabelValues("withdraw.Find", strconv.FormatBool(err == nil)).Observe(duration)
	return wt, err
}

func (w *Withdraw) FindOne(ctx context.Context, filter *model.WithdrawFilter) (*model.Withdraw, error) {
	start := time.Now()
	wt, err := w.withdraw.FindOne(ctx, filter)
	duration := time.Since(start).Seconds()

	w.histogram.WithLabelValues("withdraw.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return wt, err
}

func (w *Withdraw) Update(ctx context.Context, withdraw *model.Withdraw, filter *model.WithdrawFilter, partial *model.WithdrawPartial) error {
	start := time.Now()
	err := w.withdraw.Update(ctx, withdraw, filter, partial)
	duration := time.Since(start).Seconds()

	w.histogram.WithLabelValues("withdraw.Update", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (w *Withdraw) WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) error {
	start := time.Now()
	err := w.withdraw.WithinTransaction(ctx, txFn)
	duration := time.Since(start).Seconds()

	w.histogram.WithLabelValues("withdraw.WithinTransaction", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (w *Withdraw) WithinTransactionWithOptions(ctx context.Context, txFn func(ctx context.Context) error, opts pgx.TxOptions) error {
	start := time.Now()
	err := w.withdraw.WithinTransactionWithOptions(ctx, txFn, opts)
	duration := time.Since(start).Seconds()

	w.histogram.WithLabelValues("withdraw.WithinTransactionWithOptions", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (w *Withdraw) Runner(ctx context.Context) transactor.PgxQueryRunner {
	return w.withdraw.Runner(ctx)
}
