package worker

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/b2b/swap/internal/client"
	"code.emcdtech.com/b2b/swap/internal/repository"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	subAccountsCountingInterval  = time.Hour
	swapStatusesCountingInterval = time.Minute
	ctxTimeoutInterval           = 30 * time.Second
)

type Metrics struct {
	swapRepo           repository.Swap
	exchangeAccountCli client.ExchangeAccount
	subAccountsGauge   prometheus.Gauge
	swapStatusesGauge  *prometheus.GaugeVec
}

func NewMetrics(swapRepo repository.Swap, exchangeAccountCli client.ExchangeAccount, subAccountsGauge prometheus.Gauge, swapStatusesGauge *prometheus.GaugeVec) *Metrics {
	return &Metrics{
		swapRepo:           swapRepo,
		exchangeAccountCli: exchangeAccountCli,
		subAccountsGauge:   subAccountsGauge,
		swapStatusesGauge:  swapStatusesGauge,
	}
}

func (m *Metrics) Run(ctx context.Context) error {
	log.Debug(ctx, "subAccountCounter run")
	defer log.Debug(ctx, "subAccountCounter stopped")

	subAccountCountingTicker := time.NewTicker(subAccountsCountingInterval)
	defer subAccountCountingTicker.Stop()

	swapStatusesCountingTicker := time.NewTicker(swapStatusesCountingInterval)
	defer swapStatusesCountingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-subAccountCountingTicker.C:
			newCtx, cancel := context.WithTimeout(ctx, ctxTimeoutInterval)
			err := m.subAccountsHandle(newCtx)
			if err != nil {
				log.Error(ctx, "worker.metrics: subAccountsHandle: %s", err.Error())
				cancel()
				continue
			}
			cancel()

		case <-swapStatusesCountingTicker.C:
			newCtx, cancel := context.WithTimeout(ctx, ctxTimeoutInterval)
			err := m.swapStatusesHandle(newCtx)
			if err != nil {
				log.Error(ctx, "worker.metrics: swapStatusesHandle: %s", err.Error())
				cancel()
				continue
			}
			cancel()
		}
	}
}

func (m *Metrics) subAccountsHandle(ctx context.Context) error {
	accs, err := m.exchangeAccountCli.GetSubAccounts(ctx)
	if err != nil {
		return fmt.Errorf("getSubAccounts: %w", err)
	}

	log.Debug(ctx, "worker.metrics: count of sub accounts: %v", len(accs))
	m.subAccountsGauge.Set(float64(len(accs)))

	return nil
}

func (m *Metrics) swapStatusesHandle(ctx context.Context) error {
	statuses, err := m.swapRepo.CountSwapsByStatus(ctx)
	if err != nil {
		return fmt.Errorf("getCountSwapsByStatus: %w", err)
	}

	for status, count := range statuses {
		log.Debug(ctx, "worker.metrics: counting swap statuses: status: %s, count: %d", status.String(), count)
		m.swapStatusesGauge.WithLabelValues(status.String()).Set(float64(count))
	}

	return nil
}
