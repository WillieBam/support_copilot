package command

import (
	"context"
	"strings"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types"
)

// CommandHandler is the function signature for slash command handlers.
type CommandHandler func(ctx context.Context, prompt string) (*types.CommandResult, error)

type CommandInterceptor struct {
	handlers map[string]CommandHandler
}

func NewCommandInterceptor() interfaces.ICommandInterceptor {
	ci := &CommandInterceptor{
		handlers: make(map[string]CommandHandler),
	}
	ci.RegisterCommand("/quit", handleQuitCommand) // register command here, add on when needed
	return ci
}

func (ci *CommandInterceptor) RegisterCommand(command string, handler CommandHandler) {
	ci.handlers[strings.ToLower(command)] = handler
}

// Intercept is the function to check prompt against all registered command
func (ci *CommandInterceptor) Intercept(ctx context.Context, prompt string) (*types.CommandResult, error) {
	trimmed := strings.ToLower(strings.TrimSpace(prompt))
	for cmd, handler := range ci.handlers {
		if trimmed == cmd || strings.HasPrefix(trimmed, cmd+" ") || strings.HasPrefix(trimmed, cmd) {
			return handler(ctx, prompt)
		}
	}
	return &types.CommandResult{Handled: false}, nil
}

// function to bind with the RegisterCommand
func handleQuitCommand(ctx context.Context, prompt string) (*types.CommandResult, error) {
	return &types.CommandResult{
		Handled: true,
		Message: "LLM processing stopped by /quit command.",
	}, nil
}
