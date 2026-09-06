package config

import (
	"log/slog"
	"path"

	_ "embed"

	"github.com/IridiumNan/project-todo/internal/models"
)

const emptyPath = ""

var (
	APPConfigDir   = getAppConfigDir()
	ConfigFileName = "config.toml"
)

func getXDGConfigDir() string {
	home, err := models.GetHomeDir()
	if err != nil {
		slog.Error("when get user home dir for get xdg_config_dir", "err", err)
		return emptyPath
	}

	return path.Join(home, ".config")
}

func getAppConfigDir() string {
	return path.Join(getXDGConfigDir(), models.AppName)
}

//go:embed global_config.toml
var DefaultGlobalConfigTemplate string

type GlobalConfig struct {
	// DefaultViewer is set on config.toml file globally
	// when init, default = editorViewer
	DefaultViewer models.ViewerType `toml:"default_viewer"`

	// EditorViewCommands contains commands when run editor viewer
	EditorViewCommands []string `toml:"editor_view_commands"`
}
