package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/mustafmst/ftuck/internal/cli"
	"github.com/mustafmst/ftuck/internal/config"
	"github.com/mustafmst/ftuck/internal/filesync"
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

	slog.Info("initialing FTUCK", "config path", confPath, "working dir", cwd)

	conf, err := config.OpenConfigFile(confPath)
	if err != nil {
		return err
	}
	defer conf.Close()

	cwdfiles, err := os.ReadDir(cwd)
	if err != nil {
		return err
	}

	for _, de := range cwdfiles {
		if de.IsDir() {
			continue
		}
		if de.Name() == filesync.SYNC_FILE_NAME {
			conf.Config.SyncFile = path.Join(cwd, de.Name())
			return nil
		}
	}
	return fmt.Errorf("no file syncfile (named=%s) found in current working directory (%s)", filesync.SYNC_FILE_NAME, cwd)
}

func CreateInitCommand(ctx context.Context) *cli.Command {
	ic := &initCommand{
		ctx: ctx,
	}
	return cli.NewCommandWithFunc(
		"init", "Initialize FTUCK", ic.exec,
		cli.RegisterFlag(CONF_FLAG, CONF_DESC, cli.StringArg, CONF_DEFAULT, "c"),
		cli.RegisterFlag(WD_FLAG, WD_DESC, cli.StringArg, WD_DEFAULt, "wd"),
	)
}
