package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateMode string

const (
	UpdateModeDefault  UpdateMode = ""
	UpdateModeForceAll UpdateMode = "force_all"
	UpdateModeAll      UpdateMode = "all"
)

func NewUpdateMode(s string) (UpdateMode, error) {
	switch s {
	case string(UpdateModeDefault):
		return UpdateModeDefault, nil
	case string(UpdateModeForceAll):
		return UpdateModeForceAll, nil
	case string(UpdateModeAll):
		return UpdateModeAll, nil
	}
	return "", errors.New("invalid update mode")
}

func (u UpdateMode) String() string {
	return string(u)
}

type ReferralSettings struct {
	ReferralUUID uuid.UUID // Это UUID пользователя, из таблицы profile.users
	Preferences  []ReferralPreference
}

type ReferralPreference struct {
	Product     string
	Coin        string
	Fee         float64
	ReferralFee float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SettingForCoinAndProduct struct {
	Product string
	Coin    string
	Fee     decimal.Decimal
}
