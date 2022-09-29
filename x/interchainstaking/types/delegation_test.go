package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestRoundtripDelegationMarshalToUnmarshal(t *testing.T) {
	del1 := types.NewDelegation(
		"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
		"cosmos1ssrxxe4xsls57ehrkswlkhlkcverf0p0fpgyhzqw0hfdqj92ynxsw29r6e",
		sdk.NewCoin("uqck", sdk.NewInt(300)),
	)

	wantDelAddr := (sdk.AccAddress)([]byte{0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc})
	require.Equal(t, wantDelAddr, del1.GetDelegatorAddr(), "mismatch in delegator address")

	wantValAddr := (sdk.ValAddress)([]byte{
		0x84, 0x06, 0x63, 0x66, 0xa6, 0x87, 0xe1, 0x4f, 0x66, 0xe3, 0xb4,
		0x1d, 0xfb, 0x5f, 0xf6, 0xc3, 0x32, 0x34, 0xbc, 0x2f, 0x48, 0x50,
		0x4b, 0x88, 0x0e, 0x7d, 0xd2, 0xd0, 0x48, 0xaa, 0x24, 0xcd,
	})
	require.Equal(t, wantValAddr, del1.GetValidatorAddr(), "mismatch in validator address")

	marshaledDelBytes := types.MustMarshalDelegation(types.ModuleCdc, del1)
	unmarshaledDel := types.MustUnmarshalDelegation(types.ModuleCdc, marshaledDelBytes)
	require.Equal(t, del1, unmarshaledDel, "Roundtripping: marshal->unmarshal should produce the same delegation")

	// Finally ensure that the 2nd round marshaled bytes equal the original ones.
	marshalDelBytes2ndRound := types.MustMarshalDelegation(types.ModuleCdc, unmarshaledDel)
	require.Equal(t, marshaledDelBytes, marshalDelBytes2ndRound, "all the marshaled bytes should be equal!")
}

func TestSetForValoper(t *testing.T) {
	v1 := utils.GenerateValAddressForTest().String()
	v2 := utils.GenerateValAddressForTest().String()
	intents := types.ValidatorIntents{
		{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(10, 1)},
		{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(90, 1)},
	}

	intents = intents.SetForValoper(v1, &types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(40, 1)})
	intents = intents.SetForValoper(v2, &types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(60, 1)})

	require.Equal(t, sdk.NewDecWithPrec(40, 1), intents.MustGetForValoper(v1).Weight)
	require.Equal(t, sdk.NewDecWithPrec(60, 1), intents.MustGetForValoper(v2).Weight)
}
