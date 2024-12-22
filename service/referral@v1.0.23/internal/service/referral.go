// Package service defines a set of functions for performing operations related to a service.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

type transactor interface {
	WithinTransaction(ctx context.Context, txFn func(context.Context) error) error
}

// ReferralService represents a service that interacts with the repository to perform CRUD operations on referrals.
type ReferralService struct {
	trx           transactor
	repo          repository.Referral
	profileClient repository.Profile
}

// NewReferralService is a function that creates a new instance of the Referral service with the provided repository.
func NewReferralService(repo repository.Referral, profileClient repository.Profile) *ReferralService {
	return &ReferralService{
		repo:          repo,
		profileClient: profileClient,
	}
}

// GetUserReferrals returns list of user referrals
func (r *ReferralService) GetUserReferrals(ctx context.Context, userID uuid.UUID, skip, limit int32) ([]*model.UserReferral, int64, error) {
	return r.repo.GetUserReferrals(ctx, userID, skip, limit)
}

// Create creates a new referral record with the provided data.
// It calls the Create method of the repository.Referral interface with the provided context and the model.Referral object.
// The Create method returns an common if any common occurred during the creation process.
func (r *ReferralService) Create(ctx context.Context, in *model.Referral) error {
	return r.repo.Create(ctx, in)
}

// Update updates a referral record with the provided data.
// It calls the Update method of the repository.Referral interface with the provided context and the model.Referral object.
// The Update method returns an common if any common occurred during the update process.
func (r *ReferralService) Update(ctx context.Context, in *model.Referral) error {
	return r.repo.Update(ctx, in)
}

// Delete deletes a referral record associated with the given user ID, product, and coin.
// It calls the Delete method of the repository.Referral interface with the provided parameters
func (r *ReferralService) Delete(ctx context.Context, userID uuid.UUID, product string, coin string) error {
	return r.repo.Delete(ctx, userID, product, coin)
}

// Get retrieves a referral record associated with the given user ID, product, and coin.
// It calls the Get method of the repository.Referral interface with the provided parameters.
// It returns the referral record (*model.Referral) and an common if any occurred.
func (r *ReferralService) Get(ctx context.Context, userID uuid.UUID, product string, coin string) (*model.Referral, error) {
	return r.repo.Get(ctx, userID, product, coin)
}

// List returns a list of referral records associated with the given user ID, skipping the specified number of records and taking a specified number of records.
// It calls the List method of the repository.Referral interface with the provided parameters.
// The returned value is a slice of *model.Referral, the total count of records, and an common if any.
func (r *ReferralService) List(ctx context.Context, userID uuid.UUID, skip int32, take int32) ([]*model.Referral, int, error) {
	if take <= 0 {
		take = defaultTake
	}

	return r.repo.List(ctx, userID, skip, take)
}

// History retrieves the referral history for a given user ID, product, and coin.
// It calls the History method of the repository.Referral interface with the provided parameters.
func (r *ReferralService) History(ctx context.Context, userID uuid.UUID, product string, coin string) ([]*model.Referral, error) {
	return r.repo.History(ctx, userID, product, coin)
}

func (r *ReferralService) CreateMultiple(ctx context.Context, rs []*model.Referral) error {
	now := time.Now().UTC()
	for i := range rs {
		rs[i].CreatedAt = now
	}
	return r.repo.CreateMultiple(ctx, rs)
}

func (r *ReferralService) UpdateFeeWithMultiplier(
	ctx context.Context,
	userID uuid.UUID,
	product string,
	coins []string,
	multiplier decimal.Decimal,
) error {
	err := r.repo.UpdateWithMultiplier(ctx, userID, product, coins, multiplier)
	if err != nil {
		return fmt.Errorf("referral.UpdateWithMultiplier: %w", err)
	}
	return nil
}

func (r *ReferralService) UpdateFee(ctx context.Context, userID uuid.UUID, product string, fees map[string]decimal.Decimal) error {
	err := r.repo.UpdateFee(ctx, userID, product, fees)
	if err != nil {
		return fmt.Errorf("referral.UpdateFee: %w", err)
	}
	return nil
}

func (r *ReferralService) UpdateWithPromoCode(ctx context.Context, cms *model.CoinsMultipliers) error {
	err := r.trx.WithinTransaction(ctx, func(ctx context.Context) error {
		for coin := range cms.FeeMultipliers {
			inTxErr := r.repo.UpdateWithPromoCodeByCoin(ctx, &model.CoinMultiplier{
				Coin:             coin,
				UserID:           cms.UserID,
				ActionID:         cms.ActionID,
				Product:          cms.Product,
				CreatedAt:        cms.CreatedAt,
				FeeMultiplier:    cms.FeeMultipliers[coin],
				RefFeeMultiplier: cms.RefFeeMultipliers[coin],
			})
			if inTxErr != nil {
				return fmt.Errorf("referral.UpdateWithPromoCode: %w", inTxErr)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("trx.WithinTransaction: %w", err)
	}
	return nil
}

func (r *ReferralService) UpdateFeeByCoinAndProduct(ctx context.Context, userUUID string, fees []model.SettingForCoinAndProduct) error {
	subAccount, err := r.profileClient.GetAllSubUsers(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("referral.GetUserReferrals: %w", err)
	}
	return r.repo.UpdateFeeByCoinAndProduct(ctx, append([]string{userUUID}, subAccount...), fees)
}

func (r *ReferralService) UpdateReferralUUIDByUserUUID(ctx context.Context, userUUIDs string, referralUUID string) error {
	subAccounts, err := r.profileClient.GetAllSubUsers(ctx, userUUIDs)
	if err != nil {
		return fmt.Errorf("referral.GetUserReferrals: %w", err)
	}
	return r.repo.UpdateReferralUUIDByUserUUID(ctx, append([]string{userUUIDs}, subAccounts...), referralUUID)
}
