package metrics

import (
	"context"
	"strconv"
	"time"

	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/prometheus/client_golang/prometheus"
)

type InternalTransfer struct {
	transfer  repository.Transfer
	histogram *prometheus.HistogramVec
}

func NewInternalTransfer(transfer repository.Transfer, histogram *prometheus.HistogramVec) *InternalTransfer {
	return &InternalTransfer{
		transfer:  transfer,
		histogram: histogram,
	}
}

func (i *InternalTransfer) Add(ctx context.Context, transfer *model.InternalTransfer) error {
	start := time.Now()
	err := i.transfer.Add(ctx, transfer)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("transfer.Add", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}

func (i *InternalTransfer) Find(ctx context.Context, filter *model.InternalTransferFilter) (model.InternalTransfers, error) {
	start := time.Now()
	tr, err := i.transfer.Find(ctx, filter)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("transfer.Find", strconv.FormatBool(err == nil)).Observe(duration)
	return tr, err
}

func (i *InternalTransfer) FindOne(ctx context.Context, filter *model.InternalTransferFilter) (*model.InternalTransfer, error) {
	start := time.Now()
	tr, err := i.transfer.FindOne(ctx, filter)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("transfer.FindOne", strconv.FormatBool(err == nil)).Observe(duration)
	return tr, err
}

func (i *InternalTransfer) Update(ctx context.Context, internalTransfer *model.InternalTransfer, filter *model.InternalTransferFilter, partial *model.InternalTransferPartial) error {
	start := time.Now()
	err := i.transfer.Update(ctx, internalTransfer, filter, partial)
	duration := time.Since(start).Seconds()

	i.histogram.WithLabelValues("transfer.Update", strconv.FormatBool(err == nil)).Observe(duration)
	return err
}
