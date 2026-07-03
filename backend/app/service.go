package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/WillieBam/support_copilot/backend/app/config"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
)

type AppService struct {
	client      *appClient
	repo        *appRepository
	AuthService interfaces.IAuthService
}
type StreamEvent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func newAppService(client *appClient, repo *appRepository, authService interfaces.IAuthService) *AppService {
	return &AppService{
		client:      client,
		repo:        repo,
		AuthService: authService,
	}
}

func (s *AppService) Query(ctx context.Context, prompt string) (string, error) {
	return s.client.Query(ctx, prompt)
}

// QueryStream connects to the Ollama API to run a query with streaming enabled.
// It sends progress events and token streams back to the client via streamChan.
// It returns an error if marshalling the payload, creating the request, or the connection/decoding fails.
func (s *AppService) QueryStream(ctx context.Context, prompt string, streamChan chan<- StreamEvent) error {

	streamChan <- StreamEvent{
		Type:    "reasoning",
		Content: "Analyzing user prompt...\n ",
	}

	ollamaURL := config.Get().Ollama.BaseURL + "/api/chat"

	payload := map[string]interface{}{
		"model":  config.Get().Ollama.Model,
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
		return fmt.Errorf("[ERROR]: failed to marshal Ollama payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaURL, bytes.NewReader(payLoadBytes))
	if err != nil {
		return fmt.Errorf("[ERROR]: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	streamChan <- StreamEvent{
		Type:    "reasoning",
		Content: "Connecting to Llama 3.2...\n",
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		// check for stream abortioni
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

	// read and process the streaming JSON response chunks from Ollama
	for {
		// structure of a single response chunk from Ollama
		var chunk struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Done  bool   `json:"done"`
			Error string `json:"error"`
		}

		// decode the next JSON object in the stream
		if err := decoder.Decode(&chunk); err == io.EOF {
			break // streaming ended normally
		} else if err != nil {
			return fmt.Errorf("error decoding ollama chunk: %w", err)
		}

		// checl if Ollama API returned an error message in the chunk
		if chunk.Error != "" {
			return fmt.Errorf("[ERROR]: ollama API error: %s", chunk.Error)
		}

		// if new text fragment available, forward it to the stream channel
		if chunk.Message.Content != "" {
			streamChan <- StreamEvent{
				Type:    "text",
				Content: chunk.Message.Content,
			}
		}

		// stop looping when Ollama signals that the response generation is complete.
		if chunk.Done {
			slog.Info("[STREAM]: Ollama Done response")
			break
		}
	}
	return nil
}
