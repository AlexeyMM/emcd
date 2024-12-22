package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type AppMetrics struct {
	ByBitOrderBookWebsocketGauge   prometheus.Gauge
	ByBitReconnectWebsocketCounter *prometheus.CounterVec // Обрывы соединения
	ByBitRequestTimeHistogram      *prometheus.HistogramVec
	RepositoryRequestTimeHistogram *prometheus.HistogramVec
	SubAccountsGauge               prometheus.Gauge
	SwapStatusesGauge              *prometheus.GaugeVec
}

func New() *AppMetrics {
	m := &AppMetrics{
		ByBitOrderBookWebsocketGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "bybit_orderbook_websocket",
			Help: "Measure a count of websocket connections for listen orderbook",
		}),
		// TODO добавить в графану. Не добавил, потому что не было событий
		ByBitReconnectWebsocketCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "bybit_reconnect_websocket_total",
			Help: "Measure websocket abnormal connection count with reconnection",
		}, []string{"websocket"}),
		ByBitRequestTimeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "bybit_request_time_seconds",
			Help:        "Measure a request time of byBit",
			ConstLabels: nil,
			Buckets:     prometheus.DefBuckets,
		}, []string{"method", "success"}),
		RepositoryRequestTimeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "repository_request_time_seconds",
			Help:    "Measure a request time of repository",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "success"}),
		SubAccountsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "sub_accounts",
			Help: "Measure sub accounts",
		}),
		SwapStatusesGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "swap_statuses",
			Help: "Measure swap statuses",
		}, []string{"status"}),
	}
	return m
}

func (m *AppMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.ByBitOrderBookWebsocketGauge.Describe(ch)
	m.ByBitReconnectWebsocketCounter.Describe(ch)
	m.ByBitRequestTimeHistogram.Describe(ch)
	m.RepositoryRequestTimeHistogram.Describe(ch)
	m.SubAccountsGauge.Describe(ch)
	m.SwapStatusesGauge.Describe(ch)
}

func (m *AppMetrics) Collect(ch chan<- prometheus.Metric) {
	m.ByBitOrderBookWebsocketGauge.Collect(ch)
	m.ByBitReconnectWebsocketCounter.Collect(ch)
	m.ByBitRequestTimeHistogram.Collect(ch)
	m.RepositoryRequestTimeHistogram.Collect(ch)
	m.SubAccountsGauge.Collect(ch)
	m.SwapStatusesGauge.Collect(ch)
}
