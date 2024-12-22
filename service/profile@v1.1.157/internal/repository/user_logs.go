package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

type UserLogs interface {
	Create(ctx context.Context, ul *model.UserLog) error
	CreateWithoutToken(ctx context.Context, ul *model.UserLog) error
	DeactivateByType(ctx context.Context, userID int, changeType string) error
	Get(ctx context.Context, token, changeType string, userID int32, active bool) (*model.UserLog, error)
}

type userLogs struct {
	trx pgTx.PgxTransactor
}

func NewUserLogs(trx pgTx.PgxTransactor) *userLogs {
	return &userLogs{
		trx: trx,
	}
}

func (u *userLogs) Create(ctx context.Context, ul *model.UserLog) error {
	query := `INSERT INTO emcd.user_logs (user_id, change_type, ip, token, old_value, value, active, is_segment_sended) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := u.trx.Runner(ctx).Exec(ctx, query, ul.UserID, ul.ChangeType, ul.IP, ul.Token, ul.OldValue,
		ul.Value, ul.Active, ul.IsSegmentSent)
	if err != nil {
		return fmt.Errorf("user logs: create: %w", err)
	}
	return nil
}

func (u *userLogs) CreateWithoutToken(ctx context.Context, ul *model.UserLog) error {
	query := `INSERT INTO emcd.user_logs(user_id, change_type, old_value, value, ip, active) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := u.trx.Runner(ctx).Exec(ctx, query, ul.UserID, ul.ChangeType,
		ul.OldValue, ul.Value, ul.IP, ul.Active)
	if err != nil {
		return fmt.Errorf("user logs: create without token: %w", err)
	}
	return nil
}

func (u *userLogs) DeactivateByType(ctx context.Context, userID int, changeType string) error {
	query := `UPDATE emcd.user_logs SET active = false WHERE change_type = $1 AND user_id = $2`
	_, err := u.trx.Runner(ctx).Exec(ctx, query, changeType, userID)
	if err != nil {
		return fmt.Errorf("user logs: deactivate by type: %w", err)
	}
	return nil
}

func (u *userLogs) Get(ctx context.Context, token, changeType string, userID int32, active bool) (*model.UserLog, error) {
	query := `SELECT id,user_id, change_type,old_value, value, ip, active FROM emcd.user_logs
        WHERE token = $1 AND change_type = $2 AND used = false AND active = $3 AND user_id = $4`
	var ul model.UserLog
	err := u.trx.Runner(ctx).QueryRow(ctx, query, token, changeType, active, userID).
		Scan(&ul.ID, &ul.UserID, &ul.ChangeType, &ul.OldValue, &ul.Value, &ul.IP, &ul.Active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("user logs: get: %w", err)
	}
	return &ul, nil
}
