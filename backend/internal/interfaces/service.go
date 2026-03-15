package interfaces

import "context"

type ISupportCopilotService interface {
	Query(ctx context.Context, input string) (string, error)
}
