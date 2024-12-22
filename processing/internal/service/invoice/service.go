package invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/b2b/processing/internal/repository"
	"code.emcdtech.com/b2b/processing/internal/service"
	"code.emcdtech.com/b2b/processing/model"
)

type ServiceConfig struct {
	InvoiceTTL time.Duration
}

type Service struct {
	config          *ServiceConfig
	merchantRepo    repository.Merchant
	invoiceRepo     repository.Invoice
	addressPool     service.DepositAddressPool
	transactionRepo repository.Transaction
	coinRepository  repository.Coin
}

func NewService(
	config *ServiceConfig,
	merchantRepo repository.Merchant,
	invoiceRepo repository.Invoice,
	addressPool service.DepositAddressPool,
	transactionRepo repository.Transaction,
	coinRepository repository.Coin,
) *Service {
	return &Service{
		config:          config,
		merchantRepo:    merchantRepo,
		invoiceRepo:     invoiceRepo,
		addressPool:     addressPool,
		transactionRepo: transactionRepo,
		coinRepository:  coinRepository,
	}
}

func (s *Service) CreateInvoice(ctx context.Context, req *model.CreateInvoiceRequest) (*model.Invoice, error) {
	_, err := s.coinRepository.GetNetwork(req.CoinID, req.NetworkID)
	if err != nil {
		return nil, fmt.Errorf("getNetwork: %w", err)
	}

	merchant, err := s.merchantRepo.Get(ctx, req.MerchantID)
	if err != nil {
		return nil, fmt.Errorf("merchantRepo.Get: %w", err)
	}

	if req.Amount.LessThan(merchant.Tariff.MinPay) {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: "amount must be greater than or equal to " + merchant.Tariff.MinPay.String(),
		}
	}

	if req.Amount.GreaterThan(merchant.Tariff.MaxPay) {
		return nil, &model.Error{
			Code:    model.ErrorCodeInvalidArgument,
			Message: "amount must be less than or equal to " + merchant.Tariff.MaxPay.String(),
		}
	}
	// Use provided expiration time or calculate from TTL
	expiresAt := req.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(s.config.InvoiceTTL)
	}

	invoice := &model.Invoice{
		ID:              uuid.New(),
		MerchantID:      merchant.ID,
		ExternalID:      req.ExternalID,
		Title:           req.Title,
		Description:     req.Description,
		ExpiresAt:       expiresAt,
		CoinID:          req.CoinID,
		NetworkID:       req.NetworkID,
		PaymentAmount:   req.Amount,
		BuyerFee:        req.Amount.Mul(merchant.Tariff.UpperFee),
		MerchantFee:     req.Amount.Mul(merchant.Tariff.LowerFee),
		RequiredPayment: req.Amount.Add(req.Amount.Mul(merchant.Tariff.UpperFee)),
		PaidAmount:      decimal.Zero,
		BuyerEmail:      req.BuyerEmail,
		CheckoutURL:     req.CheckoutURL,
		Status:          model.InvoiceStatusWaitingForDeposit,
		DepositAddress:  "", // TODO: get from address pool below
	}

	err = s.invoiceRepo.WithinTransactionWithOptions(ctx, func(ctx context.Context) error {
		depositAddress, err := s.addressPool.GetOrCreate(ctx, merchant.ID, invoice.NetworkID, invoice.ID)
		if err != nil {
			return fmt.Errorf("addressPool.GetOrCreate: %w", err)
		}

		invoice.DepositAddress = depositAddress

		if err := s.invoiceRepo.CreateInvoice(ctx, invoice); err != nil {
			return fmt.Errorf("invoiceRepo.CreateInvoice: %w", err)
		}

		return nil
	}, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *Service) CreateInvoiceForm(ctx context.Context, form *model.InvoiceForm) (*model.InvoiceForm, error) {
	if err := s.invoiceRepo.CreateInvoiceForm(ctx, form); err != nil {
		return nil, fmt.Errorf("createInvoiceForm: %w", err)
	}

	return form, nil
}

func (s *Service) GetInvoiceForm(ctx context.Context, id uuid.UUID) (*model.InvoiceForm, error) {
	form, err := s.invoiceRepo.GetInvoiceForm(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cetInvoiceForm: %w", err)
	}

	return form, nil
}

func (s *Service) SubmitInvoiceForm(ctx context.Context, submission *model.InvoiceForm) (*model.Invoice, error) {
	existingForm, err := s.invoiceRepo.GetInvoiceForm(ctx, submission.ID)
	if err != nil {
		return nil, fmt.Errorf("getInvoiceForm: %w", err)
	}

	// Check expiration for personalized forms
	if existingForm.ExternalID != nil && existingForm.ExpiresAt != nil {
		if time.Now().After(*existingForm.ExpiresAt) {
			return nil, &model.Error{
				Code:    model.ErrorCodeInvoiceFormExpired,
				Message: fmt.Sprintf("invoice form %s has expired", existingForm.ID),
			}
		}
	}

	// Merge form data with submission data
	req := &model.CreateInvoiceRequest{
		MerchantID:  existingForm.MerchantID,
		CheckoutURL: existingForm.CheckoutURL,
		Title:       coalesceOrZero(existingForm.Title, submission.Title),
		Description: coalesceOrZero(existingForm.Description, submission.Description),
		CoinID:      coalesceOrZero(existingForm.CoinID, submission.CoinID),
		NetworkID:   coalesceOrZero(existingForm.NetworkID, submission.NetworkID),
		Amount:      coalesceOrZero(existingForm.Amount, submission.Amount),
		BuyerEmail:  coalesceOrZero(existingForm.BuyerEmail, submission.BuyerEmail),
		ExternalID:  coalesceOrZero(existingForm.ExternalID, nil),
		ExpiresAt:   coalesceOrZero(existingForm.ExpiresAt, nil),
	}

	return s.CreateInvoice(ctx, req)
}

func (s *Service) GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	invoice, err := s.invoiceRepo.GetInvoice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getInvoice: %w", err)
	}

	invoice.Transactions, err = s.transactionRepo.GetInvoiceTransactions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getInvoiceTransactions: %w", err)
	}

	return invoice, nil
}

func (s *Service) CalculateInvoicePayment(
	ctx context.Context,
	merchantID uuid.UUID,
	rawPayment decimal.Decimal,
) (decimal.Decimal, decimal.Decimal, error) {
	merchant, err := s.merchantRepo.Get(ctx, merchantID)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("get: %w", err)
	}

	buyerFee := rawPayment.Mul(merchant.Tariff.UpperFee)

	return rawPayment.Add(buyerFee), buyerFee, nil
}

func coalesceOrZero[T any](existing, submitted *T) T {
	if existing != nil {
		return *existing
	}

	if submitted != nil {
		return *submitted
	}

	var zero T

	return zero
}
