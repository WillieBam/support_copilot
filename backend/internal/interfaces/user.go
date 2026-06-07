package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByFirebaseUID(ctx context.Context, firebaseUid string) (*models.User, error)
}
