package endpoint

import (
	"log/slog"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/WillieBam/support_copilot/backend/app"
	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	apps *app.AppService
}

func NewHandler(a *app.AppService) *Handler {
	return &Handler{apps: a}
}

func (h *Handler) FirebaseLogin(c *echo.Context) error {
	var req types.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.IDToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "firebase id_token is required"})
	}

	// Verify credentials and register/login user using the AuthService layer
	user, err := h.apps.AuthService.LoginOrRegister(c.Request().Context(), req.IDToken)
	if err != nil {
		// Handle verification failure (expired token, revoked token, invalid credentials)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "authentication failed: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, types.LoginResponse{
		Message: "Authentication successful",
		UID:     user.FirebaseUID,
		Email:   user.Email,
		Name:    user.DisplayName,
	})
}

type queryRequest struct {
	Input string `json:"input"`
}

type queryResponse struct {
	Output string `json:"output"`
}

// Query handles POST /query/chat
func (h *Handler) Query(c *echo.Context) error {
	var req queryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Input == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "input is required"})
	}

	result, err := h.apps.Query(c.Request().Context(), req.Input)
	if err != nil {
		errMsg := err.Error()
		slog.Error("query failed", "err", err)

		if strings.Contains(errMsg, "API error (429)") || strings.Contains(errMsg, "RESOURCE_EXHAUSTED") {
			if retryAfter := extractRetryAfterSeconds(errMsg); retryAfter > 0 {
				c.Response().Header().Set(echo.HeaderRetryAfter, strconv.Itoa(retryAfter))
			}

			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "The model provider is rate limited. Please retry shortly.",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to process query"})
	}

	return c.JSON(http.StatusOK, queryResponse{Output: result})
}

func extractRetryAfterSeconds(msg string) int {
	retryDelayJSON := regexp.MustCompile(`retryDelay"\s*:\s*"(\d+)s"`)
	if m := retryDelayJSON.FindStringSubmatch(msg); len(m) == 2 {
		v, err := strconv.Atoi(m[1])
		if err == nil && v > 0 {
			return v
		}
	}

	retryDelayText := regexp.MustCompile(`Please retry in ([0-9]+(?:\.[0-9]+)?)s`)
	if m := retryDelayText.FindStringSubmatch(msg); len(m) == 2 {
		v, err := strconv.ParseFloat(m[1], 64)
		if err == nil && v > 0 {
			return int(math.Ceil(v))
		}
	}

	return 0
}
