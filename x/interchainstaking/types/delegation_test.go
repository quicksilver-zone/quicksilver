package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func TestRoundtripDelegationMarshalToUnmarshal(t *testing.T) {
	del1 := types.NewDelegation(
		"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
		"cosmos1ssrxxe4xsls57ehrkswlkhlkcverf0p0fpgyhzqw0hfdqj92ynxsw29r6e",
		sdk.NewCoin("uqck", sdkmath.NewInt(300)),
	)

	wantDelAddr := sdk.AccAddress([]byte{0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc})
	require.Equal(t, wantDelAddr, del1.GetDelegatorAddr(), "mismatch in delegator address")

	wantValAddr := sdk.ValAddress([]byte{
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

	// ensure error is returned for 0 length
	_, err := types.UnmarshalDelegation(types.ModuleCdc, []byte{})
	require.Error(t, err)
}

func TestSetForValoper(t *testing.T) {
	v1 := addressutils.GenerateValAddressForTest().String()
	v2 := addressutils.GenerateValAddressForTest().String()
	v3 := addressutils.GenerateValAddressForTest().String()

	intents := types.ValidatorIntents{
		{ValoperAddress: v1, Weight: sdkmath.LegacyNewDecWithPrec(10, 1)},
		{ValoperAddress: v2, Weight: sdkmath.LegacyNewDecWithPrec(90, 1)},
	}

	intents = intents.SetForValoper(v1, &types.ValidatorIntent{ValoperAddress: v1, Weight: sdkmath.LegacyNewDecWithPrec(40, 1)})
	intents = intents.SetForValoper(v2, &types.ValidatorIntent{ValoperAddress: v2, Weight: sdkmath.LegacyNewDecWithPrec(60, 1)})

	require.Equal(t, sdkmath.LegacyNewDecWithPrec(40, 1), intents.MustGetForValoper(v1).Weight)
	require.Equal(t, sdkmath.LegacyNewDecWithPrec(60, 1), intents.MustGetForValoper(v2).Weight)

	// check failed return
	actual := intents.MustGetForValoper(v3)
	require.Equal(t, sdkmath.LegacyZeroDec(), actual.Weight)
}

func TestNormalizeValidatorIntentsDeterminism(t *testing.T) {
	v1 := addressutils.GenerateValAddressForTest().String()
	v2 := addressutils.GenerateValAddressForTest().String()
	v3 := addressutils.GenerateValAddressForTest().String()
	v4 := addressutils.GenerateValAddressForTest().String()

	cases := []struct {
		name    string
		intents types.ValidatorIntents
	}{
		{
			name: "case 1",
			intents: types.ValidatorIntents{
				{ValoperAddress: v1, Weight: sdkmath.LegacyNewDecWithPrec(10, 1)},
				{ValoperAddress: v2, Weight: sdkmath.LegacyNewDecWithPrec(90, 1)},
			},
		},
		{
			name: "case 2",
			intents: types.ValidatorIntents{
				{ValoperAddress: v1, Weight: sdkmath.LegacyNewDecWithPrec(10, 1)},
				{ValoperAddress: v2, Weight: sdkmath.LegacyNewDecWithPrec(90, 1)},
				{ValoperAddress: v3, Weight: sdkmath.LegacyNewDecWithPrec(90, 1)},
				{ValoperAddress: v4, Weight: sdkmath.LegacyNewDecWithPrec(90, 1)},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cached := tc.intents.Normalize()
			// verify that sort is deterministic
			for i := 1; i < 3; i++ {
				normalized := tc.intents.Normalize()
				require.Equal(t, cached, normalized)

			}
		})
	}
}

func TestDetermineAllocationsForDelegation(t *testing.T) {
	// we auto generate the validator addresses in these tests. any dust gets allocated to the first validator in the list
	// once sorted alphabetically on valoper.

	vals := addressutils.GenerateValidatorsDeterministic(4)

	tc := []struct {
		current  map[string]sdkmath.Int
		target   types.ValidatorIntents
		inAmount sdk.Coins
		expected map[string]sdkmath.Int
		dust     sdkmath.Int
	}{
		{
			current: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(350000),
				vals[1]: sdkmath.NewInt(650000),
				vals[2]: sdkmath.NewInt(75000),
			},
			target: types.ValidatorIntents{ValidatorIntents: []*types.ValidatorIntent{
				{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDecWithPrec(30, 2)},
				{ValoperAddress: vals[1], Weight: sdkmath.LegacyNewDecWithPrec(63, 2)},
				{ValoperAddress: vals[2], Weight: sdkmath.LegacyNewDecWithPrec(7, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdkmath.NewInt(50000))),
			expected: map[string]sdkmath.Int{
				vals[1]: sdkmath.NewInt(47000),
				vals[2]: sdkmath.NewInt(3000),
			},
		},
		{
			current: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(52),
				vals[1]: sdkmath.NewInt(24),
				vals[2]: sdkmath.NewInt(20),
				vals[3]: sdkmath.NewInt(4),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDecWithPrec(50, 2)},
				{ValoperAddress: vals[1], Weight: sdkmath.LegacyNewDecWithPrec(25, 2)},
				{ValoperAddress: vals[2], Weight: sdkmath.LegacyNewDecWithPrec(15, 2)},
				{ValoperAddress: vals[3], Weight: sdkmath.LegacyNewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdkmath.NewInt(20))),
			expected: map[string]sdkmath.Int{
				vals[3]: sdkmath.NewInt(7),
				vals[1]: sdkmath.NewInt(5),
				vals[0]: sdkmath.NewInt(8),
			},
		},
		{
			current: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(52),
				vals[1]: sdkmath.NewInt(24),
				vals[2]: sdkmath.NewInt(20),
				vals[3]: sdkmath.NewInt(3),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: vals[0], Weight: sdkmath.LegacyNewDecWithPrec(50, 2)},
				{ValoperAddress: vals[1], Weight: sdkmath.LegacyNewDecWithPrec(25, 2)},
				{ValoperAddress: vals[2], Weight: sdkmath.LegacyNewDecWithPrec(15, 2)},
				{ValoperAddress: vals[3], Weight: sdkmath.LegacyNewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdkmath.NewInt(50))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(27),
				vals[1]: sdkmath.NewInt(12),
				vals[2]: sdkmath.NewInt(1),
				vals[3]: sdkmath.NewInt(10),
			},
		},
	}

	for caseNumber, val := range tc {
		t.Run(fmt.Sprint(caseNumber), func(t *testing.T) {
			sum := sdkmath.ZeroInt()
			for _, amount := range val.current {
				sum = sum.Add(amount)
			}
			allocations, err := types.DetermineAllocationsForDelegation(val.current, sum, val.target, val.inAmount, make(map[string]sdkmath.Int))
			require.NoError(t, err)
			require.Equal(t, len(val.expected), len(allocations))
			for valoper := range val.expected {
				ex, ok := val.expected[valoper]
				require.True(t, ok)
				ac, ok := allocations[valoper]
				require.True(t, ok)
				require.True(t, ex.Equal(ac), fmt.Sprintf("Test Case #%d failed; allocations did not equal expected output - expected %s, got %s.", caseNumber, val.expected[valoper], allocations[valoper]))
			}
		})
	}
}
