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
	h := endpoint.NewHandler(a.Service)

	s := utilserver.New(config.NewServerConfig("support-copilot"))

	if err := s.Start(ctx, func(e *echo.Echo) {
		e.Use(middlewares.RecoveryMiddleware())
		e.Use(middlewares.CORSMiddleware())

		// authentication login endpoint
		e.POST("/login", h.FirebaseLogin)

		g := e.Group("/query/chat")
		if config.Get().Auth.Enabled {
			g.Use(echo.WrapMiddleware(middlewares.AuthMiddleware))
		}
		g.POST("", h.Query)
	}); err != nil {
		slog.Error("server stopped", "err", err)
	}
}
