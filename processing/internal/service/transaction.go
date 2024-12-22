package service

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
)

type Transaction interface {
	ProcessTransaction(ctx context.Context, transaction *model.Transaction) error
}
