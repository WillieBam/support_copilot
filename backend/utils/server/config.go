package server

import (
	"strconv"
	"time"
)

type IServer interface {
	Name() string
	Port() string
	GetShutdownTimeOutDuration() time.Duration
}

type ServerConfig struct {
	name                  string
	portGetter            func() int
	shutDownTimeOutGetter func() time.Duration
}

type ServerConfigOptions struct {
	Name                  string
	PortGetter            func() int
	ShutDownTimeOutGetter func() time.Duration
}

func NewServerConfig(opts ServerConfigOptions) *ServerConfig {
	return &ServerConfig{
		name:                  opts.Name,
		portGetter:            opts.PortGetter,
		shutDownTimeOutGetter: opts.ShutDownTimeOutGetter,
	}
}

func (c *ServerConfig) Name() string {
	return c.name
}

func (c *ServerConfig) Port() string {
	p := c.portGetter()
	if p == 0 {
		return "8080"
	}

	return strconv.Itoa(p)
}

func (c *ServerConfig) GetShutdownTimeOutDuration() time.Duration {
	d := c.shutDownTimeOutGetter()
	if d == 0 {
		return 10 * time.Second
	}
	return d
}
