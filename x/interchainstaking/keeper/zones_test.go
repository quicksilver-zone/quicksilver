package keeper_test

import (
	"fmt"
	"io"
	"sort"
	"testing"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/v7/app"
	"github.com/quicksilver-zone/quicksilver/v7/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
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
		app.EmptyAppOptions{},
		true,
		false,
		app.GetWasmOpts(app.EmptyAppOptions{}),
	)
}

func TestKeeperWithZonesRoundTrip(t *testing.T) {
	quicksilver := newQuicksilver(t)

	chainID := "quicksilver-1"
	kpr := quicksilver.InterchainstakingKeeper
	ctx := quicksilver.NewContext(true)

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
					sdk.NewCoin("qck", sdkmath.NewInt(100)),
					sdk.NewCoin("uqck", sdkmath.NewInt(700000)),
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
			sdk.NewCoin("qck", sdkmath.NewInt(800)),
			sdk.NewCoin("uqck", sdkmath.NewInt(900000)),
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
	zeroDelAmt := sdk.NewCoin(firstZone.BaseDenom, sdkmath.NewInt(0))
	require.Equal(t, zeroDelAmt, gotDelAmt, "expecting 0")

	// 7.2. Set some delegations.
	del1 := types.Delegation{
		Amount:            sdk.NewCoin(firstZone.BaseDenom, sdkmath.NewInt(17000)),
		DelegationAddress: firstZone.DelegationAddress.Address,
		Height:            10,
		ValidatorAddress:  "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
	}

	kpr.SetDelegation(ctx, firstZone.ChainId, del1)

	// 7.3. Retrieve the delegation now, it should be set.
	gotDelAmt = kpr.GetDelegatedAmount(ctx, &firstZone)
	require.NotEqual(t, zeroDelAmt, gotDelAmt, "expecting a match in delegation amounts")
	wantDelAmt := sdk.NewCoin(firstZone.BaseDenom, sdkmath.NewInt(17000))
	require.Equal(t, wantDelAmt, gotDelAmt, "expecting 17000 as the delegation amount")

	// Zone for delegation account.
	zone4Del, found := kpr.GetZoneForDelegateAccount(ctx, del1.DelegationAddress)
	require.True(t, found)
	require.NotNil(t, zone4Del, "expecting a non-nil zone back")
	require.Equal(t, &firstZone, zone4Del, "expectign equivalent zones")
}

func (suite *KeeperTestSuite) TestRemoveZoneAndAssociatedRecords() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	chainID := "quicksilver-1"

	// Set zone
	quicksilver.InterchainstakingKeeper.SetZone(ctx, &types.Zone{
		ConnectionId: "connection-test",
		ChainId:      chainID,
		LocalDenom:   "uqck",
		BaseDenom:    "qck",
		DelegationAddress: &types.ICAAccount{
			Address: addressutils.GenerateAddressForTestWithPrefix("quicksilver"),
		},
		PerformanceAddress: &types.ICAAccount{
			Address: addressutils.GenerateAddressForTestWithPrefix("quicksilver"),
		},
	})
	// Check set zone
	zone, ok := quicksilver.InterchainstakingKeeper.GetZone(ctx, chainID)
	suite.True(ok, "expected to retrieve a zone")
	suite.NotEqual(types.Zone{}, zone, "Expecting a non-blank zone")

	// set val
	val0 := types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdkmath.LegacyMustNewDecFromStr("1"), VotingPower: sdkmath.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err := quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val0)
	suite.NoError(err)

	val1 := types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdkmath.LegacyMustNewDecFromStr("1"), VotingPower: sdkmath.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val1)
	suite.NoError(err)

	val2 := types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdkmath.LegacyMustNewDecFromStr("1"), VotingPower: sdkmath.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val2)
	suite.NoError(err)

	val3 := types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdkmath.LegacyMustNewDecFromStr("1"), VotingPower: sdkmath.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val3)
	suite.NoError(err)

	vals := quicksilver.InterchainstakingKeeper.GetValidators(ctx, chainID)
	// create unbonding
	quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx,
		types.UnbondingRecord{
			ChainId:       chainID,
			EpochNumber:   1,
			Validator:     vals[0].ValoperAddress,
			RelatedTxhash: []string{"ABC012"},
		})
	// create redelegations
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx,
		types.RedelegationRecord{
			ChainId:     chainID,
			EpochNumber: 1,
			Source:      vals[1].ValoperAddress,
			Destination: vals[0].ValoperAddress,
			Amount:      10000000,
		})
	// create delegation
	delegation := types.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  vals[1].ValoperAddress,
		Amount:            sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(1000)),
	}
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
	// create pert delegation
	performanceAddress := zone.PerformanceAddress
	perfDelegation := types.Delegation{
		DelegationAddress: performanceAddress.Address,
		ValidatorAddress:  vals[1].ValoperAddress,
		Amount:            sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(1000)),
	}
	quicksilver.InterchainstakingKeeper.SetPerformanceDelegation(ctx, zone.ChainId, perfDelegation)
	// create receipt
	cutOffTime := ctx.BlockTime().AddDate(0, 0, -1).Add(-2 * time.Hour)
	rcpt := types.Receipt{
		ChainId: chainID,
		Sender:  addressutils.GenerateAccAddressForTest().String(),
		Txhash:  "TestDeposit01",
		Amount: sdk.NewCoins(
			sdk.NewCoin(
				zone.BaseDenom,
				sdkmath.NewIntFromUint64(2000000),
			),
		),
		FirstSeen: &cutOffTime,
	}
	quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt)
	// create withdrawal record
	record := types.WithdrawalRecord{
		ChainId:   chainID,
		Delegator: zone.DelegationAddress.Address,
		Distribution: []*types.Distribution{
			{
				Valoper: vals[1].ValoperAddress,
				Amount:  1000000,
			},
			{
				Valoper: vals[2].ValoperAddress,
				Amount:  1000000,
			},
		},
		Recipient:      addressutils.GenerateAccAddressForTest().String(),
		Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(4000000))),
		BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdkmath.NewInt(4000000)),
		Txhash:         "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
		Status:         types.WithdrawStatusUnbond,
		CompletionTime: time.Now().UTC().Add(time.Hour),
	}
	quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)

	// Handle
	quicksilver.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, chainID)

	// check unbondings
	_, found := quicksilver.InterchainstakingKeeper.GetUnbondingRecord(ctx, chainID, vals[0].ValoperAddress, 1)
	suite.False(found, "Not found unbonding record stored in the keeper")

	// check redelegations
	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, chainID, vals[1].ValoperAddress, vals[0].ValoperAddress, 1)
	suite.False(found, "Not found redelegation record stored in the keeper")

	// check delegation
	_, found = quicksilver.InterchainstakingKeeper.GetDelegation(ctx, chainID, delegation.DelegationAddress, delegation.ValidatorAddress)
	suite.False(found, "Not found delegation stored in the keeper")

	// check pert delegation
	_, found = quicksilver.InterchainstakingKeeper.GetPerformanceDelegation(ctx, chainID, performanceAddress, perfDelegation.ValidatorAddress)
	suite.False(found, "Not found pert delegation stored in the keeper")

	// check receipts
	_, found = quicksilver.InterchainstakingKeeper.GetReceipt(ctx, chainID, rcpt.Txhash)
	suite.False(found, "Not found receipts stored in the keeper")

	// check withdrawal records
	_, found = quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, chainID, record.Txhash, record.Status)
	suite.False(found, "Not found withdrawal records stored in the keeper")

	// check validators
	vals = quicksilver.InterchainstakingKeeper.GetValidators(ctx, chainID)
	suite.Equal(len(vals), 0)

	// check zone
	zone, found = quicksilver.InterchainstakingKeeper.GetZone(ctx, chainID)
	suite.False(found, "No zone stored in the keeper")
	suite.Equal(types.Zone{}, zone, "Expecting the blank zone")
}

// TODO: convert to keeper tests

/* func TestZone_GetBondedValidatorAddressesAsSlice(t *testing.T) {
	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators, &types.Validator{
		ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
		CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
		VotingPower:    sdkmath.NewInt(2000),
		Status:         stakingtypes.BondStatusUnbonded,
	},
		&types.Validator{
			ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusUnbonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
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
		CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
		VotingPower:    sdkmath.NewInt(2000),
		Status:         stakingtypes.BondStatusUnbonded,
	},
		&types.Validator{
			ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusUnbonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(3000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
		&types.Validator{
			ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
			CommissionRate: sdkmath.LegacyMustNewDecFromStr("0.2"),
			VotingPower:    sdkmath.NewInt(2000),
			Status:         stakingtypes.BondStatusBonded,
		},
	)

	expected := types.ValidatorIntents{
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy",
			Weight:         sdkmath.LegacyNewDecWithPrec(25, 2),
		},
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			Weight:         sdkmath.LegacyNewDecWithPrec(25, 2),
		},
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
			Weight:         sdkmath.LegacyNewDecWithPrec(25, 2),
		},
		&types.ValidatorIntent{
			ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7",
			Weight:         sdkmath.LegacyNewDecWithPrec(25, 2),
		},
	}
	actual := zone.GetAggregateIntentOrDefault()
	require.Equal(t, expected, actual)
} */
