package service

import (
	sdkLog "code.emcdtech.com/emcd/sdk/log"
	"code.emcdtech.com/emcd/service/accounting/internal/repository"
	"code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/pkg/kafka"
	accountingPb "code.emcdtech.com/emcd/service/accounting/protocol/accounting"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

type OutboxTransaction interface {
	ExportToKafka(ctx context.Context) error
}

type outboxTransaction struct {
	pool                 *pgxpool.Pool
	writer               kafka.RawWriter
	balance              repository.Balance
	outboxTransactions   repository.OutboxTransactions
	iterationRecordLimit uint
}

func NewOutboxTransaction(
	pool *pgxpool.Pool,
	balance repository.Balance,
	writer kafka.RawWriter,
	outboxTransactions repository.OutboxTransactions,
	iterationRecordLimit uint,
) OutboxTransaction {
	return &outboxTransaction{
		pool:                 pool,
		balance:              balance,
		writer:               writer,
		outboxTransactions:   outboxTransactions,
		iterationRecordLimit: iterationRecordLimit,
	}
}

func (s *outboxTransaction) ExportToKafka(ctx context.Context) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			sdkLog.Error(ctx, "recover from panic: %v", p)
			if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				sdkLog.Error(ctx, err.Error())

			}
		}
	}()

	err = s.exportToKafkaWithTx(ctx, tx)
	if err != nil {
		errInt := tx.Rollback(ctx)
		if errInt != nil && !errors.Is(errInt, pgx.ErrTxClosed) {
			sdkLog.Error(ctx, "rollback transaction: %v", errInt)
		}
		return fmt.Errorf("export to kafka: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}

func (s *outboxTransaction) exportToKafkaWithTx(ctx context.Context, tx pgx.Tx) error {
	ids, err := s.outboxTransactions.List(ctx, tx, s.iterationRecordLimit)
	if err != nil {
		return fmt.Errorf("listing outbox transactions: %w", err)
	}
	if len(ids) == 0 {
		return nil
	}
	transactions, err := s.balance.GetTransactionByIDs(ctx, tx, ids...)
	if err != nil {
		return fmt.Errorf("get transaction by ids: %w", err)
	}
	messages, err := s.trxsToMessages(transactions)
	if err != nil {
		return fmt.Errorf("trx to message: %w", err)
	}
	err = s.writer.WriteMessages(ctx, messages...)
	if err != nil {
		return fmt.Errorf("write messages: %w", err)
	}
	err = s.outboxTransactions.Delete(ctx, tx, ids...)
	if err != nil {
		return fmt.Errorf("delete outbox transactions: %w", err)
	}
	return nil
}

func (s *outboxTransaction) trxsToMessages(transactions []*model.Transaction) ([]kafka.Message, error) {
	opKafka, err := kafka.GetKafka[accountingPb.TransactionCreatedEvent]()
	if err != nil {
		return nil, fmt.Errorf("get option field: %w", err)
	}
	result := make([]kafka.Message, 0, len(transactions))
	for _, t := range transactions {
		pbTrx := s.trxToPBTrx(t)

		b, err := protojson.Marshal(pbTrx)
		if err != nil {
			return nil, fmt.Errorf("marshal outboxTransaction: %w", err)
		}
		m := kafka.Message{
			ID:      uuid.UUID{},
			Key:     []byte(strconv.FormatInt(t.ID, 10)),
			Topic:   opKafka.Topic,
			Headers: nil,
			Value:   b,
		}
		result = append(result, m)
	}
	return result, nil
}

func (s *outboxTransaction) trxToPBTrx(t *model.Transaction) *accountingPb.TransactionCreatedEvent {
	return &accountingPb.TransactionCreatedEvent{
		Transaction: &accountingPb.Transaction{
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
		},
	}
}
