package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	apps        interfaces.IAppService
	authService interfaces.IAuthService
}

func NewHandler(a interfaces.IAppService, authService interfaces.IAuthService) *Handler {
	return &Handler{
		apps:        a,
		authService: authService,
	}
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

	verified, claims, err := h.authService.ExchangeToken(c.Request().Context(), req.FirebaseToken)
	if err != nil {
		if err.Error() == "mfa_required" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error":   "mfa_required",
				"message": "TOTP verification required",
			})
		}
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	var expires time.Time
	if claims != nil && claims.ExpiresAt != nil {
		expires = claims.ExpiresAt.Time
	} else {
		expires = time.Now().Add(1 * time.Hour)
	}

	cookie := &http.Cookie{
		Name:     "support_copilot_session",
		Value:    verified,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	slog.Info("Successfully created and attached HttpOnly session cookie",
		"user_uid", claims.FirebaseUID,
		"expires_at", expires.Format(time.RFC3339),
	)
	return c.JSON(http.StatusOK, map[string]string{"status": "authenticated"})
}

func (h *Handler) Me(c *echo.Context) error {
	uidVal := c.Get("user_uid")
	appUID, ok := uidVal.(string)
	if !ok || appUID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized session"})
	}

	emailVal := c.Get("user_email")
	email, _ := emailVal.(string)

	// return info about who is authenticated to revive React UI client state
	return c.JSON(http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user_uid":      appUID,
		"user_email":    email,
	})
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

	log.Printf("[LOG] Successfully authenticated user UID: %s processing query stream.", appUID)
	var req queryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Input == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "input is required"})
	}

	resp := c.Response()
	resp.Header().Set("Content-Type", "text/event-stream")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "keep-alive")
	resp.WriteHeader(http.StatusOK)

	flusher, ok := resp.(http.Flusher)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Streaming unsupported"})
	}
	flusher.Flush()

	streamChan := make(chan types.StreamEvent)
	errorChan := make(chan error, 1)

	go func() {
		// Pass the channel into the service so it can push events!
		err := h.apps.QueryStream(c.Request().Context(), req.Input, streamChan)
		if err != nil {
			errorChan <- err
		}
		// Always close the channel when the service is done
		close(streamChan)
	}()

	for {
		select {
		case event, ok := <-streamChan:
			if !ok {
				return nil
			}
			eventJSON, _ := json.Marshal(event)
			fmt.Fprintf(resp, "data: %s\n\n", eventJSON)
			flusher.Flush()

		case err := <-errorChan:
			slog.Error("[STREAM ERROR]: query stream failed", "err", err)
			errEvent := types.StreamEvent{
				Type: "text",
				// always use fmt.Sprintf to build json
				Content: fmt.Sprintf("\n\n**Error** %s", err.Error()),
			}
			// always marshal with json.Marshal
			eventJSON, _ := json.Marshal(errEvent)

			fmt.Fprintf(resp, "data: %s\n\n", eventJSON)
			flusher.Flush()
			return nil

		case <-c.Request().Context().Done():
			log.Println("[STREAM]: Client Disconnected (prompt edited or stopped). Aborting stream gracefully.")
			return nil

		}
	}

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
