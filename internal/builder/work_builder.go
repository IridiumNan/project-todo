// Package builder provide the work builder
// work builder is used for build struct [models.Work]
package builder

import (
	_ "embed"
	"time"

	"github.com/IridiumNan/project-todo/internal/utils"
)

// generateTimestampID return a hashed string which generated from timestamp and length = 8
func generateTimestampID() string {
	now := time.Now().UnixNano()

	return utils.HashByTimestamp(now, 8)
}

type WorkBuilder struct {
	// WorkDir is the /path/to/.project-todo which store the context and data.toml

	DataDir string
}
