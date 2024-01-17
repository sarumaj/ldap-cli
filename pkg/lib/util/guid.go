package util

import (
	"crypto/rand"
	"encoding/hex"
)

// NewGUID generates a new GUID
func NewGUID() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	guid := hex.EncodeToString(bytes)

	return guid[:8] + "-" +
		guid[8:12] + "-" +
		guid[12:16] + "-" +
		guid[16:20] + "-" +
		guid[20:]
}
