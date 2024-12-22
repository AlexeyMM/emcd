package handler

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	referralPb "code.emcdtech.com/emcd/service/accounting/protocol/referral"
	"context"
	"fmt"
)

type ReferralHandler struct {
	referralService service.Referral
	referralPb.UnimplementedAccountingReferralServiceServer
}

func NewReferralHandler(
	referralService service.Referral,
) *ReferralHandler {
	return &ReferralHandler{
		referralService: referralService,
		UnimplementedAccountingReferralServiceServer: referralPb.UnimplementedAccountingReferralServiceServer{},
	}
}

func (h *ReferralHandler) GetReferralsStatistic(ctx context.Context, req *referralPb.GetReferralsStatisticRequest) (*referralPb.GetReferralsStatisticResponse, error) {
	resp, err := h.referralService.GetReferralsStatistic(ctx, req)
	if err != nil {
		sdkLog.Error(ctx, "getReferralsStatistic referral handler: %v", err)
		return nil, fmt.Errorf("wallets handleer: %w", err)
	}
	return resp, nil
}
