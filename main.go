package main

import (
	"log/slog"

	"github.com/mustafmst/ftuck/internal/cli"
	"github.com/mustafmst/ftuck/internal/commands"
)

func main() {
	cmd := cli.NewCommandWithSubcommands(
		"app",
		"aplication root",
		commands.CreateInitCommand(),
	)
	err := cmd.ExecuteAsRootCommand()
	if err != nil {
		slog.Error("root command execution", "error", err)
	}
}
