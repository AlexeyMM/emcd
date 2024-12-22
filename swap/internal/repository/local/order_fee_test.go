package local

import (
	"context"
	"testing"

	businessError "code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestFee_GetFee_Success(t *testing.T) {
	ctx := context.Background()

	feeData := map[string]*model.Fee{
		"BTCUSD": {MakerFee: decimal.NewFromFloat(0.001), TakerFee: decimal.NewFromFloat(0.002)},
	}

	fee := NewFee()
	err := fee.UpdateAll(ctx, feeData)
	require.NoError(t, err)

	result, err := fee.GetFee(ctx, "BTCUSD")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, decimal.NewFromFloat(0.001), result.MakerFee)
	require.Equal(t, decimal.NewFromFloat(0.002), result.TakerFee)
}

func TestFee_GetFee_NotFound(t *testing.T) {
	// Создаем контекст
	ctx := context.Background()

	// Создаем объект Fee с пустыми данными
	fee := NewFee()

	// Тестируем получение комиссии для несуществующего символа
	result, err := fee.GetFee(ctx, "ETHUSD")
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, businessError.FeeNotFoundErr, err)
}

func TestFee_UpdateAll(t *testing.T) {
	ctx := context.Background()

	initialFee := map[string]*model.Fee{
		"BTCUSD": {MakerFee: decimal.NewFromFloat(0.001), TakerFee: decimal.NewFromFloat(0.002)},
	}

	fee := NewFee()
	err := fee.UpdateAll(ctx, initialFee)
	require.NoError(t, err)

	newFee := map[string]*model.Fee{
		"ETHUSD": {MakerFee: decimal.NewFromFloat(0.003), TakerFee: decimal.NewFromFloat(0.004)},
	}

	err = fee.UpdateAll(ctx, newFee)
	require.NoError(t, err)

	result, err := fee.GetFee(ctx, "ETHUSD")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, decimal.NewFromFloat(0.003), result.MakerFee)
	require.Equal(t, decimal.NewFromFloat(0.004), result.TakerFee)

	// Проверяем, что старых данных больше нет
	result, err = fee.GetFee(ctx, "BTCUSD")
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, businessError.FeeNotFoundErr, err)
}
