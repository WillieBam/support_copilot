package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/google/uuid"
)

type IAppService interface {
	QueryStreamWithTools(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error
	IngestAlert(ctx context.Context, incidentID uuid.UUID, serviceName, severity, metrics string) error
	RetrieveAlert(ctx context.Context, id uuid.UUID) (*models.Alert, error)
	Intercept(ctx context.Context, prompt string) (*types.CommandResult, error)
}
