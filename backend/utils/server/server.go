package server

import (
	"context"

	"github.com/labstack/echo/v5"
)

type Server struct {
	Echo   *echo.Echo
	config IServer
}

func New(config IServer) *Server {
	e := echo.New()
	return &Server{
		Echo:   e,
		config: config,
	}
}

// Start registers routes via setup, then runs the server. It gracefully
// shuts down when ctx is cancelled.
func (s *Server) Start(ctx context.Context, setup func(*echo.Echo)) error {
	setup(s.Echo)

	sc := echo.StartConfig{
		Address:       ":" + s.config.Port(),
		HideBanner:    true,
		GracefulTimeout: s.config.GetShutdownTimeOutDuration(),
	}

	return sc.Start(ctx, s.Echo)
}
