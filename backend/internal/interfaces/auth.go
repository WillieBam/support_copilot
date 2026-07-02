package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types"
)

type IAuthService interface {
	ExchangeToken(ctx context.Context, firebaseToken string) (string, *types.Claims, error)
	ParseAndValidateAuthToken(ctx context.Context, tokenString string) (*types.Claims, error)
}
