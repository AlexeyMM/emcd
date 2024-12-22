package fee

import (
	"context"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	accountingmodel "code.emcdtech.com/emcd/service/accounting/model"
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	accounting "code.emcdtech.com/emcd/service/accounting/repository"
	profilepb "code.emcdtech.com/emcd/service/profile/protocol/profile"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"code.emcdtech.com/b2b/processing/model"
)

type accountingRepos struct {
	userAccountRepo accounting.UserAccountRepository
	balanceRepo     accounting.AccountingRepository
}

type Service struct {
	profileClient profilepb.ProfileServiceClient
	accounting    accountingRepos

	feeCollectorUserID uuid.UUID
}

func NewService(
	profileClient profilepb.ProfileServiceClient,
	userAccountRepo accounting.UserAccountRepository,
	balanceRepo accounting.AccountingRepository,
	feeCollectorUserID uuid.UUID,
) *Service {
	return &Service{
		profileClient: profileClient,
		accounting: accountingRepos{
			userAccountRepo: userAccountRepo,
			balanceRepo:     balanceRepo,
		},
		feeCollectorUserID: feeCollectorUserID,
	}
}

func (s *Service) ChargeFeeForInvoice(ctx context.Context, invoice *model.Invoice) error {
	ctx, span := otel.Tracer("").Start(ctx, "charge fee for invoice")
	defer span.End()

	// check invoice status in case of some retry processing
	if invoice.Status != model.InvoiceStatusPaymentAccepted {
		log.SWarn(
			ctx,
			"invoice fee cannot be charged in current status, must be InvoiceStatusPaymentAccepted",
			map[string]any{"invoice": invoice},
		)

		return nil
	}

	// TODO: we can request both accounts in parallel

	merchantAccount, err := s.getAccount(ctx, invoice.MerchantID, invoice.CoinID)
	if err != nil {
		return fmt.Errorf("getAccount for merchant: %w", err)
	}

	// TODO: this can be optimized if we store account for repeated use after we retrieve it first time
	feeCollectorAccount, err := s.getAccount(ctx, s.feeCollectorUserID, invoice.CoinID)
	if err != nil {
		return fmt.Errorf("getAccount for feeCollector: %w", err)
	}

	err = s.accounting.balanceRepo.ChangeBalance(ctx, []*accountingmodel.Transaction{
		{
			Type:              accountingmodel.ProcessingInvoiceFee,
			CreatedAt:         time.Now(),
			SenderAccountID:   int64(merchantAccount.ID),
			ReceiverAccountID: int64(feeCollectorAccount.ID),
			CoinID:            int64(merchantAccount.CoinID),
			Amount:            invoice.MerchantFee.Add(invoice.BuyerFee),
			Comment:           "fee for invoice " + invoice.ID.String(),
			ActionID:          invoice.ID.String(), // idempotency key
		},
	})
	if err != nil {
		return fmt.Errorf("changeBalance: %w", err)
	}

	return nil
}

func (s *Service) getAccount(
	ctx context.Context,
	userID uuid.UUID,
	coinID string,
) (*accountingmodel.UserAccount, error) {
	ctx, span := otel.Tracer("coinwatch-consumer").Start(
		ctx,
		"get user account",
		oteltrace.WithAttributes(semconv.UserID(userID.String())),
		oteltrace.WithAttributes(attribute.Key("coin.id").String(coinID)),
	)
	defer span.End()

	// TODO: this must be handled on the side of accounting, push this change. For now it is here
	oldIDResp, err := s.profileClient.GetOldIDByID(ctx, &profilepb.GetOldIDByIDRequest{
		Id: userID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("getOldIDByID: %w", err)
	}

	oldID := oldIDResp.GetOldId()

	acc, err := s.accounting.userAccountRepo.GetOrCreateUserAccountByArgs(
		ctx,
		oldID,
		userID,
		coinID,
		enum.WalletAccountTypeID,
		0,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("getOrCreateUserAccountByArgs: %w", err)
	}

	return acc, nil
}
