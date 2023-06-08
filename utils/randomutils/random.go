package randomutils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will panic if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic("unable to read random bytes")
	}

	return b
}

func GenerateRandomHashAsHex(n int) string {
	return hex.EncodeToString(GenerateRandomBytes(n))
}
