package middlewares

import (
	"net/http"
	"strings"

	"github.com/WillieBam/support_copilot/backend/app/config"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/labstack/echo/v5"
)

// FirebaseAuthMiddleware handles token inspection and hooks securely into Echo's context engine
func AuthMiddleware(authSvc interfaces.IAuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cfg := config.Get()

			if !cfg.Auth.Enabled {
				return next(c)
			}

			token := getBearerToken(c.Request().Header.Get("Authorization"))
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: Missing or invalid authorization",
				})
			}

			claims, err := authSvc.ParseAndValidateAuthToken(c.Request().Context(), token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized session context:" + err.Error(),
				})
			}

			c.Set("user_uid", claims.FirebaseUID)
			return next(c)

		}
	}
}

func getBearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		token := header[len("bearer "):]
		return strings.TrimSpace(token)
	}

	return ""
}
