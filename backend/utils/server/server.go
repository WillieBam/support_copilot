package server

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo   *echo.Echo
	config IServer
}

func New(config IServer) *Server {
	e := echo.New()
	e.HideBanner = true
	return &Server{
		Echo:   e,
		config: config,
	}
}

// Start registers routes via setup, then runs the server. It gracefully
// shuts down when ctx is cancelled.
func (s *Server) Start(ctx context.Context, setup func(*echo.Echo)) error {
	setup(s.Echo)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.GetShutdownTimeOutDuration())
		defer cancel()
		s.Echo.Shutdown(shutdownCtx) //nolint:errcheck
	}()

	return s.Echo.Start(":" + s.config.Port())
}
