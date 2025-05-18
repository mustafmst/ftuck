package config

import (
	"fmt"
	"os"
	"path"

	"github.com/mustafmst/ftuck/internal/filesync"
)

func MaybeUpdateAndSaveConfig(conf *ConfigFile, cwdfiles []os.DirEntry, cwd string) error {
	defer conf.Save()
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
