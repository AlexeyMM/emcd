package bybit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"code.emcdtech.com/emcd/sdk/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestByBit_GetDepositRecords(t *testing.T) {
	bb := NewByBit(104020955, "https://api-testnet.bybit.com", "iwewYbo8QEHk8LD95p", "GarsMDtOxUhqwZ2rSbORtUt0RBLpXLgeAZ1c")

	depositRecords, err := bb.GetDepositRecords(context.Background(), "TON", time.Now().Add(-24*30*time.Hour), bb.masterApiKey, bb.masterApiSecret)
	require.NoError(t, err)

	for _, d := range depositRecords {
		log.Info(context.Background(), "%+v", d)
	}
}

func TestByBit_GetWithdrawRecords(t *testing.T) {
	bb := NewByBit(104020955, "https://api-testnet.bybit.com", "iwewYbo8QEHk8LD95p", "GarsMDtOxUhqwZ2rSbORtUt0RBLpXLgeAZ1c")
	status, err := bb.GetOrderStatus(context.Background(), uuid.MustParse("4133948d-d6a6-4199-8552-edac46499850"), bb.masterApiKey, bb.masterApiSecret)
	require.NoError(t, err)
	log.Info(context.Background(), "%d", status)
}

func TestDecimalWrapper_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected decimal.Decimal
		hasError bool
	}{
		{
			name:     "Valid decimal",
			input:    `"123.456"`,
			expected: decimal.NewFromFloat(123.456),
			hasError: false,
		},
		{
			name:     "Empty string",
			input:    `""`,
			expected: decimal.Zero,
			hasError: false,
		},
		{
			name:     "Invalid decimal",
			input:    `"abc"`,
			expected: decimal.Zero,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d decimalWrapper
			err := json.Unmarshal([]byte(tt.input), &d)
			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, decimal.Decimal(d).Equal(tt.expected))
			}
		})
	}
}

func TestIntWrapper_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int64
		hasError bool
	}{
		{
			name:     "Valid like a string",
			input:    `"123"`,
			expected: 123,
			hasError: false,
		},
		{
			name:     "Valid like an int",
			input:    `12`,
			expected: 12,
			hasError: false,
		},
		{
			name:     "Empty string",
			input:    `""`,
			expected: 0,
			hasError: false,
		},
		{
			name:     "Invalid",
			input:    `"abc"`,
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d intWrapper
			var data []byte
			var err error

			// Преобразуем input в JSON, если это не строка
			switch v := tt.input.(type) {
			case string:
				data = []byte(v) // передаем строку напрямую в []byte
			default:
				data, err = json.Marshal(v) // маршалим числовое значение
				require.NoError(t, err)
			}

			err = json.Unmarshal(data, &d)
			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, int64(d), tt.expected)
			}
		})
	}
}
