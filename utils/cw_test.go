package utils_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

func TestDecodeCwNamespacedKeyGood(t *testing.T) {
	test := "03412880801f5ddf20eca469bfb4a748a6c334b0693363b4e0edb916f8bbbcdac40009706f736974696f6e736f736d6f31787878787878787878787878"
	unhex, err := hex.DecodeString(test)
	require.NoError(t, err)
	address, parts, err := utils.DecodeCwNamespacedKey(unhex, 2)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, addressutils.MustEncodeAddressToBech32("osmo", address), "osmo1gy5gpqqlth0jpm9ydxlmff6g5mpnfvrfxd3mfc8dhyt03waumtzqt8exxr")
	require.Equal(t, parts, [][]byte{[]byte("positions"), []byte("osmo1xxxxxxxxxxxx")})
}

func TestDecodeCwNamespacedKeyBadNoPrefix(t *testing.T) {
	test := "412880801f5ddf20eca469bfb4a748a6c334b0693363b4e0edb916f8bbbcdac40009706f736974696f6e736f736d6f31787878787878787878787878"
	unhex, err := hex.DecodeString(test)
	require.NoError(t, err)
	address, parts, err := utils.DecodeCwNamespacedKey(unhex, 2)
	require.ErrorContains(t, err, "invalid prefix")
	require.Nil(t, address)
	require.Nil(t, parts)
}

func TestDecodeCwNamespacedKeyBadNoNullTerminator(t *testing.T) {
	test := "03412880801f5ddf20eca469bfb4a748a6c334b0693363b4e0edb916f8bbbcdac409706f736974696f6e736f736d6f31787878787878787878787878"
	unhex, err := hex.DecodeString(test)
	require.NoError(t, err)
	address, parts, err := utils.DecodeCwNamespacedKey(unhex, 2)
	require.ErrorContains(t, err, "expected null terminator after address")
	require.Nil(t, address)
	require.Nil(t, parts)
}

func TestDecodeCwNamespacedKeyBadKeyTooShort(t *testing.T) {
	test := "03412880801f5ddf20eca469bfb4a74f"
	unhex, err := hex.DecodeString(test)
	require.NoError(t, err)
	address, parts, err := utils.DecodeCwNamespacedKey(unhex, 2)
	require.ErrorContains(t, err, "invalid key length")
	require.Nil(t, address)
	require.Nil(t, parts)
}
