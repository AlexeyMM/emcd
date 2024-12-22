package admin

import (
	"context"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/model"
)

type MerchantService struct {
	repo repository.MerchantAdmin
}

func NewMerchantService(repo repository.MerchantAdmin) *MerchantService {
	return &MerchantService{repo: repo}
}

func (s *MerchantService) CreateMerchant(ctx context.Context, m *model.Merchant) error {
	return s.repo.SaveMerchant(ctx, m)
}
