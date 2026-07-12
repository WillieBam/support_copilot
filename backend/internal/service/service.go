package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
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
	// compact the JSON metrics to strip all whitespace, \r, and empty lines
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(metrics)); err != nil {
		// Metrics is not valid JSON  will log a warning and store the raw value
		slog.Warn("metrics field is not valid JSON, storing raw value", "err", err)
		buf.WriteString(metrics)
	}

	alert := &models.Alert{
		ID:          uuid.New(),
		IncidentID:  incidentID,
		ServiceName: serviceName,
		Severity:    severity,
		Metrics:     buf.String(),
		ReceivedAt:  time.Now(),
	}
	return s.alertRepo.StoreAlert(ctx, alert)
}

func (s *AppService) QueryStream(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error {
	return s.ollamaClient.QueryStream(ctx, prompt, streamChan)
}
