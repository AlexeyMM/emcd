package repository

import (
	"context"

	"code.emcdtech.com/b2b/endpoint/internal/model"
)

type Client interface {
	Add(ctx context.Context, client *model.Client) error
}
