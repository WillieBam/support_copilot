package config

import (
	"time"

	"github.com/WillieBam/support_copilot/backend/utils/server"
)

type IServer = server.IServer

func NewServerConfig(name string) *server.ServerConfig {
	return server.NewServerConfig(server.ServerConfigOptions{
		Name: name,
		PortGetter: func() int {
			return Get().Http.Port
		},
		ShutDownTimeOutGetter: func() time.Duration {
			return Get().Http.ShutdownTimeOut
		},
	})
}
