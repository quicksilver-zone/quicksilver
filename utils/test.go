package utils

import (
	"encoding/hex"
	"math/rand"
)

func GenerateRandomHash() []byte {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return bytes
}

func GenerateRandomHashAsHex() string {
	return hex.EncodeToString(GenerateRandomHash())
}
