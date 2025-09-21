package config

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/mustafmst/ftuck/internal/filesync"
)

func MaybeUpdateAndSaveConfig(conf *ConfigFile, cwdfiles []os.DirEntry, cwd string) error {
	defer func() {
		if err := conf.Save(); err != nil {
			slog.Error("failed to save config", "error", err)
		}
	}()
	
	slog.Debug("searching for sync file", "sync_file_name", filesync.SYNC_FILE_NAME, "directory", cwd)
	
	for _, de := range cwdfiles {
		if de.IsDir() {
			continue
		}
		if de.Name() == filesync.SYNC_FILE_NAME {
			syncFilePath := path.Join(cwd, de.Name())
			slog.Info("sync file found", "path", syncFilePath)
			conf.Config.SyncFile = syncFilePath
			return nil
		}
	}
	
	return fmt.Errorf("no file syncfile (named=%s) found in current working directory (%s)", filesync.SYNC_FILE_NAME, cwd)
}
