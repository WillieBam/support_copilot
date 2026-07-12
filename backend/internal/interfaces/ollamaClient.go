package interfaces

import (
	context "context"

	"github.com/WillieBam/support_copilot/backend/types"
)

type IOllamaClient interface {
	QueryStream(ctx context.Context, prompt string, streamChan chan<- types.StreamEvent) error
}
