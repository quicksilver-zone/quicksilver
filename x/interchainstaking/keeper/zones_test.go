package keeper_test

import (
	"fmt"
	"io"
	"sort"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"golang.org/x/exp/maps"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func newQuicksilver(t *testing.T) *app.Quicksilver {
	return app.NewQuicksilver(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		io.Discard,
		true,
		map[int64]bool{},
		t.TempDir(),
		5,
		app.MakeEncodingConfig(),
		simapp.EmptyAppOptions{},
		app.GetWasmOpts(simapp.EmptyAppOptions{}),
	)
}

func TestKeeperWithZonesRoundTrip(t *testing.T) {
	app := newQuicksilver(t)

	chainID := "quicksilver-1"
	kpr := app.InterchainstakingKeeper
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// 1. Check for a zone without having stored anything.
	zone, ok := kpr.GetZone(ctx, chainID)
	require.False(t, ok, "No zone stored in the keeper")
	require.Equal(t, types.Zone{}, zone, "Expecting the blank zone")

	// 2. Now set a zone and ensure it is retrieved.
	zone = types.Zone{
		ConnectionId: "conn-test",
		ChainId:      chainID,
		LocalDenom:   "uqck",
		BaseDenom:    "qck",
	}
	kpr.SetZone(ctx, &zone)
	gotZone, ok := kpr.GetZone(ctx, chainID)
	require.True(t, ok, "expected to retrieve a zone")
	require.NotEqual(t, types.Zone{}, gotZone, "Expecting a non-blank zone")
	require.Equal(t, zone, gotZone, "Expecting the stored zone")

	// 3. Delete the zone then try to retrieve it.
	kpr.DeleteZone(ctx, chainID)
	zone, ok = kpr.GetZone(ctx, chainID)
	require.False(t, ok, "No zone stored anymore in the keeper")
	require.Equal(t, types.Zone{}, zone, "Expecting the blank zone")

	// 4. Store many zones in the keeper, then retrieve them by chainID.
	nzones := 10
	chainIDPrefix := "quicksilver-"
	indexToZone := make(map[int64]types.Zone, nzones)
	for i := 0; i < nzones; i++ {
		chainID := fmt.Sprintf("%s-%d", chainIDPrefix, i)
		zone := types.Zone{
			ConnectionId: "conn-test",
			ChainId:      chainID,
			LocalDenom:   "qck",
			BaseDenom:    "qck",
			DelegationAddresses: []*types.ICAAccount{
				{
					Address: "cosmos1ssrxxe4xsls57ehrkswlkhlkcverf0p0fpgyhzqw0hfdqj92ynxsw29r6e",
					Balance: sdk.NewCoins(
						sdk.NewCoin("qck", sdk.NewInt(100)),
						sdk.NewCoin("uqck", sdk.NewInt(700000)),
					),
				},
			},
		}
		kpr.SetZone(ctx, &zone)
		gotZone, ok := kpr.GetZone(ctx, chainID)
		require.True(t, ok, "expected to retrieve the correct zone")
		require.NotEqual(t, types.Zone{}, gotZone, "Expecting a non-blank zone")
		require.Equal(t, zone, gotZone, "Expecting the stored zone")

		// Save the zone for later comparisons.
		indexToZone[int64(i)] = zone
	}
	require.Equal(t, nzones, len(indexToZone), "expecting unique nzones")

	// 5.1. Invoke Iterate on zones.
	gotZonesMapping := make(map[int64]types.Zone, nzones)
	kpr.IterateZones(ctx, func(index int64, zone types.Zone) bool {
		gotZonesMapping[index] = zone
		return false
	})

	require.Equal(t, nzones, len(gotZonesMapping), "expecting unique nzones")
	require.Equal(t, indexToZone, gotZonesMapping, "expecting the zones mapped by index to match")

	// 5.2. List all zones.
	gotAllZones := kpr.AllZones(ctx)
	wantAllZones := maps.Values(gotZonesMapping)
	require.Equal(t, nzones, len(gotAllZones), "expecting unique nzones")
	// Sort them for determinism
	sort.Slice(gotAllZones, func(i, j int) bool {
		zi, zj := gotAllZones[i], gotAllZones[j]
		return zi.ChainId < zj.ChainId
	})
	sort.Slice(wantAllZones, func(i, j int) bool {
		zi, zj := wantAllZones[i], wantAllZones[j]
		return zi.ChainId < zj.ChainId
	})
	require.Equal(t, wantAllZones, gotAllZones, "expecting the zones to match")

	// 6. Test performance accounts.
	perfAcctZone := indexToZone[4]
	perfAcctZone.PerformanceAddress = &types.ICAAccount{
		Address: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
		Balance: sdk.NewCoins(
			sdk.NewCoin("qck", sdk.NewInt(800)),
			sdk.NewCoin("uqck", sdk.NewInt(900000)),
		),
	}
	kpr.SetZone(ctx, &perfAcctZone)
	gotPerfAcctZone := kpr.GetZoneForPerformanceAccount(ctx, perfAcctZone.PerformanceAddress.Address)
	require.Equal(t, &perfAcctZone, gotPerfAcctZone, "expecting a match in performance accounts")
	// Try with a non-existent performance address, it should return nil.
	gotPerfAcctZone = kpr.GetZoneForPerformanceAccount(ctx, "non-existent")
	require.Nil(t, gotPerfAcctZone, "expecting no match in the performance account")
	// Try with a non-existent performance address but that of the performance zone.
	gotPerfAcctZone = kpr.GetZoneForPerformanceAccount(ctx, perfAcctZone.DelegationAddresses[0].Address)
	require.Nil(t, gotPerfAcctZone, "expecting no match in the performance account")

	// 7.1. Test delegated amounts.
	firstZone := gotAllZones[0]
	gotDelAmt := kpr.GetDelegatedAmount(ctx, &firstZone)
	// No delegations set so nothing expected.
	zeroDelAmt := sdk.NewCoin(firstZone.BaseDenom, sdk.NewInt(0))
	require.Equal(t, zeroDelAmt, gotDelAmt, "expecting 0")

	// 7.2. Set some delegations.
	del1 := types.Delegation{
		Amount:            sdk.NewCoin(firstZone.BaseDenom, sdk.NewInt(17000)),
		DelegationAddress: "cosmos1ssrxxe4xsls57ehrkswlkhlkcverf0p0fpgyhzqw0hfdqj92ynxsw29r6e",
		Height:            10,
		ValidatorAddress:  "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
	}
	kpr.SetDelegation(ctx, &firstZone, del1)

	// 7.3. Retrieve the delegation now, it should be set.
	gotDelAmt = kpr.GetDelegatedAmount(ctx, &firstZone)
	require.NotEqual(t, zeroDelAmt, gotDelAmt, "expecting a match in delegation amounts")
	wantDelAmt := sdk.NewCoin(firstZone.BaseDenom, sdk.NewInt(17000))
	require.Equal(t, wantDelAmt, gotDelAmt, "expecting 17000 as the delegation amount")

	// Zone for delegation account.
	zone4Del := kpr.GetZoneForDelegateAccount(ctx, del1.DelegationAddress)
	require.NotNil(t, zone4Del, "expecting a non-nil zone back")
	require.Equal(t, &firstZone, zone4Del, "expectign equivalent zones")
}

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
