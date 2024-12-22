package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

func TestUserLogs_Create(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	ul := model.UserLog{
		ID:         1,
		UserID:     2,
		ChangeType: "type",
		IP:         "ip",
		Token:      "token",
		OldValue:   "old value",
		Value:      "value",
		Active:     true,
	}

	rep := NewUserLogs(transactor)
	err := rep.Create(ctx, &ul)
	require.NoError(t, err)

	exec, err := dbPool.Exec(ctx, `UPDATE emcd.user_logs SET used = false, id = $1 WHERE user_id = $2`, ul.ID, ul.UserID)
	require.NoError(t, err)
	require.Equal(t, int64(1), exec.RowsAffected())

	tx, err := dbPool.Begin(ctx)
	defer tx.Rollback(ctx)
	userLog, err := rep.Get(ctx, ul.Token, ul.ChangeType, ul.UserID, ul.Active)
	require.NoError(t, err)
	userLog.Token = ul.Token
	require.Equal(t, &ul, userLog)
}

func TestUserLogs_CreateWithoutToken(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	ul := model.UserLog{
		UserID:     2,
		ChangeType: "type",
		IP:         "ip",
		OldValue:   "old value",
		Value:      "value",
		Active:     true,
	}

	rep := NewUserLogs(transactor)
	err := rep.CreateWithoutToken(ctx, &ul)
	require.NoError(t, err)

	var usl model.UserLog
	err = dbPool.QueryRow(ctx, `SELECT user_id, change_type,old_value, value, ip, active FROM emcd.user_logs
        WHERE user_id = $1`, ul.UserID).Scan(&usl.UserID, &usl.ChangeType, &usl.OldValue, &usl.Value, &usl.IP, &usl.Active)
	require.NoError(t, err)
	require.Equal(t, ul, usl)
}

func TestUserLogs_DeactivateByType(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	ul := model.UserLog{
		ID:         1,
		UserID:     2,
		ChangeType: "type",
		IP:         "ip",
		Token:      "token",
		OldValue:   "old value",
		Value:      "value",
		Active:     true,
	}

	rep := NewUserLogs(transactor)
	err := rep.Create(ctx, &ul)
	require.NoError(t, err)

	err = rep.DeactivateByType(ctx, int(ul.UserID), ul.ChangeType)
	require.NoError(t, err)

	val := true
	err = dbPool.QueryRow(ctx, `SELECT active FROM emcd.user_logs WHERE user_id = $1`, ul.UserID).Scan(&val)
	require.NoError(t, err)
	require.False(t, val)
}
