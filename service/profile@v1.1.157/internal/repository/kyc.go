package repository

import (
	"context"
	"fmt"
	"time"

	pgTx "code.emcdtech.com/emcd/sdk/pg"
)

type Kyc interface {
	GetUserStatus(ctx context.Context, userID int) (int, time.Time, string, string, string, error)
	SetUserStatus(ctx context.Context, userID, status int) error
	InsertHistory(ctx context.Context, userID int, data []byte) error
}

type kyc struct {
	trx pgTx.PgxTransactor
}

func NewKyc(trx pgTx.PgxTransactor) *kyc {
	return &kyc{
		trx: trx,
	}
}

func (k *kyc) GetUserStatus(ctx context.Context, userID int) (int, time.Time, string, string, string, error) {
	query := `SELECT u.kyc_idenfy_status, 
        COALESCE(kyc.created_at, '1970-01-01 00:00:00'),
       	COALESCE(kyc.data->'status'->>'overall', '') as overall,
       	COALESCE(kyc.data->'status'->>'manualDocument', kyc.data->'status'->>'autoDocument', '') as doc_check,
       	COALESCE(kyc.data->'status'->>'manualFace', kyc.data->'status'->>'autoFace', '') as face_check
		FROM emcd.users u
		LEFT JOIN emcd.users_kyc_idenfy_history kyc ON kyc.user_id = u.id
		WHERE u.id = $1
		ORDER BY kyc.created_at DESC LIMIT 1`

	var (
		status    int
		lastTryAt time.Time
		overall   string
		docCheck  string
		faceCheck string
	)
	if err := k.trx.Runner(ctx).QueryRow(ctx, query, userID).Scan(&status, &lastTryAt, &overall, &docCheck, &faceCheck); err != nil {
		return 0, time.Time{}, "", "", "", fmt.Errorf("queryRow: %w", err)
	}

	return status, lastTryAt, overall, docCheck, faceCheck, nil
}

func (k *kyc) SetUserStatus(ctx context.Context, userID, status int) error {
	query := `
	UPDATE emcd.users
	SET kyc_idenfy_status = $1
	WHERE id = $2`

	if _, err := k.trx.Runner(ctx).Exec(ctx, query, status, userID); err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (k *kyc) InsertHistory(ctx context.Context, userID int, data []byte) error {
	query := `
	INSERT INTO emcd.users_kyc_idenfy_history (user_id, data)
	VALUES ($1, $2)`

	if _, err := k.trx.Runner(ctx).Exec(ctx, query, userID, data); err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
