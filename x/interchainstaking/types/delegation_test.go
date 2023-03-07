package types_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
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

	// ensure error is returned for 0 length
	_, err := types.UnmarshalDelegation(types.ModuleCdc, []byte{})
	require.Error(t, err)
}

func TestSetForValoper(t *testing.T) {
	v1 := utils.GenerateValAddressForTest(r).String()
	v2 := utils.GenerateValAddressForTest(r).String()
	v3 := utils.GenerateValAddressForTest(r).String()

	intents := types.ValidatorIntents{
		{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(10, 1)},
		{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(90, 1)},
	}

	intents = intents.SetForValoper(v1, &types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(40, 1)})
	intents = intents.SetForValoper(v2, &types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(60, 1)})

	require.Equal(t, sdk.NewDecWithPrec(40, 1), intents.MustGetForValoper(v1).Weight)
	require.Equal(t, sdk.NewDecWithPrec(60, 1), intents.MustGetForValoper(v2).Weight)

	// check failed return
	actual := intents.MustGetForValoper(v3)
	require.Equal(t, sdk.ZeroDec(), actual.Weight)
}

func TestNormalizeValidatorIntentsDeterminism(t *testing.T) {
	v1 := utils.GenerateValAddressForTest(r).String()
	v2 := utils.GenerateValAddressForTest(r).String()
	v3 := utils.GenerateValAddressForTest(r).String()
	v4 := utils.GenerateValAddressForTest(r).String()

	cases := []struct {
		name    string
		intents types.ValidatorIntents
	}{
		{
			name: "case 1",
			intents: types.ValidatorIntents{
				{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(10, 1)},
				{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(90, 1)},
			},
		},
		{
			name: "case 2",
			intents: types.ValidatorIntents{
				{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(10, 1)},
				{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(90, 1)},
				{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(90, 1)},
				{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(90, 1)},
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

	val1 := utils.GenerateValAddressForTest(r)
	val2 := utils.GenerateValAddressForTest(r)
	val3 := utils.GenerateValAddressForTest(r)
	val4 := utils.GenerateValAddressForTest(r)

	tc := []struct {
		current  map[string]sdkmath.Int
		target   types.ValidatorIntents
		inAmount sdk.Coins
		expected map[string]sdkmath.Int
		dust     sdkmath.Int
	}{
		{
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(350000),
				val2.String(): sdk.NewInt(650000),
				val3.String(): sdk.NewInt(75000),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(30, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(63, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(7, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(50000))),
			expected: map[string]sdkmath.Int{
				val1.String(): sdk.ZeroInt(),
				val2.String(): sdk.NewInt(33182),
				val3.String(): sdk.NewInt(16818),
			},
		},
		{
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(52),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(4),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(20))),
			expected: map[string]sdkmath.Int{
				val4.String(): sdk.NewInt(11),
				val3.String(): sdk.ZeroInt(),
				val2.String(): sdk.NewInt(6),
				val1.String(): sdk.NewInt(3),
			},
		},
		{
			current: map[string]sdkmath.Int{
				val1.String(): sdk.NewInt(52),
				val2.String(): sdk.NewInt(24),
				val3.String(): sdk.NewInt(20),
				val4.String(): sdk.NewInt(4),
			},
			target: types.ValidatorIntents{
				{ValoperAddress: val1.String(), Weight: sdk.NewDecWithPrec(50, 2)},
				{ValoperAddress: val2.String(), Weight: sdk.NewDecWithPrec(25, 2)},
				{ValoperAddress: val3.String(), Weight: sdk.NewDecWithPrec(15, 2)},
				{ValoperAddress: val4.String(), Weight: sdk.NewDecWithPrec(10, 2)},
			},
			inAmount: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(50))),
			expected: map[string]sdkmath.Int{
				val4.String(): sdk.NewInt(20),
				val2.String(): sdk.NewInt(13),
				val1.String(): sdk.NewInt(10),
				val3.String(): sdk.NewInt(7),
			},
		},
	}

	for caseNumber, val := range tc {
		sum := sdkmath.ZeroInt()
		for _, amount := range val.current {
			sum = sum.Add(amount)
		}
		allocations := types.DetermineAllocationsForDelegation(val.current, sum, val.target, val.inAmount)
		require.Equal(t, len(val.expected), len(allocations))
		for valoper := range val.expected {
			ex, ok := val.expected[valoper]
			require.True(t, ok)
			ac, ok := allocations[valoper]
			require.True(t, ok)
			require.True(t, ex.Equal(ac), fmt.Sprintf("Test Case #%d failed; allocations did not equal expected output - expected %s, got %s.", caseNumber, val.expected[valoper], allocations[valoper]))
		}
	}
}
