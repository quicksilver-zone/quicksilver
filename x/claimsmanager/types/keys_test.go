package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
<<<<<<< HEAD
=======
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
>>>>>>> origin/develop
)

func TestKeys(t *testing.T) {
	address := utils.GenerateAccAddressForTest()

	// zone
	prefixClaim := types.GetPrefixClaim("testzone-1")
	require.Equal(t, append(types.KeyPrefixClaim, []byte("testzone-1")...), prefixClaim)

	expected := types.KeyPrefixClaim
	expected = append(expected, []byte("testzone-1")...)
	expected = append(expected, byte(0x00))
	expected = append(expected, []byte(address.String())...)

	// zone + user
	prefixUserClaim := types.GetPrefixUserClaim("testzone-1", address.String())
	require.Equal(t, expected, prefixUserClaim)

	expected = append(expected, []byte{0x00, 0x00, 0x00, 0x02}...)
	expected = append(expected, []byte("testzone-2")...)

	// zone + user + claimType + srcZone
	keyClaim := types.GetKeyClaim("testzone-1", address.String(), types.ClaimTypeOsmosisPool, "testzone-2")
	require.Equal(t, expected, keyClaim)
}

func TestLastEpochKeys(t *testing.T) {
	address := utils.GenerateAccAddressForTest()

	// zone
	prefixClaim := types.GetPrefixLastEpochClaim("testzone-1")
	require.Equal(t, append(types.KeyPrefixLastEpochClaim, []byte("testzone-1")...), prefixClaim)

	expected := types.KeyPrefixLastEpochClaim
	expected = append(expected, []byte("testzone-1")...)
	expected = append(expected, byte(0x00))
	expected = append(expected, []byte(address.String())...)

	// zone + user
	prefixUserClaim := types.GetPrefixLastEpochUserClaim("testzone-1", address.String())
	require.Equal(t, expected, prefixUserClaim)

	expected = append(expected, []byte{0x00, 0x00, 0x00, 0x02}...)
	expected = append(expected, []byte("testzone-2")...)

	// zone + user + claimType + srcZone
	keyClaim := types.GetKeyLastEpochClaim("testzone-1", address.String(), types.ClaimTypeOsmosisPool, "testzone-2")
	require.Equal(t, expected, keyClaim)
}
