package errors_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	businessErr "code.emcdtech.com/emcd/sdk/error"
)

func TestErrorMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	require.NoError(t, err)

	e := echo.New()
	c := e.NewContext(req, rec)

	testTable := []struct {
		name        string
		err         error
		expectedErr error
	}{
		{
			name:        "No error",
			err:         nil,
			expectedErr: nil,
		},
		{
			name:        "Not business error",
			err:         errors.New("not business error"),
			expectedErr: echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"),
		},
		{
			name:        "Specified by user error",
			err:         echo.NewHTTPError(http.StatusUnauthorized, "not business error"),
			expectedErr: echo.NewHTTPError(http.StatusUnauthorized, "not business error"),
		},
		{
			name:        "Business error",
			err:         businessErr.NewError("101", "user does not have enough money"),
			expectedErr: nil,
		},
		{
			name: "deeply wrapped business error",
			err: fmt.Errorf(
				"wrap level 3: %w",
				fmt.Errorf(
					"wrap level 2: %w",
					fmt.Errorf(
						"wrap level 1: %w",
						businessErr.NewError("101", "user does not have enough money"),
					),
				),
			),
			expectedErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(c echo.Context) error {
				return tc.err
			}

			err = businessErr.ErrorMiddleware(businessErr.ErrorToJSON)(handler)(c)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
