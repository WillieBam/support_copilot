package endpoint

import (
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service interfaces.ISupportCopilotService
}

func NewHandler(service interfaces.ISupportCopilotService) *Handler {
	return &Handler{service: service}
}

type queryRequest struct {
	Input string `json:"input"`
}

type queryResponse struct {
	Output string `json:"output"`
}

// Query handles POST /query/sc
func (h *Handler) Query(c *echo.Context) error {
	var req queryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Input == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "input is required"})
	}

	result, err := h.service.Query(c.Request().Context(), req.Input)
	if err != nil {
		errMsg := err.Error()
		log.Printf("query failed: %v", err)

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
