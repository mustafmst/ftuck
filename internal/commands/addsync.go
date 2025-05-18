package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/mustafmst/ftuck/internal/cli"
	"github.com/mustafmst/ftuck/internal/config"
	"github.com/mustafmst/ftuck/internal/filesync"
)

var (
	ErrNotInit        error = errors.New("config was not initiated")
	ErrObligatoryFlag error = errors.New("obligatory flag was not provided")
)

const ADDSYNK_DESC string = "Add new file sync to current config"

// FLAGS
const (
	TARGET_FLAG string = "target"
	SOURCE_FLAG string = "source"
)

// DEFAULTS
const (
	SRC_TRG_DEFAULT_VALUE string = "path not given"
)

// DESCRIPTIONS
const (
	TARGET_DESC string = "(Obligatory) This flag specifies where FTUCK will create symlink for a file"
	SOURCE_DESC string = "(Obligatory) This flag specifies what is the source of created symlink"
)

type addSyncCommand struct {
	ctx context.Context
}

func (as *addSyncCommand) exec(ctx cli.CommandContext) error {
	// get configuration path
	confPath, err := ctx.GetString(CONF_FLAG)
	if err != nil {
		return err
	}

	// get target and handle errors
	trg, err := ctx.GetString(TARGET_FLAG)
	if err != nil {
		return err
	}
	if trg == SRC_TRG_DEFAULT_VALUE {
		return fmt.Errorf("(flag = %s) %w", TARGET_FLAG, ErrObligatoryFlag)
	}

	// get source and handle errors
	src, err := ctx.GetString(SOURCE_FLAG)
	if err != nil {
		return err
	}
	if src == SRC_TRG_DEFAULT_VALUE {
		return fmt.Errorf("(flag = %s) %w", SOURCE_FLAG, ErrObligatoryFlag)
	}

	// read configuration
	conf, err := config.OpenConfigFile(confPath)
	if err != nil {
		return err
	}

	return filesync.MaybeCreateAndUpdateSyncFile(&conf.Config, src, trg)
}

func CreateAddSyncCommand(ctx context.Context) *cli.Command {
	as := &addSyncCommand{
		ctx: ctx,
	}
	return cli.NewCommandWithFunc("addsync", ADDSYNK_DESC, as.exec,
		cli.RegisterFlag(
			TARGET_FLAG,
			TARGET_DESC,
			cli.StringFlag,
			SRC_TRG_DEFAULT_VALUE, "t",
		),
		cli.RegisterFlag(
			SOURCE_FLAG,
			SOURCE_DESC,
			cli.StringFlag,
			SRC_TRG_DEFAULT_VALUE, "s",
		),
		cli.RegisterFlag(
			CONF_FLAG,
			CONF_DESC,
			cli.StringFlag,
			CONF_DEFAULT,
			"c",
		),
	)
}
