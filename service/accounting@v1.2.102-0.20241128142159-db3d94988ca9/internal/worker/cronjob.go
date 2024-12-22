package worker

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"context"
	"time"
)

type TaskCronWorker func(ctx context.Context) error

type CronWorker struct {
	task     TaskCronWorker
	interval time.Duration
}

func NewCronWorker(
	interval time.Duration,
	task TaskCronWorker,
) *CronWorker {
	return &CronWorker{
		task:     task,
		interval: interval,
	}
}

func (c *CronWorker) Run(ctx context.Context) error {
	for ctx.Err() == nil {
		c.cronExecute(ctx)
		c.asyncSleep(ctx)
	}
	return nil
}

func (c *CronWorker) cronExecute(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err := c.task(ctx)
	if err != nil {
		sdkLog.Error(ctx, "execute task in cron worker: %v", err)
	}
}

func (c *CronWorker) asyncSleep(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, c.interval)
	defer cancel()
	<-ctx.Done()
}
