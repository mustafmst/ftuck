package filesync

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// func TestMaybeCreateAndUpdateSyncFile(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for target function.
// 		conf    syncFileGetter
// 		src     string
// 		trg     string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotErr := MaybeCreateAndUpdateSyncFile(tt.conf, tt.src, tt.trg)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("MaybeCreateAndUpdateSyncFile() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("MaybeCreateAndUpdateSyncFile() succeeded unexpectedly")
// 			}
// 		})
// 	}
// }

type confMock struct {
	srcPath string
}

// GetSyncFilePath implements syncFileGetter.
func (c *confMock) GetSyncFilePath() string {
	return c.srcPath
}

func TestSchema_SyncAllEntries(t *testing.T) {
	// os.CreateTemp(dir string, pattern string)
	tmpDir := os.TempDir()
	const testDir = "test_sync_all"
	srcPath := path.Join(tmpDir, testDir, "src")
	destPath := path.Join(tmpDir, testDir, "dest")

	// Create src and dest dirs for test
	_ = os.MkdirAll(srcPath, 0x666)
	_ = os.MkdirAll(destPath, 0x666)

	srcF1Name := path.Join(srcPath, "file1")
	srcF1, _ := os.Create(srcF1Name)
	defer srcF1.Close()

	// Create all files needed for tests

	defer func(dirToClean string) {
		// Test cleanup
		os.RemoveAll(dirToClean)
	}(tmpDir)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		data *Schema
		// Named input parameters for target function.
		conf      syncFileGetter
		wantErr   bool
		checkFunc func() error
	}{
		{
			name: "non abs source is changed to abs",
			data: &Schema{
				{
					Source:      "file1",
					Destination: path.Join(destPath, "dFile1"),
				},
			},
			conf:    &confMock{srcPath},
			wantErr: false,
			checkFunc: func() error {
				f, err := os.Lstat(path.Join(destPath, "dFile1"))
				if err != nil {
					return err
				}
				if f.Mode()&os.ModeSymlink == 1 {
					return fmt.Errorf("File exist but is not a symlink (name: %s)", f.Name())
				}
				if resolvedLink, err := filepath.EvalSymlinks(f.Name()); err != nil {
					return err
				} else if resolvedLink != srcF1Name {
					return fmt.Errorf("link points to wrong file (link:%s, target:%s)", f.Name(), resolvedLink)
				}

				return nil
			},
		},
		{
			name: "abs path is not changed",
			data: &Schema{
				{
					Source:      srcF1Name,
					Destination: path.Join(destPath, "dFile1abs"),
				},
			},
			conf:    &confMock{""},
			wantErr: false,
			checkFunc: func() error {
				f, err := os.Lstat(path.Join(destPath, "dFile1abs"))
				if err != nil {
					return err
				}
				if f.Mode()&os.ModeSymlink == 1 {
					return fmt.Errorf("File exist but is not a symlink (name: %s)", f.Name())
				}
				if resolvedLink, err := filepath.EvalSymlinks(f.Name()); err != nil {
					return err
				} else if resolvedLink != srcF1Name {
					return fmt.Errorf("link points to wrong file (link:%s, target:%s)", f.Name(), resolvedLink)
				}

				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.data.SyncAllEntries(tt.conf)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SyncAllEntries() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SyncAllEntries() succeeded unexpectedly")
			}
			err := tt.checkFunc()
			if err != nil {
				t.Fatal("check failed", err)
			}
		})
	}
}
