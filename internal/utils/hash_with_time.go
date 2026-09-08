// Package utils provide utils function
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// HashByTimestamp receive the timestamp (always [time.Now().UnixNano]) and set length of this hashed string
func HashByTimestamp(timestamp int64, length int) string {
	data := []byte(strconv.FormatInt(timestamp, 10))

	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:length])
}
