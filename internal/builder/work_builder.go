// Package builder provide the work builder
// work builder is used for build struct [models.Work]
package builder

import (
	_ "embed"
)

type WorkBuilder struct {
	// WorkDir is the /path/to/.project-todo which store the context and data.toml

	DataDir string
}
