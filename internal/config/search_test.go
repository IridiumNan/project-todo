package config

import (
	"os"
	"path"
	"testing"

	"github.com/IridiumNan/project-todo/internal/models"
)

var testDir = "/tmp/go-test"

// prepareTestDir create test dir for data dir search
func prepareTestDir(t *testing.T) {
	targetDir := path.Join(testDir, models.DataDirName)

	err := os.MkdirAll(testDir, 0o755)
	if err != nil {
		t.Errorf("fail to create test dir: %s, err: %s", testDir, err.Error())
	}

	err = os.MkdirAll(targetDir, 0o755)
	if err != nil {
		t.Errorf("fail to create test dir: %s, err: %s", testDir, err.Error())
	}
}

func clearTestDir(t *testing.T) {
	err := os.RemoveAll(testDir)
	if err != nil {
		t.Errorf("error when remove the test dir: %s", err.Error())
	}
}

// TestSearchOnCurrentDir this function just test on linux platform which needs the /tmp dir
func TestSearchOnCurrentDir(t *testing.T) {
	prepareTestDir(t)

	err := os.Chdir(testDir)
	if err != nil {
		t.Errorf("error when change dir, target dir: %s, err: %s", testDir, err.Error())
	}

	foundPath, err := searchDataDir()
	if err != nil {
		t.Errorf("error when search path: %s", err.Error())
	}

	if foundPath != testDir {
		t.Errorf("found dir which contains data dir: want %s, got %s", testDir, foundPath)
	}

	clearTestDir(t)
}

func TestSearchOnParentDir(t *testing.T) {
	prepareTestDir(t)

	depthDir := path.Join(testDir, "dep1", "dep2", "dep3", "dep4")

	err := os.MkdirAll(depthDir, 0o755)
	if err != nil {
		t.Errorf("fail to create test dir: %s, err: %s", testDir, err.Error())
	}

	err = os.Chdir(depthDir)
	if err != nil {
		t.Errorf("error when change dir, target dir: %s, err: %s", depthDir, err.Error())
	}

	foundPath, err := searchDataDir()
	if err != nil {
		t.Errorf("error when search with depth 3, err: %s", err.Error())
	}

	if foundPath != testDir {
		t.Errorf("found dir which contains data dir: want %s, got %s", testDir, foundPath)
	}

	clearTestDir(t)
}
