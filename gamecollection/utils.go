package gamecollection

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
