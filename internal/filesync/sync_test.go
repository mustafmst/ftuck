package filesync

import "testing"

func TestMaybeCreateAndUpdateSyncFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		conf    syncFileGetter
		src     string
		trg     string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := MaybeCreateAndUpdateSyncFile(tt.conf, tt.src, tt.trg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MaybeCreateAndUpdateSyncFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MaybeCreateAndUpdateSyncFile() succeeded unexpectedly")
			}
		})
	}
}
