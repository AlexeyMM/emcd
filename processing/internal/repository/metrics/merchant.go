package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/model"
)

type Merchant struct {
	repo      repository.Merchant
	histogram *prometheus.HistogramVec
}

func NewMerchant(repo repository.Merchant, histogram *prometheus.HistogramVec) *Merchant {
	return &Merchant{
		repo:      repo,
		histogram: histogram,
	}
}

func (m *Merchant) Get(ctx context.Context, id uuid.UUID) (*model.Merchant, error) {
	start := time.Now()
	merchant, err := m.repo.Get(ctx, id)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("merchant.Get", strconv.FormatBool(err == nil)).Observe(duration)

	return merchant, err
}
