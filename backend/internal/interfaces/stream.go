package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types"
)

type IAppService interface {
	QueryStream(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error
}
