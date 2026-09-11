package builder

import (
	"testing"

	"github.com/IridiumNan/project-todo/internal/store"
)

func TestWorkTomlBuilder(t *testing.T) {
	dataDir, err := store.InitDataDirOnCurrentDir()
	if err != nil {
		t.Errorf("error when init data dir, err: %s", err.Error())
	}

	wb := NewWorkTomlBuilder(dataDir)

	err = wb.Build()
	if err != nil {
		t.Errorf("error when build a new work, err: %s", err.Error())
	}
}
