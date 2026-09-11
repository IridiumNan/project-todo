package models

// MDTomlConfig struct is used for UnMarshal raw markdown content
// This config define for build a new work
//
// it's different from the global configuration
//
// This toml config is parsed by from a tmp file then used for building a new work
type MDTomlConfig struct {
	Title string `toml:"title"`

	Energy Energy `toml:"energy"`

	Viewer ViewerType `toml:"viewer"`

	DependenciesID []string `toml:"dependencies_id"`
}
