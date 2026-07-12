package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WillieBam/support_copilot/backend/app/config"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types"
)

type ollamaClient struct {
	httpClient    *http.Client
	ollamaBaseUrl string
	ollamaModel   string
}

func NewOllamaClient(cfg *config.Config) interfaces.IOllamaClient {
	baseUrl := strings.TrimRight(cfg.Ollama.BaseURL, "/")
	if baseUrl == "" {
		baseUrl = "http://localhost:11434"
	}

	model := strings.TrimSpace(cfg.Ollama.Model)
	if model == "" {
		model = "llama3.2"
	}

	return &ollamaClient{
		httpClient:    &http.Client{Timeout: 3 * time.Minute},
		ollamaBaseUrl: baseUrl,
		ollamaModel:   model,
	}

}

func (c *ollamaClient) QueryStream(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error {
	streamChan <- types.StreamEvent{
		Type:    "reasoning",
		Content: "Analyzing user prompt...\n ",
	}

	ollamaURL := c.ollamaBaseUrl + "/api/chat"

	payload := map[string]interface{}{
		"model":  c.ollamaModel,
		"stream": true,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	payLoadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Ollama payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaURL, bytes.NewReader(payLoadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	streamChan <- types.StreamEvent{
		Type:    "reasoning",
		Content: "Connecting to Llama 3.2...\n",
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// check for stream abortion
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("[STREAM ERROR]: client aborted stream")
		}
		// others stream issue
		return fmt.Errorf("[STREAM ERROR]: Failed to connect Ollama: %w", err)
	}

	// close resp body after ensure no error
	defer resp.Body.Close()

	// check for ollama status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("[STREAM ERROR]: Ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)

	// Read and process the streaming JSON response chunks from Ollama.
	for {
		// Define the structure of a single response chunk from Ollama.
		var chunk struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Done  bool   `json:"done"`
			Error string `json:"error"`
		}

		// Decode the next JSON object in the stream.
		if err := decoder.Decode(&chunk); err == io.EOF {
			break // streaming ended normally
		} else if err != nil {
			return fmt.Errorf("error decoding ollama chunk: %w", err)
		}

		// Check if Ollama API returned an error message in the chunk.
		if chunk.Error != "" {
			return fmt.Errorf("ollama API error: %s", chunk.Error)
		}

		// If a new text fragment is available, forward it to the stream channel.
		if chunk.Message.Content != "" {
			streamChan <- types.StreamEvent{
				Type:    "text",
				Content: chunk.Message.Content,
			}
		}

		// Stop looping when Ollama signals that the response generation is complete.
		if chunk.Done {
			slog.Debug("Ollama Done response")
			break
		}
	}
	return nil
}
