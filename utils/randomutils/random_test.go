package randomutils_test

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/quicksilver-zone/quicksilver/utils/randomutils"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomHash(t *testing.T) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := random.Intn(128)
	byteslice := randomutils.GenerateRandomBytes(b)
	require.Equal(t, b, len(byteslice))
}

func TestGenerateRandomHashAsHex(t *testing.T) {
	hexHash := randomutils.GenerateRandomHashAsHex(32)
	require.Equal(t, 64, len(hexHash))
	byteslice, err := hex.DecodeString(hexHash)
	require.NoError(t, err)
	require.Equal(t, 32, len(byteslice))
}
