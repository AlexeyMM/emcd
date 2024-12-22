package repository_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"hash/fnv"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestUserAccount(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer truncateTables(
		ctx, t, "emcd.users", "emcd.account_types", "emcd.coins", "emcd.users_accounts")

	repo := repository.NewUserAccountRepo(dbPool)

	// user_id
	var users []*model.User
	for i := 1; i <= 2; i++ {
		user := &model.User{
			ID:       int64(i),
			Username: strconv.Itoa(i),
			Password: strconv.Itoa(i),
			NewID:    uuid.New(),
		}

		users = append(users, user)
		writeUser(ctx, dbPool, t, *user)

	}

	accountTypes := []enum.AccountTypeIdWrapper{
		{AccountTypeId: enum.WalletAccountTypeID},
		{AccountTypeId: enum.CoinholdAccountTypeID},
		// {enum.MiningAccountTypeID},
		// {enum.ReferralAccountTypeID},
		// {enum.BlockUserAccountTypeID},
	}

	writeAccountType(ctx, dbPool, t, accountTypes[0])
	writeAccountType(ctx, dbPool, t, accountTypes[1])
	// writeAccountType(ctx, dbPool, t, accountTypes[2])
	// writeAccountType(ctx, dbPool, t, accountTypes[3])
	// writeAccountType(ctx, dbPool, t, accountTypes[4])

	// coin_id
	var coins []*model.Coin
	for i := 1; i <= 11; i++ {
		coin := &model.Coin{
			ID:          int64(i),
			Name:        strconv.Itoa(i),
			Description: strconv.Itoa(i),
			Code:        strconv.Itoa(i),
			Rate:        float64(i),
		}

		coins = append(coins, coin)
		writeCoin(ctx, dbPool, t, *coin)

	}

	hashMap := make(map[string]struct{})

	t.Run("CreateUserAccounts", func(t *testing.T) {
		var userAccounts model.UserAccounts

		userAccountCountCoinhold := 0
		userAccountCountCommonUniq := 0
		userAccountCountCommonTotal := 0
		userAccountCountUserIdNewEmpty := 0
		// hashEmptyMap := make(map[string]struct{})

		for i := 0; i < 100; i++ {
			user := users[rand.Intn(len(users))]
			accountType := accountTypes[rand.Intn(len(accountTypes))]
			coin := coins[rand.Intn(len(coins))]

			userAccount := model.UserAccount{
				ID:            0,
				UserID:        int32(user.ID),
				CoinID:        int32(coin.ID),
				AccountTypeID: accountType,
				Minpay:        0,
				Address:       sql.NullString{},
				ChangedAt:     sql.NullTime{},
				Img1:          sql.NullFloat64{},
				Img2:          sql.NullFloat64{},
				IsActive:      sql.NullBool{},
				CreatedAt:     sql.NullTime{},
				UpdatedAt:     sql.NullTime{},
				Fee:           sql.NullFloat64{},
				UserIDNew:     uuid.NullUUID{UUID: user.NewID, Valid: true},
				CoinNew:       sql.NullString{String: coin.Code, Valid: true},
			}

			userAccounts = append(userAccounts, &userAccount)

			if userAccount.AccountTypeID.AccountTypeId == enum.CoinholdAccountTypeID {
				userAccountCountCoinhold++

			} else {
				userAccountCountCommonTotal++

				hash := fmt.Sprintf("%d:%d:%d", userAccount.UserID, userAccount.AccountTypeID, userAccount.CoinID)

				if _, ok := hashMap[hash]; !ok {
					hashMap[hash] = struct{}{}
					userAccountCountCommonUniq++

				}

				fnvHash := fnv.New32a()
				if _, err := fnvHash.Write([]byte(hash)); err != nil {
					require.NoError(t, err)

					// } else if fnvHash.Sum32()%100 < 0 {
					//	userAccount.UserIDNew = uuid.NullUUID{UUID: uuid.UUID{}, Valid: false}
					//	if _, ok := hashEmptyMap[hash]; !ok {
					//		hashEmptyMap[hash] = struct{}{}
					//		userAccountCountUserIdNewEmpty++
					//
					//	}
				}
			}
		}

		// userAccounts = GenerateDefaultUsersAccounts(uuid.New(), 1, 0.015)

		writeConcurrencyChunked(ctx, t, repo, 20, userAccounts)

		filterCoinhold := &model.UserAccountFilter{
			UserID:        nil,
			AccountTypeID: enum.CoinholdAccountTypeID.ToPtr(),
			IsActive:      utils.BoolToPtr(true),
			UserIDNew:     nil,
			Pagination:    nil,
		}

		_, resCoinhold, err := repo.FindUserAccountByFilter(ctx, filterCoinhold)
		require.NoError(t, err)
		require.Len(t, resCoinhold, userAccountCountCoinhold)
		t.Logf("userAccountCountCoinhold = %d", userAccountCountCoinhold)

		filterCommon := &model.UserAccountFilter{
			UserID:        nil,
			AccountTypeID: enum.WalletAccountTypeID.ToPtr(),
			IsActive:      utils.BoolToPtr(true),
			UserIDNew:     nil,
			Pagination:    nil,
		}

		_, resCommon, err := repo.FindUserAccountByFilter(ctx, filterCommon)
		require.NoError(t, err)
		require.Len(t, resCommon, userAccountCountCommonUniq)
		t.Logf("userAccountCountCommonUniq = %d", userAccountCountCommonUniq)
		t.Logf("userAccountCountCommonTotal = %d", userAccountCountCommonTotal)

		filterPagination := &model.UserAccountFilter{
			UserID:        nil,
			AccountTypeID: nil,
			IsActive:      utils.BoolToPtr(true),
			UserIDNew:     nil,
			Pagination: &model.Pagination{
				Limit:  999999999,
				Offset: 0,
			},
		}

		totalReturn, resPagination, err := repo.FindUserAccountByFilter(ctx, filterPagination)
		require.NoError(t, err)
		require.Equal(t, uint64(len(resPagination)), *totalReturn)
		t.Logf("userAccountCountTotal = %d", *totalReturn)

		nullUid := uuid.UUID{}
		emptyStr := ""
		emptyStruct := struct{}{}

		filterMigrate := &model.UserAccountFilter{
			UserID:          nil,
			AccountTypeID:   nil,
			UserIDNew:       nil,
			IsActive:        nil,
			Pagination:      nil,
			UserIDNewIsNull: &emptyStruct,
			CoinNew:         nil,
		}

		t.Logf("userAccountCountUserIdNewEmpty = %d", userAccountCountUserIdNewEmpty)
		_, resMigrate, err := repo.FindUserAccountByFilter(ctx, filterMigrate)
		require.NoError(t, err)
		require.Len(t, resMigrate, userAccountCountUserIdNewEmpty)

		partial := &model.UserAccountPartial{
			UserIDNew: &nullUid,
			CoinNew:   &emptyStr,
		}

		for _, userAccount := range resMigrate {
			if rand.Float64() < 0.5 {
				if err := repo.UpdateUserAccountByPartial(ctx, userAccount, partial); err != nil {
					require.NoError(t, err)
				}
			}
		}

		if len(resMigrate) > 0 {
			if err := repo.UpdateUserAccountForMigrateUserIdNew(ctx, resMigrate); err != nil {
				require.NoError(t, err)
			}

		}

		_, resMigrateAfter, err := repo.FindUserAccountByFilter(ctx, filterMigrate)
		require.NoError(t, err)
		require.Empty(t, resMigrateAfter)

	})

	t.Run("filter", func(t *testing.T) {
		var userAccounts model.UserAccounts

		userAccountCountIsActiveFalse := 0

		for i := 0; i < 100; i++ {
			user := users[rand.Intn(len(users))]
			accountType := enum.WalletAccountTypeID
			coin := coins[rand.Intn(len(coins))]

			var isActive sql.NullBool

			if rand.Float64() < 0.5 {
				isActive = sql.NullBool{Valid: true, Bool: false}

			} else {
				isActive = sql.NullBool{}

			}

			userAccount := model.UserAccount{
				ID:            0,
				UserID:        int32(user.ID),
				CoinID:        int32(coin.ID),
				AccountTypeID: enum.NewAccountTypeIdWrapper(accountType),
				Minpay:        0,
				Address:       sql.NullString{},
				ChangedAt:     sql.NullTime{},
				Img1:          sql.NullFloat64{},
				Img2:          sql.NullFloat64{},
				IsActive:      isActive,
				CreatedAt:     sql.NullTime{},
				UpdatedAt:     sql.NullTime{},
				Fee:           sql.NullFloat64{},
				UserIDNew:     uuid.NullUUID{UUID: user.NewID, Valid: true},
				CoinNew:       sql.NullString{String: coin.Code, Valid: true},
			}

			userAccounts = append(userAccounts, &userAccount)

			hash := fmt.Sprintf("%d:%d:%d", userAccount.UserID, userAccount.AccountTypeID, userAccount.CoinID)

			if _, ok := hashMap[hash]; !ok {
				hashMap[hash] = struct{}{}
				if userAccount.IsActive.Valid && !userAccount.IsActive.Bool {
					userAccountCountIsActiveFalse++

				}
			}
		}

		writeConcurrencyChunked(ctx, t, repo, 20, userAccounts)
		time.Sleep(100 * time.Millisecond)

		filterIsActiveFalse := &model.UserAccountFilter{
			UserID:          nil,
			AccountTypeID:   nil,
			UserIDNew:       nil,
			IsActive:        utils.BoolToPtr(false),
			Pagination:      nil,
			UserIDNewIsNull: nil,
			CoinNew:         nil,
		}

		_, resIsActiveFalse, err := repo.FindUserAccountByFilter(ctx, filterIsActiveFalse)
		require.NoError(t, err)
		require.Len(t, resIsActiveFalse, userAccountCountIsActiveFalse)

	})

}

func writeConcurrencyChunked(ctx context.Context, t *testing.T, userAccountRepo repository.UserAccountRepo, chunkCap int, userAccounts model.UserAccounts) {
	wg := new(sync.WaitGroup)

	splitMap, splitMapId := splitUserAccountsByUser(userAccounts)

	for userIdNew, userAccountsSplitList := range splitMap {
		userId := splitMapId[userIdNew]

		wg.Add(1)
		go func(userId int32, userIdNew uuid.UUID, userAccountsSplitList model.UserAccounts) {
			defer func() {
				wg.Done()

			}()

			for _, userAccountChunkList := range chunkUserAccounts(userAccountsSplitList, chunkCap) {
				func(t *testing.T) {
					ctxLocal, cancel := context.WithCancel(ctx)
					defer cancel()

					if err := userAccountRepo.WithinTransaction(ctxLocal, func(ctx context.Context) error {
						if err := userAccountRepo.AddUserAccounts(ctx, userId, userIdNew, userAccountChunkList); err != nil {
							t.Error(err)
							t.Fail()

							return err
						} else {
							for _, u := range userAccountChunkList {
								if u.ID == 0 {
									t.Errorf("user_id is zero")
									t.Fail()

									return errors.New("user_id is zero")
								} else if u.UserIDNew.UUID == uuid.New() {
									t.Errorf("user_id_new is zero")
									t.Fail()

									return errors.New("user_id_new is zero")
								}

								// require.NotEqual(t, 0, u.ID)
								// require.NotEqual(t, u.UserIDNew.UUID, uuid.Nil)
							}

							return nil
						}
					}); err != nil {
						t.Errorf(err.Error())
						t.Fail()

						return

					}
				}(t)
			}
		}(userId, userIdNew, userAccountsSplitList)
	}

	wg.Wait()
}

func splitUserAccountsByUser(userAccounts model.UserAccounts) (map[uuid.UUID]model.UserAccounts, map[uuid.UUID]int32) {
	splitMap1 := make(map[uuid.UUID]model.UserAccounts)
	splitMap2 := make(map[uuid.UUID]int32)
	for _, userAccount := range userAccounts {
		splitMap1[userAccount.UserIDNew.UUID] = append(splitMap1[userAccount.UserIDNew.UUID], userAccount)
		splitMap2[userAccount.UserIDNew.UUID] = userAccount.UserID

	}

	return splitMap1, splitMap2
}

func chunkUserAccounts(userAccounts model.UserAccounts, chunkCap int) []model.UserAccounts {
	batch := make([]model.UserAccounts, 0, 1+len(userAccounts)/chunkCap)

	var items model.UserAccounts
	for _, elem := range userAccounts {
		if len(items) == 0 {
			items = append(items, elem)

		} else if len(items) < chunkCap {
			items = append(items, elem)

		}

		if len(items) == chunkCap {
			batch = append(batch, items)
			items = nil

		}
	}

	if items != nil {
		batch = append(batch, items)

	}

	return batch
}
