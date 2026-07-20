package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types/requests"
	"github.com/WillieBam/support_copilot/backend/types/responses"
	"github.com/google/uuid"
)

type orchestratorService struct {
	alertRepo  interfaces.IAlertRepository
	mcpClient1 interfaces.IMCPClient
}

func NewOrchestratorService(repo interfaces.IAlertRepository, mcpClient1 interfaces.IMCPClient) interfaces.IOrchestratorService {
	return &orchestratorService{
		alertRepo:  repo,
		mcpClient1: mcpClient1,
	}
}

// ExecuteValidateAlert fetches alert metrics from Postgres and predicts anomalies via Python MCP.
func (s *orchestratorService) ExecuteValidateAlert(ctx context.Context, alertID uuid.UUID) (*responses.CombinedValidationResult, error) {
	slog.Info("[ORCHESTRATOR] Fetching alert from database", "alert_id", alertID.String())

	// fetch alert from Postgres
	alertRecord, err := s.alertRepo.RetrieveAlertbyID(ctx, alertID)
	if err != nil {
		slog.Error("[ORCHESTRATOR] Failed to fetch alert from DB", "alert_id", alertID.String(), "err", err)
		return nil, fmt.Errorf("failed to fetch alert #%s from database: %w", alertID, err)
	}

	slog.Info("[ORCHESTRATOR] Alert retrieved from DB", "service", alertRecord.ServiceName, "severity", alertRecord.Severity)

	// unmarshal alert JSON metrics string into AnomalyDetectionRequest struct
	var metrics requests.AnomalyDetectionRequest
	if err := json.Unmarshal([]byte(alertRecord.Metrics), &metrics); err != nil {
		slog.Error("[ORCHESTRATOR] Failed to unmarshal alert metrics JSON", "err", err)
		return nil, fmt.Errorf("failed to parse alert metrics JSON: %w", err)
	}

	slog.Info("[ORCHESTRATOR] Invoking Python MCP Server detect_anomalies", "cpu", metrics.CpuUsage, "latency", metrics.ResponseLatency)

	// directly invoke Python MCP Server (detect_anomalies)
	mcpResp, err := s.mcpClient1.DetectAnomalies(ctx, metrics)
	if err != nil {
		slog.Error("[ORCHESTRATOR] MCP anomaly detection failed", "err", err)
		return nil, fmt.Errorf("failed to analyze metrics via MCP server: %w", err)
	}

	slog.Info("[ORCHESTRATOR] MCP anomaly detection succeeded", "label", mcpResp.Label, "status", mcpResp.Status)

	// assemble combined payload package
	return &responses.CombinedValidationResult{
		AlertID:      alertRecord.ID.String(),
		ServiceName:  alertRecord.ServiceName,
		Severity:     alertRecord.Severity,
		ReceivedAt:   alertRecord.ReceivedAt,
		Metrics:      metrics,
		MLPrediction: *mcpResp,
	}, nil
}

// ExecuteValidateAlertRaw parses raw LLM JSON arguments and delegates execution.
func (s *orchestratorService) ExecuteValidateAlertRaw(ctx context.Context, rawArgs string) (string, error) {
	slog.Info("[ORCHESTRATOR] ExecuteValidateAlertRaw triggered", "rawArgs", rawArgs)

	var args struct {
		AlertID string `json:"alert_id"`
	}
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		slog.Error("[ORCHESTRATOR] Failed to parse raw tool arguments", "rawArgs", rawArgs, "err", err)
		return "", fmt.Errorf("invalid tool arguments: %w", err)
	}

	cleanAlertID := strings.TrimSpace(args.AlertID)
	if cleanAlertID == "" || cleanAlertID == "null" || cleanAlertID == "none" || cleanAlertID == "undefined" {
		slog.Warn("[ORCHESTRATOR] Empty or dummy alert_id provided", "alertID", args.AlertID)
		return "", fmt.Errorf("no valid alert_id provided: %q", args.AlertID)
	}

	alertUUID, err := uuid.Parse(cleanAlertID)
	if err != nil {
		slog.Error("[ORCHESTRATOR] Invalid alert UUID", "alertID", args.AlertID, "err", err)
		return "", fmt.Errorf("invalid alert id %q: %w", args.AlertID, err)
	}

	result, err := s.ExecuteValidateAlert(ctx, alertUUID)
	if err != nil {
		return "", err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		slog.Error("[ORCHESTRATOR] Failed to marshal validation result", "err", err)
		return "", fmt.Errorf("failed to marshal validation result: %w", err)
	}

	slog.Info("[ORCHESTRATOR] Combined validation package built successfully", "result", string(resultBytes))
	return string(resultBytes), nil
}
