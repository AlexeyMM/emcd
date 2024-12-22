package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	transactor "code.emcdtech.com/emcd/sdk/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"code.emcdtech.com/b2b/processing/internal/repository/pg/sqlc"
	"code.emcdtech.com/b2b/processing/model"
	"code.emcdtech.com/b2b/processing/pkg/gokit"
)

const foreignKeyViolationCode = "23503"

type Invoice struct {
	transactor.PgxTransactor
}

func NewInvoice(pool *pgxpool.Pool) *Invoice {
	return &Invoice{PgxTransactor: transactor.NewPgxTransactor(pool)}
}

func (i *Invoice) CreateInvoice(ctx context.Context, invoice *model.Invoice) error {
	q := sqlc.New(i.Runner(ctx))

	invoiceStatus, err := convertInvoiceStatusToDB(invoice.Status)
	if err != nil {
		return fmt.Errorf("convertInvoiceStatusToDB: %w", err)
	}

	err = q.CreateInvoice(ctx, &sqlc.CreateInvoiceParams{
		ID:              invoice.ID,
		MerchantID:      invoice.MerchantID,
		CoinID:          invoice.CoinID,
		NetworkID:       invoice.NetworkID,
		DepositAddress:  invoice.DepositAddress,
		Amount:          invoice.PaymentAmount,
		BuyerFee:        invoice.BuyerFee,
		MerchantFee:     invoice.MerchantFee,
		Title:           invoice.Title,
		Description:     invoice.Description,
		CheckoutUrl:     invoice.CheckoutURL,
		Status:          invoiceStatus,
		ExpiresAt:       invoice.ExpiresAt,
		ExternalID:      invoice.ExternalID,
		BuyerEmail:      invoice.BuyerEmail,
		RequiredPayment: invoice.RequiredPayment,
	})
	if err != nil {
		return fmt.Errorf("createInvoice: %w", err)
	}

	return nil
}

var invoiceStatusMap = map[model.InvoiceStatus]sqlc.InvoiceStatus{
	model.InvoiceStatusWaitingForDeposit:   sqlc.InvoiceStatusWaitingForDeposit,
	model.InvoiceStatusPaymentConfirmation: sqlc.InvoiceStatusPaymentConfirmation,
	model.InvoiceStatusPartiallyPaid:       sqlc.InvoiceStatusPartiallyPaid,
	model.InvoiceStatusPaymentAccepted:     sqlc.InvoiceStatusPaymentAccepted,
	model.InvoiceStatusFinished:            sqlc.InvoiceStatusFinished,
	model.InvoiceStatusExpired:             sqlc.InvoiceStatusExpired,
	model.InvoiceStatusCancelled:           sqlc.InvoiceStatusCancelled,
}

func convertInvoiceStatusToDB(status model.InvoiceStatus) (sqlc.InvoiceStatus, error) {
	if convertedStatus, ok := invoiceStatusMap[status]; ok {
		return convertedStatus, nil
	}

	return "", fmt.Errorf("invalid invoice status: %v", status)
}

func (i *Invoice) GetInvoice(ctx context.Context, invoiceID uuid.UUID) (*model.Invoice, error) {
	q := sqlc.New(i.Runner(ctx))

	invoiceRow, err := q.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, &model.Error{
			Code:    model.ErrorCodeNoSuchInvoice,
			Message: fmt.Sprintf("no such invoice: %s", invoiceID),
			Inner:   err,
		}
	}

	return convertDBInvoiceToModel(invoiceRow)
}

func convertDBInvoiceToModel(invoiceRow sqlc.GetInvoiceRow) (*model.Invoice, error) {
	invoiceStatus, err := convertInvoiceStatusFromDB(invoiceRow.Status)
	if err != nil {
		return nil, fmt.Errorf("convertInvoiceStatusFromDB: %w", err)
	}

	var finishedAt time.Time

	if invoiceRow.FinishedAt != nil {
		finishedAt = *invoiceRow.FinishedAt
	}

	return &model.Invoice{
		ID:              invoiceRow.ID,
		MerchantID:      invoiceRow.MerchantID,
		CoinID:          invoiceRow.CoinID,
		NetworkID:       invoiceRow.NetworkID,
		DepositAddress:  invoiceRow.DepositAddress,
		PaymentAmount:   invoiceRow.Amount,
		BuyerFee:        invoiceRow.BuyerFee,
		MerchantFee:     invoiceRow.MerchantFee,
		PaidAmount:      invoiceRow.PaidAmount,
		Title:           invoiceRow.Title,
		Description:     invoiceRow.Description,
		CheckoutURL:     invoiceRow.CheckoutUrl,
		Status:          invoiceStatus,
		ExpiresAt:       invoiceRow.ExpiresAt,
		ExternalID:      invoiceRow.ExternalID,
		BuyerEmail:      invoiceRow.BuyerEmail,
		RequiredPayment: invoiceRow.RequiredPayment,
		CreatedAt:       invoiceRow.CreatedAt,
		FinishedAt:      finishedAt,
	}, nil
}

var dbInvoiceStatusMap = gokit.Invert(invoiceStatusMap)

func convertInvoiceStatusFromDB(status sqlc.InvoiceStatus) (model.InvoiceStatus, error) {
	if convertedStatus, ok := dbInvoiceStatusMap[status]; ok {
		return convertedStatus, nil
	}

	return "", fmt.Errorf("invalid db invoice status: %v", status)
}

func (i *Invoice) CreateInvoiceForm(ctx context.Context, form *model.InvoiceForm) error {
	q := sqlc.New(i.Runner(ctx))

	err := q.CreateInvoiceForm(ctx, &sqlc.CreateInvoiceFormParams{
		ID:          form.ID,
		MerchantID:  form.MerchantID,
		Title:       form.Title,
		Description: form.Description,
		CoinID:      form.CoinID,
		NetworkID:   form.NetworkID,
		Amount:      form.Amount,
		BuyerEmail:  form.BuyerEmail,
		CheckoutUrl: form.CheckoutURL,
		ExternalID:  form.ExternalID,
		ExpiresAt:   form.ExpiresAt,
	})
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.ConstraintName == "invoice_form_merchant_id_fkey" && pqErr.Code == foreignKeyViolationCode { //nolint:lll
			return &model.Error{
				Code:    model.ErrorCodeNoSuchMerchant,
				Message: fmt.Sprintf("no such merchant: %s", form.MerchantID),
				Inner:   pqErr,
			}
		}

		return fmt.Errorf("createInvoiceForm: %w", err)
	}

	return nil
}

func (i *Invoice) GetInvoiceForm(ctx context.Context, id uuid.UUID) (*model.InvoiceForm, error) {
	q := sqlc.New(i.Runner(ctx))

	form, err := q.GetInvoiceForm(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Error{
				Code:    model.ErrorCodeNoSuchInvoiceForm,
				Message: fmt.Sprintf("no such invoice form: %s", id),
				Inner:   err,
			}
		}

		return nil, fmt.Errorf("getInvoiceForm: %w", err)
	}

	return &model.InvoiceForm{
		ID:          form.ID,
		MerchantID:  form.MerchantID,
		Title:       form.Title,
		Description: form.Description,
		CoinID:      form.CoinID,
		NetworkID:   form.NetworkID,
		Amount:      form.Amount,
		BuyerEmail:  form.BuyerEmail,
		CheckoutURL: form.CheckoutUrl,
		ExternalID:  form.ExternalID,
		ExpiresAt:   form.ExpiresAt,
	}, nil
}

func (i *Invoice) GetActiveInvoiceByDepositAddressForUpdate(
	ctx context.Context,
	address string,
) (*model.Invoice, error) {
	q := sqlc.New(i.Runner(ctx))

	invoice, err := q.GetActiveInvoiceByDepositAddressForUpdate(ctx, address)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Error{Code: model.ErrorCodeNoSuchInvoice}
		}

		return nil, fmt.Errorf("getActiveInvoiceByDepositAddressForUpdate: %w", err)
	}

	// NOTE: a bit hacky solution, on the other side it guarantees that both requests return the same struct
	return convertDBInvoiceToModel(sqlc.GetInvoiceRow(invoice))
}

func (i *Invoice) UpdateStatus(ctx context.Context, id uuid.UUID, status model.InvoiceStatus) error {
	q := sqlc.New(i.Runner(ctx))

	dbStatus, err := convertInvoiceStatusToDB(status)
	if err != nil {
		return fmt.Errorf("convertInvoiceStatusToDB: %w", err)
	}

	err = q.UpdateStatus(ctx, &sqlc.UpdateStatusParams{
		ID:     id,
		Status: dbStatus,
	})
	if err != nil {
		return fmt.Errorf("updateStatus: %w", err)
	}

	return nil
}

func (i *Invoice) SetInvoicesExpired(ctx context.Context) error {
	q := sqlc.New(i.Runner(ctx))

	affected, err := q.SetInvoicesExpired(ctx)
	if err != nil {
		return fmt.Errorf("setInvoicesExpired: %w", err)
	}

	log.SInfo(ctx, "expired invoices updated", map[string]any{
		"affected_count": affected,
	})

	return nil
}

func (i *Invoice) CountInvoiceByStatus(ctx context.Context) (map[model.InvoiceStatus]int, error) {
	q := sqlc.New(i.Runner(ctx))

	rawStatuses, err := q.CountInvoiceByStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("countInvoiceByStatus: %w", err)
	}

	statuses, err := convertDBCountInvoiceByStatusRow(rawStatuses)
	if err != nil {
		return nil, fmt.Errorf("convertDBCountInvoiceByStatusRow: %w", err)
	}

	return statuses, nil
}

func convertDBCountInvoiceByStatusRow(
	statuses []sqlc.CountInvoiceByStatusRow,
) (map[model.InvoiceStatus]int, error) {
	m := make(map[model.InvoiceStatus]int, len(statuses))

	for _, row := range statuses {
		invoiceStatus, err := convertInvoiceStatusFromDB(row.Status)
		if err != nil {
			return nil, fmt.Errorf("convertInvoiceStatusFromDB: %w", err)
		}

		m[invoiceStatus] = int(row.InvoiceCount)
	}

	return m, nil
}
