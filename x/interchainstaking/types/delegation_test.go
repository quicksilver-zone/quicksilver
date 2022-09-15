package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

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

func TestAllocationsMethods(t *testing.T) {
	allocs := types.Allocations{
		{
			"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("qck", sdk.NewInt(300)),
				sdk.NewCoin("uqck", sdk.NewInt(400000)),
			),
		},
		{
			"cosmosvaloper1pyfmqnramtg7ewxqwwrwyxfgc4n5ef9p2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("qck", sdk.NewInt(100)),
				sdk.NewCoin("uqck", sdk.NewInt(700000)),
			),
		},
		{
			"cosmosvaloper1ajllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("mqck", sdk.NewInt(300)),
				sdk.NewCoin("pqck", sdk.NewInt(400000)),
			),
		},
	}

	// 1. Check that the sums for various denominations are correct.
	qckSum := allocs.SumForDenom("qck")
	wantQCKSum := sdk.NewInt(300 + 100)
	require.Equal(t, wantQCKSum, qckSum, "qck sum mismatch")
	uqckSum := allocs.SumForDenom("uqck")
	wantuQCKSum := sdk.NewInt(400000 + 700000)
	require.Equal(t, wantuQCKSum, uqckSum, "uqck sum mismatch")

	// 2. Test that sorted returns values sorted by addresses.
	unsortedAllocs := make(types.Allocations, len(allocs))
	copy(unsortedAllocs, allocs)
	gotAllocsSortedByAddrs := unsortedAllocs.Sorted()
	require.NotEqual(t, allocs, gotAllocsSortedByAddrs, "allocs sorted by addrs mismatch")
	wantAllocsSortedByAddrs := types.Allocations{
		{
			"cosmosvaloper1ajllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("mqck", sdk.NewInt(300)),
				sdk.NewCoin("pqck", sdk.NewInt(400000)),
			),
		},
		{
			"cosmosvaloper1pyfmqnramtg7ewxqwwrwyxfgc4n5ef9p2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("qck", sdk.NewInt(100)),
				sdk.NewCoin("uqck", sdk.NewInt(700000)),
			),
		},
		{
			"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("qck", sdk.NewInt(300)),
				sdk.NewCoin("uqck", sdk.NewInt(400000)),
			),
		},
	}
	require.Equal(t, wantAllocsSortedByAddrs, gotAllocsSortedByAddrs, "allocs sorted by addrs mismatch")

	// 3. Test that sortedByAmount sorts values firstly by addr then by value.
	copy(unsortedAllocs, allocs)
	gotAllocsSortedByAmount := unsortedAllocs.SortedByAmount()
	require.NotEqual(t, allocs, gotAllocsSortedByAmount, "allocs sorted by amount mismatch")
	wantAllocsSortedByAmount := types.Allocations{
		{
			"cosmosvaloper1ajllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("mqck", sdk.NewInt(300)),
				sdk.NewCoin("pqck", sdk.NewInt(400000)),
			),
		},
		{
			"cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("qck", sdk.NewInt(300)),
				sdk.NewCoin("uqck", sdk.NewInt(400000)),
			),
		},
		{
			"cosmosvaloper1pyfmqnramtg7ewxqwwrwyxfgc4n5ef9p2lcnj0",
			sdk.NewCoins(
				sdk.NewCoin("qck", sdk.NewInt(100)),
				sdk.NewCoin("uqck", sdk.NewInt(700000)),
			),
		},
	}
	require.Equal(t, wantAllocsSortedByAmount, gotAllocsSortedByAmount, "allocs sorted by amount mismatch")

	// 4. Evaluate the total sum but bucketed by denominations.
	wantSum := sdk.Coins{
		sdk.NewCoin("mqck", sdk.NewInt(300)),
		sdk.NewCoin("pqck", sdk.NewInt(400000)),
		sdk.NewCoin("qck", sdk.NewInt(400)),
		sdk.NewCoin("uqck", sdk.NewInt(1100000)),
	}
	gotSum := gotAllocsSortedByAmount.Sum()
	require.Equal(t, wantSum, gotSum, "mismatch in sum bucketed by denomination")

	// 5. Evalute the total sum regardless of sum.
	gotSumAll := allocs.SumAll()
	wantSumAll := sdk.NewInt(100 + 300 + 300 + 400000 + 400000 + 700000)
	require.Equal(t, wantSumAll, gotSumAll, "mismatch in sumAll values")

	// 6. Evaluate threshold.
	gotThreshold := gotAllocsSortedByAmount.DetermineThreshold()
	// 6.1. Find the 1/3 item with all its sums, which would be the 0th item in wantAllocsSortedByAmount
	// and then its .SumAll value of: 300 + 400000.
	wantThreshold := sdk.NewInt(300 + 400000)
	require.Equal(t, wantThreshold, gotThreshold, "mismatch in threshold values")

	// 7. Evaluate smallest bin.
	gotSBin1 := gotAllocsSortedByAmount.SmallestBin()
	wantSBin := *gotAllocsSortedByAmount[0]
	require.Equal(t, wantSBin, gotSBin1, "mismatch in smallest bin values")

	// 7.5. Smallest bin value for an empty allocations object, shoould not panic.
	got0SBin := (types.Allocations{}).SmallestBin()
	want0SBin := types.Allocation{}
	require.Equal(t, want0SBin, got0SBin, "mismatch in smallest bin values for empty allocations object")
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
