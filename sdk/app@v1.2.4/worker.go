package app

import (
	"context"
)

type Worker interface {
	Run(ctx context.Context) error
}

type WorkerFn func(ctx context.Context) error

func (w WorkerFn) Run(ctx context.Context) error {
	return w(ctx)
}

type Workers []Worker
