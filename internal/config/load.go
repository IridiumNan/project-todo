package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/project-todo/internal/utils"
)

// DefaultGlobalConfig confirm this will not encounter error before call it
// NOTE: Must make sure this function work properly
func defaultGlobalConfig() *GlobalConfig {
	var config GlobalConfig

	err := toml.Unmarshal([]byte(DefaultGlobalConfigTemplate), &config)
	if err != nil {
		slog.Error("while unmarshal DefaultGlobalConfigTemplate", "err", err, "template", DefaultGlobalConfigTemplate)
		log.Fatal("config load failed, exit... function name: DefaultGlobalConfig")
	}

	return &config
}

// LoadGlobalConfig the main function for configuration setting for now
// It will not create a isolated configuration for single project
func loadGlobalConfig() *GlobalConfig {
	configFilePath := path.Join(APPConfigDir, ConfigFileName)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		slog.Warn("global configuration file not found, init a new one", "config_path", configFilePath)

		if err = initGlobalConfig(configFilePath); err != nil {
			slog.Error("while create a global config file, use default configuration value", "err", err)

			return defaultGlobalConfig()
		}
	}

	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		slog.Error("while reading configuration file", "config file path", configFilePath, "err", err)

		return defaultGlobalConfig()
	}

	var config GlobalConfig
	err = toml.Unmarshal(configFile, &config)
	if err != nil {
		return defaultGlobalConfig()
	}

	return &config
}

// initGlobalConfig Write the [DefaultGlobalConfigTemplate] into default configuration file path
// It enable user to customize their configuration
func initGlobalConfig(configFilePath string) error {
	err := utils.EnsureFileExist(configFilePath)
	if err != nil {
		return fmt.Errorf("func: initGlobalConfig, error when ensure file exist, err: %w", err)
	}
	if err := os.WriteFile(configFilePath, []byte(DefaultGlobalConfigTemplate), 0o644); err != nil {
		return fmt.Errorf("func: initGlobalConfig, err: %w", err)
	}
	return nil
}
