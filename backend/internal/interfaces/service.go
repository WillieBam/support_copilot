package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types"
	"github.com/google/uuid"
)

type IAppService interface {
	QueryStream(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error
	IngestAlert(ctx context.Context, incidentID uuid.UUID, serviceName, severity, metrics string) error
	ProcessAlert(ctx context.Context, rawMetrics string, streamChan chan<- types.StreamEvent) error
}
