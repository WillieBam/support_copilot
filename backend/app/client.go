package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/WillieBam/support_copilot/backend/app/config"
)

type appClient struct {
	httpClient  *http.Client
	ollamaBase  string
	ollamaModel string
}

func newAppClient() *appClient {
	cfg := config.Get()

	ollamaBaseURL := strings.TrimRight(cfg.Ollama.BaseURL, "/")
	if ollamaBaseURL == "" {
		ollamaBaseURL = "http://localhost:11434"
	}

	ollamaModel := strings.TrimSpace(cfg.Ollama.Model)
	if ollamaModel == "" {
		ollamaModel = "llama3.2"
	}

	return &appClient{
		httpClient:  &http.Client{Timeout: 3 * time.Minute},
		ollamaBase:  ollamaBaseURL,
		ollamaModel: ollamaModel,
	}
}

type ollamaChatRequest struct {
	Model    string `json:"model"`
	Stream   bool   `json:"stream"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ollamaChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Response string `json:"response"`
}

func (c *appClient) Query(ctx context.Context, prompt string) (string, error) {
	return c.queryOllama(ctx, prompt)
}

func (c *appClient) queryOllama(ctx context.Context, prompt string) (string, error) {
	if c.ollamaBase == "" {
		return "", errors.New("missing OLLAMA_BASE_URL")
	}

	if !strings.HasPrefix(c.ollamaBase, "http://") && !strings.HasPrefix(c.ollamaBase, "https://") {
		return "", fmt.Errorf("invalid OLLAMA_BASE_URL: %s", c.ollamaBase)
	}

	if c.ollamaModel == "" {
		return "", errors.New("missing OLLAMA_MODEL")
	}

	var body ollamaChatRequest
	body.Model = c.ollamaModel
	body.Stream = false
	body.Messages = []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		{Role: "user", Content: prompt},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("%s/api/chat", c.ollamaBase)
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
		return "", fmt.Errorf("ollama API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp ollamaChatResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return "", err
	}

	if text := strings.TrimSpace(ollamaResp.Message.Content); text != "" {
		return text, nil
	}

	if text := strings.TrimSpace(ollamaResp.Response); text != "" {
		return text, nil
	}

	return "", errors.New("ollama returned no content")
}
