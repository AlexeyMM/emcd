package log

import (
	"context"
	"github.com/labstack/echo/v4"
)

// TracerResponseFunc represents any func for processing error returned by handler.
type TracerResponseFunc func(c echo.Context, err error) error

// TracerMiddleware middleware runs preprocessing of error returned from handler.
func TracerMiddleware(errFn TracerResponseFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

// ValueRepresentationUserIDMiddleware middleware for add value representation user_id how user.id
func ValueRepresentationUserIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid := c.Get(oldUserId)
			if uid != nil {
				c.Set(userID, uid)
			}
			return next(c)
		}
	}
}

// ContextNamedValueMiddleware add context value to request for service name logging
func ContextNamedValueMiddleware(serviceName string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(ctx echo.Context) error {
			ctxEndpoint := context.WithValue(ctx.Request().Context(), serviceNameStruct{}, serviceName)
			ctx.SetRequest(ctx.Request().WithContext(ctxEndpoint))

			return next(ctx)
		}
	}
}
