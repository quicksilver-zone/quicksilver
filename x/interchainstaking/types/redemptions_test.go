package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

type detRedemptionTest struct {
	Name               string                 `json:"name"`
	CurrentAllocations map[string]sdkmath.Int `json:"cur_allocs"`
	Unlocked           map[string]sdkmath.Int `json:"Unlocked"`
	TargetAllocations  types.ValidatorIntents `json:"targ_allocs"`
	Amount             sdk.Coins              `json:"amount"`
	expected           map[string]sdkmath.Int
}

func detRedemptionTests(vals []string) []*detRedemptionTest {
	return []*detRedemptionTest{
		{
			Name: "case 0: equal delegations, equal intents; no locked",
			CurrentAllocations: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000),
				vals[1]: sdk.NewInt(1000),
				vals[2]: sdk.NewInt(1000),
				vals[3]: sdk.NewInt(1000),
			},
			Unlocked: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000),
				vals[1]: sdk.NewInt(1000),
				vals[2]: sdk.NewInt(1000),
				vals[3]: sdk.NewInt(1000),
			},
			TargetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(250),
				vals[1]: sdk.NewInt(250),
				vals[2]: sdk.NewInt(250),
				vals[3]: sdk.NewInt(250),
			},
		},
		{
			Name: "case 1: unequal delegations, equal intents; no locked",
			CurrentAllocations: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000), // + 25
				vals[1]: sdk.NewInt(950),  // -25
				vals[2]: sdk.NewInt(1200), // + 225
				vals[3]: sdk.NewInt(750),  // -225
				// 250; 0, -25, 0, -225; 225, 200, 225, 0; 650 (275, 225, 475, 25)
			},
			Unlocked: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(1000),
				vals[1]: sdk.NewInt(950),
				vals[2]: sdk.NewInt(1200),
				vals[3]: sdk.NewInt(750),
			},
			TargetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				// vals[0]: sdk.NewInt(275),
				// vals[1]: sdk.NewInt(225),
				// vals[2]: sdk.NewInt(475),
				// vals[3]: sdk.NewInt(25),
				vals[0]: sdk.NewInt(150),
				vals[1]: sdk.NewInt(150),
				vals[2]: sdk.NewInt(350),
				vals[3]: sdk.NewInt(350),
			},
		},
		{
			Name: "case 2: unequal delegations, unequal intents; no locked",
			CurrentAllocations: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(5000), // +500
				vals[1]: sdk.NewInt(1800), // 0
				vals[2]: sdk.NewInt(1200), // -150
				vals[3]: sdk.NewInt(1000), // + 550
			},
			Unlocked: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(5000),
				vals[1]: sdk.NewInt(1800),
				vals[2]: sdk.NewInt(1200),
				vals[3]: sdk.NewInt(1000),
			},
			TargetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(50, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(20, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(15, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(5, 2)},
			},
			Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(300))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdk.NewInt(143),
				vals[1]: sdk.NewInt(0), // 0
				vals[2]: sdk.NewInt(0), // 0
				vals[3]: sdk.NewInt(157),
			},
		},
		{
			Name: "case 3: unequal delegations, unequal intents, big discrepancy in intent (should not take from underallocated); no locked",
			CurrentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1234), // 252  +982
				vals[1]: sdkmath.NewInt(675),  // 504   +171
				vals[2]: sdkmath.NewInt(210),  // 756   - 546
				vals[3]: sdkmath.NewInt(401),  // 1008  - 607
			},
			Unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1234),
				vals[1]: sdkmath.NewInt(675),
				vals[2]: sdkmath.NewInt(210),
				vals[3]: sdkmath.NewInt(401),
			},
			TargetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(10, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(20, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(30, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(40, 2)},
			},
			Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(86),
				vals[1]: sdkmath.NewInt(14),
				vals[2]: sdkmath.NewInt(0),
				vals[3]: sdkmath.NewInt(0),
			},
		},
		{
			Name: "case 4: equal delegations, equal intents; vals[0] partially locked",
			CurrentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1000),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			Unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(500),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			TargetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(250),
				vals[1]: sdkmath.NewInt(250),
				vals[2]: sdkmath.NewInt(250),
				vals[3]: sdkmath.NewInt(250),
			},
		},
		{
			Name: "case 5: equal delegations, equal intents; vals[0] partially locked #2",
			CurrentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1000),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			Unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(200),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			TargetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(200),
				vals[1]: sdkmath.NewInt(300),
				vals[2]: sdkmath.NewInt(250),
				vals[3]: sdkmath.NewInt(250),
			},
		},
		{
			Name: "case 6: equal delegations, equal intents; vals[0] completely locked",
			CurrentAllocations: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(1000),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			Unlocked: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(0),
				vals[1]: sdkmath.NewInt(1000),
				vals[2]: sdkmath.NewInt(1000),
				vals[3]: sdkmath.NewInt(1000),
			},
			TargetAllocations: types.ValidatorIntents{
				&types.ValidatorIntent{ValoperAddress: vals[0], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[1], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[2], Weight: sdk.NewDecWithPrec(25, 2)},
				&types.ValidatorIntent{ValoperAddress: vals[3], Weight: sdk.NewDecWithPrec(25, 2)},
			},
			Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
			expected: map[string]sdkmath.Int{
				vals[0]: sdkmath.NewInt(0),
				vals[1]: sdkmath.NewInt(334),
				vals[2]: sdkmath.NewInt(333),
				vals[3]: sdkmath.NewInt(333),
			},
		},
	}
}

func sum(in map[string]sdkmath.Int) sdkmath.Int {
	out := sdk.ZeroInt()
	for _, i := range in {
		out = out.Add(i)
	}
	return out
}

func TestDetermineAllocationsForUndelegation(t *testing.T) {
	tests := detRedemptionTests(vals)

	for _, tt := range tests {
		allocations, err := types.DetermineAllocationsForUndelegation(tt.CurrentAllocations, map[string]bool{}, sum(tt.CurrentAllocations), tt.TargetAllocations, tt.Unlocked, tt.Amount)
		require.NoError(t, err)
		for valoper := range allocations {
			require.Equal(t, tt.expected[valoper].Int64(), allocations[valoper].Int64(), fmt.Sprintf("%s / %s", tt.Name, valoper))
		}
		// validate that the amount withdrawn always matches total amount.
		require.Equal(t, tt.Amount[0].Amount, sum(allocations))
	}
}

// The function should correctly calculate allocations for undelegation when there are overallocated validators.
func TestOverAndUnderAllocatedValidators(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(5000),
		"validator2": sdkmath.NewInt(1800),
		"validator3": sdkmath.NewInt(1200),
		"validator4": sdkmath.NewInt(1000),
	}
	lockedAllocations := map[string]bool{}
	currentSum := sdkmath.NewInt(9000)
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(50, 2)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(20, 2)},
		{ValoperAddress: "validator3", Weight: sdk.NewDecWithPrec(15, 2)},
		{ValoperAddress: "validator4", Weight: sdk.NewDecWithPrec(5, 2)},
	}.Normalize()
	availablePerValidator := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(5000),
		"validator2": sdkmath.NewInt(1800),
		"validator3": sdkmath.NewInt(1200),
		"validator4": sdkmath.NewInt(1000),
	}
	amount := sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(300))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator4": sdkmath.NewInt(300),
	}

	allocations, err := types.DetermineAllocationsForUndelegation(currentAllocations, lockedAllocations, currentSum, targetAllocations, availablePerValidator, amount)
	require.NoError(t, err)
	require.Equal(t, expectedAllocations, allocations)
}

// func TestOverAndUnderAllocatedValidatorsDiv3(t *testing.T) {
// 	currentAllocations := map[string]sdkmath.Int{
// 		"validator1": sdkmath.NewInt(5000),
// 		"validator2": sdkmath.NewInt(4000),
// 		"validator3": sdkmath.NewInt(3000),
// 		"validator4": sdkmath.NewInt(2000),
// 	}
// 	lockedAllocations := map[string]bool{}
// 	currentSum := sdkmath.NewInt(14000)
// 	targetAllocations := types.ValidatorIntents{
// 		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(40, 2)},
// 		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(30, 2)},
// 		{ValoperAddress: "validator3", Weight: sdk.NewDecWithPrec(20, 2)},
// 		{ValoperAddress: "validator4", Weight: sdk.NewDecWithPrec(10, 2)},
// 	}.Normalize()
// 	availablePerValidator := map[string]sdkmath.Int{
// 		"validator1": sdkmath.NewInt(4000),
// 		"validator2": sdkmath.NewInt(3000),
// 		"validator3": sdkmath.NewInt(2000),
// 		"validator4": sdkmath.NewInt(1000),
// 	}
// 	amount := sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(300))}

// 	expectedAllocations := map[string]sdkmath.Int{
// 		"validator1": sdkmath.NewInt(0), /// << why?
// 		"validator4": sdkmath.NewInt(300),
// 	}

// 	allocations := types.DetermineAllocationsForUndelegation(CurrentAllocations, lockedAllocations, currentSum, TargetAllocations, availablePerValidator, amount)

// 	fmt.Println(allocations)
// 	require.Equal(t, expectedAllocations, allocations)
// }
