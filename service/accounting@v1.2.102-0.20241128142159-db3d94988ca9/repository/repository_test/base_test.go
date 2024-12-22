package repository_test

import (
	"code.emcdtech.com/emcd/service/accounting/repository"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestGetValFromNullString(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := "Hello"
		stringValue := wrapperspb.StringValue{Value: value}
		result := repository.GetValFromNullString(&stringValue)
		require.Equal(t, value, *result)
	})

	t.Run("null", func(t *testing.T) {
		var stringValue *wrapperspb.StringValue
		result := repository.GetValFromNullString(stringValue)
		require.Nil(t, result)
	})
}

func TestGetValFromNullInt64(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := int64(1)
		int64Value := wrapperspb.Int64Value{Value: value}
		result := repository.GetValFromNullInt64(&int64Value)
		require.Equal(t, int(value), *result)
	})

	t.Run("null", func(t *testing.T) {
		var int64Value *wrapperspb.Int64Value
		result := repository.GetValFromNullInt64(int64Value)
		require.Nil(t, result)
	})
}

func TestGetValFromNullFloat(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := "1.1"
		floatValue := wrapperspb.StringValue{Value: value}
		result, err := repository.GetValFromNullFloat(&floatValue)
		require.NoError(t, err)

		valueAsFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err)
		require.InEpsilon(t, valueAsFloat, *result, 1e-20)
	})

	t.Run("null", func(t *testing.T) {
		var value *wrapperspb.StringValue
		result, err := repository.GetValFromNullFloat(value)
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestGetValFromNullBool(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		value := false
		boolValue := wrapperspb.BoolValue{Value: value}
		result := repository.GetValFromNullBool(&boolValue)
		require.Equal(t, value, *result)
	})

	t.Run("null", func(t *testing.T) {
		var value *wrapperspb.BoolValue
		result := repository.GetValFromNullBool(value)
		require.Nil(t, result)
	})
}
