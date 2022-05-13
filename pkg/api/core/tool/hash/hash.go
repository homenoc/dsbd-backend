package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func Generate(data string) string {
	hash := sha256.Sum256([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func Verify(data string, hashedData string) bool {
	hashedInputByte := sha256.Sum256([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hashedInputByte[:])) == strings.ToUpper(hashedData)
}
