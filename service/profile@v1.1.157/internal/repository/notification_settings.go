package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	pgTx "code.emcdtech.com/emcd/sdk/pg"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

var defaultNotificationLang = "ru"

type NotificationSettings interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.NotificationSettings, error)
	GetByUsername(ctx context.Context, username string) (*model.NotificationSettings, error)
	Save(ctx context.Context, s *model.ChangeableNotificationSettings) error
}

type notificationSettings struct {
	trx pgTx.PgxTransactor
}

func NewNotificationSettings(trx pgTx.PgxTransactor) *notificationSettings {
	return &notificationSettings{trx: trx}
}

func (n *notificationSettings) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.NotificationSettings, error) {
	query := `
SELECT u.id, coalesce(u2.email, u.email),coalesce(u.language,$1),is_email_notifications_on,is_tg_notifications_on,
       tg_id,u.whitelabel_id, is_push_notifications_on
  FROM users u LEFT JOIN 
       notification_settings ns ON u.id=ns.user_id LEFT JOIN 
       users u2 ON u.parent_id=u2.id 
 WHERE u.id=$2
 `
	var (
		ns      model.NotificationSettings
		tgOn    sql.NullBool
		pushOn  sql.NullBool
		emailOn sql.NullBool
		tgID    sql.NullInt64
	)
	err := n.trx.Runner(ctx).
		QueryRow(ctx, query, defaultNotificationLang, userID).
		Scan(&ns.UserID, &ns.Email, &ns.Language, &emailOn,
			&tgOn, &tgID, &ns.WhiteLabelID, &pushOn)
	if err != nil {
		return nil, fmt.Errorf("queryRow: scan: %w", err)
	}
	if tgOn.Valid {
		ns.IsTgNotificationsOn = tgOn.Bool
	}
	if pushOn.Valid {
		ns.IsPushNotificationsOn = pushOn.Bool
	}
	if emailOn.Valid {
		ns.IsEmailNotificationsOn = emailOn.Bool
	}
	if tgID.Valid {
		ns.TgID = tgID.Int64
	}
	return &ns, nil
}

func (n *notificationSettings) GetByUsername(ctx context.Context, username string) (*model.NotificationSettings, error) {
	query := `
SELECT u.id, coalesce(u2.email, u.email),coalesce(u.language,$1),is_email_notifications_on,is_tg_notifications_on,
       tg_id,u.whitelabel_id,is_push_notifications_on
  FROM users u LEFT JOIN 
       notification_settings ns ON u.id=ns.user_id LEFT JOIN 
       users u2 ON u.parent_id=u2.id 
 WHERE lower(u.username)=lower($2)
`
	var (
		ns      model.NotificationSettings
		tgOn    sql.NullBool
		pushOn  sql.NullBool
		emailOn sql.NullBool
		tgID    sql.NullInt64
	)
	err := n.trx.Runner(ctx).
		QueryRow(ctx, query, defaultNotificationLang, username).
		Scan(&ns.UserID, &ns.Email, &ns.Language, &emailOn,
			&tgOn, &tgID, &ns.WhiteLabelID, &pushOn)
	if err != nil {
		return nil, fmt.Errorf("queryRow: scan: %w", err)
	}
	if tgOn.Valid {
		ns.IsTgNotificationsOn = tgOn.Bool
	}
	if pushOn.Valid {
		ns.IsPushNotificationsOn = pushOn.Bool
	}
	if emailOn.Valid {
		ns.IsEmailNotificationsOn = emailOn.Bool
	}
	if tgID.Valid {
		ns.TgID = tgID.Int64
	}
	return &ns, nil
}

func (n *notificationSettings) Save(ctx context.Context, s *model.ChangeableNotificationSettings) error {
	query := `
INSERT INTO notification_settings (user_id,is_email_notifications_on,is_tg_notifications_on, tg_id,is_push_notifications_on) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (user_id)
DO UPDATE SET is_email_notifications_on=$2,is_tg_notifications_on=$3, tg_id=$4, is_push_notifications_on=$5`
	_, err := n.trx.Runner(ctx).
		Exec(ctx, query, s.UserID, s.IsEmailNotificationsOn, s.IsTgNotificationsOn, s.TgID, s.IsPushNotificationsOn)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
