package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestSumAllAllocation(t *testing.T) {
	a := types.Allocation{Address: "cosmos109kkuly5hyxz6kjqs68vazcknc3wj7hvgp4vgypz6z6c4j0y4t9qckmua8", Amount: sdk.Coins{sdk.Coin{Denom: "uatom", Amount: sdk.ZeroInt()}}}
	a.SumAll()
}

func TestAllocationsAllocate(t *testing.T) {
	a := types.Allocations{}
	a = a.Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}})
	a = a.Allocate("test2", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}})
	a = a.Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}})
	for _, i := range a {
		fmt.Printf("%v", i)
	}
}

func TestAllocationsSub(t *testing.T) {
	testCases := []struct {
		Allocations types.Allocations
		ToSub       types.Allocation
		Result      types.Allocations
		Remainder   sdk.Coins
	}{
		{
			types.Allocations{}.
				Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}).
				Allocate("test2", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}),
			types.Allocation{"test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}},
			types.Allocations{}.
				Allocate("test2", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}),
			sdk.NewCoins(),
		},
		{
			types.Allocations{}.
				Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.NewInt(16)}}).
				Allocate("test2", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}),
			types.Allocation{"test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}},
			types.Allocations{}.
				Allocate("test2", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}).Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.NewInt(15)}}),
			sdk.NewCoins(),
		},
		{
			types.Allocations{}.
				Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.NewInt(16)}}),
			types.Allocation{"test2", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.OneInt()}}},
			types.Allocations{}.
				Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.NewInt(16)}}),
			sdk.NewCoins(sdk.NewCoin("testCoin", sdk.OneInt())),
		},
		{
			types.Allocations{}.
				Allocate("test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.NewInt(16)}}),
			types.Allocation{"test1", sdk.Coins{sdk.Coin{Denom: "testCoin", Amount: sdk.NewInt(20)}}},
			types.Allocations{}.
				Allocate("test1", sdk.Coins{}),
			sdk.NewCoins(sdk.NewCoin("testCoin", sdk.NewInt(4))),
		},
	}
	for _, tc := range testCases {
		out, remainder := tc.Allocations.Sub(tc.ToSub.Amount, tc.ToSub.Address)
		for _, a := range tc.Result {
			if !out.Get(a.Address).Amount.IsEqual(a.Amount) {
				t.Errorf("allocation mismatch between expected %s and actual %s", a.Amount, out.Get(a.Address).Amount)
			}
		}
		if !remainder.IsEqual(tc.Remainder) {
			t.Errorf("remainder mismatch between expected %s and actual %s", tc.Remainder, remainder)
		}
	}
}

// func TestFindAccountForDelegation(t *testing.T) {
// 	testCases := []struct {
// 		bins      types.DelegationBins
// 		validator string
// 		expected  string
// 	}{
// 		{
// 			// no buckets in use, use first
// 			bins: types.DelegationBins{
// 				"cosmos1": types.DelegationBin{},
// 				"cosmos2": types.DelegationBin{},
// 				"cosmos3": types.DelegationBin{},
// 				"cosmos4": types.DelegationBin{},
// 				"cosmos5": types.DelegationBin{},
// 			},
// 			validator: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
// 			expected:  "cosmos1",
// 		},
// 		{
// 			// given the existing delegations are large, use the next free bucket
// 			bins: types.DelegationBins{
// 				"cosmos1": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(30000)},
// 				"cosmos2": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(45000)},
// 				"cosmos3": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(15000)},
// 				"cosmos4": types.DelegationBin{},
// 				"cosmos5": types.DelegationBin{},
// 			},
// 			validator: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
// 			expected:  "cosmos4",
// 		},
// 		{
// 			// previously unseen validator, use the smallest bucket
// 			bins: types.DelegationBins{
// 				"cosmos1": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(100)},
// 				"cosmos2": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(3000)},
// 				"cosmos3": types.DelegationBin{"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll": sdk.NewInt(6000)},
// 				"cosmos4": types.DelegationBin{"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewInt(3000)},
// 				"cosmos5": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(15000)},
// 				"cosmos6": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(26000)},
// 				"cosmos7": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(11000)},
// 			},
// 			validator: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
// 			expected:  "cosmos1",
// 		},
// 		{
// 			// previously seen, but by far largest buckets, use smallest bucket
// 			bins: types.DelegationBins{
// 				"cosmos1": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(100)},
// 				"cosmos2": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(3000)},
// 				"cosmos3": types.DelegationBin{"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll": sdk.NewInt(6000)},
// 				"cosmos4": types.DelegationBin{"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewInt(3000)},
// 				"cosmos5": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(15000)},
// 				"cosmos6": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(26000)},
// 				"cosmos7": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(11000)},
// 			},
// 			validator: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
// 			expected:  "cosmos1",
// 		},
// 		{
// 			// previously seen, but first bucket is outside of threshold, use second
// 			bins: types.DelegationBins{
// 				"cosmos1": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(10000)},
// 				"cosmos2": types.DelegationBin{"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0": sdk.NewInt(3000)},
// 				"cosmos3": types.DelegationBin{"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll": sdk.NewInt(6000)},
// 				"cosmos4": types.DelegationBin{"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy": sdk.NewInt(3000)},
// 				"cosmos5": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(15000)},
// 				"cosmos6": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(26000)},
// 				"cosmos7": types.DelegationBin{"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf": sdk.NewInt(11000)},
// 			},
// 			validator: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
// 			expected:  "cosmos2",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		selected, _ := tc.bins.FindAccountForDelegation(tc.validator, sdk.NewCoin("raa", sdk.OneInt()))
// 		if tc.expected != selected {
// 			t.Errorf("Expected %s, but got %s (tc: %v)", tc.expected, selected, tc)
// 		}
// 	}

// }
