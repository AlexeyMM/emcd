// Package model provides data models for the application.
package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"code.emcdtech.com/emcd/service/referral/protocol/default_settings"
	"code.emcdtech.com/emcd/service/referral/protocol/default_whitelabel_settings"
	"code.emcdtech.com/emcd/service/referral/protocol/referral"
)

// DefaultSettings represents the default settings for a product in the system.
type DefaultSettings struct {
	Product     string
	Coin        string
	Fee         decimal.Decimal
	ReferralFee decimal.Decimal
	CreatedAt   time.Time
}

// FromProto updates the fields of DefaultSettings based on the values of a given default_settings.Settings object.
func (ds *DefaultSettings) FromProto(p *default_settings.Settings) (err error) {
	ds.Product = p.Product
	ds.Coin = p.Coin

	ds.Fee, err = decimal.NewFromString(p.Fee)
	if err != nil {
		return fmt.Errorf("parseFee: %w", err)
	}

	ds.ReferralFee, err = decimal.NewFromString(p.ReferralFee)
	if err != nil {
		return fmt.Errorf("parseReferralFee: %w", err)
	}

	return
}

// DefaultWhitelabelSettings represents the default settings for a whitelabel product in the system.
type DefaultWhitelabelSettings struct {
	WhitelabelID uuid.UUID
	Product      string
	Coin         string
	Fee          decimal.Decimal
	ReferralFee  decimal.Decimal
	CreatedAt    time.Time
}

// FromProto updates the fields of DefaultWhitelabelSettings based on the values of a given default_whitelabel_settings.Settings object.
func (ds *DefaultWhitelabelSettings) FromProto(p *default_whitelabel_settings.Settings) (err error) {
	ds.Product = p.Product
	ds.Coin = p.Coin

	ds.WhitelabelID, err = uuid.Parse(p.WhitelabelId)
	if err != nil {
		return fmt.Errorf("parseWhitelabelID: %w", err)
	}

	ds.Fee, err = decimal.NewFromString(p.Fee)
	if err != nil {
		return fmt.Errorf("parseFee: %w", err)
	}

	ds.ReferralFee, err = decimal.NewFromString(p.ReferralFee)
	if err != nil {
		return fmt.Errorf("parseReferralFee: %w", err)
	}

	return
}

// Referral represents a referral entity in the system.
type Referral struct {
	UserID        uuid.UUID
	Product       string
	Coin          string
	WhitelabelID  uuid.UUID
	Fee           decimal.Decimal
	WhitelabelFee decimal.Decimal
	ReferralFee   decimal.Decimal
	ReferralID    uuid.UUID
	CreatedAt     time.Time
}

type UserReferral struct {
	UserID uuid.UUID
}

// FromProto updates the fields of Referral based on the values of a given referral.Settings object.
func (r *Referral) FromProto(p *referral.Referral) error {
	var err error
	r.Product = p.Product
	r.Coin = p.Coin

	r.UserID, err = uuid.Parse(p.UserId)
	if err != nil {
		return fmt.Errorf("parseUserID: %w", err)
	}

	r.WhitelabelID, err = uuid.Parse(p.WhitelabelId)
	if err != nil {
		return fmt.Errorf("parseWhitelabelID: %w", err)
	}

	r.ReferralID, err = uuid.Parse(p.ReferralId)
	if err != nil {
		return fmt.Errorf("parseReferralId: %w", err)
	}

	r.Fee, err = decimal.NewFromString(p.Fee)
	if err != nil {
		return fmt.Errorf("parseFee: %w", err)
	}

	r.WhitelabelFee, err = decimal.NewFromString(p.WhitelabelFee)
	if err != nil {
		return fmt.Errorf("parseWhitelabelFee: %w", err)
	}

	r.ReferralFee, err = decimal.NewFromString(p.ReferralFee)
	if err != nil {
		return fmt.Errorf("parseReferralFee: %w", err)
	}

	return nil
}

// Transaction represents a financial transaction made by a user.
type Transaction struct {
	UserID uuid.UUID
	Type   string
	Amount decimal.Decimal
}

type DefaultWhitelabelSettingsV2 struct {
	WhitelabelID  uuid.UUID
	Product       string
	Coin          string
	Fee           decimal.Decimal
	ReferralFee   decimal.Decimal
	CreatedAt     time.Time
	WhiteLabelFee decimal.Decimal
}

// FromProto updates the fields of DefaultWhitelabelSettings based on the values of a given default_whitelabel_settings.Settings object.
func (ds *DefaultWhitelabelSettingsV2) FromProto(p *default_whitelabel_settings.SettingsV2) (err error) {
	ds.Product = p.Product
	ds.Coin = p.Coin

	ds.WhitelabelID, err = uuid.Parse(p.WhitelabelId)
	if err != nil {
		return fmt.Errorf("parseWhitelabelID: %w", err)
	}

	ds.Fee, err = decimal.NewFromString(p.Fee)
	if err != nil {
		return fmt.Errorf("parseFee: %w", err)
	}

	ds.ReferralFee, err = decimal.NewFromString(p.ReferralFee)
	if err != nil {
		return fmt.Errorf("parseReferralFee: %w", err)
	}

	ds.WhiteLabelFee, err = decimal.NewFromString(p.WhitelabelFee)
	if err != nil {
		return fmt.Errorf("parseWhiteLabelFee: %w", err)
	}
	return
}

type CoinMultiplier struct {
	UserID           uuid.UUID
	Coin             string
	FeeMultiplier    decimal.Decimal
	RefFeeMultiplier decimal.Decimal
	CreatedAt        time.Time
	Product          string
	ActionID         uuid.UUID
}

type CoinsMultipliers struct {
	UserID            uuid.UUID
	CreatedAt         time.Time
	Product           string
	ActionID          uuid.UUID
	FeeMultipliers    map[string]decimal.Decimal
	RefFeeMultipliers map[string]decimal.Decimal
}
