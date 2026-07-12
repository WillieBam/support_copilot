package service

import (
	"context"
	"time"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/google/uuid"
)

type AppService struct {
	alertRepo    interfaces.IAlertRepository
	ollamaClient interfaces.IOllamaClient
}

func NewAppService(alertRepo interfaces.IAlertRepository, ollamaClient interfaces.IOllamaClient) interfaces.IAppService {
	return &AppService{
		alertRepo:    alertRepo,
		ollamaClient: ollamaClient,
	}
}

func (s *AppService) IngestAlert(ctx context.Context, incidentID uuid.UUID, serviceName, severity, metrics string) error {
	alert := &models.Alert{
		ID:          uuid.New(),
		IncidentID:  incidentID,
		ServiceName: serviceName,
		Severity:    severity,
		Metrics:     metrics,
		ReceivedAt:  time.Now(),
	}
	return s.alertRepo.StoreAlert(ctx, alert)
}

func (s *AppService) QueryStream(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error {
	return s.ollamaClient.QueryStream(ctx, prompt, streamChan)
}
