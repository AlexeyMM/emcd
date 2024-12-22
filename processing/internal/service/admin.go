package service

import (
	"context"

	"code.emcdtech.com/b2b/processing/model"
)

type MerchantAdminService interface {
	CreateMerchant(ctx context.Context, m *model.Merchant) error
}
