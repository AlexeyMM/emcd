package service

import (
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"context"
	"github.com/google/uuid"
)

// AccountingUserAccount deprecated, use UserAccountService
type AccountingUserAccount interface {
	// GetUserAccount deprecated, use UserAccountService.GetUserAccountById
	GetUserAccount(ctx context.Context, id int32) (*model.UserAccount, error)
	// GetUserAccounts deprecated, use UserAccountService.GetUserAccountsByUuid
	GetUserAccounts(ctx context.Context, userID uuid.UUID) (model.UserAccounts, error)
}

type AccountingUserAccountService struct {
	repo repository.UserAccountRepo
}

func NewAccountingUserAccountService(
	repo repository.UserAccountRepo,
) *AccountingUserAccountService {
	return &AccountingUserAccountService{
		repo: repo,
	}
}

// GetUserAccount deprecated, use UserAccountService.GetUserAccountById
func (s *AccountingUserAccountService) GetUserAccount(ctx context.Context, id int32) (*model.UserAccount, error) {

	return s.repo.GetUserAccountByIdLegacy(ctx, id)
}

// GetUserAccounts deprecated, use UserAccountService.GetUserAccountsByUuid
func (s *AccountingUserAccountService) GetUserAccounts(ctx context.Context, userID uuid.UUID) (model.UserAccounts, error) {

	return s.repo.FindUserAccountByUserIdLegacy(ctx, userID)
}
