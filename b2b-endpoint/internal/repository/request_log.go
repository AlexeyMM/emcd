package repository

import (
	"context"

	"code.emcdtech.com/b2b/endpoint/internal/model"
)

type RequestLog interface {
	Add(ctx context.Context, log *model.RequestLog) error
	Find(ctx context.Context, filter *model.RequestLogFilter) ([]*model.RequestLog, error)
	FindOne(ctx context.Context, filter *model.RequestLogFilter) (*model.RequestLog, error)
}
