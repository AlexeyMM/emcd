package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type AppMetrics struct {
	RepositoryRequestTimeHistogram *prometheus.HistogramVec
	InvoiceStatusesGauge           *prometheus.GaugeVec
	InvoiceExecutionHistogram      *prometheus.HistogramVec
}

func New() *AppMetrics {
	return &AppMetrics{
		RepositoryRequestTimeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "repository_request_time_seconds",
			Help:    "Measure a request time of repository",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "success"}),
		InvoiceStatusesGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "invoice_statuses",
			Help: "Measure invoice statuses",
		}, []string{"status"}),
		InvoiceExecutionHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "invoice_execution_time_seconds",
			Help: "Measure a invoice execution time",
			Buckets: []float64{
				60, 120, 180, 240, 300, 360, 420, 480, 540, 600,
			},
		}, []string{"finished_status"}),
	}
}

func (m *AppMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.RepositoryRequestTimeHistogram.Describe(ch)
	m.InvoiceStatusesGauge.Describe(ch)
	m.InvoiceExecutionHistogram.Describe(ch)
}

func (m *AppMetrics) Collect(ch chan<- prometheus.Metric) {
	m.RepositoryRequestTimeHistogram.Collect(ch)
	m.InvoiceStatusesGauge.Collect(ch)
	m.InvoiceExecutionHistogram.Collect(ch)
}
