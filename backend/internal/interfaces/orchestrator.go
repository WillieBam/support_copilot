package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types/responses"
	"github.com/google/uuid"
)

type IOrchestratorService interface {
	ExecuteValidateAlert(ctx context.Context, alertID uuid.UUID) (*responses.CombinedValidationResult, error)
	ExecuteValidateAlertRaw(ctx context.Context, rawArgs string) (string, error)
}
