package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError_Is(t *testing.T) {
	t.Parallel()

	errInnerError := errors.New("inner error")

	tests := []struct {
		name     string
		err      *Error
		target   error
		expected bool
	}{
		{
			name: "same error code and inner error",
			err: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errInnerError,
			},
			target: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errInnerError,
			},
			expected: true,
		},
		{
			name: "different error code",
			err: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errors.New("inner error"),
			},
			target: &Error{
				Code:    ErrorCodeInternal,
				Message: "internal error",
				Inner:   errors.New("inner error"),
			},
			expected: false,
		},
		{
			name: "same error code but different inner error",
			err: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errors.New("inner error 1"),
			},
			target: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errors.New("inner error 2"),
			},
			expected: false,
		},
		{
			name: "same error code and nil inner error",
			err: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
			},
			target: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
			},
			expected: true,
		},
		{
			name: "one error has Inner and other hasn't",
			err: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errors.New("inner error"),
			},
			target: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
			},
			expected: true, // inner errors are checked only when both have it
		},
		{
			name: "target is common error and is Inner",
			err: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errInnerError,
			},
			target:   errInnerError,
			expected: true,
		},
		{
			name: "target is common error and is not Inner",
			err: &Error{
				Code:    ErrorCodeNoSuchMerchant,
				Message: "no such merchant",
				Inner:   errInnerError,
			},
			target:   errors.New("other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.expected {
				require.ErrorIs(t, tt.err, tt.target)
			} else {
				require.NotErrorIs(t, tt.err, tt.target)
			}
		})
	}
}
