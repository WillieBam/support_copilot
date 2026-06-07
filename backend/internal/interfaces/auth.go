package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types/models"
)

type IAuthService interface {
	IsUserExists(ctx context.Context, idToken string) bool
	LoginOrRegister(ctx context.Context, idToken string) (*models.User, error)
}
