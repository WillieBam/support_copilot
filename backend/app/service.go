package app

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
)

type AppService struct {
	client      *appClient
	repo        *appRepository
	AuthService interfaces.IAuthService
}

func newAppService(client *appClient, repo *appRepository, authService interfaces.IAuthService) *AppService {
	return &AppService{
		client:      client,
		repo:        repo,
		AuthService: authService,
	}
}

func (s *AppService) Query(ctx context.Context, prompt string) (string, error) {
	return s.client.Query(ctx, prompt)
}
