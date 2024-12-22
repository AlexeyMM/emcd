package repository

import (
	businessErr "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

type HistoryRepository interface {
	GetHistory(ctx context.Context, client accountingPb.AccountingServiceClient, data *model.HistoryInput) (*model.HistoryOutput, error)
}

type historyRepository struct{}

func NewHistoryRepository() HistoryRepository {
	return &historyRepository{}
}

func (*historyRepository) GetHistory(ctx context.Context, client accountingPb.AccountingServiceClient,
	data *model.HistoryInput) (*model.HistoryOutput, error) {
	var req = accountingPb.GetHistoryRequest{
		Type:                string(data.Type),
		CoinCode:            data.CoinCode,
		From:                data.From,
		To:                  data.To,
		Limit:               data.Limit,
		Offset:              data.Offset,
		CoinholdID:          data.CoinholdID,
		UserID:              data.UserID,
		TransactionTypesIDs: data.TransactionTypesIDs,
		AccountTypeIDs:      data.AccountTypeIDs,
		CoinsIDs:            data.CoinsIDs,
	}

	resp, err := client.GetHistory(ctx, &req)
	if err != nil {
		return nil, businessErr.NewError("", fmt.Sprintf("accounting: %s", err.Error()))
	}

	result, err := getHistoryOutput(resp, model.HistoryType(req.Type))
	if err != nil {
		return nil, businessErr.NewError("", fmt.Sprintf("getHistoryOutput: %s", err.Error()))
	}
	return result, nil
}

func getHistoryOutput(resp *accountingPb.GetHistoryResponse, historyType model.HistoryType) (*model.HistoryOutput, error) {
	var (
		incomesSum *decimal.Decimal
		payoutsSum *decimal.Decimal
	)
	if resp.IncomesSum != "" {
		val, err := decimal.NewFromString(resp.IncomesSum)
		if err != nil {
			return nil, fmt.Errorf("invalid IncomesSum: %w", err)
		}
		incomesSum = &val
	}
	if resp.PayoutsSum != "" {
		val, err := decimal.NewFromString(resp.PayoutsSum)
		if err != nil {
			return nil, fmt.Errorf("invalid PayoutsSum: %w", err)
		}
		payoutsSum = &val
	}
	incomes, err := getIncomes(resp.Incomes, historyType)
	if err != nil {
		return nil, fmt.Errorf("getIncomes: %w", err)
	}
	payouts, err := getPayouts(resp.Payouts, historyType)
	if err != nil {
		return nil, fmt.Errorf("getPayouts: %w", err)
	}
	wallets, err := getWallets(resp.Wallets, historyType)
	if err != nil {
		return nil, fmt.Errorf("getPayouts: %w", err)
	}

	result := &model.HistoryOutput{
		TotalCount:    int(resp.TotalCount),
		IncomesSum:    incomesSum,
		PayoutsSum:    payoutsSum,
		HasNewIncome:  GetValFromNullBool(resp.HasNewIncome),
		HasNewPayouts: GetValFromNullBool(resp.HasNewPayouts),
		Incomes:       incomes,
		Payouts:       payouts,
		Wallets:       wallets,
	}

	return result, nil
}

func getIncomes(data []*accountingPb.Income, historyType model.HistoryType) ([]*model.Income, error) {
	if data == nil {
		if historyType == model.HistoryIncome {
			return []*model.Income{}, nil
		}
		return nil, nil
	}

	var (
		err     error
		incomes = make([]*model.Income, len(data))
	)

	for i, v := range data {
		incomes[i] = new(model.Income)

		incomes[i].Diff = int(v.Diff)
		incomes[i].ChangePercent, err = strconv.ParseFloat(v.ChangePercent, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ChangePercent: %w", err)
		}
		incomes[i].Time, err = strconv.ParseFloat(v.Time, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid IncomeDesktop.Time: %w", err)
		}
		incomes[i].Income, err = decimal.NewFromString(v.Income)
		if err != nil {
			return nil, fmt.Errorf("invalid IncomeDesktop.Income: %w", err)
		}
		incomes[i].Code = int(v.Code)
		incomes[i].HashRate = GetValFromNullInt64(v.HashRate)
	}

	return incomes, nil
}

func getPayouts(data []*accountingPb.Payout, historyType model.HistoryType) ([]*model.Payout, error) {
	if data == nil {
		if historyType == model.HistoryPayout {
			return []*model.Payout{}, nil
		}
		return nil, nil
	}

	var (
		err     error
		payouts = make([]*model.Payout, len(data))
	)

	for i, v := range data {
		payouts[i] = new(model.Payout)

		payouts[i].Time, err = strconv.ParseFloat(v.Time, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Time: %w", err)
		}
		payouts[i].Amount, err = strconv.ParseFloat(v.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Amount: %w", err)
		}
		payouts[i].Tx = GetValFromNullString(v.Tx)
		payouts[i].TxId = GetValFromNullString(v.TxID)
	}

	return payouts, nil
}

func getWallets(data []*accountingPb.Wallet, historyType model.HistoryType) ([]*model.Wallet, error) {
	if data == nil {
		if historyType == model.HistoryWallet || historyType == model.HistoryCoinhold {
			return []*model.Wallet{}, nil
		}
		return nil, nil
	}

	var (
		err     error
		wallets = make([]*model.Wallet, len(data))
	)

	for i, v := range data {
		wallets[i] = new(model.Wallet)

		wallets[i].TxID = GetValFromNullString(v.TxID)
		wallets[i].FiatStatus = GetValFromNullString(v.FiatStatus)
		wallets[i].Address = GetValFromNullString(v.Address)
		wallets[i].Comment = GetValFromNullString(v.Comment)
		wallets[i].CoinholdType = v.CoinholdType
		wallets[i].ExchangeToCoinID = GetValFromNullInt64(v.ExchangeToCoinID)
		wallets[i].CoinholdID = GetValFromNullInt64(v.CoinholdID)
		wallets[i].OrderID = GetValFromNullInt64(v.OrderID)
		wallets[i].CreatedAt = GetValFromNullInt64(v.CreatedAt)
		if wallets[i].Amount, err = GetValFromNullFloat(v.Amount); err != nil && !errors.Is(err, errNoError) {
			return nil, fmt.Errorf("getValFromNullFloat Amount: %w", err)
		}
		if wallets[i].Fee, err = GetValFromNullFloat(v.Fee); err != nil && !errors.Is(err, errNoError) {
			return nil, fmt.Errorf("getValFromNullFloat Fee: %w", err)
		}
		if wallets[i].FiatAmount, err = GetValFromNullFloat(v.FiatAmount); err != nil && !errors.Is(err, errNoError) {
			return nil, fmt.Errorf("getValFromNullFloat FiatAmount: %w", err)
		}
		if wallets[i].ExchangeAmountReceive, err = GetValFromNullFloat(v.ExchangeAmountReceive); err != nil && !errors.Is(err, errNoError) {
			return nil, fmt.Errorf("getValFromNullFloat ExchangeAmountReceive: %w", err)
		}
		if wallets[i].ExchangeAmountSent, err = GetValFromNullFloat(v.ExchangeAmountSent); err != nil && !errors.Is(err, errNoError) {
			return nil, fmt.Errorf("getValFromNullFloat ExchangeAmountSent: %w", err)
		}
		if wallets[i].ExchangeRate, err = GetValFromNullFloat(v.ExchangeRate); err != nil && !errors.Is(err, errNoError) {
			return nil, fmt.Errorf("getValFromNullFloat ExchangeRate: %w", err)
		}
		wallets[i].ExchangeIsSuccess = GetValFromNullBool(v.ExchangeIsSuccess)
		wallets[i].Date = v.Date.AsTime()
		wallets[i].TokenID = int(v.TokenID)
		wallets[i].CoinID = int(v.CoinID)
		wallets[i].Status = int(v.Status)
		wallets[i].Type = int(v.Type)
		wallets[i].ID = int(v.Id)
		wallets[i].P2PStatus = int(v.P2PStatus)
		wallets[i].P2POrderID = int(v.P2POrderID)
		wallets[i].ReferralEmail = GetValFromNullString(v.ReferralEmail)
		wallets[i].ReferralType = GetValFromNullInt64(v.ReferralType)
		wallets[i].NetworkID = v.NetworkID
		wallets[i].CoinStrID = v.CoinStrID
	}

	return wallets, nil
}
