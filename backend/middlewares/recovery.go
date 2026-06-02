package middlewares

import (
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// RecoveryMiddleware returns an Echo middleware that recovers from panics.
func RecoveryMiddleware() echo.MiddlewareFunc {
	return echoMiddleware.Recover()
}
