package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"code.emcdtech.com/emcd/service/whitelabel/internal/model"
	"code.emcdtech.com/emcd/service/whitelabel/internal/repository"
)

type WhiteLabel interface {
	Create(ctx context.Context, wl *model.WhiteLabel) (uuid.UUID, error)
	GetAll(ctx context.Context, skip, take int, orderBy string, asc bool) ([]*model.WhiteLabel, int, error)
	Update(ctx context.Context, wl *model.WhiteLabel) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.WhiteLabel, error)
	GetBySegmentID(ctx context.Context, segmentID int) (*model.WhiteLabel, error)
	GetByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error)
	GetByOrigin(ctx context.Context, origin string) (*model.WhiteLabel, error)
	CheckByUserID(ctx context.Context, userID int) (bool, error)
	CheckByUserIDAndOrigin(ctx context.Context, userID int, origin string) (bool, error)
	GetV2WLs(ctx context.Context) ([]*model.WhiteLabel, error)
	GetConfigByOrigin(ctx context.Context, origin string) (*model.WlConfig, error)
	SetConfigByRefID(ctx context.Context, conf *model.WlConfig) error
	SetAllowOrigin(ctx context.Context, req *model.AllowOrigin) error
	GetAllowOrigins(ctx context.Context) ([]*model.AllowOrigin, error)
	SetStratum(ctx context.Context, req *model.Stratum) error
	GetFullByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error)
	GetCoins(ctx context.Context, wlID uuid.UUID) ([]*model.WLCoins, error)
	AddCoin(ctx context.Context, wlID uuid.UUID, coinID string) error
	DeleteCoin(ctx context.Context, wlID uuid.UUID, coinID string) error
	GetWLStratumList(ctx context.Context, refID int32) ([]model.Stratum, error)
}

type whiteLabel struct {
	whiteLabelRepo repository.WhiteLabel
}

func NewWhiteLabel(whiteLabelRepo repository.WhiteLabel) WhiteLabel {
	return &whiteLabel{
		whiteLabelRepo: whiteLabelRepo,
	}
}

func (w *whiteLabel) Create(ctx context.Context, wl *model.WhiteLabel) (uuid.UUID, error) {
	wl.ID = uuid.New()
	err := w.whiteLabelRepo.Create(ctx, wl)
	if err != nil {
		return uuid.Nil, fmt.Errorf("whitelabel create: %w", err)
	}
	return wl.ID, nil
}

func (w *whiteLabel) GetAll(ctx context.Context, skip, take int, orderBy string, asc bool) ([]*model.WhiteLabel, int, error) {
	wls, count, err := w.whiteLabelRepo.GetAll(ctx, skip, take, orderBy, asc)
	if err != nil {
		return nil, 0, fmt.Errorf("whitelabel get all: %w", err)
	}

	return wls, count, nil
}

func (w *whiteLabel) Update(ctx context.Context, wl *model.WhiteLabel) error {
	err := w.whiteLabelRepo.Update(ctx, wl)
	if err != nil {
		return fmt.Errorf("whitelabel update: %w", err)
	}
	return nil
}

func (w *whiteLabel) Delete(ctx context.Context, id uuid.UUID) error {
	err := w.whiteLabelRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("whitelabel delete: %w", err)
	}
	return nil
}

func (w *whiteLabel) GetByID(ctx context.Context, id uuid.UUID) (*model.WhiteLabel, error) {
	wl, err := w.whiteLabelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("whitelabel: get by id: %w", err)
	}
	return wl, nil
}

func (w *whiteLabel) GetBySegmentID(ctx context.Context, segmentID int) (*model.WhiteLabel, error) {
	wl, err := w.whiteLabelRepo.GetBySegmentID(ctx, segmentID)
	if err != nil {
		return nil, fmt.Errorf("whitelabel: get by segment_id: %w", err)
	}

	return wl, nil
}

func (w *whiteLabel) GetByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error) {
	wl, err := w.whiteLabelRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("whitelabel: get by user_id: %w", err)
	}

	return wl, nil
}

func (w *whiteLabel) GetByOrigin(ctx context.Context, origin string) (*model.WhiteLabel, error) {
	wl, err := w.whiteLabelRepo.GetByOrigin(ctx, origin)
	if err != nil {
		return nil, fmt.Errorf("whitelabel: get by origin: %w", err)
	}

	return wl, nil
}

func (w *whiteLabel) CheckByUserID(ctx context.Context, userID int) (bool, error) {
	res, err := w.whiteLabelRepo.CheckByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("whitelabel: check by user_id: %w", err)
	}

	return res, nil
}

func (w *whiteLabel) CheckByUserIDAndOrigin(ctx context.Context, userID int, origin string) (bool, error) {
	res, err := w.whiteLabelRepo.CheckByUserIDAndOrigin(ctx, userID, origin)
	if err != nil {
		return false, fmt.Errorf("whitelabel: check by user_id and origin: %w", err)
	}

	return res, nil
}

func (w *whiteLabel) GetV2WLs(ctx context.Context) ([]*model.WhiteLabel, error) {
	res, err := w.whiteLabelRepo.GetV2WLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("whitelabel: get v2: %w", err)
	}

	return res, nil
}

func (w *whiteLabel) GetConfigByOrigin(ctx context.Context, origin string) (*model.WlConfig, error) {
	res, err := w.whiteLabelRepo.GetConfigByOrigin(ctx, origin)
	if err != nil {
		return nil, fmt.Errorf("whitelabel: get config: %w", err)
	}
	if res.RefID == "" {
		return nil, fmt.Errorf("whitelabel: get config: empty ref_id")
	}
	segmentID, err := strconv.Atoi(res.RefID)
	if err != nil {
		return nil, err
	}
	wl, err := w.GetBySegmentID(ctx, segmentID)
	if err != nil {
		return nil, err
	}
	res.WhitelabelID = wl.ID.String()
	return res, nil
}

func (w *whiteLabel) SetConfigByRefID(ctx context.Context, conf *model.WlConfig) error {
	err := w.whiteLabelRepo.SetConfigByRefID(ctx, conf)
	if err != nil {
		return fmt.Errorf("whitelabel: set config: %w", err)
	}

	return nil
}

func (w *whiteLabel) SetAllowOrigin(ctx context.Context, req *model.AllowOrigin) error {
	err := w.whiteLabelRepo.SetAllowOrigin(ctx, req)
	if err != nil {
		return fmt.Errorf("set allow origin: %w", err)
	}

	return nil
}

func (w *whiteLabel) GetAllowOrigins(ctx context.Context) ([]*model.AllowOrigin, error) {
	res, err := w.whiteLabelRepo.GetAllowOrigins(ctx)
	if err != nil {
		return nil, fmt.Errorf("get allow origins: %w", err)
	}

	return res, nil
}

func (w *whiteLabel) SetStratum(ctx context.Context, req *model.Stratum) error {
	err := w.whiteLabelRepo.SetStratum(ctx, req)
	if err != nil {
		return fmt.Errorf("set stratum: %w", err)
	}

	return nil
}

func (w *whiteLabel) GetFullByUserID(ctx context.Context, userID int) (*model.WhiteLabel, error) {
	return w.whiteLabelRepo.GetFullByUserID(ctx, userID)
}

func (w *whiteLabel) GetCoins(ctx context.Context, wlID uuid.UUID) ([]*model.WLCoins, error) {
	coins, err := w.whiteLabelRepo.GetCoins(ctx, wlID)
	if err != nil {
		return nil, fmt.Errorf("GetCoins: %w", err)
	}
	return coins, nil
}

func (w *whiteLabel) AddCoin(ctx context.Context, wlID uuid.UUID, coinID string) error {
	err := w.whiteLabelRepo.AddCoin(ctx, wlID, coinID)
	if err != nil {
		return fmt.Errorf("add coin: %w", err)
	}
	return nil
}

func (w *whiteLabel) DeleteCoin(ctx context.Context, wlID uuid.UUID, coinID string) error {
	err := w.whiteLabelRepo.DeleteCoin(ctx, wlID, coinID)
	if err != nil {
		return fmt.Errorf("delete coin: %w", err)
	}
	return nil
}

func (w *whiteLabel) GetWLStratumList(ctx context.Context, refID int32) ([]model.Stratum, error) {
	return w.whiteLabelRepo.GetWLStratumListV2(ctx, refID)
}
