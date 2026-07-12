package service_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/google/uuid"

	"github.com/WillieBam/support_copilot/backend/app/config"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/internal/mocks"
	"github.com/WillieBam/support_copilot/backend/internal/service"
	"github.com/WillieBam/support_copilot/backend/types"
)

var _ = Describe("AppService (Streaming & Alerts)", func() {
	var (
		appSvc        interfaces.IAppService
		mockAlertRepo *mocks.IAlertRepository
		ctx           context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockAlertRepo = &mocks.IAlertRepository{}
		appSvc = service.NewAppService(mockAlertRepo)
	})

	Context("QueryStream", func() {
		It("should connect to Ollama server and stream token events correctly", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.URL.Path).To(Equal("/api/chat"))
				Expect(r.Method).To(Equal(http.MethodPost))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				// Write mock streaming chunks from Ollama
				_, _ = w.Write([]byte(`{"message":{"content":"Hello"},"done":false}`))
				_, _ = w.Write([]byte(`{"message":{"content":" world!"},"done":true}`))
			}))
			defer server.Close()

			config.Get().Ollama.BaseURL = server.URL
			config.Get().Ollama.Model = "llama3.2"

			streamChan := make(chan types.StreamEvent, 10)

			err := appSvc.QueryStream(ctx, "hello test", streamChan)
			Expect(err).NotTo(HaveOccurred())
			close(streamChan)

			var events []types.StreamEvent
			for ev := range streamChan {
				events = append(events, ev)
			}

			// Expected result:
			// 1. "reasoning" - Analyzing user prompt...
			// 2. "reasoning" - Connecting to Llama 3.2...
			// 3. "text" - Hello
			// 4. "text" -  world!
			Expect(len(events)).To(BeNumerically(">=", 3))
			Expect(events[0].Type).To(Equal("reasoning"))
			Expect(events[2].Type).To(Equal("text"))
			Expect(events[2].Content).To(Equal("Hello"))
			Expect(events[3].Type).To(Equal("text"))
			Expect(events[3].Content).To(Equal(" world!"))
		})

		It("should return an error if the server returns non-200 status code", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("ollama internal error"))
			}))
			defer server.Close()

			config.Get().Ollama.BaseURL = server.URL

			streamChan := make(chan types.StreamEvent, 10)

			err := appSvc.QueryStream(ctx, "hello test", streamChan)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Ollama returned status 500"))
		})

		It("should return an error if client cancels the context", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Wait or write chunk
				time.Sleep(100 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			config.Get().Ollama.BaseURL = server.URL

			cancelCtx, cancel := context.WithCancel(ctx)
			cancel() // cancel immediately

			streamChan := make(chan types.StreamEvent, 10)

			err := appSvc.QueryStream(cancelCtx, "hello test", streamChan)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("IngestAlert", func() {
		It("should successfully store alert in repository", func() {
			incidentID := uuid.New()
			mockAlertRepo.On("StoreAlert", mock.Anything, mock.Anything).Return(nil)

			err := appSvc.IngestAlert(ctx, incidentID, "auth-service", "critical", "cpu_util > 90%")
			Expect(err).NotTo(HaveOccurred())
			mockAlertRepo.AssertExpectations(GinkgoT())
		})

		It("should return error if repository fails to store alert", func() {
			incidentID := uuid.New()
			mockAlertRepo.On("StoreAlert", mock.Anything, mock.Anything).Return(errors.New("db error"))

			err := appSvc.IngestAlert(ctx, incidentID, "auth-service", "critical", "cpu_util > 90%")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("db error"))
			mockAlertRepo.AssertExpectations(GinkgoT())
		})
	})
})
