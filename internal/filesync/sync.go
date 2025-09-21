package filesync

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
)

type syncFileGetter interface {
	GetSyncFilePath() string
}

func MaybeCreateAndUpdateSyncFile(conf syncFileGetter, src string, trg string) error {
	syncFile := conf.GetSyncFilePath()
	// return error if sync file not set
	if syncFile == "" {
		cwd, _ := os.Getwd()
		syncFile = path.Join(cwd, SYNC_FILE_NAME)
		slog.Debug("sync file not configured, using default", "sync_file", syncFile)
	} else {
		slog.Debug("using configured sync file", "sync_file", syncFile)
	}

	// read sync definitions
	slog.Debug("reading sync file", "path", syncFile)
	d, err := ReadOrCreate(syncFile)
	if err != nil {
		return fmt.Errorf("%w : run init again", err)
	}

	s, err := ReadSchema(d)
	if err != nil {
		return err
	}

	// add new definition
	slog.Info("adding sync definition", "source", src, "destination", trg)
	s.Append(SyncDefinition{
		Source:      src,
		Destination: trg,
	})

	return s.WriteToFile(syncFile)
}

func (s *Schema) SyncAllEntries(conf syncFileGetter) error {
	return s.ForEach(func(sd SyncDefinition) error {
		source := sd.Source
		if !filepath.IsAbs(source) {
			source = filepath.Join(conf.GetSyncFilePath(), source)
		}
		fi, err := os.Lstat(sd.Destination)
		if err != nil && !os.IsNotExist(err) {
			slog.Error("syncing", "error", err, "target", sd.Destination)
			return err
		}
		// Creating link if it does not exist
		if err != nil && os.IsNotExist(err) {
			slog.Info("creating link", "source", source, "target", sd.Destination)
			os.Symlink(source, sd.Destination)
			return nil
		}
		// ommiting if file exists
		if fi.Mode()&os.ModeSymlink == 1 {
			err := fmt.Errorf("(path = %s) file exists", sd.Destination)
			slog.Error("target file already exists and is not a Symlink", "error", err)
			return nil
		}
		existingSource, err := os.Readlink(sd.Destination)
		if err != nil {
			slog.Error("reading link", "error", err)
			return err
		}

		// updating if source is different
		if existingSource != source {
			slog.Error("link different", "source", source, "link", existingSource)
			err := os.Remove(sd.Destination)
			if err != nil {
				return err
			}
			err = os.Symlink(source, sd.Destination)
			if err != nil {
				return err
			}
			return nil
		}

		// nothing to be done here
		slog.Info("nothing to do", "target", sd.Destination)
		return nil
	})
}
