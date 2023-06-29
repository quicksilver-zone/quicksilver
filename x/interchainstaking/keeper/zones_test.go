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
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func newQuicksilver(t *testing.T) *app.Quicksilver {
	t.Helper()

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
	quicksilver := newQuicksilver(t)

	chainID := "quicksilver-1"
	kpr := quicksilver.InterchainstakingKeeper
	ctx := quicksilver.NewContext(true, tmproto.Header{Height: quicksilver.LastBlockHeight()})

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
		chainID := fmt.Sprintf("%s%d", chainIDPrefix, i)
		delegationAddr := addressutils.GenerateAddressForTestWithPrefix("cosmos")
		zone := types.Zone{
			ConnectionId: "conn-test",
			ChainId:      chainID,
			LocalDenom:   "qck",
			BaseDenom:    "qck",
			DelegationAddress: &types.ICAAccount{
				Address: delegationAddr,
				Balance: sdk.NewCoins(
					sdk.NewCoin("qck", sdk.NewInt(100)),
					sdk.NewCoin("uqck", sdk.NewInt(700000)),
				),
			},
			Is_118: true,
		}
		kpr.SetAddressZoneMapping(ctx, delegationAddr, zone.ChainId)
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
	kpr.IterateZones(ctx, func(index int64, zone *types.Zone) bool {
		gotZonesMapping[index] = *zone
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
	kpr.SetAddressZoneMapping(ctx, "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", perfAcctZone.ChainId)
	kpr.SetZone(ctx, &perfAcctZone)
	gotPerfAcctZone, found := kpr.GetZoneForPerformanceAccount(ctx, perfAcctZone.PerformanceAddress.Address)
	require.True(t, found)
	require.Equal(t, &perfAcctZone, gotPerfAcctZone, "expecting a match in performance accounts")

	// Try with a non-existent performance address, it should return nil.
	gotPerfAcctZone, found = kpr.GetZoneForPerformanceAccount(ctx, "non-existent")
	require.False(t, found)
	require.Nil(t, gotPerfAcctZone, "expecting no match in the performance account")

	// Try with a non-existent performance address but that of the performance zone.
	gotPerfAcctZone, found = kpr.GetZoneForPerformanceAccount(ctx, perfAcctZone.DelegationAddress.Address)
	require.False(t, found)
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
		DelegationAddress: firstZone.DelegationAddress.Address,
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
	zone4Del, found := kpr.GetZoneForDelegateAccount(ctx, del1.DelegationAddress)
	require.True(t, found)
	require.NotNil(t, zone4Del, "expecting a non-nil zone back")
	require.Equal(t, &firstZone, zone4Del, "expectign equivalent zones")
}

// TODO: convert to keeper tests

/*func TestZone_GetBondedValidatorAddressesAsSlice(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators, &types.Validator{
		ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
		CommissionRate: sdk.MustNewDecFromStr("0.2"),
		VotingPower:    sdk.NewInt(2000),
		Status:         stakingtypes.BondStatusUnbonded,
	},
		&types.Validator{
			ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusUnbonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
	)

	// sorted list
	expected := []string{
		"cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy",
		"cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
		"cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
		"cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7",
	}
	require.Equal(t, expected, zone.GetBondedValidatorAddressesAsSlice())
}

func TestZone_GetAggregateIntentOrDefault(t *testing.T) {
	// empty
	zone := types.Zone{}
	require.Equal(t, types.ValidatorIntents(nil), zone.GetAggregateIntentOrDefault())

	zone = types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators, &types.Validator{
		ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
		CommissionRate: sdk.MustNewDecFromStr("0.2"),
		VotingPower:    sdk.NewInt(2000),
		Status:         stakingtypes.BondStatusUnbonded,
	},
		&types.Validator{
			ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusUnbonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(3000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
			CommissionRate: sdk.MustNewDecFromStr("0.2"),
			VotingPower:    sdk.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
	)

	expected := types.ValidatorIntents{
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy",
			Weight:         sdk.NewDecWithPrec(25, 2),
		},
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			Weight:         sdk.NewDecWithPrec(25, 2),
		},
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
			Weight:         sdk.NewDecWithPrec(25, 2),
		},
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7",
			Weight:         sdk.NewDecWithPrec(25, 2),
		},
	}
	actual := zone.GetAggregateIntentOrDefault()
	require.Equal(t, expected, actual)
}*/
