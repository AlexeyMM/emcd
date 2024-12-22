package repository

import (
	"context"

	"code.emcdtech.com/b2b/swap/model"
)

// Account отвечает за хранение суб-аккаунтов
type Account interface {
	Add(ctx context.Context, account *model.Account) error
	Find(ctx context.Context, filter *model.AccountFilter) (model.Accounts, error)
	FindOne(ctx context.Context, filter *model.AccountFilter) (*model.Account, error)
}
