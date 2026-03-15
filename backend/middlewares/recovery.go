package middlewares

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// RecoveryMiddleware returns an Echo middleware that recovers from panics.
func RecoveryMiddleware() echo.MiddlewareFunc {
	return echoMiddleware.Recover()
}
