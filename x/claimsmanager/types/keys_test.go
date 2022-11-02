package types

import (
	"testing"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	address := utils.GenerateAccAddressForTest()
	keyClaim := GetKeyClaim("testzone-1", address.String(), ClaimTypeOsmosisPool, "testzone-2")
	prefixClaim := GetPrefixClaim("testzone-1")

	expected := append(
		[]byte{0x00},
		append(
			[]byte("testzone-1"),
			append(
				[]byte{0x00},
				append(
					[]byte(address.String()),
					append(
						[]byte{0x00, 0x00, 0x00, 0x02},
						[]byte("testzone-2")...,
					)...,
				)...,
			)...,
		)...,
	)

	require.Equal(t, append([]byte{0x00}, []byte("testzone-1")...), prefixClaim)
	require.Equal(t, expected, keyClaim)
}

func TestLastEpochKeys(t *testing.T) {
	address := utils.GenerateAccAddressForTest()
	keyClaim := GetKeyLastEpochClaim("testzone-1", address.String(), ClaimTypeOsmosisPool, "testzone-2")
	prefixClaim := GetPrefixLastEpochClaim("testzone-1")

	expected := append(
		[]byte{0x01},
		append(
			[]byte("testzone-1"),
			append(
				[]byte{0x00},
				append(
					[]byte(address.String()),
					append(
						[]byte{0x00, 0x00, 0x00, 0x02},
						[]byte("testzone-2")...,
					)...,
				)...,
			)...,
		)...,
	)

	require.Equal(t, append([]byte{0x01}, []byte("testzone-1")...), prefixClaim)
	require.Equal(t, expected, keyClaim)
}
