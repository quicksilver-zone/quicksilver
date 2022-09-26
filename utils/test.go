package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomHash() []byte {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("unable to read random bytes")
	}
	return bytes
}

func GenerateRandomHashAsHex() string {
	return hex.EncodeToString(GenerateRandomHash())
}
