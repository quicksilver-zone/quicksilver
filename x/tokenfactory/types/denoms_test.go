package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	cmdcfg "github.com/ingenuity-build/quicksilver/cmd/config"
	"github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

func TestDeconstructDenom(t *testing.T) {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)

	for _, tc := range []struct {
		desc             string
		denom            string
		expectedSubdenom string
		err              error
	}{
		{
			desc:  "empty is invalid",
			denom: "",
			err:   types.ErrInvalidDenom,
		},
		{
			desc:             "normal",
			denom:            "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/bitcoin",
			expectedSubdenom: "bitcoin",
		},
		{
			desc:             "multiple slashes in subdenom",
			denom:            "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/bitcoin/1",
			expectedSubdenom: "bitcoin/1",
		},
		{
			desc:             "no subdenom",
			denom:            "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/",
			expectedSubdenom: "",
		},
		{
			desc:  "incorrect prefix",
			denom: "ibc/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/bitcoin",
			err:   types.ErrInvalidDenom,
		},
		{
			desc:             "subdenom of only slashes",
			denom:            "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/////",
			expectedSubdenom: "////",
		},
		{
			desc:  "too long name",
			denom: "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			err:   types.ErrInvalidDenom,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			expectedCreator := "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf"
			creator, subdenom, err := types.DeconstructDenom(tc.denom)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, expectedCreator, creator)
				require.Equal(t, tc.expectedSubdenom, subdenom)
			}
		})
	}
}

func TestGetTokenDenom(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		creator  string
		subdenom string
		valid    bool
	}{
		{
			desc:     "normal",
			creator:  "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
			subdenom: "bitcoin",
			valid:    true,
		},
		{
			desc:     "multiple slashes in subdenom",
			creator:  "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
			subdenom: "bitcoin/1",
			valid:    true,
		},
		{
			desc:     "no subdenom",
			creator:  "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
			subdenom: "",
			valid:    true,
		},
		{
			desc:     "subdenom of only slashes",
			creator:  "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
			subdenom: "/////",
			valid:    true,
		},
		{
			desc:     "too long name",
			creator:  "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppfffffffffffffffffffffffffffffffffff",
			subdenom: ".",
			valid:    false,
		},
		{
			desc:     "subdenom is exactly max length",
			creator:  "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
			subdenom: "bitcoinfsadfsdfeadfsafwefsefsefsdfsdafasefsf",
			valid:    true,
		},
		{
			desc:     "creator is exactly max length",
			creator:  "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppfkhanhanchaucascascascas",
			subdenom: "bitcoin",
			valid:    true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := types.GetTokenDenom(tc.creator, tc.subdenom)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
