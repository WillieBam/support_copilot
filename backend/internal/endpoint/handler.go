package endpoint

import (
	"log"
	"log/slog"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/WillieBam/support_copilot/backend/app"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	apps *app.AppService
}

func NewHandler(a *app.AppService) *Handler {
	return &Handler{apps: a}
}

type queryRequest struct {
	Input string `json:"input"`
}

type queryResponse struct {
	Output string `json:"output"`
}

// Define the payload structures
type TokenExchangeRequest struct {
	FirebaseToken string `json:"firebase_token"`
}

type TokenExchangeResponse struct {
	Token string `json:"token"`
}

// TokenExchangeHandler converts a validated Firebase token into a JWT session token
func (h *Handler) TokenExchangeHandler(c *echo.Context) error {
	var req TokenExchangeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing request payload"})
	}

	if req.FirebaseToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing firebase request"})
	}

	verified, err := h.apps.AuthService.ExchangeToken(c.Request().Context(), req.FirebaseToken)
	if err != nil {
		if err.Error() == "mfa_required" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error":   "mfa_required",
				"message": "TOTP verification required",
			})
		}
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, TokenExchangeResponse{Token: verified})
}

// Query handles POST /query/chat
func (h *Handler) Query(c *echo.Context) error {
	uidVal := c.Get("user_uid")
	appUID, ok := uidVal.(string)
	if !ok || appUID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Authentication required: missing support copilot token context",
		})
	}

	log.Printf("[DEBUG] Successfully authenticated user UID: %s processing query stream.", appUID)
	var req queryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Input == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "input is required"})
	}

	// route query into service
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
