package model

import (
	"fmt"
	"time"

	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID    uuid.UUID
	OldID int32

	// атрибуты - поля, изменяющиеся напрямую без дополнительной логики
	Username                 string
	ParentID                 uuid.UUID
	Language                 string
	WhiteLabelID             uuid.UUID
	PoolType                 string
	OldParentID              *int32
	WasReferralLinkGenerated bool // не обновляется из-за обратной совместимости
	IsAmbassador             bool

	// поля, требующие вызова дополнительной логики
	RefID    int
	NewRefID uuid.UUID

	Email    string
	Password string
	ApiKey   string
	SecretKey    string

	IsActive  bool // доступ разрешен
	Suspended bool // частичный доступ, заморозка денег

	AppleID    string
	TgID       string
	TgUsername string

	// хранится не в сервисе profile
	Vip       bool
	SegmentID int

	CreatedAt time.Time
}

type UserAttributes struct {
	Username                 *string
	ParentID                 *uuid.UUID
	Language                 *string
	WhiteLabelID             *uuid.UUID
	PoolType                 *string
	WasReferralLinkGenerated *bool
	IsAmbassador             *bool
}

func (u *User) SetAttributes(attrs UserAttributes) {
	if attrs.Username != nil {
		u.Username = *attrs.Username
	}
	if attrs.ParentID != nil {
		u.ParentID = *attrs.ParentID
	}
	if attrs.Language != nil {
		u.Language = *attrs.Language
	}
	if attrs.WhiteLabelID != nil {
		u.WhiteLabelID = *attrs.WhiteLabelID
	}
	if attrs.PoolType != nil {
		u.PoolType = *attrs.PoolType
	}
	if attrs.WasReferralLinkGenerated != nil {
		u.WasReferralLinkGenerated = *attrs.WasReferralLinkGenerated
	}
	if attrs.IsAmbassador != nil {
		u.IsAmbassador = *attrs.IsAmbassador
	}
}

type Profile struct {
	User *User
}

func (p *User) SetRandomPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password bcrypt: %w", err)
	}

	p.Password = string(hash)
	return nil
}

type Referral struct {
	UserID    uuid.UUID
	Username  string
	Email     string
	CreatedAt time.Time
}

type UsernameAndHashrate struct {
	Username string
	Hashrate decimal.Decimal
}

type Address struct {
	Coin          string
	AccountTypeId userAccountModelEnum.AccountTypeIdWrapper
	MinPay        decimal.Decimal
	WalletAddress string
	MiningAddress string
}

type UserShortInfo struct {
	UserUUID  string
	UserName  string
	Email     string
	CreatedAt time.Time
}

type CoinAndAddress struct {
	Coin    string
	Address string
}

type UserPromoCode struct {
	ID                 int32
	CreatedAt          time.Time
	BonusFeeDaysAmount int32
	BonusFeeHasNoLimit bool
	Fee                decimal.Decimal
	IsSummable         bool
	CoinID             string
	ExpiresAt          time.Time
}

type ChangeWalletAddressConfirmResult struct {
	Address string
	UserID  int32
	CoinID  string
}

var IncreasedFeeRefMap = map[int32]float64{
	121356:          0.03,
	BmrGnsRefID:     BmrGnsFee,
	Next87RefID:     Next87Fee,
	BesitoRefID:     BesitoFee,
	Vladimir72RefID: Vladimir72Fee,
	ParrotRefID:     ParrotFee,
	Nata71RefID:     Nata71Fee,
	GarantexRefID:   GarantexFee,
}

const (
	Next87RefID = 132090
	Next87Fee   = 0.02

	BmrGnsFee   = 0.025
	BmrGnsRefID = 27524

	BesitoRefID = 132956
	BesitoFee   = 0.02

	Vladimir72RefID = 132507
	Vladimir72Fee   = 0.015

	ParrotRefID = 5017
	ParrotFee   = 0.02

	// Referee of Pavel071(id = 108257)
	Nata71RefID = 134222
	Nata71Fee   = 0.025

	GarantexRefID  = 111595
	GarantexFee    = 0.005
	GarantexUserID = 111595
)
