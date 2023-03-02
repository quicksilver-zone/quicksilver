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
	v1 := "cosmosvaloper1f9xf8pl9ev5e6amuec9680scqhpyd85ee25q5y90uqpxdcdfmufsal5nd0"
	v2 := "cosmosvaloper1345w0hf2x4c5vdcm4wmj38pn26z0pmgvvjj8s0tplj87cxegw34qu8f6l1"
	v3 := "cosmosvaloper1fslnzgde7z8mexm9y3evcy9a8t9km0lshtr92a2jsq3mhtku789qy62jg2"
	v4 := "cosmosvaloper1vx4el3dyqzucc3jhc6leaxt8gg9sar78elrce7wvyqsk9uu2d99slwapz3"
	tests := []struct {
		name               string
		currentAllocations map[string]sdkmath.Int
		unlocked           map[string]sdkmath.Int
		targetAllocations  types.ValidatorIntents
		amount             sdk.Coins
		expected           map[string]sdkmath.Int
	}{
		{
			name: "equal delegations, equal intents; no locked",
			currentAllocations: map[string]sdkmath.Int{
				v1: sdk.NewInt(1000),
				v2: sdk.NewInt(1000),
				v3: sdk.NewInt(1000),
				v4: sdk.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				v1: sdk.NewInt(1000),
				v2: sdk.NewInt(1000),
				v3: sdk.NewInt(1000),
				v4: sdk.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				v1: sdk.NewInt(250),
				v2: sdk.NewInt(250),
				v3: sdk.NewInt(250),
				v4: sdk.NewInt(250),
			},
		},
		{
			name: "unequal delegations, equal intents; no locked",
			currentAllocations: map[string]sdkmath.Int{
				v1: sdk.NewInt(1000), // + 25
				v2: sdk.NewInt(950),  // -25
				v3: sdk.NewInt(1200), // + 225
				v4: sdk.NewInt(750),  // -225
				// 250; 0, -25, 0, -225; 225, 200, 225, 0; 650 (275, 225, 475, 25)
			},
			unlocked: map[string]sdkmath.Int{
				v1: sdk.NewInt(1000),
				v2: sdk.NewInt(950),
				v3: sdk.NewInt(1200),
				v4: sdk.NewInt(750),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				v1: sdk.NewInt(275),
				v2: sdk.NewInt(225),
				v3: sdk.NewInt(475),
				v4: sdk.NewInt(25),
			},
		},
		{
			name: "unequal delegations, unequal intents; no locked",
			currentAllocations: map[string]sdkmath.Int{
				v1: sdk.NewInt(5000),
				v2: sdk.NewInt(1800),
				v3: sdk.NewInt(1200),
				v4: sdk.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				v1: sdk.NewInt(5000),
				v2: sdk.NewInt(1800),
				v3: sdk.NewInt(1200),
				v4: sdk.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(50, 2)},
				&types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(20, 2)},
				&types.ValidatorIntent{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(15, 2)},
				&types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(5, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(300))),
			expected: map[string]sdkmath.Int{
				v1: sdk.NewInt(154),
				v2: sdk.NewInt(14),
				v3: sdk.NewInt(0),
				v4: sdk.NewInt(132),
			},
		},
		{
			name: "unequal delegations, unequal intents, big discrepancy in intent (should not take from underallocated); no locked",
			currentAllocations: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(1234),
				v2: sdkmath.NewInt(675),
				v3: sdkmath.NewInt(210),
				v4: sdkmath.NewInt(401),
			},
			unlocked: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(1234),
				v2: sdkmath.NewInt(675),
				v3: sdkmath.NewInt(210),
				v4: sdkmath.NewInt(401),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(10, 2)},
				&types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(20, 2)},
				&types.ValidatorIntent{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(30, 2)},
				&types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(40, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100))),
			expected: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(84),
				v2: sdkmath.NewInt(16),
				v3: sdkmath.NewInt(0),
				v4: sdkmath.NewInt(0),
			},
		},
		{
			name: "equal delegations, equal intents; v1 partially locked",
			currentAllocations: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(1000),
				v2: sdkmath.NewInt(1000),
				v3: sdkmath.NewInt(1000),
				v4: sdkmath.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(500),
				v2: sdkmath.NewInt(1000),
				v3: sdkmath.NewInt(1000),
				v4: sdkmath.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(250),
				v2: sdkmath.NewInt(250),
				v3: sdkmath.NewInt(250),
				v4: sdkmath.NewInt(250),
			},
		},
		{
			name: "equal delegations, equal intents; v1 partially locked #2",
			currentAllocations: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(1000),
				v2: sdkmath.NewInt(1000),
				v3: sdkmath.NewInt(1000),
				v4: sdkmath.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(200),
				v2: sdkmath.NewInt(1000),
				v3: sdkmath.NewInt(1000),
				v4: sdkmath.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(200),
				v2: sdkmath.NewInt(276),
				v3: sdkmath.NewInt(262),
				v4: sdkmath.NewInt(262),
			},
		},
		{
			name: "equal delegations, equal intents; v1 completely locked",
			currentAllocations: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(1000),
				v2: sdkmath.NewInt(1000),
				v3: sdkmath.NewInt(1000),
				v4: sdkmath.NewInt(1000),
			},
			unlocked: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(0),
				v2: sdkmath.NewInt(1000),
				v3: sdkmath.NewInt(1000),
				v4: sdkmath.NewInt(1000),
			},
			targetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v3, Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(25, 2)},
			},
			amount: sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				v1: sdkmath.NewInt(0),
				v2: sdkmath.NewInt(376),
				v3: sdkmath.NewInt(312),
				v4: sdkmath.NewInt(312),
			},
		},
	}
	for testidx, tt := range tests {
		sum := func(in map[string]sdkmath.Int) sdkmath.Int {
			out := sdk.ZeroInt()
			for _, i := range in {
				out = out.Add(i)
			}
			return out
		}
		fmt.Printf("=============== case %d ===============\n", testidx)
		allocations := types.DetermineAllocationsForUndelegation(tt.currentAllocations, sum(tt.currentAllocations), tt.targetAllocations, tt.unlocked, tt.amount)
		for valoper := range allocations {
			require.Equal(t, tt.expected[valoper].Int64(), allocations[valoper].Int64(), fmt.Sprintf("%s (%d) / %s", tt.name, testidx, valoper))
		}
		// validate that the amount withdrawn always matches total amount.
		require.Equal(t, tt.amount[0].Amount, sum(allocations))
	}
}
