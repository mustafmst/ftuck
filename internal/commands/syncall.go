package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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

	return s.ForEach(func(sd filesync.SyncDefinition) error {
		// TODO: Move logic to a seperate sync handler.
		// Command definitions should only include getting proper
		// data from command context and passing them to handler.
		fi, err := os.Lstat(sd.Target)
		if err != nil && !os.IsNotExist(err) {
			slog.Error("syncing", "error", err, "target", sd.Target)
			return err
		}
		// Creating link if it does not exist
		if err != nil && os.IsNotExist(err) {
			slog.Info("creating link", "source", sd.Source, "target", sd.Target)
			os.Symlink(sd.Source, sd.Target)
			return nil
		}
		// ommiting if file exists
		if fi.Mode()&os.ModeSymlink == 1 {
			err := fmt.Errorf("(path = %s) file exists", sd.Target)
			slog.Error("target file already exists and is not a Symlink", "error", err)
			return nil
		}
		existingSource, err := os.Readlink(sd.Target)
		if err != nil {
			slog.Error("reading link", "error", err)
			return err
		}

		// updating if source is different
		if existingSource != sd.Source {
			slog.Error("link different", "source", sd.Source, "link", existingSource)
			err := os.Remove(sd.Target)
			if err != nil {
				return err
			}
			err = os.Symlink(sd.Source, sd.Target)
			if err != nil {
				return err
			}
			return nil
		}

		// nothing to be done here
		slog.Info("nothing to do", "target", sd.Target)
		return nil
	})
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
