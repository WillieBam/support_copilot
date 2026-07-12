package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/google/uuid"
)

type IAlertRepository interface {
	StoreAlert(ctx context.Context, alert *models.Alert) error
	RetrieveAlert(ctx context.Context, id uuid.UUID) (*models.Alert, error)
}
