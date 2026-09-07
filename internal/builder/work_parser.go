package builder

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// generateTimestampID return a hashed string which generated from timestamp and length = 8
func generateTimestampID() string {
	now := time.Now().UnixNano()

	// Convert integer timestamp to string bytes
	data := []byte(strconv.FormatInt(now, 10))

	// Hash using SHA-256
	hash := sha256.Sum256(data)

	// Encode hash as a hexadecimal string
	return hex.EncodeToString(hash[:8])
}
