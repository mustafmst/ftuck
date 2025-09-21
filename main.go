package main

import (
	"context"

	"github.com/mustafmst/ftuck/internal/cli"
	"github.com/mustafmst/ftuck/internal/commands"
	"github.com/mustafmst/ftuck/internal/logging"
)

func main() {
	// Initialize structured logging
	logConfig := logging.DefaultConfig()
	logger := logging.InitLogger(logConfig)
	
	ctx := context.Background()
	cmd := cli.NewCommandWithSubcommands(
		"ftuck",
		"File synchronization tool for managing dotfiles",
		commands.CreateInitCommand(ctx),
		commands.CreateAddSyncCommand(ctx),
		commands.CreateSyncAllCommand(ctx),
	)
	
	err := cmd.ExecuteAsRootCommand()
	if err != nil {
		// Don't log help requests as errors
		errStr := err.Error()
		if errStr == "command not found" || errStr == "" {
			// These are expected conditions (help, unknown commands), don't log as errors
			return
		}
		logger.Error("application error", "error", err)
	}
}
