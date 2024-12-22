package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	InvoiceID   uuid.UUID
	Address     string // receiver-address
	Amount      decimal.Decimal
	Hash        string // blockchain tx hash
	CoinID      string
	IsConfirmed bool
	CreatedAt   time.Time
}

type InvoiceStatus string

const (
	InvoiceStatusUnknown             InvoiceStatus = "UNKNOWN"
	InvoiceStatusWaitingForDeposit   InvoiceStatus = "WAITING_FOR_DEPOSIT"
	InvoiceStatusPaymentConfirmation InvoiceStatus = "PAYMENT_CONFIRMATION"
	InvoiceStatusPartiallyPaid       InvoiceStatus = "PARTIALLY_PAID"
	InvoiceStatusPaymentAccepted     InvoiceStatus = "PAYMENT_ACCEPTED"
	InvoiceStatusFinished            InvoiceStatus = "FINISHED"
	InvoiceStatusCancelled           InvoiceStatus = "CANCELLED"
	InvoiceStatusExpired             InvoiceStatus = "EXPIRED"
)

type Invoice struct {
	ID              uuid.UUID
	MerchantID      uuid.UUID
	ExternalID      string
	Title           string
	Description     string
	ExpiresAt       time.Time
	CoinID          string
	NetworkID       string
	PaymentAmount   decimal.Decimal
	BuyerFee        decimal.Decimal // from merchant's tariff upper fee
	MerchantFee     decimal.Decimal // from merchant's tariff lower fee
	RequiredPayment decimal.Decimal // this is how much a buyer will pay
	PaidAmount      decimal.Decimal
	BuyerEmail      string
	CheckoutURL     string
	Status          InvoiceStatus
	DepositAddress  string
	Transactions    []*Transaction
	CreatedAt       time.Time
	FinishedAt      time.Time
}

type CreateInvoiceRequest struct {
	ExternalID  string
	Title       string
	Description string
	CoinID      string
	NetworkID   string
	Amount      decimal.Decimal
	BuyerEmail  string
	CheckoutURL string
	MerchantID  uuid.UUID
	ExpiresAt   time.Time
}
