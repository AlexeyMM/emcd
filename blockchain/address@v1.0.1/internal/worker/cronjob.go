package worker

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	sdkLog "code.emcdtech.com/emcd/sdk/log"
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

	tracer := otel.Tracer("cron-worker")
	ctx, span := tracer.Start(ctx, "cron-job-execute",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	span.SetAttributes(attribute.String("transaction.type", "background"))
	defer span.End()

	err := c.task(ctx)
	if err != nil {
		sdkLog.Error(ctx, "execute task in cron worker: %s", err.Error())
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("event", "error"),
			attribute.String("message", err.Error()),
		)
		span.SetStatus(codes.Error, err.Error())
	}
}

func (c *CronWorker) asyncSleep(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, c.interval)
	defer cancel()
	<-ctx.Done()
}
