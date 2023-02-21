package keeper_test

import (
	"fmt"
	"io"
	"sort"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"golang.org/x/exp/maps"

	"github.com/ingenuity-build/quicksilver/app"
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
		wasm.EnableAllProposals,
		app.EmptyAppOptions{},
		app.GetWasmOpts(app.EmptyAppOptions{}),
		true,
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
			DelegationAddress: &types.ICAAccount{
				Address: "cosmos1ssrxxe4xsls57ehrkswlkhlkcverf0p0fpgyhzqw0hfdqj92ynxsw29r6e",
				Balance: sdk.NewCoins(
					sdk.NewCoin("qck", sdk.NewInt(100)),
					sdk.NewCoin("uqck", sdk.NewInt(700000)),
				),
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
	gotPerfAcctZone = kpr.GetZoneForPerformanceAccount(ctx, perfAcctZone.DelegationAddress.Address)
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
