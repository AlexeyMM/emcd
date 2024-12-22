package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type AppMetrics struct {
	CpuTemp    prometheus.Gauge
	HDFailures *prometheus.CounterVec
}

func NewAppMetrics() *AppMetrics {
	cpuTemp := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        "cpu_temperature_celsius",
		Help:        "Current temperature of the CPU.",
		ConstLabels: nil,
	})

	hdFailures := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "hd_errors_total",
			Help:        "Number of hard-disk errors.",
			ConstLabels: nil,
		},
		[]string{"device"},
	)

	return &AppMetrics{
		CpuTemp:    cpuTemp,
		HDFailures: hdFailures,
	}
}

func (m *AppMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.CpuTemp.Describe(ch)
	m.HDFailures.Describe(ch)
}

func (m *AppMetrics) Collect(ch chan<- prometheus.Metric) {
	m.CpuTemp.Collect(ch)
	m.HDFailures.Collect(ch)
}
