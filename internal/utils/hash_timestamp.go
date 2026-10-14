package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

const defautHashLength = 10

func HashByTimestampWithLen(timestamp int64, length int) string {
	data := []byte(strconv.FormatInt(timestamp, 10))

	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:length])
}

// HashByTimestamp receive the timestamp (always [time.Now().UnixNano]) and set length of this hashed string
// with [defautHashLength]
func HashByTimestamp(timestamp int64) string {
	return HashByTimestampWithLen(timestamp, defautHashLength)
}
