package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

func TestProfile_Create(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			NewRefID:     uuid.New(),
		},
	}
	err := profileRepo.Create(ctx, &expected)
	require.NoError(t, err)
	var (
		u   model.User
		pwd []byte
	)
	err = dbPool.QueryRow(ctx, `SELECT id,username,ref_id,email,password,created_at,whitelabel_id,api_key,new_ref_id
		FROM users where id = $1`, expected.User.ID).Scan(&u.ID, &u.Username, &u.RefID, &u.Email, &pwd,
		&u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.NewRefID)
	require.NoError(t, err)
	u.Password = string(pwd)
	require.Equal(t, expected.User, &u)
}

func TestProfile_SaveInsert(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			RefID:        312,
			ParentID:     uuid.New(),
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			AppleID:      uuid.New().String(),
			IsActive:     true,
			NewRefID:     uuid.New(),
			Language:     "en",
			PoolType:     "emcd",
		},
	}
	err := profileRepo.SaveV3(ctx, &expected)
	require.NoError(t, err)
	var (
		u   model.User
		pwd []byte
	)

	err = dbPool.QueryRow(ctx, `SELECT id,username,ref_id,email,password,created_at,whitelabel_id,api_key,is_active,apple_id,parent_id,new_ref_id,pool_type,language
		FROM users where id = $1`, expected.User.ID).Scan(&u.ID, &u.Username, &u.RefID, &u.Email, &pwd,
		&u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.IsActive, &u.AppleID, &u.ParentID, &u.NewRefID, &u.PoolType, &u.Language)
	require.NoError(t, err)
	u.Password = string(pwd)
	require.Equal(t, expected.User, &u)
}

func TestProfile_SaveUpdate(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			RefID:        312,
			ParentID:     uuid.New(),
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			AppleID:      uuid.New().String(),
			IsActive:     true,
			NewRefID:     uuid.New(),
			Language:     "ru",
		},
	}
	err := profileRepo.SaveV3(ctx, &expected)
	require.NoError(t, err)
	expected.User.Password = "new password"
	expected.User.Username = "new username"
	expected.User.RefID = 666
	expected.User.ApiKey = "new apikey"
	err = profileRepo.SaveV3(ctx, &expected)
	require.NoError(t, err)
	var (
		u   model.User
		pwd []byte
	)
	err = dbPool.QueryRow(ctx, `SELECT id,username,ref_id,email,password,created_at,whitelabel_id,api_key,is_active,apple_id,parent_id,new_ref_id,language
		FROM users where id = $1`, expected.User.ID).Scan(&u.ID, &u.Username, &u.RefID, &u.Email, &pwd,
		&u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.IsActive, &u.AppleID, &u.ParentID, &u.NewRefID, &u.Language)
	require.NoError(t, err)
	u.Password = string(pwd)
	require.Equal(t, expected.User, &u)
}

func TestProfile_UpdatePassword(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			Vip:          true,
			SegmentID:    234,
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			ParentID:     uuid.New(),
			IsActive:     true,
		},
	}
	err := profileRepo.SaveV3(ctx, &expected)
	require.NoError(t, err)
	newPassword := "new password 228"
	err = profileRepo.UpdatePassword(ctx, "email", newPassword)
	require.NoError(t, err)
	expected.User.Password = newPassword
	oldID := 123
	query := `INSERT INTO emcd.users (id, new_id, username,api_key) VALUES ($1,$2,$3,$4)`
	_, err = dbPool.Exec(ctx, query, oldID, expected.User.ID, "uuu", "apiapiapi")
	require.NoError(t, err)
	query = `INSERT INTO emcd.vip_users (user_id) VALUES ($1)`
	_, err = dbPool.Exec(ctx, query, oldID)
	require.NoError(t, err)
	query = `INSERT INTO histories.segment_userids (user_id,segment_id) VALUES ($1,$2)`
	_, err = dbPool.Exec(ctx, query, oldID, expected.User.SegmentID)
	require.NoError(t, err)
	u, count, err := profileRepo.GetAllUsers(ctx, 0, 22, "username", true, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(u))
	require.Equal(t, 1, count)
	require.Equal(t, *expected.User, *u[0])
}

func TestProfile_GetUserByApiKey(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	oldID := 15863
	e := model.User{
		ID:        uuid.New(),
		SegmentID: 998877,
		Username:  "username",
		RefID:     312,
		Email:     "email123",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey",
		Vip:       true,
		ParentID:  uuid.New(),
		IsActive:  true,
		NewRefID:  uuid.New(),
	}
	err := insertUser(ctx, &e, oldID)
	require.NoError(t, err)
	u, err := profileRepo.GetUserByApiKey(ctx, e.ApiKey)
	require.NoError(t, err)
	require.Equal(t, &e, u)
}

func TestProfile_GetAllUsers(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	oldID := 15863
	e := model.User{
		ID:        uuid.New(),
		SegmentID: 998877,
		Username:  "username",
		RefID:     312,
		Email:     "email123",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey",
		Vip:       true,
		ParentID:  uuid.New(),
		IsActive:  true,
		NewRefID:  uuid.New(),
	}
	err := insertUser(ctx, &e, oldID)
	require.NoError(t, err)

	oldID2 := 64761
	e2 := model.User{
		ID:        uuid.New(),
		SegmentID: 461313,
		Username:  "ususus",
		RefID:     846,
		Email:     "email456",
		Password:  "password123!",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey2",
		Vip:       true,
		ParentID:  uuid.New(),
		IsActive:  true,
		NewRefID:  uuid.New(),
	}
	err = insertUser(ctx, &e2, oldID2)
	require.NoError(t, err)
	users, count, err := profileRepo.GetAllUsers(ctx, 0, 10, "email", true, "")
	require.NoError(t, err)
	require.Equal(t, len(users), 2)
	require.Equal(t, 2, count)

	e.PoolType = "emcd"  // default
	e2.PoolType = "emcd" // default
	require.Equal(t, &e, users[0])
	require.Equal(t, &e2, users[1])
}

func TestProfile_GetAllUsersWithoutSegmentID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	e := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     312,
		Email:     "email123",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey",
		ParentID:  uuid.New(),
		IsActive:  true,
		NewRefID:  uuid.New(),
	}
	query := `INSERT INTO users (id,username,ref_id,email,password,created_at,whitelabel_id,api_key,parent_id,is_active,new_ref_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := dbPool.Exec(ctx, query, e.ID, e.Username, e.RefID, e.Email, e.Password, e.CreatedAt, uuid.Nil, e.ApiKey, e.ParentID, e.IsActive, e.NewRefID)
	require.NoError(t, err)
	users, count, err := profileRepo.GetAllUsers(ctx, 0, 10, "email", true, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(users))
	require.Equal(t, 1, count)

	e.PoolType = "emcd" // default
	require.Equal(t, &e, users[0])
}

func insertUser(ctx context.Context, e *model.User, oldID int) error {
	query := `INSERT INTO users (id,username,ref_id,email,password,created_at,whitelabel_id,api_key,parent_id,is_active,new_ref_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := dbPool.Exec(ctx, query, e.ID, e.Username, e.RefID, e.Email, e.Password, e.CreatedAt, uuid.Nil, e.ApiKey, e.ParentID, e.IsActive, e.NewRefID)
	if err != nil {
		return err
	}
	query = `INSERT INTO emcd.users (id, new_id, username,api_key) VALUES ($1,$2,$3,$4)`
	_, err = dbPool.Exec(ctx, query, oldID, e.ID, "uuu", "apiapiapi")
	if err != nil {
		return err
	}
	query = `INSERT INTO emcd.vip_users (user_id) VALUES ($1)`
	_, err = dbPool.Exec(ctx, query, oldID)
	if err != nil {
		return err
	}
	query = `INSERT INTO histories.segment_userids (user_id,segment_id) VALUES ($1,$2)`
	_, err = dbPool.Exec(ctx, query, oldID, e.SegmentID)
	if err != nil {
		return err
	}
	return nil
}

func TestProfile_GetAllUsersByWlID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	oldID := 15863
	e := model.User{
		ID:        uuid.New(),
		SegmentID: 998877,
		Username:  "username",
		RefID:     312,
		Email:     "email123",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey",
		Vip:       true,
		ParentID:  uuid.New(),
		IsActive:  true,
		NewRefID:  uuid.New(),
	}
	err := insertUser(ctx, &e, oldID)
	require.NoError(t, err)

	oldID2 := 64761
	e2 := model.User{
		ID:        uuid.New(),
		SegmentID: 461313,
		Username:  "ususus",
		RefID:     846,
		Email:     "email456",
		Password:  "password123!",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey2",
		Vip:       true,
		ParentID:  uuid.New(),
		IsActive:  true,
		NewRefID:  uuid.New(),
	}
	err = insertUser(ctx, &e2, oldID2)
	require.NoError(t, err)
	users, count, err := profileRepo.GetAllUsersByWlID(ctx, 0, 10, "email", true, "", uuid.Nil)
	require.NoError(t, err)
	require.Equal(t, len(users), 2)
	require.Equal(t, 2, count)

	e.PoolType = "emcd"  // default
	e2.PoolType = "emcd" // default
	require.Equal(t, &e, users[0])
	require.Equal(t, &e2, users[1])
}

func TestProfile_Get(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			Vip:          true,
			SegmentID:    666,
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			ParentID:     uuid.New(),
			IsActive:     true,
			NewRefID:     uuid.New(),
			Language:     "ru",
		},
	}
	err := profileRepo.Create(ctx, &expected)
	require.NoError(t, err)
	query := `UPDATE users SET parent_id=$1,is_active=$2,language=$3 WHERE id = $4`
	_, err = dbPool.Exec(ctx, query, expected.User.ParentID, expected.User.IsActive, expected.User.Language, expected.User.ID)
	require.NoError(t, err)
	oldID := 124
	query = `INSERT INTO emcd.users (id, new_id, username,api_key) VALUES ($1,$2,$3,$4)`
	_, err = dbPool.Exec(ctx, query, oldID, expected.User.ID, "uuu", "apiapiapi")
	require.NoError(t, err)
	query = `INSERT INTO emcd.vip_users (user_id) VALUES ($1)`
	_, err = dbPool.Exec(ctx, query, oldID)
	require.NoError(t, err)
	query = `INSERT INTO histories.segment_userids (user_id,segment_id) VALUES ($1,$2)`
	_, err = dbPool.Exec(ctx, query, oldID, expected.User.SegmentID)
	require.NoError(t, err)
	actual, err := profileRepo.GetByUserID(ctx, expected.User.ID)
	require.NoError(t, err)
	require.Equal(t, *expected.User, *actual.User)

	profiles, err := profileRepo.GetByUsernames(ctx, []string{"username"})
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	require.Equal(t, expected.User, profiles[0].User)
}

func TestProfile_GetUserIsActive(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	u := model.User{
		ID:           uuid.New(),
		Username:     "username",
		Vip:          true,
		SegmentID:    666,
		RefID:        312,
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		ParentID:     uuid.Nil,
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	err := insertUser(ctx, &u, 100)
	require.NoError(t, err)
	active, err := profileRepo.GetUserIsActive(ctx, u.ID)
	require.NoError(t, err)
	require.True(t, active)
}

func TestProfile_GetUserIsActiveByParentID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	parent := model.User{
		ID:           uuid.New(),
		Username:     "username1",
		Vip:          true,
		SegmentID:    666,
		RefID:        312,
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		ParentID:     uuid.Nil,
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	child := model.User{
		ID:           uuid.New(),
		Username:     "username2",
		Vip:          true,
		SegmentID:    111,
		RefID:        2312,
		Email:        "email1",
		Password:     "password123",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey12123",
		ParentID:     parent.ID,
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	err := insertUser(ctx, &parent, 100)
	require.NoError(t, err)
	err = insertUser(ctx, &child, 101)
	require.NoError(t, err)
	active, err := profileRepo.GetUserIsActive(ctx, child.ID)
	require.NoError(t, err)
	require.True(t, active)
}

func TestProfile_GetUserIsActiveByParentID2(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	parent := model.User{
		ID:           uuid.New(),
		Username:     "username1",
		Vip:          true,
		SegmentID:    666,
		RefID:        312,
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		ParentID:     uuid.Nil,
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	child := model.User{
		ID:           uuid.New(),
		Username:     "username2",
		Vip:          true,
		SegmentID:    111,
		RefID:        2312,
		Email:        "email1",
		Password:     "password123",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey12123",
		ParentID:     parent.ID,
		IsActive:     false,
		NewRefID:     uuid.New(),
	}
	err := insertUser(ctx, &parent, 100)
	require.NoError(t, err)
	err = insertUser(ctx, &child, 101)
	require.NoError(t, err)
	active, err := profileRepo.GetUserIsActive(ctx, child.ID)
	require.NoError(t, err)
	require.True(t, active)
}

func TestProfile_GetUserIsActiveByParentID3(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	parent := model.User{
		ID:           uuid.New(),
		Username:     "username1",
		Vip:          true,
		SegmentID:    666,
		RefID:        312,
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		ParentID:     uuid.Nil,
		IsActive:     false,
		NewRefID:     uuid.New(),
	}
	child := model.User{
		ID:           uuid.New(),
		Username:     "username2",
		Vip:          true,
		SegmentID:    111,
		RefID:        2312,
		Email:        "email1",
		Password:     "password123",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey12123",
		ParentID:     parent.ID,
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	err := insertUser(ctx, &parent, 100)
	require.NoError(t, err)
	err = insertUser(ctx, &child, 101)
	require.NoError(t, err)
	active, err := profileRepo.GetUserIsActive(ctx, child.ID)
	require.NoError(t, err)
	require.False(t, active)
}

func TestProfile_GetUserIsActiveByParentID4(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	parent := model.User{
		ID:           uuid.New(),
		Username:     "username1",
		Vip:          true,
		SegmentID:    666,
		RefID:        312,
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		ParentID:     uuid.Nil,
		IsActive:     false,
		NewRefID:     uuid.New(),
	}
	child := model.User{
		ID:           uuid.New(),
		Username:     "username2",
		Vip:          true,
		SegmentID:    111,
		RefID:        2312,
		Email:        "email1",
		Password:     "password123",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey12123",
		ParentID:     parent.ID,
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	child2 := model.User{
		ID:           uuid.New(),
		Username:     "username23",
		Vip:          true,
		SegmentID:    11112,
		RefID:        231212,
		Email:        "email111",
		Password:     "password12312",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey121231212",
		ParentID:     parent.ID,
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	err := insertUser(ctx, &parent, 100)
	require.NoError(t, err)
	err = insertUser(ctx, &child, 101)
	require.NoError(t, err)
	err = insertUser(ctx, &child2, 102)
	require.NoError(t, err)
	active, err := profileRepo.GetUserIsActive(ctx, child2.ID)
	require.NoError(t, err)
	require.False(t, active)
}

func TestProfile_ExistAppleID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	ex := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			Vip:          false,
			SegmentID:    666,
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			AppleID:      "apple_id",
		},
	}

	err := profileRepo.SaveV3(ctx, &ex)
	require.NoError(t, err)

	exist, em, err := profileRepo.ExistAppleID(ctx, ex.User.AppleID)
	require.NoError(t, err)
	require.True(t, exist)
	require.Equal(t, ex.User.Email, em)
}

func TestProfile_UpdateAppleID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	ex := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			Vip:          false,
			SegmentID:    666,
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
		},
	}

	err := profileRepo.SaveV3(ctx, &ex)
	require.NoError(t, err)

	ex.User.AppleID = "apple"

	exist, _, err := profileRepo.ExistAppleID(ctx, ex.User.AppleID)
	require.NoError(t, err)
	require.False(t, exist)

	rowsAffected, err := profileRepo.UpdateAppleID(ctx, ex.User.Email, ex.User.AppleID)
	require.NoError(t, err)
	require.Equal(t, 1, rowsAffected)

	exist, _, err = profileRepo.ExistAppleID(ctx, ex.User.AppleID)
	require.NoError(t, err)
	require.True(t, exist)
}

func TestProfile_UpdateAppleIDNothingUpdated(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)

	rowsAffected, err := profileRepo.UpdateAppleID(ctx, "email", "apple_id")
	require.NoError(t, err)
	require.Equal(t, 0, rowsAffected)
}

func TestProfile_CheckAppleAccountDoesntExist(t *testing.T) {
	ctx := context.Background()
	exist, em, err := profileRepo.ExistAppleID(ctx, "apple_id")
	require.NoError(t, err)
	require.False(t, exist)
	require.Equal(t, "", em)
}

func TestProfile_GetUserByEmailAndWl(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	oldID := 15863
	e := model.User{
		ID:        uuid.New(),
		SegmentID: 998877,
		Username:  "username",
		RefID:     312,
		Email:     "email123",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey",
		ParentID:  uuid.New(),
		IsActive:  true,
		Vip:       true,
		NewRefID:  uuid.New(),
		Suspended: true,
	}
	query := `INSERT INTO users (id,username,ref_id,email,password,created_at,whitelabel_id,api_key,parent_id,is_active,new_ref_id,suspended) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := dbPool.Exec(ctx, query, e.ID, e.Username, e.RefID, e.Email, e.Password, e.CreatedAt, uuid.Nil, e.ApiKey, e.ParentID, e.IsActive, e.NewRefID, time.Now())
	require.NoError(t, err)
	query = `INSERT INTO emcd.users (id, new_id, username,api_key) VALUES ($1,$2,$3,$4)`
	_, err = dbPool.Exec(ctx, query, oldID, e.ID, "uuu", "apiapiapi")
	require.NoError(t, err)
	query = `INSERT INTO emcd.vip_users (user_id) VALUES ($1)`
	_, err = dbPool.Exec(ctx, query, oldID)
	require.NoError(t, err)
	query = `INSERT INTO histories.segment_userids (user_id,segment_id) VALUES ($1,$2)`
	_, err = dbPool.Exec(ctx, query, oldID, e.SegmentID)
	require.NoError(t, err)
	u, err := profileRepo.GetUserByEmailAndWl(ctx, e.Email, uuid.Nil)
	require.NoError(t, err)
	require.Equal(t, &e, u)
}

func TestProfile_GetUserByEmailAndWlWithoutSegment(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	oldID := 15863
	e := model.User{
		ID: uuid.New(),
		// SegmentID: 998877,
		Username:  "username",
		RefID:     312,
		Email:     "email123",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		ApiKey:    "apiKey",
		ParentID:  uuid.New(),
		IsActive:  true,
		Vip:       true,
		NewRefID:  uuid.New(),
		Suspended: true,
	}
	query := `INSERT INTO users (id,username,ref_id,email,password,created_at,whitelabel_id,api_key,parent_id,is_active,new_ref_id,suspended) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := dbPool.Exec(ctx, query, e.ID, e.Username, e.RefID, e.Email, e.Password, e.CreatedAt, uuid.Nil, e.ApiKey, e.ParentID, e.IsActive, e.NewRefID, time.Now())
	require.NoError(t, err)
	query = `INSERT INTO emcd.users (id, new_id, username,api_key) VALUES ($1,$2,$3,$4)`
	_, err = dbPool.Exec(ctx, query, oldID, e.ID, "uuu", "apiapiapi")
	require.NoError(t, err)
	query = `INSERT INTO emcd.vip_users (user_id) VALUES ($1)`
	_, err = dbPool.Exec(ctx, query, oldID)
	require.NoError(t, err)
	u, err := profileRepo.GetUserByEmailAndWl(ctx, e.Email, uuid.Nil)
	require.NoError(t, err)
	require.NotEqual(t, 0, u.SegmentID)
	e.SegmentID = u.SegmentID
	require.Equal(t, &e, u)
}

func TestProfile_RelatedAccountsMainSub(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	e := model.User{
		ID:       uuid.New(),
		Username: "username",
		Email:    "email123",
		Password: "password",
		ApiKey:   "apiKey",
		ParentID: uuid.Nil,
	}

	e2 := model.User{
		ID:       uuid.New(),
		Username: "username-2",
		Email:    "email123-2",
		Password: "password",
		ApiKey:   "apiKey",
		ParentID: e.ID,
	}
	fmt.Println(e, e2)
	query := `INSERT INTO users (id,username,email,password,api_key,parent_id) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := dbPool.Exec(ctx, query, e.ID, e.Username, e.Email, e.Password, e.ApiKey, e.ParentID)
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, query, e2.ID, e2.Username, e2.Email, e2.Password, e2.ApiKey, e2.ParentID)
	require.NoError(t, err)

	related, err := profileRepo.RelatedUsers(ctx, e.ID, e2.ID)
	require.NoError(t, err)
	require.True(t, related)
}

func TestProfile_RelatedAccountsSubSub(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	parentID := uuid.New()
	e := model.User{
		ID:       uuid.New(),
		Username: "username",
		Email:    "email123",
		Password: "password",
		ApiKey:   "apiKey",
		ParentID: parentID,
	}

	e2 := model.User{
		ID:       uuid.New(),
		Username: "username-2",
		Email:    "email123-2",
		Password: "password",
		ApiKey:   "apiKey",
		ParentID: parentID,
	}
	fmt.Println(e, e2)
	query := `INSERT INTO users (id,username,email,password,api_key,parent_id) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := dbPool.Exec(ctx, query, e.ID, e.Username, e.Email, e.Password, e.ApiKey, e.ParentID)
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, query, e2.ID, e2.Username, e2.Email, e2.Password, e2.ApiKey, e2.ParentID)
	require.NoError(t, err)

	related, err := profileRepo.RelatedUsers(ctx, e.ID, e2.ID)
	require.NoError(t, err)
	require.True(t, related)
}

func TestProfile_RelatedAccountsSubMain(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	e := model.User{
		ID:       uuid.New(),
		Username: "username",
		Email:    "email123",
		Password: "password",
		ApiKey:   "apiKey",
		ParentID: uuid.Nil,
	}

	e2 := model.User{
		ID:       uuid.New(),
		Username: "username-2",
		Email:    "email123-2",
		Password: "password",
		ApiKey:   "apiKey",
		ParentID: e.ID,
	}
	fmt.Println(e, e2)
	query := `INSERT INTO users (id,username,email,password,api_key,parent_id) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := dbPool.Exec(ctx, query, e.ID, e.Username, e.Email, e.Password, e.ApiKey, e.ParentID)
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, query, e2.ID, e2.Username, e2.Email, e2.Password, e2.ApiKey, e2.ParentID)
	require.NoError(t, err)

	related, err := profileRepo.RelatedUsers(ctx, e2.ID, e.ID)
	require.NoError(t, err)
	require.True(t, related)
}

func TestAccounts_GetAllSubUsers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)

	apiKey := "1"
	parentID := uuid.New()
	parentOldID := 100
	_, err := dbPool.Exec(ctx, "INSERT INTO emcd.users (id, new_id, username, api_key) VALUES ($1,$2,$3,$4)", parentOldID, parentID, "parent", apiKey)
	require.NoError(t, err)

	u1 := model.User{
		Username: "username1",
		OldID:    11,
		ID:       uuid.New(),
	}
	u2 := model.User{
		Username: "username2",
		OldID:    22,
		ID:       uuid.New(),
	}
	u3 := model.User{
		Username: "username3",
		OldID:    33,
		ID:       uuid.New(),
	}
	for _, u := range []*model.User{
		&u1, &u2, &u3,
	} {
		_, err := dbPool.Exec(ctx, "INSERT INTO emcd.users (id, new_id, username, parent_id, api_key) VALUES ($1,$2,$3,$4,$5)",
			u.OldID, u.ID, u.Username, parentOldID, apiKey)
		require.NoError(t, err)
	}

	subUsers, err := profileRepo.GetAllSubUsers(ctx, parentID)
	require.NoError(t, err)
	require.Equal(t, 3, len(subUsers))
	var count int
	for _, u := range subUsers {
		switch u.ID {
		case u1.ID:
			require.Equal(t, u1.Username, u.Username)
			count++
		case u2.ID:
			require.Equal(t, u2.Username, u.Username)
			count++
		case u3.ID:
			require.Equal(t, u3.Username, u.Username)
			count++
		}
	}
	require.Equal(t, 3, count)
}

func TestProfile_GetAllUserIDsByUsername(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)

	u1 := model.User{
		Username: "username1",
		ID:       uuid.New(),
	}
	u2 := model.User{
		Username: "username2",
		ID:       uuid.New(),
	}
	u3 := model.User{
		Username: "username3",
		ID:       uuid.New(),
	}
	for _, u := range []*model.User{
		&u1, &u2, &u3,
	} {
		_, err := dbPool.Exec(ctx, "INSERT INTO users (id,username) VALUES ($1,$2)", u.ID, u.Username)
		require.NoError(t, err)
	}
	m, err := profileRepo.GetAllUserIDsByUsername(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, len(m))
	require.Equal(t, u1.ID, m[u1.Username])
	require.Equal(t, u2.ID, m[u2.Username])
	require.Equal(t, u3.ID, m[u3.Username])
}

func TestProfile_GetReferrals(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)

	now := time.Now()
	u1 := model.User{
		Username:  "username1",
		Email:     "email1",
		ID:        uuid.New(),
		CreatedAt: now,
	}
	u2 := model.User{
		Username:  "username2",
		ID:        uuid.New(),
		Email:     "email2",
		NewRefID:  u1.ID,
		CreatedAt: now,
	}
	u3 := model.User{
		Username:  "username3",
		ID:        uuid.New(),
		Email:     "email3",
		NewRefID:  u1.ID,
		CreatedAt: now,
	}
	other := model.User{
		Username:  "username4",
		ID:        uuid.New(),
		Email:     "email4",
		NewRefID:  uuid.New(),
		CreatedAt: now,
	}
	sortField := "created_at"

	for _, u := range []*model.User{
		&u1, &u2, &u3, &other,
	} {
		_, err := dbPool.Exec(ctx, "INSERT INTO users (id,username,email,new_ref_id, created_at) VALUES ($1,$2,$3,$4,$5)", u.ID, u.Username, u.Email, u.NewRefID, u.CreatedAt)
		require.NoError(t, err)
	}

	users, count, err := profileRepo.GetReferrals(ctx, u1.ID, 0, 5, sortField, false)
	require.NoError(t, err)
	require.Equal(t, 2, count)
	require.Equal(t, 2, len(users))
	// u2.ParentID, u3.ParentID = uuid.Nil, uuid.Nil
	require.Equal(t, u2.ID.String(), users[0].ID.String())
	require.Equal(t, u3.ID.String(), users[1].ID.String())
}

func TestProfile_GetReferrals2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)

	u1 := model.User{
		Username:  "username1",
		Email:     "email1",
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
	u2 := model.User{
		Username:  "username2",
		ID:        uuid.New(),
		Email:     "email2",
		NewRefID:  u1.ID,
		CreatedAt: time.Now(),
	}
	u3 := model.User{
		Username:  "username3",
		ID:        uuid.New(),
		Email:     "email3",
		NewRefID:  u1.ID,
		CreatedAt: time.Now(),
	}
	u4 := model.User{
		Username:  "username4",
		ID:        uuid.New(),
		Email:     "email4",
		NewRefID:  u1.ID,
		CreatedAt: time.Now(),
	}
	u5 := model.User{
		Username:  "username5",
		ID:        uuid.New(),
		Email:     "email5",
		NewRefID:  u1.ID,
		CreatedAt: time.Now(),
	}
	other := model.User{
		Username:  "username6",
		ID:        uuid.New(),
		Email:     "email6",
		NewRefID:  uuid.New(),
		CreatedAt: time.Now(),
	}
	sortField := "created_at"
	for _, u := range []*model.User{
		&u1, &u2, &u3, &u4, &u5, &other,
	} {
		_, err := dbPool.Exec(ctx, "INSERT INTO users (id,username,email,new_ref_id,created_at) VALUES ($1,$2,$3,$4,$5)", u.ID, u.Username, u.Email, u.NewRefID, u.CreatedAt)
		require.NoError(t, err)
	}

	users, count, err := profileRepo.GetReferrals(ctx, u1.ID, 2, 1, sortField, true)
	require.NoError(t, err)
	require.Equal(t, 4, count)
	require.Equal(t, 1, len(users))
	u4.ParentID = uuid.Nil
	require.Equal(t, u4.ID, users[0].ID)
}

func TestOldUsers_SetGetSuspended(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)

	id := uuid.New()
	id2 := uuid.New()

	_, err := dbPool.Exec(ctx, `INSERT INTO users (id,username,api_key) VALUES ($1,$2,$3)`, id, "user", "api")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, `INSERT INTO users (id,username,api_key) VALUES ($1,$2,$3)`, id2, "user2", "api2")
	require.NoError(t, err)

	err = profileRepo.SetSuspended(ctx, []uuid.UUID{id, id2}, true)
	require.NoError(t, err)

	sus, err := profileRepo.GetSuspended(ctx, id)
	require.NoError(t, err)
	require.True(t, sus)
	sus, err = profileRepo.GetSuspended(ctx, id2)
	require.NoError(t, err)
	require.True(t, sus)
}

func TestOldUsers_NullSuspended(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)

	id := uuid.New()

	_, err := dbPool.Exec(ctx, `INSERT INTO users (id,username,api_key) VALUES ($1,$2,$3)`, id, "user", "api")
	require.NoError(t, err)

	sus, err := profileRepo.GetSuspended(ctx, id)
	require.NoError(t, err)
	require.False(t, sus)
}

func truncateAllProfile(t *testing.T, ctx context.Context) {
	_, err := dbPool.Exec(ctx, "TRUNCATE users CASCADE")
	require.NoError(t, err)
}

func TestProfile_GetUsernamesByIDs(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)

	query := `INSERT INTO public.users (id, username) values ($1,$2)`
	users := make([]struct {
		ID       uuid.UUID
		Username string
	}, 0)
	ids := make([]uuid.UUID, 0)
	for i := 0; i < 10; i++ {
		users = append(users, struct {
			ID       uuid.UUID
			Username string
		}{ID: uuid.New(), Username: fmt.Sprintf("username%d", i)})
		ids = append(ids, users[i].ID)
		_, err := dbPool.Exec(ctx, query, users[i].ID, users[i].Username)
		require.NoError(t, err)
	}
	_, err := dbPool.Exec(ctx, query, uuid.New(), "1234")
	require.NoError(t, err)
	usernames, err := profileRepo.GetUsernamesByIDs(ctx, ids)
	require.NoError(t, err)

	for i := range users {
		require.Equal(t, users[i].Username, usernames[users[i].ID])
	}
}

func TestProfile_GetEmailsByIDs(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)

	query := `INSERT INTO public.users (id, email) values ($1,$2)`
	users := make([]struct {
		ID    uuid.UUID
		Email string
	}, 0)
	ids := make([]uuid.UUID, 0)
	for i := 0; i < 10; i++ {
		users = append(users, struct {
			ID    uuid.UUID
			Email string
		}{ID: uuid.New(), Email: fmt.Sprintf("email%d", i)})
		ids = append(ids, users[i].ID)
		_, err := dbPool.Exec(ctx, query, users[i].ID, users[i].Email)
		require.NoError(t, err)
	}
	_, err := dbPool.Exec(ctx, query, uuid.New(), "1234")
	require.NoError(t, err)
	emails, err := profileRepo.GetEmailsByIDs(ctx, ids)
	require.NoError(t, err)

	for i := range users {
		require.Equal(t, users[i].Email, emails[users[i].ID])
	}
}

func TestProfile_SaveUserInsert(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)

	expected22 := &model.User{
		ID:           uuid.New(),
		Username:     "username",
		RefID:        312,
		ParentID:     uuid.New(),
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		AppleID:      uuid.New().String(),
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	err := profileRepo.SaveUser(ctx, expected22)
	require.NoError(t, err)

	var (
		u   model.User
		pwd []byte
	)

	err = dbPool.QueryRow(ctx, `SELECT id,username,ref_id,email,password,created_at,whitelabel_id,api_key,is_active,apple_id,parent_id,new_ref_id
		FROM users where id = $1`, expected22.ID).Scan(&u.ID, &u.Username, &u.RefID, &u.Email, &pwd,
		&u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.IsActive, &u.AppleID, &u.ParentID, &u.NewRefID)
	require.NoError(t, err)
	u.Password = string(pwd)
	require.Equal(t, expected22, &u)
}

func TestProfile_SaveUserUpdate(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)

	expected22 := &model.User{
		ID:           uuid.New(),
		Username:     "username",
		RefID:        312,
		ParentID:     uuid.New(),
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		AppleID:      uuid.New().String(),
		IsActive:     true,
		NewRefID:     uuid.New(),
	}
	err := profileRepo.SaveUser(ctx, expected22)
	require.NoError(t, err)

	expected22.Password = "new password"
	expected22.Username = "new username"
	expected22.RefID = 666
	expected22.ApiKey = "new apikey"
	err = profileRepo.SaveUser(ctx, expected22)
	require.NoError(t, err)

	var (
		u   model.User
		pwd []byte
	)

	err = dbPool.QueryRow(ctx, `SELECT id,username,ref_id,email,password,created_at,whitelabel_id,api_key,is_active,apple_id,parent_id,new_ref_id
		FROM users where id = $1`, expected22.ID).Scan(&u.ID, &u.Username, &u.RefID, &u.Email, &pwd,
		&u.CreatedAt, &u.WhiteLabelID, &u.ApiKey, &u.IsActive, &u.AppleID, &u.ParentID, &u.NewRefID)
	require.NoError(t, err)
	u.Password = string(pwd)
	require.Equal(t, expected22, &u)
}

func TestProfile_UpdateRefID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			NewRefID:     uuid.New(),
		},
	}
	err := profileRepo.Create(ctx, &expected)
	require.NoError(t, err)

	var newRefID int32 = 999
	err = profileRepo.UpdateRefID(ctx, expected.User.ID, newRefID, expected.User.NewRefID)
	require.NoError(t, err)

	var actualRefID int32
	var newRefUUID uuid.UUID
	err = dbPool.QueryRow(ctx, `SELECT ref_id,new_ref_id FROM users where id = $1`, expected.User.ID).
		Scan(&actualRefID, &newRefUUID)
	require.NoError(t, err)
	require.Equal(t, newRefID, actualRefID)
	require.Equal(t, expected.User.NewRefID, newRefUUID)
}

func TestProfile_ReferralLinkGenerated(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			NewRefID:     uuid.New(),
		},
	}
	err := profileRepo.Create(ctx, &expected)
	require.NoError(t, err)

	err = profileRepo.SetFlagReferralLinkGenerated(ctx, expected.User.ID, true)
	require.NoError(t, err)
	flg, err := profileRepo.GetFlagReferralLinkGenerated(ctx, expected.User.ID)
	require.NoError(t, err)
	require.Equal(t, flg, true)
}

func TestProfile_GetAllReferrals(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)

	referrer := &model.User{
		ID:       uuid.New(),
		NewRefID: uuid.New(),
		Email:    "email",
		Username: "referrer",
	}
	referral1 := &model.User{
		ID:       uuid.New(),
		NewRefID: referrer.ID,
		Email:    "email1",
		Username: "referral1",
	}
	referral2 := &model.User{
		ID:       uuid.New(),
		NewRefID: referrer.ID,
		Email:    "email2",
		Username: "referral2",
	}
	err := insertUser(ctx, referrer, 1)
	require.NoError(t, err)

	err = insertUser(ctx, referral1, 2)
	require.NoError(t, err)

	err = insertUser(ctx, referral2, 3)
	require.NoError(t, err)

	referrals, usernames, err := profileRepo.GetAllReferrals(ctx, referrer.ID, &referral1.Username, "username", true, 1)
	require.NoError(t, err)

	require.Len(t, referrals, 1)
	require.Len(t, usernames, 1)
	require.Equal(t, referral2.Username, usernames[0])

	require.Equal(t, model.Referral{
		UserID:   referral2.ID,
		Email:    referral2.Email,
		Username: referral2.Username,
	}, *referrals[usernames[0]])
}

func TestProfile_GetAllReferralsNilRetrieveAfter(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)

	referrer := &model.User{
		ID:       uuid.New(),
		NewRefID: uuid.New(),
		Email:    "email",
		Username: "referrer",
	}
	referral1 := &model.User{
		ID:       uuid.New(),
		NewRefID: referrer.ID,
		Email:    "email1",
		Username: "referral1",
	}
	referral2 := &model.User{
		ID:       uuid.New(),
		NewRefID: referrer.ID,
		Email:    "email2",
		Username: "referral2",
	}
	err := insertUser(ctx, referrer, 1)
	require.NoError(t, err)

	err = insertUser(ctx, referral1, 2)
	require.NoError(t, err)

	err = insertUser(ctx, referral2, 3)
	require.NoError(t, err)

	referrals, usernames, err := profileRepo.GetAllReferrals(ctx, referrer.ID, nil, "username", false, 2)
	require.NoError(t, err)

	require.Len(t, referrals, 2)
	require.Len(t, usernames, 2)
	require.Equal(t, referral2.Username, usernames[0])
	require.Equal(t, referral1.Username, usernames[1])

	require.Equal(t, model.Referral{
		UserID:   referral2.ID,
		Email:    referral2.Email,
		Username: referral2.Username,
	}, *referrals[usernames[0]])
	require.Equal(t, model.Referral{
		UserID:   referral1.ID,
		Email:    referral1.Email,
		Username: referral1.Username,
	}, *referrals[usernames[1]])
}

func TestProfile_UpdateUser(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	u := &model.User{
		ID:           uuid.New(),
		Username:     "username",
		RefID:        312,
		ParentID:     uuid.New(),
		Email:        "email",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		AppleID:      uuid.New().String(),
		IsActive:     true,
		NewRefID:     uuid.New(),
		Language:     "ru",
	}
	err := profileRepo.SaveUser(ctx, u)
	require.NoError(t, err)
	err = profileRepo.UpdateUser(ctx, u.ID, func(user *model.User) error {
		user.Email = "new_email"
		return nil
	})
	require.NoError(t, err)

	var actual model.User
	err = dbPool.QueryRow(ctx, `SELECT id,username,ref_id,email,created_at,whitelabel_id,api_key,is_active,apple_id,parent_id,new_ref_id,language
		FROM users where id = $1`, u.ID).Scan(&actual.ID, &actual.Username, &actual.RefID, &actual.Email,
		&actual.CreatedAt, &actual.WhiteLabelID, &actual.ApiKey, &actual.IsActive, &actual.AppleID, &actual.ParentID, &actual.NewRefID, &actual.Language)
	require.NoError(t, err)
	require.Equal(t, u.Email, actual.Email)
}

func TestProfile_UpdateUserEmailExist(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	u1 := &model.User{
		ID:           uuid.New(),
		Username:     "username1",
		RefID:        312,
		ParentID:     uuid.New(),
		Email:        "email1",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		AppleID:      uuid.New().String(),
		IsActive:     true,
		NewRefID:     uuid.New(),
		Language:     "ru",
	}
	u2 := &model.User{
		ID:           uuid.New(),
		Username:     "username2",
		RefID:        312,
		ParentID:     uuid.New(),
		Email:        "email2",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		AppleID:      uuid.New().String(),
		IsActive:     true,
		NewRefID:     uuid.New(),
		Language:     "ru",
	}
	err := profileRepo.SaveUser(ctx, u1)
	require.NoError(t, err)
	err = profileRepo.SaveUser(ctx, u2)
	require.NoError(t, err)
	oldEmail := u2.Email
	err = profileRepo.UpdateUser(ctx, u2.ID, func(user *model.User) error {
		user.Email = u1.Email
		return nil
	})
	require.Error(t, err)

	var actual model.User
	err = dbPool.QueryRow(ctx, `SELECT id,username,ref_id,email,created_at,whitelabel_id,api_key,is_active,apple_id,parent_id,new_ref_id,language
		FROM users where id = $1`, u2.ID).Scan(&actual.ID, &actual.Username, &actual.RefID, &actual.Email,
		&actual.CreatedAt, &actual.WhiteLabelID, &actual.ApiKey, &actual.IsActive, &actual.AppleID, &actual.ParentID, &actual.NewRefID, &actual.Language)
	require.NoError(t, err)
	require.Equal(t, oldEmail, actual.Email)
}

func TestProfile_SetLanguage(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	expected := model.Profile{
		User: &model.User{
			ID:           uuid.New(),
			Username:     "username",
			RefID:        312,
			Email:        "email",
			Password:     "password",
			CreatedAt:    time.Now().Truncate(time.Second).UTC(),
			WhiteLabelID: uuid.New(),
			ApiKey:       "apiKey",
			NewRefID:     uuid.New(),
			Language:     "ru",
		},
	}
	err := profileRepo.Create(ctx, &expected)
	require.NoError(t, err)

	err = profileRepo.SetLanguage(ctx, expected.User.ID, "en")
	require.NoError(t, err)

	var lang string
	err = profileRepo.trx.Runner(ctx).QueryRow(ctx, `SELECT language FROM users WHERE id = $1`, expected.User.ID).Scan(&lang)
	require.NoError(t, err)
	require.Equal(t, "en", lang)
}

func TestProfile_UpdateUserEmptyAppleID(t *testing.T) {
	ctx := context.Background()
	defer truncateAllProfile(t, ctx)
	defer truncateOldUsers(t, ctx)
	u1 := &model.User{
		ID:           uuid.New(),
		Username:     "username1",
		RefID:        312,
		ParentID:     uuid.New(),
		Email:        "email1",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		AppleID:      "",
		IsActive:     true,
		NewRefID:     uuid.New(),
		Language:     "ru",
	}
	u2 := &model.User{
		ID:           uuid.New(),
		Username:     "username2",
		RefID:        312,
		ParentID:     uuid.New(),
		Email:        "email2",
		Password:     "password",
		CreatedAt:    time.Now().Truncate(time.Second).UTC(),
		WhiteLabelID: uuid.New(),
		ApiKey:       "apiKey",
		AppleID:      uuid.New().String(),
		IsActive:     true,
		NewRefID:     uuid.New(),
		Language:     "ru",
	}
	err := profileRepo.SaveUser(ctx, u1)
	require.NoError(t, err)
	err = profileRepo.SaveUser(ctx, u2)
	require.NoError(t, err)
	// u1 already has empty AppleID
	// try to update u2 by empty AppleID
	err = profileRepo.UpdateUser(ctx, u2.ID, func(user *model.User) error {
		user.AppleID = ""
		return nil
	})
	require.NoError(t, err)
}

func getSubUsersProfile(ctx context.Context, userID uuid.UUID) ([]*model.User, error) {
	query := `SELECT id, username, email FROM users WHERE parent_id = $1`
	rows, err := dbPool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}
