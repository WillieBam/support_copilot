package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types/requests"
)

type IToolRegistry interface {
	GetTools() []requests.OllamaTool
	Execute(ctx context.Context, name string, rawArgs string) (string, error)
}
