package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// func TestApplyDeltasToIntent(t *testing.T) {

// 	testCases := []struct {
// 		diffs       []types.Diff
// 		allocations types.Allocations
// 		output      types.Allocations
// 	}{
// 		{
// 			[]types.Diff{
// 				{Valoper: "val1", Amount: sdk.NewInt(3000)},
// 				{Valoper: "val2", Amount: sdk.NewInt(5000)},
// 				{Valoper: "val3", Amount: sdk.NewInt(9000)},
// 				{Valoper: "val4", Amount: sdk.NewInt(-16000)},
// 				{Valoper: "val5", Amount: sdk.NewInt(-1000)},
// 			},
// 			types.Allocations{}.Allocate("val5", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(200))}).Allocate("val2", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(600))}),
// 			types.Allocations{}.Allocate("val4", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(800))}),
// 		},
// 		{ // all zero, no change.
// 			[]types.Diff{
// 				{Valoper: "val1", Amount: sdk.NewInt(0)},
// 				{Valoper: "val2", Amount: sdk.NewInt(0)},
// 				{Valoper: "val3", Amount: sdk.NewInt(0)},
// 				{Valoper: "val4", Amount: sdk.NewInt(0)},
// 				{Valoper: "val5", Amount: sdk.NewInt(0)},
// 			},
// 			types.Allocations{}.Allocate("val4", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(20000))}).Allocate("val2", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(600))}),
// 			types.Allocations{}.Allocate("val4", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(20000))}).Allocate("val2", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(600))}),
// 		},
// 		{
// 			[]types.Diff{
// 				{Valoper: "val1", Amount: sdk.NewInt(26000)},
// 				{Valoper: "val2", Amount: sdk.NewInt(0)},
// 				{Valoper: "val3", Amount: sdk.NewInt(-9000)},
// 				{Valoper: "val4", Amount: sdk.NewInt(-16000)},
// 				{Valoper: "val5", Amount: sdk.NewInt(-1000)},
// 			},
// 			types.Allocations{}.Allocate("val1", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(20000))}).Allocate("val2", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(600))}),
// 			types.Allocations{}.Allocate("val4", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(16000))}).Allocate("val3", sdk.Coins{sdk.NewCoin(types.GenericToken, sdk.NewInt(4600))}),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		deltas := tc.diffs

// 		sort.SliceStable(deltas, func(i, j int) bool {
// 			return deltas[i].Amount.LT(deltas[j].Amount)
// 		})

// 		out := keeper.ApplyDeltasToIntent(tc.allocations, deltas)

// 		for _, i := range tc.output {
// 			if !out.Get(i.Address).Amount.AmountOf(types.GenericToken).Equal(i.Amount.AmountOf(types.GenericToken)) {
// 				t.Errorf("mismatch between expected tokens (%s) and actual tokens (%s)", i.Amount.AmountOf(types.GenericToken), out.Get(i.Address).Amount.AmountOf(types.GenericToken))
// 			}
// 		}
// 	}
// }

func generateTestBins() types.Allocations {
	return types.Allocations{}.
		Allocate("del1", sdk.NewCoins(sdk.NewCoin("val1", sdk.NewInt(83333)))).
		Allocate("del2", sdk.NewCoins(sdk.NewCoin("val2", sdk.NewInt(83333)))).
		Allocate("del3", sdk.NewCoins(sdk.NewCoin("val2", sdk.NewInt(300000)))).
		Allocate("del3", sdk.NewCoins(sdk.NewCoin("val3", sdk.NewInt(300000)))).
		Allocate("del3", sdk.NewCoins(sdk.NewCoin("val4", sdk.NewInt(400000)))).
		Allocate("del4", sdk.NewCoins(sdk.NewCoin("val2", sdk.NewInt(50000))))
}

func generateIntents() types.ValidatorIntents {
	return types.ValidatorIntents{
		"val1": {ValoperAddress: "val1", Weight: sdk.NewDecWithPrec(3, 1)},
		"val2": {ValoperAddress: "val2", Weight: sdk.NewDecWithPrec(25, 2)},
		"val3": {ValoperAddress: "val3", Weight: sdk.NewDecWithPrec(35, 2)},
		"val4": {ValoperAddress: "val4", Weight: sdk.NewDecWithPrec(1, 1)},
	}
}

func TestDeltasAndIntents(t *testing.T) {
	requests := types.Allocations{}.Allocate("val4", sdk.Coins{sdk.Coin{Denom: types.GenericToken, Amount: sdk.NewInt(900000)}})

	bins := generateTestBins()
	deltas := types.DetermineIntentDelta(bins, bins.SumAll(), generateIntents())
	fmt.Println(deltas)

	requests = keeper.ApplyDeltasToIntent(requests, deltas, bins)
	for _, i := range requests {
		fmt.Println(i)
	}
}
