package repository

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

type Whitelabel interface {
	GetByUserID(ctx context.Context, userID int32) (*model.WhiteLabel, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.WhiteLabel, error)
}

type Wl struct {
	whitelabelCli whitelabel.WhitelabelServiceClient
}

func NewWhitelabel(whitelabelCli whitelabel.WhitelabelServiceClient) *Wl {
	return &Wl{
		whitelabelCli: whitelabelCli,
	}
}

func (w *Wl) GetByUserID(ctx context.Context, userID int32) (*model.WhiteLabel, error) {
	resp, err := w.whitelabelCli.GetByUserID(ctx, &whitelabel.GetByUserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return w.parseProto(resp.WhiteLabel)
}

func (w *Wl) parseProto(wl *whitelabel.WhiteLabel) (*model.WhiteLabel, error) {
	if wl == nil {
		return nil, nil
	}

	var (
		id  uuid.UUID
		err error
	)
	if wl.Id != "" {
		id, err = uuid.Parse(wl.Id)
		if err != nil {
			return nil, fmt.Errorf("parse proto: parse wl id %s: %w", wl.Id, err)
		}
	}
	return &model.WhiteLabel{
		ID:                    id,
		Prefix:                wl.GetPrefix(),
		Version:               int(wl.GetVersion()),
		IsEmailConfirmEnabled: wl.IsEmailConfirmEnabled,
	}, nil
}

func (w *Wl) GetByID(ctx context.Context, id uuid.UUID) (*model.WhiteLabel, error) {
	resp, err := w.whitelabelCli.GetByID(ctx, &whitelabel.GetByIDRequest{
		Id: id.String(),
	})
	if err != nil {
		return nil, err
	}
	return w.parseProto(resp.WhiteLabel)
}
