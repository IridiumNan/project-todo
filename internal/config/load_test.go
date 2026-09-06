package config

import (
	"fmt"
	"os"
	"path"
	"slices"
	"testing"

	"github.com/IridiumNan/project-todo/internal/models"
)

// WARN: change this after you modify the default configuration file content
var expectedDefaultConfig = &GlobalConfig{
	DefaultViewer: models.ViewerEditor,
	EditorViewCommands: []string{
		"nano %s",
		"vim %s",
	},
}

// WARN: change this function if default configuration file modified
func checkConfig(want *GlobalConfig, got *GlobalConfig) error {
	if got.DefaultViewer != want.DefaultViewer {
		return fmt.Errorf("want %v, got %v", want.DefaultViewer, got.DefaultViewer)
	}

	if slices.Compare(got.EditorViewCommands, got.EditorViewCommands) != 0 {
		return fmt.Errorf("want %v, got %v", want.EditorViewCommands, got.EditorViewCommands)
	}
	return nil
}

func TestDefaultConfig(t *testing.T) {
	got := defaultGlobalConfig()
	want := expectedDefaultConfig

	if err := checkConfig(want, got); err != nil {
		t.Error(err)
	}
}

func clearOriginalConfig(t *testing.T) {
	configFilePath := path.Join(APPConfigDir, ConfigFileName)

	// if dir not exist, this [os.RemoveAll] will return nil (not err)
	err := os.RemoveAll(configFilePath)
	if err != nil {
		t.Errorf("error when remove original dir, target dir: %s, err: %s", configFilePath, err.Error())
	}
}

func TestLoadGlobalConfig(t *testing.T) {
	clearOriginalConfig(t)

	got := loadGlobalConfig()

	want := expectedDefaultConfig

	if err := checkConfig(want, got); err != nil {
		t.Error(err)
	}

	// check if configuration file automatically create and write the template content

	configFilePath := path.Join(APPConfigDir, ConfigFileName)
	b, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Errorf("error when read configuration file, file path: %s, err: %s", configFilePath, err.Error())
	}

	if string(b) != DefaultGlobalConfigTemplate {
		t.Errorf("want: %s, got: %s", DefaultGlobalConfigTemplate, string(b))
	}
}
