package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/model"
)

type MerchantAdmin struct {
	repo      repository.MerchantAdmin
	histogram *prometheus.HistogramVec
}

func NewMerchantAdmin(repo repository.MerchantAdmin, histogram *prometheus.HistogramVec) *MerchantAdmin {
	return &MerchantAdmin{
		repo:      repo,
		histogram: histogram,
	}
}

func (m *MerchantAdmin) SaveMerchant(ctx context.Context, merchant *model.Merchant) error {
	start := time.Now()
	err := m.repo.SaveMerchant(ctx, merchant)
	duration := time.Since(start).Seconds()

	m.histogram.WithLabelValues("merchantAdmin.SaveMerchant", strconv.FormatBool(err == nil)).Observe(duration)

	return err
}
