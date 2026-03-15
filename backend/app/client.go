package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/WillieBam/support_copilot/backend/app/config"
)

type appClient struct {
	httpClient  *http.Client
	geminiAPI   string
	geminiBase  string
	geminiModel string
}

func newAppClient() *appClient {
	cfg := config.Get()
	baseURL := strings.TrimRight(cfg.Gemini.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com"
	}

	return &appClient{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		geminiAPI:   cfg.Gemini.APIKey,
		geminiBase:  baseURL,
		geminiModel: cfg.Gemini.Model,
	}
}

type geminiGenerateContentRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiGenerateContentResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (c *appClient) QueryGemini(ctx context.Context, prompt string) (string, error) {
	if c.geminiAPI == "" {
		return "", errors.New("missing GEMINI_API_KEY")
	}

	if c.geminiBase == "" {
		return "", errors.New("missing GEMINI_BASE_URL")
	}

	if !strings.HasPrefix(c.geminiBase, "http://") && !strings.HasPrefix(c.geminiBase, "https://") {
		return "", fmt.Errorf("invalid GEMINI_BASE_URL: %s", c.geminiBase)
	}

	if c.geminiModel == "" {
		c.geminiModel = "gemini-2.0-flash"
	}

	var body geminiGenerateContentRequest
	body.Contents = []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{
		{
			Parts: []struct {
				Text string `json:"text"`
			}{{Text: prompt}},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", c.geminiBase, c.geminiModel, url.QueryEscape(c.geminiAPI))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("gemini API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var geminiResp geminiGenerateContentResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("gemini returned no content")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
