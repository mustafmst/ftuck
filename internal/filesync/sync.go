package filesync

import (
	"fmt"
	"os"
	"path"
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
	}

	// read sync definitions
	d, err := ReadOrCreate(syncFile)
	if err != nil {
		return fmt.Errorf("%w : run init again", err)
	}

	s, err := ReadSchema(d)
	if err != nil {
		return err
	}

	// add new definition
	s.Append(SyncDefinition{
		Source: src,
		Target: trg,
	})

	return s.WriteToFile(syncFile)
}
