package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func HashByTimestamp(timestamp int64, length int) string {
	data := []byte(strconv.FormatInt(timestamp, 10))

	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:length])
}
