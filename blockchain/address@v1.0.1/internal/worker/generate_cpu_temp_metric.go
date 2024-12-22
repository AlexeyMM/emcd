package worker

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	sdkLog "code.emcdtech.com/emcd/sdk/log"

	"github.com/prometheus/client_golang/prometheus"

	"code.emcdtech.com/emcd/sdk/app"
)

const cpuFull = 100
const cpuLow = 5

func NewGenerateCpuTempMetric(cpuTemp prometheus.Gauge) app.Worker {

	return app.WorkerFn(func(ctx context.Context) error {
		for ctx.Err() == nil {
			if r, err := rand.Int(rand.Reader, big.NewInt(cpuFull-cpuLow)); err != nil {
				sdkLog.Error(ctx, "rand generate: %v", err)

			} else {
				cpuTemp.Set(float64(r.Int64() + cpuLow))
				time.Sleep(time.Second)

			}
		}

		return nil
	})
}
