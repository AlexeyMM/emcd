package metrics

import (
	"context"
	"strconv"
	"time"

	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
)

type Swap struct {
	swap      repository.Swap
	histogram *prometheus.HistogramVec
}

func NewSwap(swap repository.Swap, histogram *prometheus.HistogramVec) *Swap {
	return &Swap{
		swap:      swap,
		histogram: histogram,
	}
}

func (s *Swap) Add(ctx context.Context, swap *model.Swap) error {
	start := time.Now()
	err := s.swap.Add(ctx, swap)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.Add", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (s *Swap) Find(ctx context.Context, filter *model.SwapFilter) (model.Swaps, error) {
	start := time.Now()
	sw, err := s.swap.Find(ctx, filter)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.Find", strconv.FormatBool(err == nil)).Observe(duration)
	return sw, err
}

func (s *Swap) FindOne(ctx context.Context, filter *model.SwapFilter) (*model.Swap, error) {
	start := time.Now()
	sw, err := s.swap.FindOne(ctx, filter)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return sw, err
}

func (s *Swap) Update(ctx context.Context, swap *model.Swap, filter *model.SwapFilter, partial *model.SwapPartial) error {
	start := time.Now()
	err := s.swap.Update(ctx, swap, filter, partial)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.Update", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (s *Swap) CountSwapsByStatus(ctx context.Context) (map[model.Status]int, error) {
	start := time.Now()
	statuses, err := s.swap.CountSwapsByStatus(ctx)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.CountSwapsByStatus", strconv.FormatBool(err == nil)).Observe(duration)
	return statuses, err
}

func (s *Swap) WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) error {
	start := time.Now()
	err := s.swap.WithinTransaction(ctx, txFn)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.WithinTransaction", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (s *Swap) WithinTransactionWithOptions(ctx context.Context, txFn func(ctx context.Context) error, opts pgx.TxOptions) error {
	start := time.Now()
	err := s.swap.WithinTransactionWithOptions(ctx, txFn, opts)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.WithinTransactionWithOptions", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (s *Swap) Runner(ctx context.Context) transactor.PgxQueryRunner {
	return s.swap.Runner(ctx)
}

func (s *Swap) CountTotalWithFilter(ctx context.Context, filter *model.SwapFilter) (int, error) {
	start := time.Now()
	count, err := s.swap.CountTotalWithFilter(ctx, filter)
	duration := time.Since(start).Seconds()

	s.histogram.WithLabelValues("swap.CountTotalWithFilter", strconv.FormatBool(err == nil)).Observe(duration)
	return count, err
}
