package app

import (
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
)

type AppRepository struct {
	Client *appClient
	User   interfaces.IUserRepository
	Alert  interfaces.IAlertRepository
}

func NewAppRepository(client *appClient, user interfaces.IUserRepository, alert interfaces.IAlertRepository) *AppRepository {
	return &AppRepository{
		Client: client,
		User:   user,
		Alert:  alert,
	}
}
