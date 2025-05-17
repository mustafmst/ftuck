package commands

import (
	"log/slog"
	"os"
	"path"

	"github.com/mustafmst/ftuck/internal/cli"
)

// FLAGS
const CONF_FLAG string = "conf"

// DEFAULTS
var CONF_DEFAULT string = path.Join(os.Getenv("HOME"), ".ftuck.yaml")

// DESCRITIOPNS
const CONF_DESC string = "Specify configuration path"

func CreateInitCommand() *cli.Command {
	return cli.NewCommandWithFunc(
		"init", "Initialize FTUCK", func(ctx cli.CommandContext) error {
			// get flag values
			conf, err := ctx.GetString("conf")
			if err != nil {
				return err
			}
			slog.Info("executing init subcommand", "conf", conf)
			return nil
		},
		cli.RegisterFlag(CONF_FLAG, "Config path", cli.StringArg, CONF_DEFAULT, "c"),
	)
}
