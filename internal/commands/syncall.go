package commands

import (
	"context"

	"github.com/mustafmst/ftuck/internal/cli"
)

type syncAllCommand struct {
	ctx context.Context
}

func (sa *syncAllCommand) exec(ctx cli.CommandContext) error {
	panic("not implemented")
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
