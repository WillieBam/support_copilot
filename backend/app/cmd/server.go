package cmd

import (
	"context"
	"log/slog"

	"github.com/WillieBam/support_copilot/backend/app"
	"github.com/WillieBam/support_copilot/backend/app/config"
	"github.com/WillieBam/support_copilot/backend/internal/endpoint"
	"github.com/WillieBam/support_copilot/backend/middlewares"
	utilserver "github.com/WillieBam/support_copilot/backend/utils/server"
	"github.com/labstack/echo/v5"
	"github.com/spf13/cobra"
)

var supportCopilotCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  `Run Support Copilot Server`,
	Run:   supportCopilotExec,
}

func init() {
	rootCmd.AddCommand(supportCopilotCmd)
}

func supportCopilotExec(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	a := app.NewApp()
	h := endpoint.NewHandler(a.Service, a.AuthService)

	s := utilserver.New(config.NewServerConfig("support-copilot"))

	if err := s.Start(ctx, func(e *echo.Echo) {
		e.Use(middlewares.RecoveryMiddleware())
		e.Use(middlewares.CORSMiddleware())

		e.POST("/auth/exchange", h.TokenExchangeHandler)

		// '/api' group endpoints
		apiGroup := e.Group("/api")
		apiGroup.Use(middlewares.AuthMiddleware(a.AuthService))
		apiGroup.GET("/auth/me", h.Me)
		apiGroup.GET("/alerts/ingest", h.IngestAlert)

		// '/query' group endpoints
		g := e.Group("/query")
		g.Use(middlewares.AuthMiddleware(a.AuthService))
		g.POST("/chat", h.Query)
	}); err != nil {
		slog.Error("server gave up", "err", err)
	}
}
