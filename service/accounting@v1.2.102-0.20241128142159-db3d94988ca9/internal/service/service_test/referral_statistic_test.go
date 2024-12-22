package service_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"code.emcdtech.com/emcd/service/accounting/internal/utils"
	repositoryMock "code.emcdtech.com/emcd/service/accounting/mocks/internal_/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	referralPb "code.emcdtech.com/emcd/service/accounting/protocol/referral"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReferralStatisticGetReferralStatistic_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	referralStatisticRepo := repositoryMock.NewMockReferralStatistic(t)
	userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)

	referralStatisticService := service.NewReferral(referralStatisticRepo, userAccountRepo)

	userIdNew := uuid.New()
	userId := int32(1)
	coinId := int32(1)
	coinNew := "btc"

	req := &referralPb.GetReferralsStatisticRequest{
		UserId:              userIdNew.String(),
		TransactionTypesIds: nil,
	}

	userAccountsModel := model.UserAccounts{{
		ID:            0,
		UserID:        userId,
		CoinID:        coinId,
		AccountTypeID: enum.AccountTypeIdWrapper{},
		Minpay:        0,
		Address:       sql.NullString{},
		ChangedAt:     sql.NullTime{},
		Img1:          sql.NullFloat64{},
		Img2:          sql.NullFloat64{},
		IsActive:      sql.NullBool{},
		CreatedAt:     sql.NullTime{},
		UpdatedAt:     sql.NullTime{},
		Fee:           sql.NullFloat64{},
		UserIDNew:     utils.UuidToUuidNull(userIdNew),
		CoinNew:       utils.StringToStringNull(coinNew),
	},
	}

	userAccountIds := []int64{0}

	reqIn := &model.ReferralsStatisticInput{
		UserID:              userIdNew,
		AccountIDs:          userAccountIds,
		TransactionTypesIDs: req.TransactionTypesIds,
	}

	reqOut := &model.ReferralsStatisticOutput{
		Income: model.AggregatedReferralIncome{
			Yesterday: nil,
			ThisMonth: nil,
		},
	}

	userAccountRepo.EXPECT().
		FindUserAccountByUserIdLegacy(ctx, userIdNew).
		Return(userAccountsModel, nil)
	referralStatisticRepo.EXPECT().
		GetReferralsStatistic(ctx, reqIn).
		Return(reqOut, nil)

	result, err := referralStatisticService.GetReferralsStatistic(ctx, req)
	require.NotNil(t, result)
	require.NoError(t, err)

}

func TestReferralStatisticGetReferralStatistic_FailCauseAccount(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	referralStatisticRepo := repositoryMock.NewMockReferralStatistic(t)
	userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)

	referralStatisticService := service.NewReferral(referralStatisticRepo, userAccountRepo)

	userIdNew := uuid.New()

	req := &referralPb.GetReferralsStatisticRequest{
		UserId:              userIdNew.String(),
		TransactionTypesIds: nil,
	}

	retError := errors.New("error")

	userAccountRepo.EXPECT().
		FindUserAccountByUserIdLegacy(ctx, userIdNew).
		Return(nil, retError)

	result, err := referralStatisticService.GetReferralsStatistic(ctx, req)
	require.Nil(t, result)
	require.Error(t, err)

	require.EqualError(t, err, fmt.Errorf("get referrals statistic service: getting account ids for user: %s: %w", userIdNew.String(), retError).Error())

}

func TestReferralStatisticGetReferralStatistic_FailCauseStatistic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	referralStatisticRepo := repositoryMock.NewMockReferralStatistic(t)
	userAccountRepo := repositoryMock.NewMockUserAccountRepo(t)

	referralStatisticService := service.NewReferral(referralStatisticRepo, userAccountRepo)

	userIdNew := uuid.New()
	userId := int32(1)
	coinId := int32(1)
	coinNew := "btc"

	req := &referralPb.GetReferralsStatisticRequest{
		UserId:              userIdNew.String(),
		TransactionTypesIds: nil,
	}

	userAccountsModel := model.UserAccounts{{
		ID:            0,
		UserID:        userId,
		CoinID:        coinId,
		AccountTypeID: enum.AccountTypeIdWrapper{},
		Minpay:        0,
		Address:       sql.NullString{},
		ChangedAt:     sql.NullTime{},
		Img1:          sql.NullFloat64{},
		Img2:          sql.NullFloat64{},
		IsActive:      sql.NullBool{},
		CreatedAt:     sql.NullTime{},
		UpdatedAt:     sql.NullTime{},
		Fee:           sql.NullFloat64{},
		UserIDNew:     utils.UuidToUuidNull(userIdNew),
		CoinNew:       utils.StringToStringNull(coinNew),
	},
	}

	userAccountIds := []int64{0}

	reqIn := &model.ReferralsStatisticInput{
		UserID:              userIdNew,
		AccountIDs:          userAccountIds,
		TransactionTypesIDs: req.TransactionTypesIds,
	}

	retError := newMockError()

	userAccountRepo.EXPECT().
		FindUserAccountByUserIdLegacy(ctx, userIdNew).
		Return(userAccountsModel, nil)
	referralStatisticRepo.EXPECT().
		GetReferralsStatistic(ctx, reqIn).
		Return(nil, retError)

	result, err := referralStatisticService.GetReferralsStatistic(ctx, req)
	require.Nil(t, result)
	require.Error(t, err)

	require.ErrorIs(t, err, retError)

}
