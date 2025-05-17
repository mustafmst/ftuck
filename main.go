package main

import (
	"context"
	"log/slog"

	"github.com/mustafmst/ftuck/internal/cli"
	"github.com/mustafmst/ftuck/internal/commands"
)

func main() {
	ctx := context.Background()
	cmd := cli.NewCommandWithSubcommands(
		"app",
		"aplication root",
		commands.CreateInitCommand(ctx),
		commands.CreateAddSyncCommand(ctx),
	)
	err := cmd.ExecuteAsRootCommand()
	if err != nil {
		slog.Error("root command execution", "error", err)
	}
}
