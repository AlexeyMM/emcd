package handler

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/model"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

func (h *AccountingHandler) ListTransactions(
	ctx context.Context,
	request *accountingPb.ListTransactionsRequest,
) (*accountingPb.ListTransactionsResponse, error) {
	const limit = 200
	err := validation.ValidateStructWithContext(ctx, request,
		validation.Field(&request.From, validation.Required),
		validation.Field(&request.To, validation.Required),
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}

	receiverAccountIds := h.sliceInt64ToInt(request.ReceiverAccountIds)
	types := h.sliceInt64ToInt(request.Types)
	from := request.From.AsTime()
	to := request.To.AsTime()
	// page_token = 5;
	trxs, lastTrxId, err := h.transactionService.ListTransactions(
		ctx,
		receiverAccountIds,
		types,
		from,
		to,
		limit,
		request.FromTransactionId,
	)
	if err != nil {
		sdkLog.Error(ctx, "service: %v", err)
		return nil, err
	}
	if len(trxs) < limit {
		lastTrxId = 0
	}

	return &accountingPb.ListTransactionsResponse{
		Transactions:      h.trxsToPBTrxs(trxs),
		LastTransactionId: lastTrxId,
	}, nil
}

func (h *AccountingHandler) sliceInt64ToInt(s []int64) []int {
	r := make([]int, len(s))
	for i := range s {
		r[i] = int(s[i])
	}
	return r
}

func (h *AccountingHandler) trxsToPBTrxs(t []*model.Transaction) []*accountingPb.Transaction {
	r := make([]*accountingPb.Transaction, 0, len(t))
	for i := range t {
		r = append(r, h.trxToPBTrx(t[i]))
	}
	return r
}

func (h *AccountingHandler) trxToPBTrx(t *model.Transaction) *accountingPb.Transaction {
	return &accountingPb.Transaction{
		Type:              t.Type.Int64(),
		CreatedAt:         timestamppb.New(t.CreatedAt),
		SenderAccountID:   t.SenderAccountID,
		ReceiverAccountID: t.ReceiverAccountID,
		CoinID:            strconv.FormatInt(t.CoinID, 10),
		Amount:            t.Amount.String(),
		Comment:           t.Comment,
		Hash:              t.Hash,
		Hashrate:          t.Hashrate,
		FromReferralId:    t.FromReferralId,
		ReceiverAddress:   t.ReceiverAddress,
		TokenID:           t.TokenID,
		ActionID:          t.ActionID,
	}
}
