package service

import (
	"context"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/referral/internal/clients/coin"
	"code.emcdtech.com/emcd/service/referral/internal/clients/promocode"
	"code.emcdtech.com/emcd/service/referral/internal/model"
	"code.emcdtech.com/emcd/service/referral/internal/repository"
)

type DefaultUsersSettings struct {
	repo            *repository.DefaultUserSettings
	promoCodeClient *promocode.Client
	coinClient      *coin.Client
}

func NewDefaultUsersSettings(
	repo *repository.DefaultUserSettings,
	promoCodeClient *promocode.Client,
	coinClient *coin.Client) *DefaultUsersSettings {
	return &DefaultUsersSettings{
		repo:            repo,
		promoCodeClient: promoCodeClient,
		coinClient:      coinClient,
	}
}

func (s *DefaultUsersSettings) CreateUsersSettings(ctx context.Context, settings []model.ReferralSettings) error {
	err := s.repo.CreateUsersSettings(ctx, settings)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultUsersSettings) UpdateUsersSettings(
	ctx context.Context,
	settings []model.ReferralSettings,
	updateMode model.UpdateMode) error {

	// пользователи с промокодами
	userAndPromoCodes := map[string]map[string][]string{} // referralUUID, Coin, list of users uuid with promoCode to this coin

	if updateMode == model.UpdateModeAll {
		users, err := s.repo.GetUserUUIDsByReferralsUUID(ctx, settings)
		if err != nil {
			return err
		}

		if len(users) > 0 {
			coins, err := s.coinClient.GetCoins(ctx)
			if err != nil {
				return err
			}

			userUUIDs := make([]string, 0, len(users))
			userAndReferralLink := make(map[string]uuid.UUID)

			for referralUUID, user := range users {
				userUUIDs = append(userUUIDs, user...)

				for _, uUUID := range user {
					userAndReferralLink[uUUID] = referralUUID
				}
			}

			userAndPromo, err := s.promoCodeClient.GetActiveUsersPromoCodes(ctx, userUUIDs)
			if err != nil {
				return err
			}

			for userUUID, promoCodes := range userAndPromo {
				for _, code := range promoCodes.PromoCodes {
					for _, coinOldID := range code.PromoCode.CoinIDs {
						coinCode, ok := coins[coinOldID]
						if !ok {
							continue
						}

						referralUUID, ok := userAndReferralLink[userUUID.String()]
						if !ok {
							continue
						}

						referralCoins, ok := userAndPromoCodes[referralUUID.String()]
						if !ok {
							userAndPromoCodes[referralUUID.String()] = map[string][]string{
								coinCode: {userUUID.String()},
							}
							continue
						}

						referralCoins[coinCode] = append(referralCoins[coinCode], userUUID.String())
						userAndPromoCodes[referralUUID.String()] = referralCoins
					}
				}
			}
		}
	}

	ctxForUpdate := context.WithoutCancel(ctx)
	err := s.repo.UpdateUsersSettings(ctxForUpdate, settings, updateMode, userAndPromoCodes)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultUsersSettings) GetUsersSettings(
	ctx context.Context,
	referrals []uuid.UUID,
	products, coins []string) (map[uuid.UUID]model.ReferralSettings, error) {
	users, err := s.repo.GetUsersSettings(ctx, referrals, products, coins)
	if err != nil {
		return nil, err
	}

	return users, nil
}
