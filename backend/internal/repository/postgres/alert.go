package postgres

import (
	"context"
	"errors"

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
	result := a.db.WithContext(ctx).Create(&alert)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *alertRepository) RetrieveAlertbyID(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	var alert models.Alert
	if err := a.db.WithContext(ctx).First(&alert, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, errors.New("Internal Server Error")
	}
	return &alert, nil
}
