package types_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestDetermineAllocationsForUndelegation(t *testing.T) {
	vals := GenerateValidatorsDeterministic(4)
	tests := []struct {
		name               string
		currentAllocations map[string]sdkmath.Int
		unlocked           map[string]sdkmath.Int
		targetAllocations  types.ValidatorIntents
		amount             sdk.Coins
		expected           map[string]sdkmath.Int
	}{
		{
			name: "case 0: equal delegations, equal intents; no locked",
			currentAllocations: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000),
				vals[1]: sdk.NewInt(1000),
				vals[2]: sdk.NewInt(1000),
				vals[3]: sdk.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000),
				vals[1]: sdk.NewInt(1000),
				vals[2]: sdk.NewInt(1000),
				vals[3]: sdk.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(250),
				vals[1]: sdk.NewInt(250),
				vals[2]: sdk.NewInt(250),
				vals[3]: sdk.NewInt(250),
			},
		},
		{
			name: "case 1: unequal delegations, equal intents; no locked",
			currentAllocations: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000), // + 25
				vals[1]: sdk.NewInt(950),  // -25
				vals[2]: sdk.NewInt(1200), // + 225
				vals[3]: sdk.NewInt(750),  // -225
				// 250; 0, -25, 0, -225; 225, 200, 225, 0; 650 (275, 225, 475, 25)
			},
			unlocked: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000),
				vals[1]: sdk.NewInt(950),
				vals[2]: sdk.NewInt(1200),
				vals[3]: sdk.NewInt(750),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(275),
				vals[1]: sdk.NewInt(225),
				vals[2]: sdk.NewInt(475),
				vals[3]: sdk.NewInt(25),
			},
		},
		{
			name: "case 2: unequal delegations, unequal intents; no locked",
			currentAllocations: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(5000), // +500
				vals[1]: sdk.NewInt(1800), // 0
				vals[2]: sdk.NewInt(1200), // -150
				vals[3]: sdk.NewInt(1000), // + 550
			},
			unlocked: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(5000),
				vals[1]: sdk.NewInt(1800),
				vals[2]: sdk.NewInt(1200),
				vals[3]: sdk.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(50, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(20, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(15, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(5, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(300))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(143),
				vals[1]: sdk.NewInt(0), // 0
				vals[2]: sdk.NewInt(0), // 0
				vals[3]: sdk.NewInt(157),
			},
		},
		{
			name: "case 3: unequal delegations, unequal intents, big discrepancy in intent (should not take from underallocated); no locked",
			currentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1234), // 252  +982
				vals[1]: sdkmath.NewInt(675),  // 504   +171
				vals[2]: sdkmath.NewInt(210),  // 756   - 546
				vals[3]: sdkmath.NewInt(401),  // 1008  - 607
			},
			unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1234),
				vals[1]: sdkmath.NewInt(675),
				vals[2]: sdkmath.NewInt(210),
				vals[3]: sdkmath.NewInt(401),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(10, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(20, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(30, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(40, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(86),
				vals[1]: sdkmath.NewInt(14),
				vals[2]: sdkmath.NewInt(0),
				vals[3]: sdkmath.NewInt(0),
			},
		},
		{
			name: "case 4: equal delegations, equal intents; vals[0] partially locked",
			currentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1000),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(500),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(250),
				vals[1]: sdkmath.NewInt(250),
				vals[2]: sdkmath.NewInt(250),
				vals[3]: sdkmath.NewInt(250),
			},
		},
		{
			name: "case 5: equal delegations, equal intents; vals[0] partially locked #2",
			currentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1000),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(200),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(200),
				vals[1]: sdkmath.NewInt(300),
				vals[2]: sdkmath.NewInt(250),
				vals[3]: sdkmath.NewInt(250),
			},
		},
		{
			name: "case 6: equal delegations, equal intents; vals[0] completely locked",
			currentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1000),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(0),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(0),
				vals[1]: sdkmath.NewInt(334),
				vals[2]: sdkmath.NewInt(333),
				vals[3]: sdkmath.NewInt(333),
			},
		},
	}
	for _, tt := range tests {
		sum := func(in map[string]sdkmath.Int) sdkmath.Int {
			out := sdk.ZeroInt()
			for _, i := range in {
				out = out.Add(i)
			}
			return out
		}
		allocations := types.DetermineAllocationsForUndelegation(tt.currentAllocations, map[string]bool{}, sum(tt.currentAllocations), tt.targetAllocations, tt.unlocked, tt.amount)
		for valoper := range allocations {
			require.Equal(t, tt.expected[valoper].Int64(), allocations[valoper].Int64(), fmt.Sprintf("%s / %s", tt.name, valoper))
		}
		// validate that the amount withdrawn always matches total amount.
		require.Equal(t, tt.amount[0].Amount, sum(allocations))
	}
}
