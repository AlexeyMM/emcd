package model

import (
	"code.emcdtech.com/emcd/service/accounting/model/enum"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	Email                 string
	Username              string
	Password              string
	Nopay                 bool
	IsEmailNotificationOn bool
	IsTglNotificationOn   bool
	TgID                  int
	AuthSecret            string
	ApiKey                string
	Language              string
	Timezone              string
	IsDonationOn          bool
	IsCoinholdEnabled     bool
	IsEmployee            bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	ID                    int64
	RefID                 int
	ParentID              int
	MasterFee             float64
	// phone                 string
	IsPhoneVerified   bool
	TwaCode           string
	IsHedgeEnabled    bool
	CompanyName       string
	IsAutopayDisabled bool
	KycStatus         string
	WbLinkID          string
	PassUpdatedAy     string
	PrimaryCurrency   string
	IsActive          bool
	FreeWithdraw      bool
	DefCoinID         int
	NewID             uuid.UUID
}

type UserBeforePayoutMining struct {
	UserID        int64
	CoinID        int64
	BlockID       int64
	AccountTypeID enum.AccountTypeId
	LastPay       time.Time
}

type UserBeforePayoutWallet struct {
	UserID         int64
	CoinID         int64
	TransactionIDs []int64
	AccountTypeID  enum.AccountTypeId
}

type UserWalletDiff struct {
	UserID  int64
	BlockID int64
	Diff    decimal.Decimal
}
