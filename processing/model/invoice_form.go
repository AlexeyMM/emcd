package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InvoiceForm struct {
	ID          uuid.UUID
	MerchantID  uuid.UUID
	Title       *string
	Description *string
	CoinID      *string
	NetworkID   *string
	Amount      *decimal.Decimal
	BuyerEmail  *string
	CheckoutURL string
	// Fields below may be filled for personal invoice form
	ExternalID *string
	ExpiresAt  *time.Time
}
