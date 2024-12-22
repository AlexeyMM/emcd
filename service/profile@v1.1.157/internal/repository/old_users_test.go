package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/model"
)

func TestOldUsers_Create(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	expected := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	oldID, id, err := oldUsersRepo.Create(ctx, &expected)
	require.NoError(t, err)
	require.Equal(t, expected.ID, id)
	require.Equal(t, oldID != 0, true)
	actual, err := oldUsersRepo.GetByEmail(ctx, expected.Email)
	expected.OldID = actual.OldID
	require.NoError(t, err)
	expected.SegmentID = actual.SegmentID
	require.Equal(t, expected, *actual)
}

func TestOldUsers_SaveInsert(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	expected := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	trx := oldUsersRepo.Begin()
	err := trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			_, err := oldUsersRepo.SaveV2(txCtx, &expected)
			require.NoError(t, err)
			actual, err := oldUsersRepo.GetByEmail(txCtx, expected.Email)
			require.NoError(t, err)
			expected.OldID = actual.OldID
			expected.SegmentID = actual.SegmentID
			require.Equal(t, expected, *actual)
			return err
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)
	require.NoError(t, err)
}

func TestOldUsers_SaveUpdate(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	expected := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	var id1 int32
	trx := oldUsersRepo.Begin()
	err := trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			id, err := oldUsersRepo.SaveV2(txCtx, &expected)
			require.NoError(t, err)
			id1 = id
			return err
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)
	require.NoError(t, err)
	expected.ID = uuid.New()
	expected.Username = "new username"
	expected.RefID = 345
	expected.ApiKey = "new apikey"
	err = trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			id2, err := oldUsersRepo.SaveV2(txCtx, &expected)
			require.NoError(t, err)
			require.Equal(t, id1, id2)
			actual, err := oldUsersRepo.GetByEmail(txCtx, expected.Email)
			require.NoError(t, err)
			expected.OldID = actual.OldID
			expected.SegmentID = actual.SegmentID
			require.Equal(t, expected, *actual)
			require.Equal(t, 1, getAllCount(t, ctx))
			return err
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)

	require.NoError(t, err)
}

func getAllCount(t *testing.T, ctx context.Context) int {
	var res int
	err := dbPool.QueryRow(ctx, `SELECT COUNT(*) FROM emcd.users`).Scan(&res)
	require.NoError(t, err)
	return res
}

func TestOldUsers_UpdatePassword(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	expected := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	_, _, err := oldUsersRepo.Create(ctx, &expected)
	require.NoError(t, err)
	newPwd := "new password"
	err = oldUsersRepo.UpdatePassword(ctx, expected.Email, newPwd)
	require.NoError(t, err)
	expected.Password = newPwd
	actual, err := oldUsersRepo.GetByEmail(ctx, expected.Email)
	require.NoError(t, err)
	expected.OldID = actual.OldID
	expected.SegmentID = actual.SegmentID
	require.Equal(t, expected, *actual)
}

func TestOldUsers_GetUserByUUID(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	newID := uuid.New()
	olID := int32(4894)
	_, err := dbPool.Exec(ctx, `INSERT INTO emcd.users (id,new_id,username,api_key,email,password) VALUES ($1,$2,$3,$4,$5,$6)`, olID, newID, "user", "api", "foo@bar.com", "password")
	require.NoError(t, err)
	u, err := oldUsersRepo.GetUserByUUID(ctx, newID)
	require.NoError(t, err)
	require.Equal(t, &model.ID{New: newID, Old: olID}, &model.ID{New: u.ID, Old: u.OldID})
}

func TestOldUsers_GetIDs(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	expected := make([]*model.ID, 0)
	newIDs := make([]uuid.UUID, 0)
	for i := 0; i < 10; i++ {
		newID := uuid.New()
		_, err := dbPool.Exec(ctx, `INSERT INTO emcd.users (id,new_id,username,api_key,email) VALUES ($1,$2,$3,$4,$5)`,
			i+1, newID, "username", "api", fmt.Sprintf("email%d", i))
		require.NoError(t, err)
		expected = append(expected, &model.ID{
			New: newID,
			Old: int32(i + 1),
		})
		newIDs = append(newIDs, newID)
	}
	_, _, err := oldUsersRepo.Create(ctx, &model.User{
		ID:    uuid.New(),
		Email: "email",
	})
	require.NoError(t, err)
	actual, err := oldUsersRepo.GetIDs(ctx, newIDs)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestOldUsers_GetByEmail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	defer truncateOldUsers(t, ctx)

	// Create main user
	rep := NewOldUsers(transactor)
	u := model.User{
		ID:        uuid.New(),
		Username:  "username",
		Vip:       true,
		SegmentID: 151,
		RefID:     3,
		Email:     "qwe@gmail.com",
		Password:  "Password",
		CreatedAt: time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second),
		ApiKey:    "vpojvovj[evjqepjvpejvevjpjrevrpovj",
		OldID:     25697,
		Suspended: true,
	}
	parentID := 150
	initSql := `INSERT INTO emcd.users (id,new_id,username,ref_id,email,password,created_at, api_key, suspended, parent_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := dbPool.Exec(ctx, initSql, u.OldID, u.ID, u.Username, u.RefID, u.Email, u.Password, u.CreatedAt, u.ApiKey, time.Now().UTC(), parentID)
	require.NoError(t, err)

	initSql = `INSERT INTO histories.segment_userids (user_id, segment_id) VALUES ($1,$2)`
	_, err = dbPool.Exec(ctx, initSql, u.OldID, u.SegmentID)
	require.NoError(t, err)

	initSql = `INSERT INTO emcd.vip_users (user_id) VALUES ($1)`
	_, err = dbPool.Exec(ctx, initSql, u.OldID)

	// Create parent
	parentUser := model.User{
		ID:        uuid.New(),
		Username:  "parent",
		Vip:       true,
		SegmentID: 190,
		RefID:     6,
		Email:     "asdasdasdasd@gmail.com",
		Password:  "Password123!!!",
		CreatedAt: time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second),
		ApiKey:    "evjpjrevrpovj",
		OldID:     int32(parentID),
		Suspended: false,
	}
	initSql = `INSERT INTO emcd.users (id,new_id,username,ref_id,email,password,created_at, api_key, suspended, parent_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err = dbPool.Exec(ctx, initSql, parentUser.OldID, parentUser.ID, parentUser.Username, parentUser.RefID, parentUser.Email, parentUser.Password, parentUser.CreatedAt, parentUser.ApiKey, time.Now().UTC(), nil)
	require.NoError(t, err)

	initSql = `INSERT INTO histories.segment_userids (user_id, segment_id) VALUES ($1,$2)`
	_, err = dbPool.Exec(ctx, initSql, parentUser.OldID, parentUser.SegmentID)
	require.NoError(t, err)

	initSql = `INSERT INTO emcd.vip_users (user_id) VALUES ($1)`
	_, err = dbPool.Exec(ctx, initSql, parentUser.OldID)

	// Set
	u.ParentID = parentUser.ID

	// Test
	user, err := rep.GetByEmail(ctx, u.Email)
	require.NoError(t, err)
	require.Equal(t, &u, user)
}

func TestOldUsers_GetDefaultAddress(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	id := 6554
	userID := 65446
	coinID := 1
	initSql := `INSERT INTO emcd.users_accounts (id, user_id, coin_id, account_type_id) VALUES ($1,$2,$3,$4)`
	_, err := dbPool.Exec(ctx, initSql, id, userID, coinID, 2)
	require.NoError(t, err)
	address := "cfc[cmpsfnvp"
	initSql = `INSERT INTO emcd.autopay_addresses (id, user_account_id, address, percent) VALUES ($1,$2,$3,$4)`
	_, err = dbPool.Exec(ctx, initSql, 489, id, address, 100)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)

	defaultAddress, err := rep.GetDefaultAddress(ctx, id)
	require.NoError(t, err)
	require.Equal(t, address, defaultAddress)
}

func TestIldUsers_GetRefIDAndEmail(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	id := 6464
	refID := int32(55)
	em := "qwe@gmail.com"
	initSql := `INSERT INTO emcd.users (id, ref_id, email, username, api_key) VALUES ($1,$2,$3,$4,$5)`
	_, err := dbPool.Exec(ctx, initSql, id, refID, em, "username", "api")
	require.NoError(t, err)

	rep := NewOldUsers(transactor)

	rID, myEmail, err := rep.GetRefIDAndEmail(ctx, id)
	require.NoError(t, err)
	require.Equal(t, refID, rID)
	require.Equal(t, em, myEmail)
}

func TestOldUsers_UpdateUserMinPay(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	userID := 684
	minpay := float32(1)
	coinID := 2
	accID := 155
	initSql := `INSERT INTO emcd.users_accounts (id, minpay, user_id, coin_id, account_type_id) VALUES ($1,$2,$3,$4,$5)`
	_, err := dbPool.Exec(ctx, initSql, accID, 0, userID, coinID, 2)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)

	err = rep.UpdateUserMinPay(ctx, accID, minpay)
	require.NoError(t, err)

	var mp float32
	err = dbPool.QueryRow(ctx, `SELECT minpay FROM emcd.users_accounts WHERE user_id = $1`, userID).Scan(&mp)
	require.NoError(t, err)

	require.Equal(t, minpay, mp)
}

func TestOldUsers_GetSegmentID(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	userID := int32(34564)
	segmentID := 8743879
	initSql := `INSERT INTO histories.segment_userids (user_id, segment_id) VALUES ($1,$2)`
	_, err := dbPool.Exec(ctx, initSql, userID, segmentID)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)
	sid, err := rep.GetSegmentID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, segmentID, sid)
}

func TestOldUsers_GetSegmentID2(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	rep := NewOldUsers(transactor)
	sid, err := rep.GetSegmentID(ctx, 100)
	require.NoError(t, err)
	require.NotEqual(t, 0, sid)
}

func TestOldUsers_GetUserByOldID(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	u := model.User{
		ID:       uuid.New(),
		Username: "user",
		// Vip:          falseЗД,
		// SegmentID:    0,
		RefID:     5,
		Email:     "ert@gmail.com",
		Password:  "password",
		CreatedAt: time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second),
		// WhiteLabelID: uuid.UUID{},
		ApiKey:   "api",
		OldID:    64848,
		Language: "ru",
	}
	initSql := `INSERT INTO emcd.users (id,new_id,username,ref_id,email,password,created_at,api_key,language)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := dbPool.Exec(ctx, initSql, u.OldID, u.ID, u.Username, u.RefID, u.Email, u.Password, u.CreatedAt, u.ApiKey, u.Language)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)

	user, err := rep.GetUserByOldID(ctx, int(u.OldID))
	require.NoError(t, err)
	require.Equal(t, &u, user)
}

func TestOldUsers_CreateDeleteAddress(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	userID := int32(6516)
	coinID := 1
	id := 648648
	address := "address"
	initSql := `INSERT INTO emcd.users_accounts (id, user_id, coin_id, account_type_id) VALUES ($1,$2,$3,$4)`
	_, err := dbPool.Exec(ctx, initSql, id, userID, coinID, 2)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)
	trx := rep.Begin()
	err = trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			return rep.CreateAddress(txCtx, userID, coinID, address)
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)

	require.NoError(t, err)

	var count int
	err = dbPool.QueryRow(ctx, `SELECT count(1) FROM emcd.autopay_addresses`).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	err = trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			return rep.DeleteAddress(txCtx, userID, coinID)
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)
	require.NoError(t, err)
	err = dbPool.QueryRow(ctx, `SELECT count(1) FROM emcd.autopay_addresses WHERE id = $1`, id).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestOldUsers_UpdateUserAccountChangedAt(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	userID := int32(65165)
	initSql := `INSERT INTO emcd.users_accounts (user_id,coin_id,account_type_id) VALUES ($1,$2,$3)`
	_, err := dbPool.Exec(ctx, initSql, userID, 1, 2)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)

	err = rep.UpdateUserAccountChangedAt(ctx, userID, 1)
	require.NoError(t, err)

	var updatedAt time.Time
	err = dbPool.QueryRow(ctx, `SELECT changed_at FROM emcd.users_accounts WHERE user_id = $1`, userID).Scan(&updatedAt)
	require.NoError(t, err)
	require.False(t, updatedAt.IsZero())
}

func TestOldUsers_GetEmailWithParentID(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	u := model.User{
		Username: "user",
		Email:    "ert@gmail.com",
		ApiKey:   "api",
		OldID:    64848,
	}
	initSql := `INSERT INTO emcd.users (id,username,email,api_key) VALUES ($1,$2,$3,$4)`
	_, err := dbPool.Exec(ctx, initSql, u.OldID, u.Username, u.Email, u.ApiKey)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)

	em, err := rep.GetEmailWithParentID(ctx, u.OldID)
	require.NoError(t, err)

	require.Equal(t, u.Email, em)
}

func TestOldUsers_GetEmailWithParentID2(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	u := model.User{
		Username: "user",
		Email:    "ert@gmail.com",
		ApiKey:   "api",
		OldID:    64848,
	}
	u2 := model.User{
		Username: "user2",
		Email:    "ertqwqw@gmail.com",
		ApiKey:   "api2",
		OldID:    6546488,
	}
	initSql := `INSERT INTO emcd.users (id,username,email,api_key) VALUES ($1,$2,$3,$4)`
	_, err := dbPool.Exec(ctx, initSql, u.OldID, u.Username, u.Email, u.ApiKey)
	require.NoError(t, err)
	initSql = `INSERT INTO emcd.users (id,username,email,api_key,parent_id) VALUES ($1,$2,$3,$4,$5)`
	_, err = dbPool.Exec(ctx, initSql, u2.OldID, u2.Username, u2.Email, u2.ApiKey, u.OldID)
	require.NoError(t, err)

	rep := NewOldUsers(transactor)

	em, err := rep.GetEmailWithParentID(ctx, u2.OldID)
	require.NoError(t, err)

	require.Equal(t, u.Email, em)
}

func TestOldUsers_GetNoPay(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)

	u := model.User{
		Username: "user",
		Email:    "ert@gmail.com",
		ApiKey:   "api",
		OldID:    64848,
	}
	newID := uuid.New()
	query := `INSERT INTO emcd.users (id,username,email,api_key,nopay,new_id) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := dbPool.Exec(ctx, query, u.OldID, u.Username, u.Email, u.ApiKey, true, newID)
	require.NoError(t, err)

	nopay, err := oldUsersRepo.GetNoPay(ctx, newID)
	require.NoError(t, err)
	require.True(t, nopay)

	query = `UPDATE emcd.users SET nopay = false WHERE new_id = $1`
	_, err = dbPool.Exec(ctx, query, newID)
	require.NoError(t, err)

	nopay, err = oldUsersRepo.GetNoPay(ctx, newID)
	require.NoError(t, err)
	require.False(t, nopay)
}

func TestOldUsers_UpdateRefID(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	initialUser := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	oldID, _, err := oldUsersRepo.Create(ctx, &initialUser)
	require.NoError(t, err)

	var newRefID int32 = 999
	err = oldUsersRepo.UpdateRefID(ctx, oldID, newRefID)
	require.NoError(t, err)
	actual, err := oldUsersRepo.GetByEmail(ctx, initialUser.Email)
	require.NoError(t, err)
	require.Equal(t, newRefID, int32(actual.RefID))
}

func TestOldUsers_GetByID(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	u := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
		OldID:     64848,
		Language:  "ru",
	}
	initSql := `INSERT INTO emcd.users (id,new_id,username,ref_id,email,password,created_at,api_key,language)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := dbPool.Exec(ctx, initSql, u.OldID, u.ID, u.Username, u.RefID, u.Email, u.Password, u.CreatedAt, u.ApiKey, u.Language)
	require.NoError(t, err)
	_, err = oldUsersRepo.GetUserByOldID(ctx, int(u.OldID))
	require.NoError(t, err)
}

func TestOldUsers_SetTimezone(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	initialUser := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	oldID, _, err := oldUsersRepo.Create(ctx, &initialUser)
	require.NoError(t, err)

	timezone := "+3"
	err = oldUsersRepo.SetTimezone(ctx, int(oldID), timezone)
	require.NoError(t, err)

	var tm string
	err = oldUsersRepo.trx.Runner(ctx).QueryRow(ctx, `SELECT timezone FROM emcd.users WHERE id = $1`, oldID).Scan(&tm)
	require.NoError(t, err)
	require.Equal(t, timezone, tm)
}

func TestOldUsers_GetByIDNoLanguage(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	u := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
		OldID:     64848,
	}
	initSql := `INSERT INTO emcd.users (id,new_id,username,ref_id,email,password,created_at,api_key)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := dbPool.Exec(ctx, initSql, u.OldID, u.ID, u.Username, u.RefID, u.Email, u.Password, u.CreatedAt, u.ApiKey)
	require.NoError(t, err)
	actual, err := oldUsersRepo.GetUserByOldID(ctx, int(u.OldID))
	require.NoError(t, err)
	require.Equal(t, "en", actual.Language)
}

func TestOldUsers_UpdateUserEmail(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	u := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	oldID, _, err := oldUsersRepo.Create(ctx, &u)
	u.OldID = oldID
	require.NoError(t, err)
	u.Email = "new_email"
	err = oldUsersRepo.UpdateUser(ctx, &u)
	require.NoError(t, err)
	actual, err := oldUsersRepo.GetUserByOldID(ctx, int(u.OldID))
	require.NoError(t, err)
	require.Equal(t, u.Email, actual.Email)
}

func TestOldUsers_UpdateUserEmailExist(t *testing.T) {
	ctx := context.Background()
	defer truncateOldUsers(t, ctx)
	u1 := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email1",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	u2 := model.User{
		ID:        uuid.New(),
		Username:  "username",
		RefID:     123,
		Email:     "email2",
		Password:  "password",
		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
		ApiKey:    "apikey",
	}
	oldID1, _, err := oldUsersRepo.Create(ctx, &u1)
	u1.OldID = oldID1
	require.NoError(t, err)
	oldID2, _, err := oldUsersRepo.Create(ctx, &u2)
	u2.OldID = oldID2
	require.NoError(t, err)
	oldEmail := u2.Email
	u2.Email = u1.Email
	err = oldUsersRepo.UpdateUser(ctx, &u2)
	require.Error(t, err)
	actual, err := oldUsersRepo.GetUserByOldID(ctx, int(u2.OldID))
	require.NoError(t, err)
	require.Equal(t, oldEmail, actual.Email)
}

func insertCoins(ctx context.Context) {
	_, err := dbPool.Exec(ctx, `INSERT INTO emcd.coins (id,code,name) VALUES (1,'btc','BTC'),
        (2,'bch', 'BCH'),(3,'bsv','BSV'),(4,'ltc','LTC'),(5,'dash','DASH'),(6,'eth','ETH'),(7,'etc','ETC'),
        (8,'doge','DOGE'),(9,'bnb','BNB'),(10,'usdt','USDT'),(11,'usdc','USDC'),(12,'trx','TRX'),(13,'ton','TON')`)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}
}

func insertAccountReferral(t *testing.T, ctx context.Context, id, coinID int, reward float64, accID int) {
	_, err := dbPool.Exec(ctx, `INSERT INTO emcd.accounts_referral (id,account_id,coin_id,referral_fee) VALUES ($1,$2,$3,$4)`,
		id, accID, coinID, reward)
	require.NoError(t, err)
}

func insertUsersAccounts(t *testing.T, ctx context.Context, id, userID int32, coinID int, accType int, fee float64) {
	_, err := dbPool.Exec(ctx, `INSERT INTO emcd.users_accounts (id,user_id,coin_id,account_type_id,fee)
VALUES ($1,$2,$3,$4,$5)`, id, userID, coinID, accType, fee)
	require.NoError(t, err)
}

func truncateOldUsers(t *testing.T, ctx context.Context) {
	_, err := dbPool.Exec(ctx, "TRUNCATE TABLE emcd.users CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE emcd.vip_users CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE histories.segment_userids CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE emcd.users_accounts CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE emcd.accounts_referral CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE histories.segment_userids CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE emcd.accounts_pool CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE emcd.autopay_addresses CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE emcd.histories CASCADE")
	require.NoError(t, err)
	_, err = dbPool.Exec(ctx, "TRUNCATE TABLE emcd.user_logs CASCADE")
	require.NoError(t, err)
}

func createUser(ctx context.Context, u *model.User, oldParentID int32) (int32, error) {
	var id int32
	query := `INSERT INTO emcd.users (username, email, password, parent_id, api_key) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := dbPool.QueryRow(ctx, query, u.Username, u.Email, u.Password, oldParentID, u.ApiKey).Scan(&id)
	return id, err
}

type subUserOldUser struct {
	ID       int32
	ParentID int32
	Username string
	Email    string
}

func getSubUsersOldUsers(ctx context.Context, oldParentID int32) ([]*subUserOldUser, error) {
	query := `SELECT id, username, email, parent_id FROM emcd.users WHERE parent_id = $1`
	rows, err := dbPool.Query(ctx, query, oldParentID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	users := make([]*subUserOldUser, 0)
	for rows.Next() {
		var user subUserOldUser
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.ParentID)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}
