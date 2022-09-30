package types

import (
	"testing"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	address := utils.GenerateAccAddressForTest()
	keyClaim := GetKeyClaim("testzone-1", address.String())
	prefixClaim := GetPrefixClaim("testzone-1")

	require.Equal(t, append([]byte{0x01}, []byte("testzone-1")...), prefixClaim)
	require.Equal(t, append([]byte{0x01}, []byte("testzone-1"+address.String())...), keyClaim)
}
