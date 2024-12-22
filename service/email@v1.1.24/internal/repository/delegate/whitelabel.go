package delegate

import (
	"context"
	"errors"

	wl "code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/email/internal/model"
)

type whitelabelRepo struct {
	cli wl.WhitelabelServiceClient
}

func NewWhitelabel(cli wl.WhitelabelServiceClient) *whitelabelRepo {
	return &whitelabelRepo{
		cli: cli,
	}
}

func (w *whitelabelRepo) GetVersion(ctx context.Context, wlID uuid.UUID) (int, error) {
	resp, err := w.cli.GetByID(ctx, &wl.GetByIDRequest{
		Id: wlID.String(),
	})
	if err != nil {
		return 0, err
	}
	if resp.WhiteLabel == nil {
		return 0, errors.New("whitelabel not found")
	}
	return int(resp.WhiteLabel.Version), nil
}

func (w *whitelabelRepo) GetWlFullDomain(ctx context.Context, wlOwnerUserID int32) (string, error) {
	originsResp, err := w.cli.GetAllowOrigins(ctx, &wl.EmptyRequest{})
	if err != nil {
		return "", err
	}
	for _, originsData := range originsResp.List {
		if originsData.UserId == wlOwnerUserID {
			return originsData.Origin, nil
		}
	}
	return "", errors.New("failed to get allow origin for whitelabel")
}

func (w *whitelabelRepo) GetWlByID(ctx context.Context, wlID uuid.UUID) (*model.Whitelabel, error) {
	resp, err := w.cli.GetByID(ctx, &wl.GetByIDRequest{Id: wlID.String()})
	if err != nil {
		return nil, err
	}
	whitelabel := resp.GetWhiteLabel()
	return &model.Whitelabel{
		Id:     whitelabel.Id,
		UserId: whitelabel.UserId,
		Origin: whitelabel.Origin,
		Prefix: whitelabel.Prefix,
	}, nil
}
