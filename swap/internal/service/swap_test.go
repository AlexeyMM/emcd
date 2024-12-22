package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	businessError "code.emcdtech.com/b2b/swap/internal/business_error"
	"code.emcdtech.com/b2b/swap/mocks/internal_/repository"
	"code.emcdtech.com/b2b/swap/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCalculateAmountToDirectionBuy(t *testing.T) {
	mockOrderBookRep := new(repository.MockOrderBook)
	swapService := &swap{orderBookRep: mockOrderBookRep}

	levels := [][2]float64{
		{5.29, 29.69}, // price = 5.29, size = 29.69
		{5.30, 20.00}, // price = 5.30, size = 20.00
	}

	mockOrderBookRep.On("GetAskTopLevels", "TONUSDT", 1).Return(levels, nil)

	amountFrom := decimal.NewFromFloat(200.00) // amountFrom (USDT)

	totalAmountTo, err := swapService.calculateAmountToDirectionBuy("TONUSDT", amountFrom, 1)

	require.NoError(t, err)

	// Ожидаемое значение amountTo (TON) для amountFrom = 100 USDT
	// Для первого уровня: 29.69 * 5.29 = 157.12 USDT -> не выкупается на первом уровне
	// Второй уровень: 200 - 157.12 = 42.88 / 5.3 = 8.09 + 29.69 = 37.78
	expectedAmountTo := decimal.NewFromFloat(37.78)

	// Определяем допустимую погрешность
	inaccuracy := decimal.NewFromFloat(0.02)

	// Проверяем, что значение находится в пределах допустимой погрешности
	diff := totalAmountTo.Sub(expectedAmountTo).Abs()
	require.LessOrEqual(t, diff.Cmp(inaccuracy), 0)
}

func TestCalculateAmountToDirectionSell(t *testing.T) {
	mockOrderBookRep := new(repository.MockOrderBook)
	swapService := &swap{orderBookRep: mockOrderBookRep}

	levels := [][2]float64{
		{5.29, 29.00}, // price = 5.29, size = 29.69
		{5.30, 20.00}, // price = 5.30, size = 20.00
	}

	mockOrderBookRep.On("GetBidTopLevels", "TONUSDT", 1).Return(levels, nil)

	amountFrom := decimal.NewFromFloat(30)

	totalAmountTo, err := swapService.calculateAmountToDirectionSell("TONUSDT", amountFrom, 1)
	require.NoError(t, err)

	// Ожидаемое значение amountTo (USDT) для amountFrom = 20 TON:
	// 1-й уровень: 29 TON * 5.29 = 153.41 USDT
	// 2-й уровень: 1 * 5.30 = 5.30 + 153.41 = 158.71
	expectedAmountTo := decimal.NewFromFloat(158.71)

	// Определяем допустимую погрешность
	inaccuracy := decimal.NewFromFloat(0.02)

	// Проверяем, что значение находится в пределах допустимой погрешности
	diff := totalAmountTo.Sub(expectedAmountTo).Abs()
	require.LessOrEqual(t, diff.Cmp(inaccuracy), 0)
}

func TestGetSwapOptions(t *testing.T) {
	mockOrderBookRep := new(repository.MockOrderBook)
	swapService := &swap{orderBookRep: mockOrderBookRep}
	ctx := context.Background()

	// Тест 1: Direct swap from USDT -> Symbol
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "BTC", model.CoinUSDT)).Return(true)
	step, indirect, err := swapService.getSwapOptions(ctx, model.CoinUSDT, "BTC")
	require.NoError(t, err)
	require.NotNil(t, step)
	require.Nil(t, indirect)
	require.Equal(t, "BTCUSDT", step.Symbol)
	require.Equal(t, model.Buy, int(step.Direction))

	// Тест 2: Direct swap Symbol -> USDT
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "ETH", model.CoinUSDT)).Return(true)
	step, indirect, err = swapService.getSwapOptions(ctx, "ETH", model.CoinUSDT)
	require.NoError(t, err)
	require.NotNil(t, step)
	require.Nil(t, indirect)
	require.Equal(t, "ETHUSDT", step.Symbol)
	require.Equal(t, model.Sell, int(step.Direction))

	// Тест 3: No direct swap, but indirect swap available
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "TON", "INJ")).Return(false)
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "INJ", "TON")).Return(false)
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "TON", model.CoinUSDT)).Return(true)
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "INJ", model.CoinUSDT)).Return(true)
	step, indirect, err = swapService.getSwapOptions(ctx, "TON", "INJ")
	require.NoError(t, err)
	require.Nil(t, step)
	require.Len(t, indirect, 2)
	require.Equal(t, "TONUSDT", indirect[0].Symbol)
	require.Equal(t, model.Sell, int(indirect[0].Direction))
	require.Equal(t, "INJUSDT", indirect[1].Symbol)
	require.Equal(t, model.Buy, int(indirect[1].Direction))

	// Тест 4: Error - no swap available
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "ADA", "INC")).Return(false)
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "INC", "ADA")).Return(false)
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "ADA", model.CoinUSDT)).Return(false)
	mockOrderBookRep.On("IsExist", fmt.Sprintf("%s%s", "INC", model.CoinUSDT)).Return(false)
	step, indirect, err = swapService.getSwapOptions(ctx, "ADA", "INC")
	require.Error(t, err)
	require.Equal(t, businessError.NoPathToSwapErr, errors.Unwrap(err))
	require.Nil(t, step)
	require.Nil(t, indirect)
}
