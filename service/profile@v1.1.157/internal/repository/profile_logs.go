package repository

import (
	"context"

	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

type ProfileLog interface {
	Log(ctx context.Context, info *model.ProfileLog) error
}

type profileLog struct {
	trx pgTx.PgxTransactor
}

func NewProfileLog(trx pgTx.PgxTransactor) *profileLog {
	return &profileLog{
		trx: trx,
	}
}

func (p *profileLog) Log(ctx context.Context, info *model.ProfileLog) error {
	query := `INSERT INTO profile_logs (originator, change_type, details) VALUES ($1, $2, $3)`
	_, err := p.trx.Runner(ctx).Exec(ctx, query, info.Originator, info.ChangeType, info.Details)

	return err
}
