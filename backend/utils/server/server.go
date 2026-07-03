package server

import (
	"context"
	"net/http"
	"time"

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
		Address:         ":" + s.config.Port(),
		HideBanner:      true,
		GracefulTimeout: s.config.GetShutdownTimeOutDuration(),
		BeforeServeFunc: func(s *http.Server) error {
			s.WriteTimeout = 0
			s.ReadTimeout = 5 * time.Minute
			s.IdleTimeout = 120 * time.Second
			return nil
		},
	}

	return sc.Start(ctx, s.Echo)
}
