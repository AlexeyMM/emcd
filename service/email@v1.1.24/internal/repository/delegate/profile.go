package delegate

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/service/profile/protocol/profile"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

func SupportLanguage(language string) string {
	if language == "ru" {
		return language
	}
	return "en"
}

type profileClient struct {
	cli profile.ProfileServiceClient
}

func NewProfileClient(cli profile.ProfileServiceClient) *profileClient {
	return &profileClient{
		cli: cli,
	}
}

func (s *profileClient) GetByUserID(ctx context.Context, userID uuid.UUID) (model.Profile, error) {
	req, err := s.cli.GetByUserID(ctx,
		&profile.GetByUserIDRequest{
			UserID: userID.String(),
		},
	)
	if err != nil {
		return model.Profile{}, fmt.Errorf("profile.GetByUserID: %w", err)
	}
	if req.GetProfile().GetUser() == nil {
		return model.Profile{}, fmt.Errorf("profile.GetByUserID: unexpected nil value")
	}
	return s.convertProtoProfile(req.Profile.User)
}

func (s *profileClient) GetNotificationSettings(
	ctx context.Context,
	userID uuid.UUID,
) (model.NotificationSettings, error) {
	resp, err := s.cli.GetNotificationSettings(ctx,
		&profile.GetNotificationSettingsRequest{
			UserId: userID.String(),
		},
	)
	if err != nil {
		return model.NotificationSettings{}, fmt.Errorf("profile.GetNotificationSettings: %w", err)
	}
	return newNotificationSettings(userID, resp)
}

func (s *profileClient) convertProtoProfile(user *profile.User) (model.Profile, error) {
	whiteLabelID, err := uuid.Parse(user.WhiteLabelID)
	if err != nil {
		return model.Profile{}, fmt.Errorf("parse WhiteLabelID: %w", err)
	}
	p := model.Profile{
		Email:        user.Email,
		WhiteLabelID: whiteLabelID,
		Language:     SupportLanguage(user.Language),
	}
	return p, nil
}

func newNotificationSettings(
	userID uuid.UUID,
	in *profile.GetNotificationSettingsResponse,
) (model.NotificationSettings, error) {
	whiteLabelID, err := uuid.Parse(in.GetSettings().WhitelabelId)
	if err != nil {
		return model.NotificationSettings{}, fmt.Errorf("parse WhiteLabelID: %w", err)
	}
	return model.NotificationSettings{
		UserID:                 userID,
		Email:                  in.GetSettings().Email,
		Language:               in.GetSettings().Language,
		IsEmailNotificationsOn: in.GetSettings().IsEmailNotificationsOn,
		WhiteLabelID:           whiteLabelID,
	}, nil
}
