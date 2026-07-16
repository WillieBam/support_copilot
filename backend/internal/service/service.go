package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/WillieBam/support_copilot/backend/types/requests"
	"github.com/google/uuid"
)

type AppService struct {
	alertRepo    interfaces.IAlertRepository
	ollamaClient interfaces.IOllamaClient
	mcpClient    interfaces.IMCPClient
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

func (s *AppService) ProcessAlert(ctx context.Context, rawMetrics string, streamChan chan<- types.StreamEvent) error {
	streamChan <- types.StreamEvent{Type: "status", Content: "Step 1/3: Standardizing multi-team alert ingestion vector..."}

	// Unmarshal your string database records into the 9-metric infrastructure model
	var metrics requests.AnomalyDetectionRequest
	if err := json.Unmarshal([]byte(rawMetrics), &metrics); err != nil {
		return fmt.Errorf("failed parsing ingestion metrics: %w", err)
	}

	streamChan <- types.StreamEvent{Type: "status", Content: "Step 2/3: Invoking IsolationForest tool on mcp_server_1..."}

	mlResult, err := s.mcpClient.DetectAnomalies(ctx, metrics)
	if err != nil {
		return fmt.Errorf("mcp tool execution failed: %w", err)
	}

	streamChan <- types.StreamEvent{Type: "status", Content: "Step 3/3: Passing payload + anomaly state to LLM Validation Agent..."}

	prompt := fmt.Sprintf(
		"System Alert Context:\n- Service Target: Datagateway Gateway\n- ML IsolationForest Classification: %s (DB Code: %d)\n\nTelemetry Raw Vector:\n- CPU: %.2f%%\n- Memory: %.2f%%\n- Latency: %.2fms\n- Availability: %.2f%%\n\nAnalyze this environment profile and output diagnostic recommendations.",
		mlResult.Label, mlResult.Status, metrics.CpuUsage, metrics.MemoryUsage, metrics.ResponseLatency, metrics.AvailabilityPercent,
	)

	err = s.ollamaClient.QueryStream(ctx, prompt, streamChan)
	if err != nil {
		return fmt.Errorf("agent validation reflection failed: %w", err)
	}

	return nil
}
