package service

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

type CreateProfileRequest struct {
	Username         string
	Email            string
	Password         string
	RefId            int32
	RefUuid          string
	WhiteLabelUuid   uuid.UUID
	AppleId          string
	Language         string
	TelegramId       string
	TelegramUserName string
}

type CreateProfileResponse struct {
	UserUUID uuid.UUID
	UserID   int32
}

func (p *profile) CreateProfile(ctx context.Context, req CreateProfileRequest) (*CreateProfileResponse, error) {
	refUUID := uuid.Nil
	var err error
	if req.RefId > 0 {
		refUUID, err = p.GetUserIDByOldID(ctx, req.RefId)
		if err != nil {
			return nil, err
		}
	}

	email := req.Email
	wlUUID := uuid.Nil
	if req.WhiteLabelUuid != uuid.Nil {
		wl, err := p.whiteLabel.GetByID(ctx, wlUUID)
		if err != nil {
			return nil, err
		}
		prefix := wl.Prefix
		if prefix != "" {
			email = prefix + email
		}
	}

	user := &model.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  req.Password,
		IsActive:  true,
		AppleID:   req.AppleId,
		CreatedAt: time.Now(),

		RefID:    int(req.RefId),
		NewRefID: refUUID,

		WhiteLabelID: wlUUID,
		Language:     req.Language,

		TgID:       req.TelegramId,
		TgUsername: req.TelegramUserName,
	}

	// при использовании альтернативных способов регистрации пользователь не вводит пароль,
	// Будем генерировать случайный. Потом пользователь сможет сбросить пароль при необходимости.
	if len(user.Password) == 0 {
		err := user.SetRandomPassword()
		if err != nil {
			return nil, fmt.Errorf("set random password:  %w", err)
		}
		log.Info(ctx, "set random password for user %s", user.Email)
	}

	id, err := p.createUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("profile: create: %w", err)
	}

	err = p.setReferralSettingsForSubAccount(ctx, user.ID, user.ParentID, user.WhiteLabelID, user.NewRefID)
	if err != nil {
		log.Error(ctx, "profile.CreateProfile: %v", err)
	}

	return &CreateProfileResponse{
		UserUUID: user.ID,
		UserID:   id,
	}, nil
}
