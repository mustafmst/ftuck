package commands

import (
	"context"
	"fmt"

	"github.com/mustafmst/ftuck/internal/cli"
	"github.com/mustafmst/ftuck/internal/config"
	"github.com/mustafmst/ftuck/internal/filesync"
)

type syncAllCommand struct {
	ctx context.Context
}

func (sa *syncAllCommand) exec(ctx cli.CommandContext) error {
	// get flag values
	confPath, err := ctx.GetString(CONF_FLAG)
	if err != nil {
		return err
	}

	conf, err := config.OpenConfigFile(confPath)
	if err != nil {
		return err
	}

	if conf.Config.SyncFile == "" {
		return ErrNotInit
	}

	// read sync definitions
	d, err := filesync.ReadOrCreate(conf.Config.SyncFile)
	if err != nil {
		return fmt.Errorf("%w : run init again", err)
	}

	s, err := filesync.ReadSchema(d)
	if err != nil {
		return err
	}
	return s.SyncAllEntries()
}

func CreateSyncAllCommand(ctx context.Context) *cli.Command {
	sa := &syncAllCommand{
		ctx: ctx,
	}
	return cli.NewCommandWithFunc(
		"sync",
		"Sync files with current configuration",
		sa.exec,
		cli.RegisterFlag(CONF_FLAG, CONF_DESC, cli.StringFlag, CONF_DEFAULT, "c"),
	)
}
