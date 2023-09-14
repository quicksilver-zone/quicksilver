package utils_test

import (
	"testing"

	utils "github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const expectedDenom = "denom"

func TestDenomFromRequestKey(t *testing.T) {
	cases := []struct {
		name string
		fn   func() (sdk.AccAddress, []byte)
		err  string
	}{
		{
			"valid",
			func() (sdk.AccAddress, []byte) {
				accAddr := addressutils.GenerateAccAddressForTest()
				key := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
				key = append(key, []byte(expectedDenom)...)
				return accAddr, key
			},
			"",
		},
		{
			"invalid - address mismatch",
			func() (sdk.AccAddress, []byte) {
				keyAddr, err := addressutils.AccAddressFromBech32("cosmos135rd8ft0dyq8fv3w3hhmaa55qu3pe668j99qh67mg747ew4ad03qsgq8vh", "cosmos")
				require.NoError(t, err)
				checkAddr, err := addressutils.AccAddressFromBech32("cosmos1ent5eg0xn3pskf3fhdw8mky88ry7t4kx628ru3pzp4nqjp6eufusphlldy", "cosmos")
				require.NoError(t, err)
				key := banktypes.CreateAccountBalancesPrefix(keyAddr.Bytes())
				key = append(key, []byte(expectedDenom)...)
				return checkAddr, key
			},
			"account mismatch; expected cosmos135rd8ft0dyq8fv3w3hhmaa55qu3pe668j99qh67mg747ew4ad03qsgq8vh, got cosmos1ent5eg0xn3pskf3fhdw8mky88ry7t4kx628ru3pzp4nqjp6eufusphlldy",
		},
		{
			"invalid - empty address",
			func() (sdk.AccAddress, []byte) {
				accAddr := sdk.AccAddress{}
				key := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
				key = append(key, []byte(expectedDenom)...)
				return accAddr, key
			},
			"invalid key",
		},
		{
			"invalid - empty denom",
			func() (sdk.AccAddress, []byte) {
				accAddr := addressutils.GenerateAccAddressForTest()
				key := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
				key = append(key, []byte("")...)
				return accAddr, key
			},
			"key contained no denom",
		},
	}

	for _, c := range cases {
		address, key := c.fn()
		actualDenom, err := utils.DenomFromRequestKey(key, address)
		if c.err == "" {
			require.NoError(t, err)
			require.Equal(t, expectedDenom, actualDenom)
		} else {
			require.Errorf(t, err, c.err)
		}
	}
}
