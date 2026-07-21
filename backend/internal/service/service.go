package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/WillieBam/support_copilot/backend/internal/command"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/internal/tools"
	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/WillieBam/support_copilot/backend/types/requests"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppService struct {
	alertRepo          interfaces.IAlertRepository
	ollamaClient       interfaces.IOllamaClient
	mcpClient          interfaces.IMCPClient
	orchestrator       interfaces.IOrchestratorService
	toolRegistry       interfaces.IToolRegistry
	commandInterceptor interfaces.ICommandInterceptor
}

func NewAppService(alertRepo interfaces.IAlertRepository, ollamaClient interfaces.IOllamaClient, mcpClient interfaces.IMCPClient, toolRegAndInterceptor ...interface{}) interfaces.IAppService {
	orchestrator := NewOrchestratorService(alertRepo, mcpClient)

	var registry interfaces.IToolRegistry
	var cmdInterceptor interfaces.ICommandInterceptor

	for _, arg := range toolRegAndInterceptor {
		if tr, ok := arg.(interfaces.IToolRegistry); ok && tr != nil {
			registry = tr
		}
		if ci, ok := arg.(interfaces.ICommandInterceptor); ok && ci != nil {
			cmdInterceptor = ci
		}
	}

	if registry == nil {
		tr := tools.NewToolRegistry()
		tools.RegisterDefaultTools(tr, orchestrator)
		registry = tr
	}

	if cmdInterceptor == nil {
		cmdInterceptor = command.NewCommandInterceptor()
	}

	return &AppService{
		alertRepo:          alertRepo,
		ollamaClient:       ollamaClient,
		mcpClient:          mcpClient,
		orchestrator:       orchestrator,
		toolRegistry:       registry,
		commandInterceptor: cmdInterceptor,
	}
}

func (s *AppService) IngestAlert(ctx context.Context, incidentID uuid.UUID, serviceName, severity, metrics string) error {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(metrics)); err != nil {
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

func isValidToolCallArgs(toolName string, args map[string]interface{}) bool {
	if toolName == "validate_alert" {
		alertIDVal, exists := args["alert_id"]
		if !exists || alertIDVal == nil {
			return false
		}
		alertID, ok := alertIDVal.(string)
		if !ok {
			return false
		}
		alertID = strings.TrimSpace(alertID)
		if alertID == "" || alertID == "null" || alertID == "none" || alertID == "undefined" || alertID == "{alert_id}" || alertID == "00000000-0000-0000-0000-000000000000" {
			return false
		}
		if _, err := uuid.Parse(alertID); err != nil {
			return false
		}
	}
	return true
}

func (s *AppService) Intercept(ctx context.Context, prompt string) (*types.CommandResult, error) {
	if s.commandInterceptor != nil {
		return s.commandInterceptor.Intercept(ctx, prompt)
	}
	return &types.CommandResult{Handled: false}, nil
}

func (s *AppService) QueryStreamWithTools(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error {
	slog.Info("[APP SERVICE] QueryStreamWithTools started", "prompt", prompt)

	res, err := s.Intercept(ctx, prompt)
	if err != nil {
		slog.Error("[APP SERVICE] Command interceptor error", "err", err)
		return err
	}
	if res != nil && res.Handled {
		slog.Info("[APP SERVICE] Prompt intercepted by command parser", "prompt", prompt)
		if res.Message != "" {
			streamChan <- types.StreamEvent{
				Type:    "text",
				Content: res.Message,
			}
		}
		return nil
	}

	systemPrompt := "You are a Support Copilot. Support Copilot is responsible to assist support engineer in resolve incidents. You have access to tools to inspect system state and alerts. Call tools ONLY when the user provides an alert ID or explicitly requests tool execution. Do NOT call any tools for general conversational input, greetings, or acknowledgments like 'ok', 'thanks', or 'thank you'."

	messages := []requests.OllamaMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: prompt},
	}

	availableTools := s.toolRegistry.GetTools()

	req := requests.OllamaChatRequest{
		Messages: messages,
		Tools:    availableTools,
	}

	// Call Ollama streaming with tool declarations dynamically provided by ToolRegistry
	assistantMsg, err := s.ollamaClient.QueryStreamWithTools(ctx, req, streamChan)
	if err != nil {
		slog.Error("[APP SERVICE] First pass QueryStreamWithTools failed", "err", err)
		return err
	}

	// Check if LLM requested a tool execution
	if assistantMsg != nil && len(assistantMsg.ToolCalls) > 0 {
		// Emit reasoning event so React UI displays it immediately in the reasoning block
		streamChan <- types.StreamEvent{
			Type:    "reasoning",
			Content: fmt.Sprintln("🔍 Accessing Database to get alert... "),
		}

		for _, toolCall := range assistantMsg.ToolCalls {
			toolName := toolCall.Function.Name
			slog.Info("[APP SERVICE] Ollama triggered tool call", "tool", toolName, "args", toolCall.Function.Arguments)

			// Guardrail: Pre-check tool arguments for dummy/invalid values BEFORE executing or emitting tool reasoning
			if !isValidToolCallArgs(toolName, toolCall.Function.Arguments) {
				slog.Warn("[APP SERVICE] Tool call skipped due to dummy or missing valid parameters", "tool", toolName, "args", toolCall.Function.Arguments)

				if strings.TrimSpace(assistantMsg.Content) == "" {
					slog.Info("[APP SERVICE] Falling back to direct text response without tools")
					fallbackReq := requests.OllamaChatRequest{
						Messages: messages,
					}
					_, fallbackErr := s.ollamaClient.QueryStreamWithTools(ctx, fallbackReq, streamChan)
					return fallbackErr
				}
				continue
			}

			// Emit reasoning event so React UI displays it immediately in the reasoning block
			streamChan <- types.StreamEvent{
				Type:    "reasoning",
				Content: fmt.Sprintf("🔍 Intercepted tool call: %s. Executing tool...\n", toolName),
			}

			argsBytes, _ := json.Marshal(toolCall.Function.Arguments)
			toolResult, err := s.toolRegistry.Execute(ctx, toolName, string(argsBytes))
			if err != nil {
				slog.Warn("[APP SERVICE] Tool execution failed via ToolRegistry", "tool", toolName, "err", err)

				// Guardrail: If tool call failed (e.g. invalid alert_id "null") and no content was streamed yet,
				// fall back to a conversational text stream without tools to avoid sending raw error noise to user
				if strings.TrimSpace(assistantMsg.Content) == "" {
					slog.Info("[APP SERVICE] Falling back to direct text response without tools")
					fallbackReq := requests.OllamaChatRequest{
						Messages: messages,
					}
					_, fallbackErr := s.ollamaClient.QueryStreamWithTools(ctx, fallbackReq, streamChan)
					return fallbackErr
				}
				toolResult = fmt.Sprintf(`{"error": "%s"}`, err.Error())
			}

			slog.Info("[APP SERVICE] Tool result retrieved", "tool", toolName, "resultLen", len(toolResult))

			messages = append(messages, *assistantMsg)
			messages = append(messages, requests.OllamaMessage{
				Role:    "tool",
				Content: toolResult,
			})

			// 2nd Pass: Stream Ollama's final synthesis based on the tool result
			secondReq := requests.OllamaChatRequest{
				Messages: messages,
			}
			_, err = s.ollamaClient.QueryStreamWithTools(ctx, secondReq, streamChan)
			return err
		}
	}

	return nil
}

func (s *AppService) RetrieveAlert(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	alert, err := s.alertRepo.RetrieveAlertbyID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alert not found")
		}
		return nil, err
	}
	return alert, nil
}
