package worker

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/emcd/sdk/app"
)

const sleepHdFailureInterval = time.Second * 5

func NewGenerateHDFailuresMetric(hdFailures *prometheus.CounterVec) app.Worker {
	return app.WorkerFn(func(ctx context.Context) error {
		for ctx.Err() == nil {
			hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
			time.Sleep(sleepHdFailureInterval)
		}

		return nil
	})
}
