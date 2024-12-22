package repository

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/service/profile/protocol/profile"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/referral/internal/model"
)

//go:generate mockery --name=Profile --structname=MockProfile --outpkg=repository --output ./ --filename profile_mock.go
type Profile interface {
	Get(ctx context.Context, userID uuid.UUID) (*model.User, error)
	GetAllSubUsers(ctx context.Context, userUUID string) ([]string, error)
}

type ProfileRepository struct {
	cli profile.ProfileServiceClient
}

func (c *ProfileRepository) Get(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	resp, err := c.cli.GetByUserID(ctx, &profile.GetByUserIDRequest{
		UserID: userID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("profile.GetByUserID (user_id: %s): %w", userID, err)
	}
	r := &model.User{
		UserID: userID,
	}
	r.ReferralID, err = uuid.Parse(resp.GetProfile().GetUser().GetNewRefId())
	if err != nil {
		return nil, fmt.Errorf("parse referral_id `%s`: %w", resp.GetProfile().GetUser().GetNewRefId(), err)
	}
	r.WhitelableID, err = uuid.Parse(resp.GetProfile().GetUser().GetWhiteLabelID())
	if err != nil {
		return nil, fmt.Errorf("parse whitelable_id `%s`: %w", resp.GetProfile().GetUser().GetWhiteLabelID(), err)
	}
	return r, nil
}

func (c *ProfileRepository) GetAllSubUsers(ctx context.Context, userUUID string) ([]string, error) {
	resp, err := c.cli.GetAllSubUsers(ctx, &profile.GetAllSubUsersRequest{
		UserId: userUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("profile.GetAllSubUsers (user_id: %s): %w", userUUID, err)
	}

	subAccUUIDs := make([]string, 0, len(resp.GetSubs()))
	for _, subs := range resp.GetSubs() {
		subAccUUIDs = append(subAccUUIDs, subs.UserId)
	}
	return subAccUUIDs, nil
}

func NewProfileRepository(
	cli profile.ProfileServiceClient,
) Profile {
	return &ProfileRepository{
		cli: cli,
	}
}
