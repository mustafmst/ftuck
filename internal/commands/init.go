package commands

import (
	"context"
	"os"
	"path"

	"github.com/mustafmst/ftuck/internal/cli"
	"github.com/mustafmst/ftuck/internal/config"
)

// FLAGS
const (
	CONF_FLAG string = "conf"
	WD_FLAG   string = "workdir"
)

// DEFAULTS
var (
	CONF_DEFAULT string = path.Join(os.Getenv("HOME"), ".ftuck.yaml")
	WD_DEFAULt   string = "not provided"
)

// DESCRITIOPN
var (
	CONF_DESC string = "Specify configuration path. (DEFAULT=" + CONF_DEFAULT + ")"
	WD_DESC   string = "Use different working directory than current."
)

type initCommand struct {
	ctx context.Context
}

func (i *initCommand) exec(ctx cli.CommandContext) error {
	// get flag values
	confPath, err := ctx.GetString(CONF_FLAG)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	conf, err := config.OpenConfigFile(confPath)
	if err != nil {
		return err
	}

	cwdfiles, err := os.ReadDir(cwd)
	if err != nil {
		return err
	}

	return config.MaybeUpdateAndSaveConfig(conf, cwdfiles, cwd)
}

func CreateInitCommand(ctx context.Context) *cli.Command {
	ic := &initCommand{
		ctx: ctx,
	}
	return cli.NewCommandWithFunc(
		"init", "Initialize FTUCK", ic.exec,
		cli.RegisterFlag(CONF_FLAG, CONF_DESC, cli.StringFlag, CONF_DEFAULT, "c"),
		cli.RegisterFlag(WD_FLAG, WD_DESC, cli.StringFlag, WD_DEFAULt, "wd"),
	)
}
