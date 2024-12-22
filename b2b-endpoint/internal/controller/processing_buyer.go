package controller

import (
	"errors"
	"net/http"
	"time"

	"code.emcdtech.com/b2b/processing/protocol/buyerpb"
	"code.emcdtech.com/b2b/processing/protocol/commonpb"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/sdk/log"
	"github.com/labstack/echo/v4"

	"code.emcdtech.com/b2b/endpoint/internal/business_error"
	"code.emcdtech.com/b2b/endpoint/pkg/gokit"
)

type ProcessingBuyerController struct {
	processingCli buyerpb.InvoiceBuyerServiceClient
}

func NewProcessingBuyerController(
	processingCli buyerpb.InvoiceBuyerServiceClient,
) *ProcessingBuyerController {
	return &ProcessingBuyerController{
		processingCli: processingCli,
	}
}

type getInvoiceResponse struct {
	ID                    string               `json:"id"`
	Title                 string               `json:"title"`
	Description           string               `json:"description"`
	ExpiresAt             time.Time            `json:"expires_at"`
	CoinID                string               `json:"coin_id"`
	NetworkID             string               `json:"network_id"`
	PaymentAmount         string               `json:"payment_amount"`
	PaidAmount            string               `json:"paid_amount"`
	BuyerEmail            string               `json:"buyer_email"`
	CheckoutURL           string               `json:"checkout_url"`
	Status                string               `json:"status"`
	DepositAddress        string               `json:"deposit_address"`
	BuyerFee              string               `json:"buyer_fee"`
	RequiredPaymentAmount string               `json:"required_payment_amount"`
	Transactions          []invoiceTransaction `json:"transactions"`
	CreatedAt             time.Time            `json:"created_at"`
	FinishedAt            *time.Time           `json:"finished_at"`
}

type invoiceTransaction struct {
	Hash        string    `json:"hash"`
	Amount      string    `json:"amount"`
	Address     string    `json:"address"`
	IsConfirmed bool      `json:"is_confirmed"`
	CreatedAt   time.Time `json:"created_at"`
}

func convertPbInvoice(resp *commonpb.Invoice) getInvoiceResponse {
	var finishedAt *time.Time
	if finishedAtTime := resp.FinishedAt.AsTime(); !finishedAtTime.IsZero() {
		finishedAt = &finishedAtTime
	}

	return getInvoiceResponse{
		ID:                    resp.Id,
		Title:                 resp.Title,
		Description:           resp.Description,
		ExpiresAt:             resp.ExpiresAt.AsTime(),
		CoinID:                resp.CoinId,
		NetworkID:             resp.NetworkId,
		PaymentAmount:         resp.PaymentAmount,
		PaidAmount:            resp.PaidAmount,
		BuyerEmail:            resp.BuyerEmail,
		CheckoutURL:           resp.CheckoutUrl,
		Status:                resp.Status.String(),
		DepositAddress:        resp.DepositAddress,
		BuyerFee:              resp.BuyerFee,
		RequiredPaymentAmount: resp.RequiredPaymentAmount,
		Transactions:          convertPbTransactions(resp.Transactions),
		CreatedAt:             resp.CreatedAt.AsTime(),
		FinishedAt:            finishedAt,
	}
}

// GetInvoice returns invoice details by ID
// @Summary Get invoice details
// @Description Get invoice details by invoice ID
// @Tags Processing
// @Produce json
// @Param invoice_id path string true "Invoice ID"
// @Success 200 {object} getInvoiceResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 404 {object} sdkError.Error "Invoice not found"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /processing/invoice/{invoice_id} [get]
func (p *ProcessingBuyerController) GetInvoice(c echo.Context) error {
	invoiceID := c.Param("invoice_id")
	if invoiceID == "" {
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	resp, err := p.processingCli.GetInvoice(c.Request().Context(), &buyerpb.GetInvoiceRequest{
		Id: invoiceID,
	})
	if err != nil {
		log.Error(c.Request().Context(), "getInvoice: %s", err.Error())
		// TODO: use sdk/error in processing instead of custom proto message (though keep custom business error struct)
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}

		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, convertPbInvoice(resp))
}

func convertPbTransactions(pbTxs []*commonpb.Transaction) []invoiceTransaction {
	txs := make([]invoiceTransaction, 0, len(pbTxs))
	for _, tx := range pbTxs {
		txs = append(txs, invoiceTransaction{
			Hash:        tx.Hash,
			Amount:      tx.Amount,
			Address:     tx.Address,
			IsConfirmed: tx.IsConfirmed,
			CreatedAt:   tx.CreatedAt.AsTime(),
		})
	}
	return txs
}

type getInvoiceFormResponse struct {
	ID          string     `json:"id"`
	MerchantID  string     `json:"merchant_id"`
	Title       *string    `json:"title" extensions:"x-nullable"`
	Description *string    `json:"description" extensions:"x-nullable"`
	CoinID      *string    `json:"coin_id" extensions:"x-nullable"`
	NetworkID   *string    `json:"network_id" extensions:"x-nullable"`
	Amount      *string    `json:"amount" extensions:"x-nullable"`
	BuyerEmail  *string    `json:"buyer_email" extensions:"x-nullable"`
	ExpiresAt   *time.Time `json:"expires_at" extensions:"x-nullable"`
}

// GetInvoiceForm returns invoice form details by ID
// @Summary Get invoice form details
// @Description Get invoice form details by invoice ID
// @Tags Processing
// @Produce json
// @Param form_id path string true "Form ID"
// @Success 200 {object} getInvoiceFormResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 404 {object} sdkError.Error "Invoice form not found"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /processing/invoice_form/{form_id} [get]
func (p *ProcessingBuyerController) GetInvoiceForm(c echo.Context) error {
	formID := c.Param("form_id")
	if formID == "" {
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	resp, err := p.processingCli.GetInvoiceForm(c.Request().Context(), &buyerpb.GetInvoiceFormRequest{
		Id: formID,
	})
	if err != nil {
		log.Error(c.Request().Context(), "getInvoiceForm: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	var expiresAt *time.Time
	if resp.ExpiresAt != nil {
		expiresAt = gokit.Ptr(resp.ExpiresAt.AsTime())
	}

	return c.JSON(http.StatusOK, getInvoiceFormResponse{
		ID:          resp.Id,
		MerchantID:  resp.MerchantId,
		Title:       resp.Title,
		Description: resp.Description,
		CoinID:      resp.CoinId,
		NetworkID:   resp.NetworkId,
		Amount:      resp.Amount,
		BuyerEmail:  resp.BuyerEmail,
		ExpiresAt:   expiresAt,
	})
}

type submitInvoiceFormRequest struct {
	ID          string  `json:"id"`
	MerchantID  string  `json:"merchant_id"`
	Title       *string `json:"title" extensions:"x-nullable"`
	Description *string `json:"description" extensions:"x-nullable"`
	CoinID      *string `json:"coin_id" extensions:"x-nullable"`
	NetworkID   *string `json:"network_id" extensions:"x-nullable"`
	Amount      *string `json:"amount" extensions:"x-nullable"`
	BuyerEmail  *string `json:"buyer_email" extensions:"x-nullable"`
	ExternalID  *string `json:"external_id" extensions:"x-nullable"`
	ExpiresAt   *string `json:"expires_at" extensions:"x-nullable"`
}

// SubmitInvoiceForm submits the invoice form with buyer details
// @Summary Submit invoice form
// @Description Submit invoice form with buyer details and payment preferences
// @Tags Processing
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param request body submitInvoiceFormRequest true "Invoice form submission request"
// @Success 200 {object} getInvoiceFormResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 404 {object} sdkError.Error "Invoice form not found"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /processing/invoice_form/{form_id} [put]
func (p *ProcessingBuyerController) SubmitInvoiceForm(c echo.Context) error {
	var req submitInvoiceFormRequest
	if err := c.Bind(&req); err != nil {
		log.Error(c.Request().Context(), "submitInvoiceForm: bind: %s", err.Error())
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	if err := c.Validate(&req); err != nil {
		log.Error(c.Request().Context(), "submitInvoiceForm: validate: %s", err.Error())
		return c.JSON(http.StatusBadRequest, businessError.InvalidRequest)
	}

	resp, err := p.processingCli.SubmitInvoiceForm(c.Request().Context(), &buyerpb.InvoiceForm{
		Id:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		CoinId:      req.CoinID,
		NetworkId:   req.NetworkID,
		Amount:      req.Amount,
		BuyerEmail:  req.BuyerEmail,
		MerchantId:  req.MerchantID,
	})
	if err != nil {
		log.Error(c.Request().Context(), "submitInvoiceForm: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, convertPbInvoice(resp))
}

type getBuyerFeeResponse struct {
	Amount   string `json:"amount"`
	BuyerFee string `json:"buyer_fee"`
}

// CalculateInvoicePayment
// @Summary Get amount, buyer fee
// @Description Calculate and return amount and buyer fee
// @Tags Processing
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Param raw_payment_amount query string true "Payment amount without fee"
// @Success 200 {object} getBuyerFeeResponse
// @Failure 400 {object} sdkError.Error "Invalid request"
// @Failure 500 {object} sdkError.Error "Internal server error"
// @Router /processing/merchant/{id}/calculate_invoice_payment [get]
func (p *ProcessingBuyerController) CalculateInvoicePayment(c echo.Context) error {
	merchantID := c.Param("id")

	rawAmount := c.QueryParam("raw_payment_amount")

	resp, err := p.processingCli.CalculateInvoicePayment(c.Request().Context(), &buyerpb.CalculateInvoicePaymentRequest{
		MerchantId: merchantID,
		RawAmount:  rawAmount,
	})
	if err != nil {
		log.Error(c.Request().Context(), "getBuyerFee: %s", err.Error())
		var e *sdkError.Error
		if errors.As(err, &e) {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusInternalServerError, businessError.Internal)
	}

	return c.JSON(http.StatusOK, getBuyerFeeResponse{
		Amount:   resp.Amount,
		BuyerFee: resp.BuyerFee,
	})
}
