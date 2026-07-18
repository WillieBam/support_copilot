package service_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/internal/mocks"
	"github.com/WillieBam/support_copilot/backend/internal/service"
	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/WillieBam/support_copilot/backend/types/models"
)

var _ = Describe("AppService (Streaming & Alerts)", func() {
	var (
		appSvc        interfaces.IAppService
		mockAlertRepo *mocks.IAlertRepository
		mockOllama    *mocks.IOllamaClient
		mockMcpOne    *mocks.IMCPClient
		ctx           context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockAlertRepo = &mocks.IAlertRepository{}
		mockOllama = &mocks.IOllamaClient{}
		mockMcpOne = &mocks.IMCPClient{}

		appSvc = service.NewAppService(mockAlertRepo, mockOllama, mockMcpOne)
	})

	Context("QueryStream", func() {
		It("should connect to Ollama server and stream token events correctly", func() {
			mockOllama.On("QueryStream", mock.Anything, "hello test", mock.Anything).
				Return(nil).
				Run(func(args mock.Arguments) {
					streamChan := args.Get(2).(chan<- types.StreamEvent)
					streamChan <- types.StreamEvent{Type: "reasoning", Content: "Analyzing user prompt...\n "}
					streamChan <- types.StreamEvent{Type: "reasoning", Content: "Connecting to Llama 3.2...\n"}
					streamChan <- types.StreamEvent{Type: "text", Content: "Hello"}
					streamChan <- types.StreamEvent{Type: "text", Content: " world!"}
				})

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
			mockOllama.AssertExpectations(GinkgoT())
		})

		It("should return an error if the server returns non-200 status code", func() {
			mockOllama.On("QueryStream", mock.Anything, "hello test", mock.Anything).
				Return(errors.New("Ollama returned status 500: ollama internal error"))

			streamChan := make(chan types.StreamEvent, 10)

			err := appSvc.QueryStream(ctx, "hello test", streamChan)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Ollama returned status 500"))
			mockOllama.AssertExpectations(GinkgoT())
		})

		It("should return an error if client cancels the context", func() {
			mockOllama.On("QueryStream", mock.Anything, "hello test", mock.Anything).
				Return(context.Canceled)

			cancelCtx, cancel := context.WithCancel(ctx)
			cancel() // cancel immediately

			streamChan := make(chan types.StreamEvent, 10)

			err := appSvc.QueryStream(cancelCtx, "hello test", streamChan)
			Expect(err).To(HaveOccurred())
			mockOllama.AssertExpectations(GinkgoT())
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

	Context("RetrieveAlert", func() {
		It("should successfully retrieve alert from repository", func() {
			alertID := uuid.New()
			expectedAlert := &models.Alert{
				ID:          alertID,
				ServiceName: "test-service",
				Severity:    "high",
			}
			mockAlertRepo.On("RetrieveAlert", mock.Anything, alertID).Return(expectedAlert, nil)

			alert, err := appSvc.RetrieveAlert(ctx, alertID)
			Expect(err).NotTo(HaveOccurred())
			Expect(alert).To(Equal(expectedAlert))
			mockAlertRepo.AssertExpectations(GinkgoT())
		})

		It("should return 'alert not found' error if repository returns gorm.ErrRecordNotFound", func() {
			alertID := uuid.New()
			mockAlertRepo.On("RetrieveAlert", mock.Anything, alertID).Return(nil, gorm.ErrRecordNotFound)

			alert, err := appSvc.RetrieveAlert(ctx, alertID)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("alert not found"))
			mockAlertRepo.AssertExpectations(GinkgoT())
			Expect(alert).To(BeNil())
		})

		It("should return other repository errors as-is", func() {
			alertID := uuid.New()
			mockAlertRepo.On("RetrieveAlert", mock.Anything, alertID).Return(nil, errors.New("Internal Server Error"))

			alert, err := appSvc.RetrieveAlert(ctx, alertID)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Internal Server Error"))
			mockAlertRepo.AssertExpectations(GinkgoT())
			Expect(alert).To(BeNil())

		})
	})
})
