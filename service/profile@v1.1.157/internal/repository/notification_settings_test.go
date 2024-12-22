package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

var (
	nsTemplate = model.NotificationSettings{
		UserID:                 uuid.New(),
		IsEmailNotificationsOn: true,
		IsTgNotificationsOn:    true,
		TgID:                   3211,
		IsPushNotificationsOn:  true,
	}

	userTemplateForNotificationSettings = model.User{
		ID:           nsTemplate.UserID,
		Email:        "email",
		Language:     "en",
		WhiteLabelID: uuid.New(),
		Username:     "username",
	}
)

func TestNotificationSettings_GetByUserID(t *testing.T) {
	ctx := context.Background()
	defer truncateNotificationSettings(ctx)

	expectedUser := userTemplateForNotificationSettings
	expectedNs := nsTemplate
	expectedNs.Email = expectedUser.Email
	expectedNs.Language = expectedUser.Language
	expectedNs.WhiteLabelID = expectedUser.WhiteLabelID

	insertUserForNotificationSettings(ctx, &expectedUser)
	err := notificationSettingsRepo.Save(ctx, toChangeable(&expectedNs))
	require.NoError(t, err)

	actualNs, err := notificationSettingsRepo.GetByUserID(ctx, expectedUser.ID)
	require.NoError(t, err)
	require.Equal(t, &expectedNs, actualNs)
}

// TODO Виталик

func TestNotificationSettings_GetByUserIDSub(t *testing.T) {
	ctx := context.Background()
	defer truncateNotificationSettings(ctx)

	parentUser := userTemplateForNotificationSettings
	parentUser.ID = uuid.New()
	parentUser.Email = "TARGET EMAIL!!!"
	parentUser.Username = "parent_username"
	expectedUser := userTemplateForNotificationSettings

	expectedUser.ParentID = parentUser.ID
	expectedNs := nsTemplate

	expectedNs.Email = parentUser.Email
	expectedNs.Language = expectedUser.Language
	expectedNs.WhiteLabelID = expectedUser.WhiteLabelID

	insertUserForNotificationSettings(ctx, &parentUser)
	insertUserForNotificationSettings(ctx, &expectedUser)
	err := notificationSettingsRepo.Save(ctx, toChangeable(&expectedNs))
	require.NoError(t, err)

	actualNs, err := notificationSettingsRepo.GetByUserID(ctx, expectedUser.ID)
	require.NoError(t, err)
	require.Equal(t, &expectedNs, actualNs)
}

func TestNotificationSettings_GetByUsername(t *testing.T) {
	ctx := context.Background()
	defer truncateNotificationSettings(ctx)

	expectedUser := userTemplateForNotificationSettings
	expectedNs := nsTemplate
	expectedNs.Email = expectedUser.Email
	expectedNs.Language = expectedUser.Language
	expectedNs.WhiteLabelID = expectedUser.WhiteLabelID

	insertUserForNotificationSettings(ctx, &expectedUser)
	err := notificationSettingsRepo.Save(ctx, toChangeable(&expectedNs))
	require.NoError(t, err)

	actualNs, err := notificationSettingsRepo.GetByUsername(ctx, expectedUser.Username)
	require.NoError(t, err)
	require.Equal(t, &expectedNs, actualNs)
}

func toChangeable(ns *model.NotificationSettings) *model.ChangeableNotificationSettings {
	return &model.ChangeableNotificationSettings{
		UserID:                 ns.UserID,
		IsEmailNotificationsOn: ns.IsEmailNotificationsOn,
		IsTgNotificationsOn:    ns.IsTgNotificationsOn,
		TgID:                   ns.TgID,
		IsPushNotificationsOn:  ns.IsPushNotificationsOn,
	}
}

func TestNotificationSettings_GetByUsernameSub(t *testing.T) {
	ctx := context.Background()
	defer truncateNotificationSettings(ctx)

	parentUser := userTemplateForNotificationSettings
	parentUser.ID = uuid.New()
	parentUser.Email = "TARGET EMAIL!!!"
	parentUser.Username = "different username"
	expectedUser := userTemplateForNotificationSettings
	expectedUser.ParentID = parentUser.ID
	expectedNs := nsTemplate
	expectedNs.Email = parentUser.Email
	expectedNs.Language = expectedUser.Language
	expectedNs.WhiteLabelID = expectedUser.WhiteLabelID

	insertUserForNotificationSettings(ctx, &parentUser)
	expectedUser.Username += "asd"
	insertUserForNotificationSettings(ctx, &expectedUser)
	err := notificationSettingsRepo.Save(ctx, toChangeable(&expectedNs))
	require.NoError(t, err)

	actualNs, err := notificationSettingsRepo.GetByUsername(ctx, expectedUser.Username)
	require.NoError(t, err)
	require.Equal(t, &expectedNs, actualNs)
}

func TestNotificationSettings_SaveUpdate(t *testing.T) {
	ctx := context.Background()
	defer truncateNotificationSettings(ctx)

	expectedUser := userTemplateForNotificationSettings
	expectedNs := nsTemplate
	expectedNs.Email = expectedUser.Email
	expectedNs.Language = expectedUser.Language
	expectedNs.WhiteLabelID = expectedUser.WhiteLabelID

	insertUserForNotificationSettings(ctx, &expectedUser)
	err := notificationSettingsRepo.Save(ctx, toChangeable(&expectedNs))
	require.NoError(t, err)
	expectedNs.TgID = 4444
	expectedNs.IsEmailNotificationsOn = false
	err = notificationSettingsRepo.Save(ctx, toChangeable(&expectedNs))
	require.NoError(t, err)

	actualNs, err := notificationSettingsRepo.GetByUserID(ctx, expectedUser.ID)
	require.NoError(t, err)
	require.Equal(t, &expectedNs, actualNs)
}

func truncateNotificationSettings(ctx context.Context) {
	query := `TRUNCATE TABLE notification_settings, users cascade`
	_, err := dbPool.Exec(ctx, query)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
}

func insertUserForNotificationSettings(ctx context.Context, user *model.User) {
	query := `INSERT INTO users (id,email,language,parent_id,whitelabel_id,username) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := dbPool.Exec(ctx, query, user.ID, user.Email, user.Language, user.ParentID, user.WhiteLabelID, user.Username)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
}
