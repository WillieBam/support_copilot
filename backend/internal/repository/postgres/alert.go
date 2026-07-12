package postgres

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type alertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) interfaces.IAlertRepository {
	return &alertRepository{db: db}
}

func (a *alertRepository) StoreAlert(ctx context.Context, alert *models.Alert) error {
	return a.db.WithContext(ctx).Create(alert).Error
}

func (a *alertRepository) RetrieveAlert(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	var alert models.Alert
	if err := a.db.WithContext(ctx).First(&alert, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}
