package app

import (
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
)

type AppRepository struct {
	User  interfaces.IUserRepository
	Alert interfaces.IAlertRepository
	LLM   interfaces.IOllamaClient
}

func NewAppRepository(ollama interfaces.IOllamaClient, user interfaces.IUserRepository, alert interfaces.IAlertRepository) *AppRepository {
	return &AppRepository{
		User:  user,
		Alert: alert,
		LLM:   ollama,
	}
}
