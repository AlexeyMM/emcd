package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
)

const (
	Day  = time.Hour * 24
	Year = time.Hour * 24 * 365

	ChangeBalance            string = "ChangeBalance"
	ChangeBalanceWithBlock   string = "ChangeBalanceWithBlock"
	ChangeBalanceWithUnblock string = "ChangeBalanceWithUnblock"
	DeprecatedAction         string = "Deprecated"
)

var TransactionsActions = map[model.TransactionType]string{
	model.IncomeBillTrTypeID:                      ChangeBalance, // 1
	model.CnhldInitialTrTypeID:                    ChangeBalance, // 2
	model.MainCoinMiningPayoutTrTypeID:            ChangeBalance, // 5
	model.MergeCoinMiningPayoutTrTypeID:           ChangeBalance, // 6
	model.UserPaysPoolComsTrTypeID:                ChangeBalance, // 7
	model.UserPaysPoolDonationsTrTypeID:           ChangeBalance, // 8
	model.PoolPaysUserCompensationTrTypeID:        ChangeBalance, // 9
	model.PayUserReferralBonusTrTypeID:            ChangeBalance, // 10
	model.FromMiningReplenishmentTrTypeID:         ChangeBalance, // 11
	model.PayoutTrTypeID:                          ChangeBalance, // 13
	model.CnhldWithdrawalTrTypeID:                 ChangeBalance, // 16
	model.PoolPaysUsersReferralsTrTypeID:          ChangeBalance, // 22
	model.PoolPaysBenefitOtherUserTrTypeID:        ChangeBalance, // 23
	model.EarlyPaymentTrTypeID:                    ChangeBalance, // 28
	model.TransferBetweenAccountsTrTypeID:         ChangeBalance, // 29
	model.WithdrawalWithCommissionTrTypeID:        ChangeBalance, // 30
	model.PayoutCommissionBillTrTypeID:            ChangeBalance, // 32
	model.CnhldDiffBalanceTrTypeID:                ChangeBalance, // 34
	model.InterestAccrualTrTypeID:                 ChangeBalance, // 35
	model.ReturnInterestsTrTypeID:                 ChangeBalance, // 37
	model.InterestCapitalizationTrTypeID:          ChangeBalance, // 38
	model.CommissionFeesTrTypeID:                  ChangeBalance, // 40
	model.ToFinanceTransferBillTrTypeID:           ChangeBalance, // 41
	model.UserUnblockTransferCommonWalletTrTypeID: ChangeBalance, // 46
	model.FundsTransferHedgeAccTrTypeID:           ChangeBalance, // 47
	model.ComRemoveHedgeAccTrTypeID:               ChangeBalance, // 48
	model.UserHedgeIncomeTrTypeID:                 ChangeBalance, // 49
	model.MiningWalletPayoutTrType:                ChangeBalance, // 51
	model.PaymentReferralAccrualTrType:            ChangeBalance, // 52
	model.BalanceTransferToAdminTrType:            ChangeBalance, // 53
	model.CnhldRefBonusTrTypeID:                   ChangeBalance, // 54
	model.PromoAuthorBonusTrTypeID:                ChangeBalance, // 55
	model.UserCompensationOrBonusTrTypeID:         ChangeBalance, // 56
	model.ExchPayoutTrTypeID:                      ChangeBalance, // 60
	model.ToKucoinTradeTransferTrTypeID:           ChangeBalance, // 61
	model.FromKucoinTradeTransferTrTypeID:         ChangeBalance, // 62
	model.KucoinTransferFeeSpendTrTypeID:          ChangeBalance, // 63
	model.WalletToWalletTrTypeID:                  ChangeBalance, // 64
	model.MiningToWalletTrTypeID:                  ChangeBalance, // 65
	model.RejectedWithdrawalTrTypeID:              ChangeBalance, // 70
	model.RejectedWithdrawalFeeTrTypeID:           ChangeBalance, // 71
	model.P2PUnblockToAdminTrTypeID:               ChangeBalance, // 72
	model.TransferFromWalletToP2PAccount:          ChangeBalance, // 74
	model.TransferFromP2PAccountToWallet:          ChangeBalance, // 75
	model.CreditingFreeBonusToAccount:             ChangeBalance, // 76
	model.CompensationForProblemsInProduct:        ChangeBalance, // 77
	model.CorrectionErrorsInBilling:               ChangeBalance, // 78
	model.IncomeForSharedMining:                   ChangeBalance, // 85
	model.FiatCardWithdraw:                        ChangeBalance, // 86
	model.RefundFiatCardWithdrawFailed:            ChangeBalance, // 87
	model.UserPaysWlComsTrTypeID:                  ChangeBalance, // 88
	model.WlPaysUserReferralsTrTypeID:             ChangeBalance, // 89
	model.OffchainWalletWithdraw:                  ChangeBalance, // 90
	model.OffchainWalletReceipt:                   ChangeBalance, // 91
	model.P2PTransfer:                             ChangeBalance, // 93
	model.FiatProviderExternalAddrTransfer:        ChangeBalance, // 94
	model.P2PCryptoHold:                           ChangeBalance, // 95
	model.P2PCryptoRelease:                        ChangeBalance, // 96
	model.P2PTransferFee:                          ChangeBalance, // 97
	model.P2PTransferFeeRef:                       ChangeBalance, // 98
	model.ProcessingInvoiceFee:                    ChangeBalance, // 99

	model.PoolPaysUsersBalanceTrTypeID:     ChangeBalanceWithBlock, // 21
	model.WalletMiningTransferTrTypeID:     ChangeBalanceWithBlock, // 31
	model.CnhldEarlyCloseTrTypeID:          ChangeBalanceWithBlock, // 36
	model.FiatWithdrawTrTypeID:             ChangeBalanceWithBlock, // 42
	model.HedgeBuyBlockTrTypeID:            ChangeBalanceWithBlock, // 45
	model.ExchBlockTrTypeID:                ChangeBalanceWithBlock, // 57
	model.P2PSellTrType:                    ChangeBalanceWithBlock, // 66
	model.P2PSellCommissionTrType:          ChangeBalanceWithBlock, // 69
	model.BlockTransferFromDepositToWallet: ChangeBalanceWithBlock, // 83

	model.AdminUnblockWithdrawalTrTypeID:          ChangeBalanceWithUnblock, // 43
	model.UnblockBackToUserFailWithdrawalTrTypeID: ChangeBalanceWithUnblock, // 44
	model.UserUnblockRefundHedgeTrTypeID:          ChangeBalanceWithUnblock, // 50
	model.ExchRollbackTrTypeID:                    ChangeBalanceWithUnblock, // 58
	model.ExchUnblockTrTypeID:                     ChangeBalanceWithUnblock, // 59
	model.P2PBuyTrType:                            ChangeBalanceWithUnblock, // 67
	model.P2PBuyFailTrType:                        ChangeBalanceWithUnblock, // 68
	model.P2PUnblockFailToSellerTrTypeID:          ChangeBalanceWithUnblock, // 73
	model.UnblockPayoutMiningToNode:               ChangeBalanceWithUnblock, // 79
	model.RefundCanceledMiningPayouts:             ChangeBalanceWithUnblock, // 80
	model.UnblockPayoutFreeWalletToNode:           ChangeBalanceWithUnblock, // 81
	model.RefundCanceledFreePayoutFromWallet:      ChangeBalanceWithUnblock, // 82
	model.UnblockTransferFromDepositToWallet:      ChangeBalanceWithUnblock, // 84
	model.OffchainMiningPayout:                    ChangeBalanceWithUnblock, // 92

	model.PoolIncomeFromOuterPoolsMainTrTypeID:  DeprecatedAction, // 3
	model.PoolIncomeFromOuterPoolsMergeTrTypeID: DeprecatedAction, // 4
	model.FromPlatformReplenishmentTrTypeID:     DeprecatedAction, // 12
	model.UserPercentCapTypeID:                  DeprecatedAction, // 14
	model.CnhldCloseTrTypeID:                    DeprecatedAction, // 15
	model.PoolTransfersCoinsExternalTrTypeID:    DeprecatedAction, // 17
	model.PoolGetsPercentsExternalCapTrTypeID:   DeprecatedAction, // 18
	model.PoolGetsPercentsExternalTrTypeID:      DeprecatedAction, // 19
	model.PoolGetsWithdrawalExternalTrTypeID:    DeprecatedAction, // 20
	model.DBTransferAccType1TrTypeID:            DeprecatedAction, // 24
	model.DBTransferAccType2TrTypeID:            DeprecatedAction, // 25
	model.DBTransferAccType3TrTypeID:            DeprecatedAction, // 26
	model.DBTransferAccType4TrTypeID:            DeprecatedAction, // 27
	model.AccrualsChineseTechAccountTrTypeID:    DeprecatedAction, // 33
	model.BaseCoinFeesTokenWithdrawalsTrTypeID:  DeprecatedAction, // 39

}

func NullString(v *string) *wrapperspb.StringValue {
	var result *wrapperspb.StringValue
	if v == nil {
		return result
	}
	result = &wrapperspb.StringValue{Value: *v}
	return result
}

func NullInt64(v *int) *wrapperspb.Int64Value {
	var result *wrapperspb.Int64Value
	if v == nil {
		return result
	}
	result = &wrapperspb.Int64Value{Value: int64(*v)}
	return result
}

func NullFloat(v *float64) *wrapperspb.StringValue {
	var result *wrapperspb.StringValue
	if v == nil {
		return result
	}
	strFromFloat := strconv.FormatFloat(*v, 'f', -1, 64)
	result = &wrapperspb.StringValue{Value: strFromFloat}
	return result
}

func NullBool(v *bool) *wrapperspb.BoolValue {
	var result *wrapperspb.BoolValue
	if v == nil {
		return result
	}
	result = &wrapperspb.BoolValue{Value: *v}
	return result
}

func ParseProtoTransaction(proto *accountingPb.Transaction) (*model.Transaction, error) {
	coinID, err := strconv.Atoi(proto.CoinID)
	if err != nil {
		return nil, fmt.Errorf("parsing coin id: %w", err)
	}

	amount, err := decimal.NewFromString(proto.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount value: %w", err)
	}

	return &model.Transaction{
		ID:                0,
		Type:              model.TransactionType(proto.Type),
		CreatedAt:         proto.CreatedAt.AsTime(),
		SenderAccountID:   proto.SenderAccountID,
		ReceiverAccountID: proto.ReceiverAccountID,
		CoinID:            int64(coinID),
		Amount:            amount,
		Comment:           proto.Comment,
		Hash:              proto.Hash,
		ReceiverAddress:   proto.ReceiverAddress,
		Hashrate:          0,
		FromReferralId:    0,
		ActionID:          "", // TODO: why empty?
		TokenID:           proto.TokenID,
		UnblockAccountId:  0,
		BlockedTill:       time.Time{},
	}, nil
}

func GetBlockTillByType(transactionType model.TransactionType) time.Time {
	switch transactionType {
	case model.CnhldDiffBalanceTrTypeID, model.CnhldEarlyCloseTrTypeID, model.BlockTransferFromDepositToWallet:
		return time.Now().Add(Day)
	default:
		return time.Now().Add(Year)
	}
}
