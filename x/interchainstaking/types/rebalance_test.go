package types_test

import (
	"sort"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/require"
)

func GenerateValidatorsDeterministic(n int) (out []string) {
	out = make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper"))
	}
	sort.Strings(out)
	return out
}

func TestDetermineAllocationsForRebalancing(t *testing.T) {
	vals := GenerateValidatorsDeterministic(5)

	type testcase struct {
		name        string
		allocations map[string]math.Int
		target      types.ValidatorIntents
		locked      map[string]bool
		expected    types.RebalanceTargets
	}

	tcs := []testcase{
		{
			name: "100% No Existing Redelegations",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(10),
				vals[1]: math.NewInt(10),
				vals[2]: math.NewInt(10),
				vals[3]: math.NewInt(10),
				vals[4]: math.NewInt(10),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[0],
					Weight:         sdk.NewDec(1),
				},
			},
			locked: map[string]bool{},
			expected: types.RebalanceTargets{
				{
					Source: vals[1],
					Target: vals[0],
					Amount: math.NewInt(10),
				},
				{
					Source: vals[2],
					Target: vals[0],
					Amount: math.NewInt(10),
				},
				{
					Source: vals[3],
					Target: vals[0],
					Amount: math.NewInt(5),
				},
			},
		},
		{
			name: "50/50 No Existing Redelegations, Constrained by total",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(10),
				vals[1]: math.NewInt(10),
				vals[2]: math.NewInt(10),
				vals[3]: math.NewInt(10),
				vals[4]: math.NewInt(10),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[0],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[1],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
			},
			locked: map[string]bool{},
			expected: types.RebalanceTargets{
				{
					Source: vals[2],
					Target: vals[0],
					Amount: math.NewInt(10),
				},
				{
					Source: vals[3],
					Target: vals[0],
					Amount: math.NewInt(5),
				},
				{
					Source: vals[3],
					Target: vals[1],
					Amount: math.NewInt(5),
				},
				{
					Source: vals[4],
					Target: vals[1],
					Amount: math.NewInt(5),
				},
			},
		},
		{
			name: "50/50 No Existing Redelegations, Unconstrained",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(10),
				vals[1]: math.NewInt(10),
				vals[2]: math.NewInt(10),
				vals[3]: math.NewInt(10),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[0],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[1],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
			},
			locked: map[string]bool{},
			expected: []*types.RebalanceTarget{
				{
					Source: vals[2],
					Target: vals[0],
					Amount: math.NewInt(10),
				},
				{
					Source: vals[3],
					Target: vals[1],
					Amount: math.NewInt(10),
				},
			},
		},
		{
			name: "Drop one validator, No Existing Redelegations",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(8),
				vals[1]: math.NewInt(8),
				vals[2]: math.NewInt(8),
				vals[3]: math.NewInt(8),
				vals[4]: math.NewInt(8),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[0],
					Weight:         sdk.NewDecWithPrec(25, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[1],
					Weight:         sdk.NewDecWithPrec(25, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[2],
					Weight:         sdk.NewDecWithPrec(25, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[3],
					Weight:         sdk.NewDecWithPrec(25, 2),
				},
			},
			locked: map[string]bool{},
			expected: types.RebalanceTargets{
				{
					Source: vals[4],
					Target: vals[0],
					Amount: math.NewInt(2),
				},
				{
					Source: vals[4],
					Target: vals[1],
					Amount: math.NewInt(2),
				},
				{
					Source: vals[4],
					Target: vals[2],
					Amount: math.NewInt(2),
				},
				{
					Source: vals[4],
					Target: vals[3],
					Amount: math.NewInt(2),
				},
			},
		},
		{
			name: "Add one validator, No Existing Redelegations",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(10),
				vals[1]: math.NewInt(10),
				vals[2]: math.NewInt(10),
				vals[3]: math.NewInt(10),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[0],
					Weight:         sdk.NewDecWithPrec(20, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[1],
					Weight:         sdk.NewDecWithPrec(20, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[2],
					Weight:         sdk.NewDecWithPrec(20, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[3],
					Weight:         sdk.NewDecWithPrec(20, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[4],
					Weight:         sdk.NewDecWithPrec(20, 2),
				},
			},
			locked: map[string]bool{},
			expected: types.RebalanceTargets{
				{
					Source: vals[0],
					Target: vals[4],
					Amount: math.NewInt(2),
				},
				{
					Source: vals[1],
					Target: vals[4],
					Amount: math.NewInt(2),
				},
				{
					Source: vals[2],
					Target: vals[4],
					Amount: math.NewInt(2),
				},
				{
					Source: vals[3],
					Target: vals[4],
					Amount: math.NewInt(2),
				},
			},
		},
		{
			name: "Attempt redelegate away from locked validator; no-op",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(10),
				vals[1]: math.NewInt(10),
				vals[2]: math.NewInt(10),
				vals[3]: math.NewInt(10),
				vals[4]: math.NewInt(10),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[0],
					Weight:         sdk.NewDecWithPrec(10, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[1],
					Weight:         sdk.NewDecWithPrec(225, 3),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[2],
					Weight:         sdk.NewDecWithPrec(225, 3),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[3],
					Weight:         sdk.NewDecWithPrec(225, 3),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[4],
					Weight:         sdk.NewDecWithPrec(225, 3),
				},
			},
			locked: map[string]bool{
				vals[0]: true,
			},
			expected: types.RebalanceTargets{},
		},
		{
			name: "Delegate away from 2; 1 locked validator; v1 -15; v2 + 10; v3 +5",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(20),
				vals[1]: math.NewInt(20),
				vals[2]: math.NewInt(20),
				vals[3]: math.NewInt(20),
				vals[4]: math.NewInt(20),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[0],
					Weight:         sdk.NewDecWithPrec(5, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[1],
					Weight:         sdk.NewDecWithPrec(5, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[2],
					Weight:         sdk.NewDecWithPrec(30, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[3],
					Weight:         sdk.NewDecWithPrec(30, 2),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[4],
					Weight:         sdk.NewDecWithPrec(30, 2),
				},
			},
			locked: map[string]bool{
				vals[0]: true,
			},
			expected: types.RebalanceTargets{
				{
					Source: vals[1],
					Target: vals[2],
					Amount: math.NewInt(10),
				},
				{
					Source: vals[1],
					Target: vals[3],
					Amount: math.NewInt(5),
				},
			},
		},
		{
			name: "v0 missing, v1 zero; one new vals. Should delegate v0: -50; v1: -25; v2: +25; v3: +50",
			allocations: map[string]math.Int{
				vals[0]: math.NewInt(50),
				vals[1]: math.NewInt(50),
				vals[2]: math.NewInt(50),
			},
			target: types.ValidatorIntents{
				&types.ValidatorIntent{
					ValoperAddress: vals[1],
					Weight:         sdk.ZeroDec(),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[2],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				&types.ValidatorIntent{
					ValoperAddress: vals[3],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
			},
			locked: map[string]bool{},
			expected: types.RebalanceTargets{
				{
					Source: vals[0],
					Target: vals[2],
					Amount: math.NewInt(25),
				},
				{
					Source: vals[0],
					Target: vals[3],
					Amount: math.NewInt(25),
				},
				{
					Source: vals[1],
					Target: vals[3],
					Amount: math.NewInt(25),
				},
			},
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			currentSum, lockedSum := func(in map[string]math.Int, locked map[string]bool) (sum, lockedSum math.Int) {
				sum = math.ZeroInt()
				lockedSum = math.ZeroInt()
				for k, v := range in {
					sum = sum.Add(v)
					if locked[k] {
						lockedSum = lockedSum.Add(v)
					}
				}
				return sum, lockedSum
			}(tt.allocations, tt.locked)

			actual := types.DetermineAllocationsForRebalancing(
				tt.allocations, tt.locked, currentSum, lockedSum, tt.target, nil,
			)

			require.ElementsMatch(t, tt.expected, actual)
		})
	}
}
