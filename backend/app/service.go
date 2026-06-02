package app

import (
	"context"
)

type appService struct {
	client     *appClient
	repository *appRepository
}

func newAppService(client *appClient, repository *appRepository) *appService {
	return &appService{
		client:     client,
		repository: repository,
	}
}

// Query satisfies interfaces.ISupportCopilotService.
func (s *appService) Query(ctx context.Context, input string) (string, error) {
	return s.client.Query(ctx, input)
}
