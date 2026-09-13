package cache

import (
	"os"
	"slices"
	"testing"
)

var testNewDirs = []string{
	"/tmp/.project-todo",
	"/tmp/test/.project-todo",

	// This fake dir will not created, it should not appear on cache file
	"/tmp/fake/.project-todo",
}

var expectedDirs = []string{
	"/tmp/.project-todo",
	"/tmp/test/.project-todo",
}

func createTestDirs(t *testing.T) {
	for _, dir := range expectedDirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Errorf("error when create test dirs, err: %s", err.Error())
		}
	}
}

func clearTestDirs(t *testing.T) {
	for _, dir := range expectedDirs {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("error when clear test dirs, err: %s", err.Error())
		}
	}

	os.RemoveAll(DefaultDirCacheFilePath)
}

func TestPushNewDir(t *testing.T) {
	createTestDirs(t)

	defer clearTestDirs(t)

	cache, err := DefaultCache()
	if err != nil {
		t.Fatal(err)
	}

	for _, dir := range testNewDirs {
		err := cache.PushNewDir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}

	dirs, err := cache.Dirs()
	if err != nil {
		t.Errorf("error when read cache dirs, err: %s", "err")
	}

	if !slices.Equal(dirs, expectedDirs) {
		t.Errorf("unexpected cache dirs, want: %v, got: %v", expectedDirs, dirs)
	}
}
