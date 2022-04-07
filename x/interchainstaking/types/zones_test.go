package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// func TestDetermineStateIntentDiff(t *testing.T) {
// 	zone := types.RegisteredZone{}
// 	d1 := []*types.Delegation{}
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1234", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(1000)})
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1235", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(500)})
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1236", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(300)})
// 	d1 = append(d1, &types.Delegation{DelegationAddress: "cosmos1user1237", ValidatorAddress: "cosmos12345667890", Amount: sdk.NewDec(200)})

// 	i1 := []types.DelegatorIntent{}
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1234", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.MustNewDecFromStr("0.5")}, {ValoperAddress: "cosmos987654321", Weight: sdk.MustNewDecFromStr("0.5")}}})
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1235", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.NewDec(1)}}})
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1236", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.NewDec(1)}}})
// 	i1 = append(i1, types.DelegatorIntent{Delegator: "cosmos1user1237", Intents: []*types.ValidatorIntent{{ValoperAddress: "cosmos12345667890", Weight: sdk.NewDec(1)}}})

// 	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmos12345667890", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewDec(2000), Delegations: d1})

// 	//fmt.Println(zone.DetermineStateIntentDiff(i1))

// 	require.Equal(t, 0, 0)
// }

func TestApplyDiffsToDistribution(t *testing.T) {
	testCases := []struct {
		distribution         map[string]sdk.Coin
		diff                 map[string]sdk.Int
		expectedDistribution map[string]sdk.Coin
		expectedRemainder    sdk.Int
	}{
		{
			distribution: map[string]sdk.Coin{
				"val1": sdk.NewInt64Coin("uatom", 3),
				"val2": sdk.NewInt64Coin("uatom", 3),
			},
			diff: map[string]sdk.Int{
				"val1": sdk.NewInt(-1),
				"val2": sdk.NewInt(1),
			},
			expectedDistribution: map[string]sdk.Coin{
				"val1": sdk.NewInt64Coin("uatom", 4),
				"val2": sdk.NewInt64Coin("uatom", 2),
			},
			expectedRemainder: sdk.ZeroInt(),
		},

		{
			distribution: map[string]sdk.Coin{
				"val1": sdk.NewInt64Coin("uatom", 1),
				"val2": sdk.NewInt64Coin("uatom", 5),
			},
			diff: map[string]sdk.Int{
				"val1": sdk.NewInt(-1),
				"val2": sdk.NewInt(1),
			},
			expectedDistribution: map[string]sdk.Coin{
				"val1": sdk.NewInt64Coin("uatom", 2),
				"val2": sdk.NewInt64Coin("uatom", 4),
			},
			expectedRemainder: sdk.ZeroInt(),
		},
		{
			distribution: map[string]sdk.Coin{
				"val1": sdk.NewInt64Coin("uatom", 1),
				"val2": sdk.NewInt64Coin("uatom", 5),
			},
			diff: map[string]sdk.Int{
				"val1": sdk.NewInt(2),
				"val2": sdk.NewInt(2),
				"val3": sdk.NewInt(-4),
				"val4": sdk.NewInt(0),
			},
			expectedDistribution: map[string]sdk.Coin{
				"val2": sdk.NewInt64Coin("uatom", 3),
			},
			expectedRemainder: sdk.NewInt(3),
		},
		{
			distribution: map[string]sdk.Coin{
				"val1": sdk.NewInt64Coin("uatom", 1),
				"val2": sdk.NewInt64Coin("uatom", 5),
				"val3": sdk.NewInt64Coin("uatom", 0),
			},
			diff: map[string]sdk.Int{
				"val1": sdk.NewInt(2),
				"val2": sdk.NewInt(2),
				"val3": sdk.NewInt(-4),
				"val4": sdk.NewInt(0),
			},
			expectedDistribution: map[string]sdk.Coin{
				"val2": sdk.NewInt64Coin("uatom", 3),
				"val3": sdk.NewInt64Coin("uatom", 3),
			},
			expectedRemainder: sdk.ZeroInt(),
		},
	}

	zone := types.RegisteredZone{}

	for idx, i := range testCases {
		result, remainder := zone.ApplyDiffsToDistribution(i.distribution, i.diff)
		for k, v := range i.expectedDistribution {
			require.Truef(t, v.IsEqual(result[k]), "case %d: distribution %v does not match expected %v", idx, result[k], v)
		}
		require.Truef(t, i.expectedRemainder.Equal(remainder), "case %d: remainder %v does not match expected %v", idx, remainder, i.expectedRemainder)
	}

}
