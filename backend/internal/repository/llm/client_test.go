package ollama_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/WillieBam/support_copilot/backend/app/config"
	ollama "github.com/WillieBam/support_copilot/backend/internal/repository/llm"
	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/WillieBam/support_copilot/backend/types/requests"
)

var _ = Describe("OllamaClient", func() {
	Context("NewOllamaClient", func() {
		It("should initialize with defaults when config fields are empty", func() {
			cfg := &config.Config{}
			client := ollama.NewOllamaClient(cfg)
			Expect(client).NotTo(BeNil())
		})

		It("should initialize with provided config values", func() {
			cfg := &config.Config{}
			cfg.Ollama.BaseURL = "http://localhost:11434/"
			cfg.Ollama.Model = "llama3.2:latest"

			client := ollama.NewOllamaClient(cfg)
			Expect(client).NotTo(BeNil())
		})
	})

	Context("QueryStreamWithTools", func() {
		It("should stream responses successfully from Ollama mock server", func() {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.URL.Path).To(Equal("/api/chat"))
				Expect(r.Method).To(Equal(http.MethodPost))
				w.Header().Set("Content-Type", "application/x-ndjson")
				w.WriteHeader(http.StatusOK)

				flusher, ok := w.(http.Flusher)
				Expect(ok).To(BeTrue())

				fmt.Fprintln(w, `{"message":{"content":"Hello "},"done":false}`)
				flusher.Flush()

				fmt.Fprintln(w, `{"message":{"content":"World!"},"done":true}`)
				flusher.Flush()
			}))
			defer mockServer.Close()

			cfg := &config.Config{}
			cfg.Ollama.BaseURL = mockServer.URL
			cfg.Ollama.Model = "test-model"

			client := ollama.NewOllamaClient(cfg)

			ch := make(chan types.StreamEvent, 10)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			req := requests.OllamaChatRequest{
				Messages: []requests.OllamaMessage{
					{Role: "user", Content: "Hello"},
				},
			}

			msg, err := client.QueryStreamWithTools(ctx, req, ch)
			close(ch)

			Expect(err).NotTo(HaveOccurred())
			Expect(msg.Content).To(Equal("Hello World!"))

			var events []types.StreamEvent
			for ev := range ch {
				events = append(events, ev)
			}

			Expect(len(events)).To(Equal(2))
			Expect(events[0].Content).To(Equal("Hello "))
			Expect(events[1].Content).To(Equal("World!"))
		})

		It("should return error when server returns non-200 status code", func() {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal model error"))
			}))
			defer mockServer.Close()

			cfg := &config.Config{}
			cfg.Ollama.BaseURL = mockServer.URL
			cfg.Ollama.Model = "test-model"

			client := ollama.NewOllamaClient(cfg)

			ch := make(chan types.StreamEvent, 10)
			req := requests.OllamaChatRequest{
				Messages: []requests.OllamaMessage{{Role: "user", Content: "Hello"}},
			}
			_, err := client.QueryStreamWithTools(context.Background(), req, ch)
			close(ch)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("status code 500"))
		})

		It("should return error when Ollama returns an in-stream API error", func() {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"error":"model not found"}`)
			}))
			defer mockServer.Close()

			cfg := &config.Config{}
			cfg.Ollama.BaseURL = mockServer.URL
			cfg.Ollama.Model = "nonexistent-model"

			client := ollama.NewOllamaClient(cfg)

			ch := make(chan types.StreamEvent, 10)
			req := requests.OllamaChatRequest{
				Messages: []requests.OllamaMessage{{Role: "user", Content: "Hello"}},
			}
			_, err := client.QueryStreamWithTools(context.Background(), req, ch)
			close(ch)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ollama API error: model not found"))
		})

		It("should return error when context is canceled", func() {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond)
			}))
			defer mockServer.Close()

			cfg := &config.Config{}
			cfg.Ollama.BaseURL = mockServer.URL
			cfg.Ollama.Model = "test-model"

			client := ollama.NewOllamaClient(cfg)

			ch := make(chan types.StreamEvent, 10)
			ctx, cancel := context.WithCancel(context.Background())
			cancel() // cancel immediately

			req := requests.OllamaChatRequest{
				Messages: []requests.OllamaMessage{{Role: "user", Content: "Hello"}},
			}
			_, err := client.QueryStreamWithTools(ctx, req, ch)
			close(ch)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("STREAM ERROR"))
		})
	})
})
