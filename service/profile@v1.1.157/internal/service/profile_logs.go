package service

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/model"
	"code.emcdtech.com/emcd/service/profile/internal/repository"
)

type ProfileLog interface {
	Log(ctx context.Context, info *model.ProfileLog) error
}

type profileLog struct {
	logRepo repository.ProfileLog
}

func NewProfileLog(profileLogRepo repository.ProfileLog) *profileLog {
	return &profileLog{
		logRepo: profileLogRepo,
	}
}

func (p *profileLog) Log(ctx context.Context, info *model.ProfileLog) error {
	const op = "service.ProfileLog.Log"

	if err := p.logRepo.Log(ctx, info); err != nil {
		log.Error(ctx, "%s: log data - %+v: %v", op, *info, err)
		return fmt.Errorf("%s: log", op)
	}
	return nil
}
