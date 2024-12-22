package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"

	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
)

type WalletsHistory interface {
	GetHistory(ctx context.Context, params *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error)
}

type walletsHistory struct {
	repo repository.WalletsHistory
}

func NewWalletsHistory(repo repository.WalletsHistory) WalletsHistory {
	return &walletsHistory{
		repo: repo,
	}
}

func (s *walletsHistory) GetHistory(ctx context.Context, params *accountingPb.GetHistoryRequest) (*accountingPb.GetHistoryResponse, error) {
	in := &model.HistoryInput{
		Type:                "", // TODO: why empty?
		CoinCode:            "", // TODO: why empty?
		From:                params.From,
		To:                  params.To,
		Limit:               params.Limit,
		Offset:              params.Offset,
		AccountTypeIDs:      params.AccountTypeIDs,
		CoinsIDs:            params.CoinsIDs,
		CoinholdID:          params.CoinholdID,
		UserID:              params.UserID,
		TransactionTypesIDs: params.TransactionTypesIDs,
	}
	items, err := s.repo.GetWalletHistory(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("GetWalletHistory: %w", err)
	}
	var history accountingPb.GetHistoryResponse
	history.TotalCount, err = s.repo.GetWalletHistoryTotal(ctx, in)
	if err != nil {
		return &history, nil
	}
	history.Wallets = getPBWallets(items.Wallets)

	return &history, nil
}

func getPBWallets(data []*model.Wallet) []*accountingPb.Wallet {
	result := make([]*accountingPb.Wallet, len(data))
	for i, v := range data {
		item := accountingPb.Wallet{
			TxID:                  NullString(v.TxID),
			FiatStatus:            NullString(v.FiatStatus),
			Address:               NullString(v.Address),
			Comment:               NullString(v.Comment),
			CoinholdType:          v.CoinholdType,
			ExchangeToCoinID:      NullInt64(v.ExchangeToCoinID),
			CoinholdID:            NullInt64(v.CoinholdID),
			OrderID:               NullInt64(v.OrderID),
			CreatedAt:             NullInt64(v.CreatedAt),
			Amount:                NullFloat(v.Amount),
			Fee:                   NullFloat(v.Fee),
			FiatAmount:            NullFloat(v.FiatAmount),
			ExchangeAmountReceive: NullFloat(v.ExchangeAmountReceive),
			ExchangeAmountSent:    NullFloat(v.ExchangeAmountSent),
			ExchangeRate:          NullFloat(v.ExchangeRate),
			ExchangeIsSuccess:     NullBool(v.ExchangeIsSuccess),
			Date:                  timestamppb.New(v.Date),
			TokenID:               int64(v.TokenID),
			CoinID:                int64(v.CoinID),
			Status:                int64(v.Status),
			Type:                  int64(v.Type),
			Id:                    int64(v.ID),
			P2PStatus:             int64(v.P2PStatus),
			P2POrderID:            int64(v.P2POrderID),
			ReferralEmail:         NullString(v.ReferralEmail),
			ReferralType:          NullInt64(v.ReferralType),
			NetworkID:             v.NetworkID,
			CoinStrID:             v.CoinStrID,
		}
		result[i] = &item
	}

	return result
}
