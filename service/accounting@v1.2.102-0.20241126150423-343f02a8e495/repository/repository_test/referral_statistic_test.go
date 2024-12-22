package repository_test

import (
	referralMock "code.emcdtech.com/emcd/service/accounting/mocks/protocol/referral"
	"code.emcdtech.com/emcd/service/accounting/model"
	referralPb "code.emcdtech.com/emcd/service/accounting/protocol/referral"
	externalRepository "code.emcdtech.com/emcd/service/accounting/repository"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientReferralRepositoryGetReferralStatistic_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	referralStatisticRepository := externalRepository.NewReferralStatisticRepository()
	protoAccountingReferralServiceClientMock := referralMock.NewMockAccountingReferralServiceClient(t)

	req := model.ReferralsStatisticInput{
		UserID:              uuid.UUID{},
		AccountIDs:          nil,
		TransactionTypesIDs: []int64{1},
	}

	protoAccountingReferralServiceClientMock.On("GetReferralsStatistic", mock.Anything, &referralPb.GetReferralsStatisticRequest{
		UserId:              req.UserID.String(),
		TransactionTypesIds: req.TransactionTypesIDs,
	}).Return(&referralPb.GetReferralsStatisticResponse{
		ThisMonthIncome: nil,
		YesterdayIncome: nil,
	}, nil)

	stats, err := referralStatisticRepository.GetReferralsStatistic(ctx, protoAccountingReferralServiceClientMock, req)
	require.NoError(t, err)
	require.NotNil(t, stats)
}

func TestClientReferralRepositoryGetReferralStatistic_FailEmptyTransactionsList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	referralStatisticRepository := externalRepository.NewReferralStatisticRepository()
	protoAccountingReferralServiceClientMock := referralMock.NewMockAccountingReferralServiceClient(t)

	req := model.ReferralsStatisticInput{
		UserID:              uuid.UUID{},
		AccountIDs:          nil,
		TransactionTypesIDs: nil,
	}

	_, err := referralStatisticRepository.GetReferralsStatistic(ctx, protoAccountingReferralServiceClientMock, req)
	// require.Nil(t, stats)
	require.Error(t, err)

	require.EqualError(t, err, "referral accounting: empty transaction types filter")

}

func TestClientReferralRepositoryGetReferralStatistic_ErrorMockReturn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	referralStatisticRepository := externalRepository.NewReferralStatisticRepository()
	protoAccountingReferralServiceClientMock := referralMock.NewMockAccountingReferralServiceClient(t)

	req := model.ReferralsStatisticInput{
		UserID:              uuid.UUID{},
		AccountIDs:          nil,
		TransactionTypesIDs: []int64{1},
	}

	retError := newMockError()

	protoAccountingReferralServiceClientMock.On("GetReferralsStatistic", mock.Anything, &referralPb.GetReferralsStatisticRequest{
		UserId:              req.UserID.String(),
		TransactionTypesIds: req.TransactionTypesIDs,
	}).Return(nil, retError)

	_, err := referralStatisticRepository.GetReferralsStatistic(ctx, protoAccountingReferralServiceClientMock, req)
	// require.Nil(t, stats)
	require.Error(t, err)

	require.ErrorAs(t, err, &retError)

}
