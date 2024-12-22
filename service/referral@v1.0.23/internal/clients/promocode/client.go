package promocode

import (
	"context"
	"errors"
	"io"

	"code.emcdtech.com/emcd/service/promocode/protocol/promocode"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/referral/internal/model"
)

type Client struct {
	c promocode.PromocodeClient
}

func NewClient(c promocode.PromocodeClient) *Client {
	return &Client{
		c: c}
}

func (c *Client) GetActiveUsersPromoCodes(ctx context.Context, userUUIDs []string) (map[uuid.UUID]model.UserAndPromoCodes, error) {
	resp, err := c.c.GetActiveUsersPromoCodes(ctx, &promocode.GetActiveUsersPromoCodesRequest{
		UserUuid: userUUIDs,
	})
	if err != nil {
		return nil, err
	}
	list := make(map[uuid.UUID]model.UserAndPromoCodes, 20)
	for {
		msg, err := resp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		for _, user := range msg.GetUsersPromoCodes() {
			if len(user.GetPromoCodesInfo()) == 0 {
				continue
			}

			userUUID, err := uuid.Parse(user.GetUserUuid())
			if err != nil {
				return nil, err
			}
			listItem := model.UserAndPromoCodes{
				UserUUID:   userUUID,
				UserID:     user.GetUserId(),
				PromoCodes: make([]model.UserPromoCodeAndPromo, 0, len(user.GetPromoCodesInfo())),
			}
			for _, promoCodeInfo := range user.GetPromoCodesInfo() {
				userPromoInfo := promoCodeInfo.GetUserPromoCode()
				usedPromoCode := promoCodeInfo.GetPromoCodes()

				promoCode := model.UserPromoCodeAndPromo{
					UserPromoCode: &model.UserPromoCode{
						UserID:      userPromoInfo.GetUserId(),
						PromoCodeID: userPromoInfo.GetPromoCodeId(),
						ExpiresAt:   userPromoInfo.GetExpiresAt().AsTime(),
						CreatedAt:   userPromoInfo.GetCreatedAt().AsTime(),
					},
					PromoCode: &model.PromoCode{
						ID:                      usedPromoCode.GetId(),
						Code:                    usedPromoCode.GetCode(),
						ValidDaysAmount:         usedPromoCode.GetValidDaysAmount(),
						HasNoLimit:              usedPromoCode.GetHasNoLimit(),
						FeePercent:              usedPromoCode.GetFeePercent(),
						ReferralEnabled:         usedPromoCode.GetReferralEnabled(),
						IsActive:                usedPromoCode.GetIsActive(),
						IsDisposable:            usedPromoCode.GetIsDisposable(),
						RefID:                   usedPromoCode.GetRefId(),
						CoinIDs:                 usedPromoCode.GetCoinIds(),
						IsSummable:              usedPromoCode.GetIsSummable(),
						IsOnlyForRegistration:   usedPromoCode.GetIsOnlyForRegistration(),
						IsOnlyForPrivateCabinet: usedPromoCode.GetIsOnlyForPrivateCabinet(),
						CreatedAt:               usedPromoCode.GetCreatedAt().AsTime(),
						ExpiresAt:               usedPromoCode.GetExpiresAt().AsTime(),
					},
				}

				listItem.PromoCodes = append(listItem.PromoCodes, promoCode)
			}

			list[userUUID] = listItem
		}
	}

	return list, nil
}
