package builder

import (
	"time"

	"github.com/IridiumNan/project-todo/internal/utils"
)

// generateTimestampID return a hashed string which generated from timestamp and length = 8
func generateTimestampID() string {
	now := time.Now().UnixNano()

	return utils.HashByTimestamp(now, 8)
}
