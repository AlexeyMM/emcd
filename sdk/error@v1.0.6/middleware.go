package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorResponseFunc represents any func for processing error returned by handler.
type ErrorResponseFunc func(c echo.Context, err error) error

// ErrorMiddleware middleware runs preprocessing of error returned from handler.
func ErrorMiddleware(errFn ErrorResponseFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				return errFn(c, err)
			}
			return nil
		}
	}
}

// ErrorToJSON respond business error in JSON format. If err is not business error type - return 500 status code.
func ErrorToJSON(c echo.Context, err error) error {
	var e *Error
	if errors.As(err, &e) {
		return c.JSON(http.StatusBadRequest, e)
	}

	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		return err
	}

	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}
