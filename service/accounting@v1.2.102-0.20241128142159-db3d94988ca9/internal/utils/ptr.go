package utils

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

func BoolToPtr(v bool) *bool {

	return &v
}

func IntToPtr(v int) *int {

	return &v
}

func Int32ToPtr(v int32) *int32 {

	return &v
}

func Int64ToPtr(v int64) *int64 {

	return &v
}

func Uint64ToPtr(v uint64) *uint64 {

	return &v
}

func Float64ToPtr(v float64) *float64 {

	return &v
}

func StringToPtr(v string) *string {

	return &v
}

func TimeToPtr(v time.Time) *time.Time {

	return &v
}

func DurationToPtr(v time.Duration) *time.Duration {

	return &v
}

func DecimalToPtr(v decimal.Decimal) *decimal.Decimal {

	return &v
}

func UuidToPtr(v uuid.UUID) *uuid.UUID {

	return &v
}
