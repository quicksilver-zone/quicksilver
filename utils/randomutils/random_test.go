package randomutils_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils/randomutils"
)

func TestGenerateRandomHash(t *testing.T) {
	var byteslice []byte = randomutils.GenerateRandomHash()
	require.Equal(t, 32, len(byteslice))
}

func TestGenerateRandomHashAsHex(t *testing.T) {
	var hexHash string = randomutils.GenerateRandomHashAsHex()
	require.Equal(t, 64, len(hexHash))
	byteslice, err := hex.DecodeString(hexHash)
	require.NoError(t, err)
	require.Equal(t, 32, len(byteslice))
}
