package repository

import (
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

var errNoError = errors.New("no error")

func GetValFromNullString(v *wrapperspb.StringValue) *string {
	if v != nil {
		return &v.Value
	}
	return nil
}

func GetValFromNullInt64(v *wrapperspb.Int64Value) *int {
	if v != nil {
		result := int(v.Value)
		return &result
	}
	return nil
}

func GetValFromNullFloat(v *wrapperspb.StringValue) (*float64, error) {
	if v == nil {

		return nil, errNoError
	} else {
		if result, err := strconv.ParseFloat(v.Value, 64); err != nil {

			return nil, fmt.Errorf("ParseFloat: %w", err)
		} else {

			return &result, nil
		}
	}
}

func GetValFromNullBool(v *wrapperspb.BoolValue) *bool {
	if v != nil {
		return &v.Value
	}
	return nil
}
