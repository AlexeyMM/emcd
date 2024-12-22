package metrics

import (
	"context"
	"strconv"
	"time"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/model"
)

type DepositAddressPool struct {
	repo      repository.DepositAddressPool
	histogram *prometheus.HistogramVec
}

func NewDepositAddressPool(repo repository.DepositAddressPool, histogram *prometheus.HistogramVec) *DepositAddressPool {
	return &DepositAddressPool{
		repo:      repo,
		histogram: histogram,
	}
}

func (d *DepositAddressPool) WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) error {
	start := time.Now()
	err := d.repo.WithinTransaction(ctx, txFn)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues(
		"depositAddressPool.WithinTransaction",
		strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (d *DepositAddressPool) WithinTransactionWithOptions(
	ctx context.Context,
	txFn func(ctx context.Context) error,
	opts pgx.TxOptions,
) error {
	start := time.Now()
	err := d.repo.WithinTransactionWithOptions(ctx, txFn, opts)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues(
		"depositAddressPool.WithinTransactionWithOptions",
		strconv.FormatBool(err == nil)).Observe(duration)

	return err
}

func (d *DepositAddressPool) Runner(ctx context.Context) transactor.PgxQueryRunner {
	return d.repo.Runner(ctx)
}

func (d *DepositAddressPool) OccupyAddress(
	ctx context.Context,
	merchantID uuid.UUID,
	networkID string,
) (*model.Address, error) {
	start := time.Now()
	addr, err := d.repo.OccupyAddress(ctx, merchantID, networkID)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues(
		"depositAddressPool.OccupyAddress",
		strconv.FormatBool(err == nil)).Observe(duration)

	return addr, err
}

func (d *DepositAddressPool) Save(ctx context.Context, address *model.Address) error {
	start := time.Now()
	err := d.repo.Save(ctx, address)
	duration := time.Since(start).Seconds()

	d.histogram.WithLabelValues("depositAddressPool.Save", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}
