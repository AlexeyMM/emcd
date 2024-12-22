package repository

import (
	"code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"context"
	"github.com/google/uuid"
)

type WhiteLabel interface {
	GetUserIDByID(ctx context.Context, wlID uuid.UUID) (int32, error)
}

type whiteLabelRepo struct {
	cli whitelabel.WhitelabelServiceClient
}

func NewWhiteLabel(cli whitelabel.WhitelabelServiceClient) WhiteLabel {
	return &whiteLabelRepo{
		cli: cli,
	}
}

func (r *whiteLabelRepo) GetUserIDByID(ctx context.Context, wlID uuid.UUID) (int32, error) {
	resp, err := r.cli.GetByID(ctx, &whitelabel.GetByIDRequest{
		Id: wlID.String(),
	})
	if err != nil {
		return 0, err
	}
	return resp.WhiteLabel.UserId, nil
}
