package main

import (
	"log/slog"

	"github.com/mustafmst/ftuck/internal/commands"
)

func main() {
	cmd := commands.NewCommandWithSubcommands(
		"app",
		"aplication root",
		commands.NewCommandWithFunc(
			"init", "Initialize FTUCK", func(ctx commands.CommandContext) error {
				// get flag values
				conf, err := ctx.GetString("conf")
				if err != nil {
					return err
				}
				user, err := ctx.GetString("user")
				if err != nil {
					return err
				}
				slog.Info("executing init subcommand", "conf", conf, "user", user)
				return nil
			},
			commands.RegisterFlag("conf", "Config path", commands.StringArg, "c"),
			commands.RegisterFlag("user", "User for the operation", commands.StringArg, "u"),
		))
	err := cmd.ExecuteAsRootCommand()
	if err != nil {
		slog.Error("root command execution error", "error", err)
	}
}
