package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types/requests"
)

type ToolHandler func(ctx context.Context, rawArgs string) (string, error)

type ToolDefinition struct {
	Tool    requests.OllamaTool
	Handler ToolHandler
}

type ToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]ToolDefinition
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]ToolDefinition),
	}
}

func (r *ToolRegistry) Register(name string, tool requests.OllamaTool, handler ToolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[name] = ToolDefinition{
		Tool:    tool,
		Handler: handler,
	}
}

func (r *ToolRegistry) GetTools() []requests.OllamaTool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]requests.OllamaTool, 0, len(r.tools))
	for _, def := range r.tools {
		result = append(result, def.Tool)
	}
	return result
}

func (r *ToolRegistry) Execute(ctx context.Context, name string, rawArgs string) (string, error) {
	r.mu.RLock()
	def, exists := r.tools[name]
	r.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("tool %q not registered in ToolRegistry", name)
	}
	return def.Handler(ctx, rawArgs)
}

// RegisterDefaultTools registers standard backend tools such as validate_alert
func RegisterDefaultTools(registry *ToolRegistry, orchestrator interfaces.IOrchestratorService) {
	registry.Register("validate_alert", requests.OllamaTool{
		Type: "function",
		Function: requests.OllamaFunction{
			Name:        "validate_alert",
			Description: "Retrieves telemetry metrics for a given alert_id from Postgres and predicts whether the system state is Anomaly or Normal using IsolationForest ML. Call 'validate_alert' ONLY when an alert ID is provided or explicit alert validation is requested. Do NOT call this tool for general conversational input, greetings, or acknowledgments.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"alert_id": map[string]interface{}{
						"type":        "string",
						"description": "The unique UUID string of the alert to validate (e.g. '550e8400-e29b-41d4-a716-446655440000')",
					},
				},
				"required": []string{"alert_id"},
			},
		},
	}, func(ctx context.Context, rawArgs string) (string, error) {
		return orchestrator.ExecuteValidateAlertRaw(ctx, rawArgs)
	})
}
