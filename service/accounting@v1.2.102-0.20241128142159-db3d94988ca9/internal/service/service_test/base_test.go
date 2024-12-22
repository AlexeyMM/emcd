package service_test

import (
	"code.emcdtech.com/emcd/service/accounting/internal/service"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"code.emcdtech.com/emcd/service/accounting/model"
)

func TestNullString(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := "Hello"
		result := service.NullString(&value)
		expected := &wrapperspb.StringValue{Value: value}
		require.Equal(t, expected, result)
	})

	t.Run("null", func(t *testing.T) {
		result := service.NullString(nil)
		require.Nil(t, result)
	})
}

func TestNullInt64(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := 1
		result := service.NullInt64(&value)
		expected := &wrapperspb.Int64Value{Value: int64(value)}
		require.Equal(t, expected, result)
	})

	t.Run("null", func(t *testing.T) {
		result := service.NullInt64(nil)
		require.Nil(t, result)
	})
}

func TestNullFloat(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := 1.1
		result := service.NullFloat(&value)
		expected := &wrapperspb.StringValue{Value: decimal.NewFromFloat(value).String()}
		require.Equal(t, expected, result)
	})

	t.Run("null", func(t *testing.T) {
		result := service.NullFloat(nil)
		require.Nil(t, result)
	})
}

func TestNullBool(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := false
		result := service.NullBool(&value)
		expected := &wrapperspb.BoolValue{Value: value}
		require.Equal(t, expected, result)
	})

	t.Run("null", func(t *testing.T) {
		result := service.NullBool(nil)
		require.Nil(t, result)
	})
}

func Test_getBlockTillByType(t *testing.T) {
	t.Run("Success: 34 type", func(t *testing.T) {
		result := service.GetBlockTillByType(model.CnhldDiffBalanceTrTypeID)
		require.LessOrEqual(t, result, time.Now().Add(service.Day+1*time.Millisecond))
	})

	t.Run("Success: 36 type", func(t *testing.T) {
		result := service.GetBlockTillByType(model.CnhldEarlyCloseTrTypeID)
		require.LessOrEqual(t, result, time.Now().Add(service.Day+1*time.Millisecond))
	})

	t.Run("Success: other type", func(t *testing.T) {
		result := service.GetBlockTillByType(model.ReturnInterestsTrTypeID)
		require.LessOrEqual(t, time.Now().Add(service.Year-1*time.Millisecond), result)
	})
}
