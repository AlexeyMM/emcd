package repository

import (
	"context"

	"code.emcdtech.com/b2b/endpoint/internal/model"
)

type IP interface {
	Add(ctx context.Context, ips []*model.IP) error
	Find(ctx context.Context, filter *model.IPFilter) ([]*model.IP, error)
	Delete(ctx context.Context, filter *model.IPFilter) error
}
