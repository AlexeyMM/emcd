// Package service defines a set of functions for performing operations related to a service.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

const (
	normal = "normal"
	fee    = "fee"
	wlFee  = "wlFee"
	refFee = "referral"
)

// Reward is an interface that defines the Calculate method for calculating rewards based on the given parameters.
// RewardServer is a struct containing an implementation of the Reward interface.
type Reward interface {
	Calculate(
		ctx context.Context,
		userID uuid.UUID,
		product string,
		coin string,
		amount decimal.Decimal,
	) ([]*model.Transaction, error)
}

// reward is a struct that represents a reward object.
// It contains two fields: settingsRps of type repository.DefaultSettings and referralRps of type repository.Referral.
// The settingsRps is used to interact with the default settings repository, while the referralRps is used to interact with the referral repository.
type reward struct {
	referralRps                  repository.Referral
	profileRps                   repository.Profile
	defaultSettingsRps           repository.DefaultSettings
	defaultWhitelabelSettingsRps repository.DefaultWhitelabelSettings
}

// NewReward creates a new instance of Reward using the provided repositories for settings and referrals. It returns the created Reward.
// The settingsRps parameter is an implementation of the DefaultSettings interface, which provides methods to create, update, delete, and retrieve default settings.
// The referralRps parameter is an implementation of the Referral interface, which provides methods to create, update, delete, retrieve, list, and retrieve referral history.
// The function initializes a reward struct with the provided repositories and returns it as a Reward interface.
func NewReward(
	referralRps repository.Referral,
	profileRps repository.Profile,
	defaultSettingsRps repository.DefaultSettings,
	defaultWhitelabelSettingsRps repository.DefaultWhitelabelSettings,
) Reward {
	return &reward{
		referralRps:                  referralRps,
		profileRps:                   profileRps,
		defaultSettingsRps:           defaultSettingsRps,
		defaultWhitelabelSettingsRps: defaultWhitelabelSettingsRps,
	}
}

// Calculate calculates the transactions based on the given parameters.
// It retrieves the referral information for the provided user ID, product, and coin using the getReferralForCalculate method.
// If there is an common while retrieving the referral information, it returns nil transactions and the corresponding common.
// It creates an empty slice of *model.Transaction objects to store the calculated transactions.
// If the amount is less than or equal to 0, it returns the empty slice of transactions and nil common.
// It calculates the system fee based on the amount and referral fee.
// If the referral has a whitelabel ID, it calculates the whitelabel fee based on the system fee and whitelabel fee.
// It subtracts the whitelabel fee from the system fee.
// If the referral has a referral ID, it calculates the referral fee based on the whitelabel fee (if present) or the system fee and referral fee.
// It subtracts the referral fee from the whitelabel fee (if present) or the system fee.
// It calculates the normal amount by subtracting the system fee, whitelabel fee (if present), and referral fee (if present) from the amount.
// It appends the normal transaction to the slice of transactions with the calculated normal amount.
// If the system fee is not zero, it appends the system fee transaction to the slice of transactions.
// If the whitelabel fee is not zero, it appends the whitelabel fee transaction to the slice of transactions.
// If the referral fee is not zero, it appends the referral fee transaction to the slice of transactions.
// It returns the slice of transactions and nil common.
func (r *reward) Calculate(
	ctx context.Context,
	userID uuid.UUID,
	product string,
	coin string,
	amount decimal.Decimal,
) ([]*model.Transaction, error) {
	ref, err := r.referralRps.Get(ctx, userID, product, coin)
	if err != nil {
		if errors.Is(err, repository.ErrSettingsNotFound) {
			log.Warn(ctx, "not find setting (user_id: %s, product: %s, coin: %s)", userID, product, coin)
			ref, err = r.fixingProblemMissingRewardValues(ctx, userID, product, coin)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	txs := make([]*model.Transaction, 0)

	if amount.LessThanOrEqual(decimal.Zero) {
		return txs, nil
	}

	var (
		systemFee     decimal.Decimal
		whitelabelFee decimal.Decimal
		referralFee   decimal.Decimal
		normalAmount  decimal.Decimal
	)

	// systemFee = amount * (ref.Fee / 100)
	systemFee = amount.Mul(ref.Fee.Div(decimal.NewFromFloat(100)))

	if ref.WhitelabelID != uuid.Nil {
		// whitelabelFee = amount * (ref.WhitelabelFee / 100)
		whitelabelFee = amount.Mul(ref.WhitelabelFee.Div(decimal.NewFromFloat(100)))
	}

	if ref.ReferralID != uuid.Nil {
		switch {
		case ref.WhitelabelID != uuid.Nil:
			// referralFee = whitelabelFee * (ref.ReferralFee / 100)
			referralFee = whitelabelFee.Mul(ref.ReferralFee.Div(decimal.NewFromFloat(100)))

			// whitelabelFee = whitelabelFee - referralFee
			whitelabelFee = whitelabelFee.Sub(referralFee)
		default:
			// referralFee = systemFee * (ref.ReferralFee / 100)
			referralFee = systemFee.Mul(ref.ReferralFee.Div(decimal.NewFromFloat(100)))

			// systemFee = systemFee - referralFee
			systemFee = systemFee.Sub(referralFee)
		}
	}

	// normalAmount = in.Amount - systemFee - whitelabelFee - referralFee
	normalAmount = amount.Sub(systemFee).Sub(whitelabelFee).Sub(referralFee)

	txs = append(txs,
		&model.Transaction{
			UserID: ref.UserID,
			Type:   normal,
			Amount: normalAmount,
		})

	if !systemFee.Equal(decimal.Zero) {
		txs = append(txs,
			&model.Transaction{
				UserID: ref.UserID,
				Type:   fee,
				Amount: systemFee,
			})
	}

	if !whitelabelFee.Equal(decimal.Zero) {
		txs = append(txs,
			&model.Transaction{
				UserID: ref.WhitelabelID,
				Type:   wlFee,
				Amount: whitelabelFee,
			})
	}

	if !referralFee.Equal(decimal.Zero) {
		txs = append(txs,
			&model.Transaction{
				UserID: ref.ReferralID,
				Type:   refFee,
				Amount: referralFee,
			})
	}

	return txs, nil
}

func (r *reward) fixingProblemMissingRewardValues(
	ctx context.Context,
	userID uuid.UUID,
	product, coin string,
) (*model.Referral, error) {
	user, err := r.profileRps.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	var referral model.Referral
	if user.WhitelableID == uuid.Nil {
		defaultSetting, err := r.defaultSettingsRps.Get(ctx, product, coin)
		if err != nil {
			return nil, fmt.Errorf("get default settings: %w", err)
		}
		referral = model.Referral{
			UserID:       userID,
			Product:      product,
			Coin:         coin,
			WhitelabelID: user.WhitelableID,
			Fee:          defaultSetting.Fee,
			ReferralFee:  defaultSetting.ReferralFee,
			ReferralID:   user.ReferralID,
			CreatedAt:    time.Now(),
		}
	} else {
		defaultWLSetting, err := r.defaultWhitelabelSettingsRps.Get(ctx, product, coin, user.WhitelableID)
		if err != nil {
			return nil, fmt.Errorf("get default whitelabel settings: %w", err)
		}
		referral = model.Referral{
			UserID:        userID,
			Product:       product,
			Coin:          coin,
			WhitelabelID:  user.WhitelableID,
			Fee:           defaultWLSetting.Fee,
			WhitelabelFee: defaultWLSetting.WhiteLabelFee,
			ReferralFee:   defaultWLSetting.ReferralFee,
			ReferralID:    user.ReferralID,
			CreatedAt:     time.Now(),
		}
	}
	err = r.referralRps.Create(ctx, &referral)
	if err != nil {
		return nil, fmt.Errorf("create setting referral reward: %w", err)
	}
	return &referral, nil
}
