package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestDefaultIntent(t *testing.T) {
	zone := types.RegisteredZone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})

	out := zone.DefaultAggregateIntents()
	for _, v := range out {
		if !v.Weight.Equal(sdk.NewDecWithPrec(2, 1)) {
			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(2, 1), v.Weight)
		}
	}
}

func TestCoinsToIntent(t *testing.T) {
	zone := types.RegisteredZone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})

	testCases := []struct {
		amount         sdk.Coins
		expectedIntent map[string]sdk.Dec
	}{
		{
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(45)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf16", sdk.NewInt(55)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(45),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(55),
			},
		},
		{
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(350)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf16", sdk.NewInt(350)),
				sdk.NewCoin("cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy6", sdk.NewInt(300)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(350),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(350),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(300),
			},
		},
		{
			amount: sdk.NewCoins(
				sdk.NewCoin("cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj01", sdk.NewInt(3900)),
				sdk.NewCoin("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf16", sdk.NewInt(5500)),
				sdk.NewCoin("cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy6", sdk.NewInt(3000)),
				sdk.NewCoin("cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll2", sdk.NewInt(500)),
			),
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(3900),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(5500),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(3000),
				"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll": sdk.NewDec(500),
			},
		},
	}

	for _, tc := range testCases {
		out := zone.ConvertCoinsToOrdinalIntents(tc.amount)
		for k, v := range out {
			if !tc.expectedIntent[k].Equal(v.Weight) {
				t.Errorf("Got %v expected %v", v.Weight, tc.expectedIntent[k])
			}
		}
	}
}

func TestBase64MemoToIntent(t *testing.T) {
	zone := types.RegisteredZone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})
	zone.Validators = append(zone.Validators, &types.Validator{ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000)})

	testCases := []struct {
		memo           string
		amount         int
		expectedIntent map[string]sdk.Dec
	}{
		{
			memo:   "WoS/+Ex92tEcuMBzhukZKMVnXKS8bqaQBJTx9zza4rrxyLiP9fwLijOc",
			amount: 100,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(45),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(55),
			},
		},
		{
			memo:   "RoS/+Ex92tEcuMBzhukZKMVnXKS8RqaQBJTx9zza4rrxyLiP9fwLijOcPK/59acWzdcBME6ub8f0LID97qWE",
			amount: 1000,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDec(350),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDec(350),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(300),
			},
		},
		{
			memo:   "ToS/+Ex92tEcuMBzhukZKMVnXKS8NKaQBJTx9zza4rrxyLiP9fwLijOcPK/59acWzdcBME6ub8f0LID97qWECuxJKXmxBM1YBQyWHSAXDiwmMY78",
			amount: 10,
			expectedIntent: map[string]sdk.Dec{
				"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewDecWithPrec(39, 1),
				"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewDecWithPrec(26, 1),
				"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewDec(3),
				"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll": sdk.NewDecWithPrec(5, 1),
			},
		},
	}

	for _, tc := range testCases {
		out := zone.ConvertMemoToOrdinalIntents(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(int64(tc.amount)))), tc.memo)
		for k, v := range out {
			if !tc.expectedIntent[k].Equal(v.Weight) {
				t.Errorf("Got %v expected %v", v.Weight, tc.expectedIntent[k])
			}
		}
	}
}

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

// 	require.Equal(t, 0, 0)
// }

// func TestApplyDiffsToDistribution(t *testing.T) {
// 	testCases := []struct {
// 		distribution         map[string]sdk.Coin
// 		diff                 map[string]sdk.Int
// 		expectedDistribution map[string]sdk.Coin
// 		expectedRemainder    sdk.Int
// 	}{
// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 3),
// 				"val2": sdk.NewInt64Coin("uatom", 3),
// 			},
// 			diff: map[string]sdk.Int{
// 				"val1": sdk.NewInt(-1),
// 				"val2": sdk.NewInt(1),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 4),
// 				"val2": sdk.NewInt64Coin("uatom", 2),
// 			},
// 			expectedRemainder: sdk.ZeroInt(),
// 		},

// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 1),
// 				"val2": sdk.NewInt64Coin("uatom", 5),
// 			},
// 			diff: map[string]sdk.Int{
// 				"val1": sdk.NewInt(-1),
// 				"val2": sdk.NewInt(1),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 2),
// 				"val2": sdk.NewInt64Coin("uatom", 4),
// 			},
// 			expectedRemainder: sdk.ZeroInt(),
// 		},
// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 1),
// 				"val2": sdk.NewInt64Coin("uatom", 5),
// 			},
// 			diff: map[string]sdk.Int{
// 				"val1": sdk.NewInt(2),
// 				"val2": sdk.NewInt(2),
// 				"val3": sdk.NewInt(-4),
// 				"val4": sdk.NewInt(0),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val2": sdk.NewInt64Coin("uatom", 3),
// 			},
// 			expectedRemainder: sdk.NewInt(3),
// 		},
// 		{
// 			distribution: map[string]sdk.Coin{
// 				"val1": sdk.NewInt64Coin("uatom", 1),
// 				"val2": sdk.NewInt64Coin("uatom", 5),
// 				"val3": sdk.NewInt64Coin("uatom", 0),
// 			},
// 			diff: map[string]sdk.Int{
// 				"val1": sdk.NewInt(2),
// 				"val2": sdk.NewInt(2),
// 				"val3": sdk.NewInt(-4),
// 				"val4": sdk.NewInt(0),
// 			},
// 			expectedDistribution: map[string]sdk.Coin{
// 				"val2": sdk.NewInt64Coin("uatom", 3),
// 				"val3": sdk.NewInt64Coin("uatom", 3),
// 			},
// 			expectedRemainder: sdk.ZeroInt(),
// 		},
// 	}

// 	zone := types.RegisteredZone{}

// 	for idx, i := range testCases {
// 		result, remainder := zone.ApplyDiffsToDistribution(i.distribution, i.diff)
// 		for k, v := range i.expectedDistribution {
// 			require.Truef(t, v.IsEqual(result[k]), "case %d: distribution %v does not match expected %v", idx, result[k], v)
// 		}
// 		require.Truef(t, i.expectedRemainder.Equal(remainder), "case %d: remainder %v does not match expected %v", idx, remainder, i.expectedRemainder)
// 	}

// }
