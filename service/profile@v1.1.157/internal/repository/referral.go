package repository

import (
	"context"
	"fmt"

	defPb "code.emcdtech.com/emcd/service/referral/protocol/default_settings"
	defWlPb "code.emcdtech.com/emcd/service/referral/protocol/default_whitelabel_settings"
	refPb "code.emcdtech.com/emcd/service/referral/protocol/referral"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Referral is an interface that defines the method for creating a referral.
// The method takes a context and a pointer to a model.Referral as parameters.
// It returns an error if any error occurs during the creation process.
type Referral interface {
	Create(ctx context.Context, userID, wlID, refID uuid.UUID) error
	SetReferralUUID(ctx context.Context, userUUID, refUUID uuid.UUID) error
	CreateForSubAccount(ctx context.Context, userUUID, parentUUID, wlID, refID uuid.UUID) error
}

// referral is a struct that represents a referral object.
// It contains a field refCli of type refPb.ReferralServiceClient
type referral struct {
	refCli   refPb.ReferralServiceClient
	defCli   defPb.DefaultSettingsServiceClient
	defWlCli defWlPb.DefaultWhitelabelSettingsServiceClient
	// whitelabels which take commissions with master-slave relationship. maybe will be needed
	// wlsWithZeroCommission []uuid.UUID
}

// NewReferral creates a new instance of `referral` with the given `refCli`
// as the `refPb.ReferralServiceClient`.
func NewReferral(
	refCli refPb.ReferralServiceClient,
	defCli defPb.DefaultSettingsServiceClient,
	defWlCli defWlPb.DefaultWhitelabelSettingsServiceClient,
) *referral {
	return &referral{
		refCli:   refCli,
		defCli:   defCli,
		defWlCli: defWlCli,
	}
}

// Create creates a new referral in the referral service using the provided context and referral information.
// It sends a CreateRequest to `r.refCli` with the referral details converted to the protobuf equivalent.
// If an error occurs during the Create request, it returns an error wrapped with additional context.
func (r *referral) Create(ctx context.Context, userID, wlID, refID uuid.UUID) error {
	var (
		req *refPb.CreateMultipleRequest
		err error
	)

	if wlID == uuid.Nil {
		req, err = r.getDefaultReferralSetting(ctx, userID, wlID, refID)
		if err != nil {
			return fmt.Errorf("referral.getDefaultReferralSetting: %v", err)
		}
	} else {
		req, err = r.getDefaultWhitelabelReferralSetting(ctx, userID, wlID, refID)
		if err != nil {
			return fmt.Errorf("referral.getDefaultWhitelabelReferralSetting: %w", err)
		}
	}
	if len(req.Referrals) == 0 {
		return nil
	}

	_, err = r.refCli.CreateMultiple(ctx, req)
	if err != nil {
		return fmt.Errorf("referral.refCli.CreateMultiple: %v", err)
	}

	return nil
}

func (r *referral) SetReferralUUID(ctx context.Context, userUUID, refUUID uuid.UUID) error {
	_, err := r.refCli.SetReferralUUID(ctx, &refPb.SetReferralUUIDRequest{
		UserUuid:     userUUID.String(),
		ReferralUuid: refUUID.String(),
	})
	if err != nil {
		return fmt.Errorf("referral.refCli.SetReferralUUID: %v", err)
	}
	return nil
}

func (r *referral) CreateForSubAccount(ctx context.Context, userUUID, parentUUID, wlID, refID uuid.UUID) error {
	req, err := r.getUserSettingByUserUUID(ctx, userUUID, parentUUID, wlID, refID)
	if err != nil {
		return fmt.Errorf("referral.getUserSettingByUserUUID: %v", err)
	}

	if len(req.Referrals) == 0 {
		return nil
	}

	_, err = r.refCli.CreateMultiple(ctx, req)
	if err != nil {
		return fmt.Errorf("referral.refCli.CreateMultiple: %v", err)
	}

	return nil
}

func (r *referral) getUserSettingByUserUUID(ctx context.Context, userUUID, parentUUID, wlID, refID uuid.UUID) (*refPb.CreateMultipleRequest, error) {
	settingsResp, err := r.refCli.List(ctx, &refPb.ListRequest{
		UserId: parentUUID.String(),
		Skip:   0,
		Take:   1000,
	})
	if err != nil {
		return nil, fmt.Errorf("referral.refCli.List: %v", err)
	}

	req := new(refPb.CreateMultipleRequest)

	for i := range settingsResp.List {
		req.Referrals = append(req.Referrals, &refPb.Referral{
			UserId:        userUUID.String(),
			Product:       settingsResp.List[i].Product,
			Coin:          settingsResp.List[i].Coin,
			Fee:           settingsResp.List[i].Fee,
			WhitelabelFee: decimal.Zero.String(),
			WhitelabelId:  wlID.String(),
			ReferralId:    refID.String(),
			ReferralFee:   settingsResp.List[i].ReferralFee,
		})
	}

	return req, nil
}

func (r *referral) getDefaultReferralSetting(
	ctx context.Context,
	userID, wlID, refID uuid.UUID,
) (*refPb.CreateMultipleRequest, error) {
	settingsResp, err := r.defCli.GetAllWithoutPagination(ctx, &defPb.GetAllWithoutPaginationRequest{
		ReferrerId: refID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("referral.defCli.GetAllWithoutPagination: %v", err)
	}

	req := new(refPb.CreateMultipleRequest)

	for i := range settingsResp.List {
		req.Referrals = append(req.Referrals, &refPb.Referral{
			UserId:        userID.String(),
			Product:       settingsResp.List[i].Product,
			Coin:          settingsResp.List[i].Coin,
			Fee:           settingsResp.List[i].Fee,
			WhitelabelFee: decimal.Zero.String(),
			WhitelabelId:  wlID.String(),
			ReferralId:    refID.String(),
			ReferralFee:   settingsResp.List[i].ReferralFee,
		})
	}

	return req, nil
}

func (r *referral) getDefaultWhitelabelReferralSetting(
	ctx context.Context,
	userID, wlID, refID uuid.UUID,
) (*refPb.CreateMultipleRequest, error) {
	settingsResp, err := r.defWlCli.GetV2(ctx, &defWlPb.GetV2Request{
		WhitelabelId: wlID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("referral.defWlCli.GetAllWithoutPagination: %v", err)
	}

	req := new(refPb.CreateMultipleRequest)

	for i := range settingsResp.Settings {
		req.Referrals = append(req.Referrals, &refPb.Referral{
			UserId:        userID.String(),
			Product:       settingsResp.Settings[i].Product,
			Coin:          settingsResp.Settings[i].Coin,
			Fee:           settingsResp.Settings[i].Fee,
			WhitelabelFee: settingsResp.Settings[i].WhitelabelFee,
			WhitelabelId:  wlID.String(),
			ReferralId:    refID.String(),
			ReferralFee:   settingsResp.Settings[i].ReferralFee,
		})
	}

	return req, nil
}
