package randomutils

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

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
