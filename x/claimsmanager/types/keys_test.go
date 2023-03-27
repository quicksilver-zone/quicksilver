package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
)

func TestKeys(t *testing.T) {
	address := utils.GenerateAccAddressForTest()

	// zone
	prefixClaim := GetPrefixClaim("testzone-1")
	require.Equal(t, append(KeyPrefixClaim, []byte("testzone-1")...), prefixClaim)

	expected := KeyPrefixClaim
	expected = append(expected, []byte("testzone-1")...)
	expected = append(expected, byte(0x00))
	expected = append(expected, []byte(address.String())...)

	// zone + user
	prefixUserClaim := GetPrefixUserClaim("testzone-1", address.String())
	require.Equal(t, expected, prefixUserClaim)

	expected = append(expected, []byte{0x00, 0x00, 0x00, 0x02}...)
	expected = append(expected, []byte("testzone-2")...)

	// zone + user + claimType + srcZone
	keyClaim := GetKeyClaim("testzone-1", address.String(), ClaimTypeOsmosisPool, "testzone-2")
	require.Equal(t, expected, keyClaim)
}

func TestLastEpochKeys(t *testing.T) {
	address := utils.GenerateAccAddressForTest()

	// zone
	prefixClaim := GetPrefixLastEpochClaim("testzone-1")
	require.Equal(t, append(KeyPrefixLastEpochClaim, []byte("testzone-1")...), prefixClaim)

	expected := KeyPrefixLastEpochClaim
	expected = append(expected, []byte("testzone-1")...)
	expected = append(expected, byte(0x00))
	expected = append(expected, []byte(address.String())...)

	// zone + user
	prefixUserClaim := GetPrefixLastEpochUserClaim("testzone-1", address.String())
	require.Equal(t, expected, prefixUserClaim)

	expected = append(expected, []byte{0x00, 0x00, 0x00, 0x02}...)
	expected = append(expected, []byte("testzone-2")...)

	// zone + user + claimType + srcZone
	keyClaim := GetKeyLastEpochClaim("testzone-1", address.String(), ClaimTypeOsmosisPool, "testzone-2")
	require.Equal(t, expected, keyClaim)
}
