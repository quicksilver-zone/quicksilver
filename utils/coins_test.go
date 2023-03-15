package utils_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	utils "github.com/ingenuity-build/quicksilver/utils"
)

func TestDenomFromRequestKey(t *testing.T) {
	cases := []struct {
		name string
		fn   func() (sdk.AccAddress, string, []byte)
		err  string
	}{
		{
			"valid",
			func() (sdk.AccAddress, string, []byte) {
				accAddr := utils.GenerateAccAddressForTest()
				prefix := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
				key := append(prefix, []byte("denom")...)
				return accAddr, "denom", key
			},
			"",
		},
		{
			"invalid - address mismatch",
			func() (sdk.AccAddress, string, []byte) {
				keyAddr, err := utils.AccAddressFromBech32("cosmos135rd8ft0dyq8fv3w3hhmaa55qu3pe668j99qh67mg747ew4ad03qsgq8vh", "cosmos")
				require.NoError(t, err)
				checkAddr, err := utils.AccAddressFromBech32("cosmos1ent5eg0xn3pskf3fhdw8mky88ry7t4kx628ru3pzp4nqjp6eufusphlldy", "cosmos")
				require.NoError(t, err)
				prefix := banktypes.CreateAccountBalancesPrefix(keyAddr.Bytes())
				key := append(prefix, []byte("denom")...)
				return checkAddr, "denom", key
			},
			"account mismatch; expected cosmos135rd8ft0dyq8fv3w3hhmaa55qu3pe668j99qh67mg747ew4ad03qsgq8vh, got cosmos1ent5eg0xn3pskf3fhdw8mky88ry7t4kx628ru3pzp4nqjp6eufusphlldy",
		},
		{
			"invalid - empty address",
			func() (sdk.AccAddress, string, []byte) {
				accAddr := sdk.AccAddress{}
				prefix := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
				key := append(prefix, []byte("denom")...)
				return accAddr, "denom", key
			},
			"invalid key",
		},
		{
			"invalid - empty denom",
			func() (sdk.AccAddress, string, []byte) {
				accAddr := utils.GenerateAccAddressForTest()
				prefix := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
				key := append(prefix, []byte("")...)
				return accAddr, "", key
			},
			"key contained no denom",
		},
	}

	for _, c := range cases {
		address, expectedDenom, key := c.fn()
		actualDenom, error := utils.DenomFromRequestKey(key, address)
		if len(c.err) == 0 {
			require.NoError(t, error)
			require.Equal(t, expectedDenom, actualDenom)
		} else {
			require.Errorf(t, error, c.err)
		}
	}
}
