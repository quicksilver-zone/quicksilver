package keeper_test

import (
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

var (
	user1 = addressutils.GenerateAccAddressForTest()
	user2 = addressutils.GenerateAccAddressForTest()
)

func (suite *KeeperTestSuite) TestKeeper_IntentStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)
	zoneValidatorAddresses := icsKeeper.GetValidators(ctx, zone.ChainId)

	// check that there are no intents
	intents := icsKeeper.AllDelegatorIntents(ctx, &zone, false)
	suite.Len(intents, 0)

	// set intents for testAddress
	icsKeeper.SetDelegatorIntent(
		ctx,
		&zone,
		icstypes.DelegatorIntent{
			Delegator: testAddress,
			Intents: icstypes.ValidatorIntents{
				{
					ValoperAddress: zoneValidatorAddresses[0].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[1].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[2].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[3].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
			},
		},
		false,
	)
	// set intents for user1
	icsKeeper.SetDelegatorIntent(
		ctx,
		&zone,
		icstypes.DelegatorIntent{
			Delegator: user1.String(),
			Intents: icstypes.ValidatorIntents{
				{
					ValoperAddress: zoneValidatorAddresses[0].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[1].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[2].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[3].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
			},
		},
		false,
	)
	// set intents for user2
	icsKeeper.SetDelegatorIntent(
		ctx,
		&zone,
		icstypes.DelegatorIntent{
			Delegator: user2.String(),
			Intents: icstypes.ValidatorIntents{
				{
					ValoperAddress: zoneValidatorAddresses[0].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.5"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[1].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.3"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[2].ValoperAddress,
					Weight:         sdk.MustNewDecFromStr("0.2"),
				},
			},
		},
		false,
	)

	// check for intents set above
	intents = icsKeeper.AllDelegatorIntents(ctx, &zone, false)
	suite.Len(intents, 3)

	// delete intent for testAddress
	icsKeeper.DeleteDelegatorIntent(ctx, &zone, testAddress, false)

	// check intents
	intents = icsKeeper.AllDelegatorIntents(ctx, &zone, false)
	suite.Len(intents, 2)

	suite.T().Logf("intents:\n%+v\n", intents)

	// update intent for user1
	err := icsKeeper.UpdateDelegatorIntent(
		ctx,
		user2,
		&zone,
		sdk.NewCoins(
			sdk.NewCoin(
				zone.BaseDenom,
				sdkmath.NewInt(7313913),
			),
		),
		nil,
	)
	suite.NoError(err)

	// load and match pointers
	intentsPointers := icsKeeper.AllDelegatorIntentsAsPointer(ctx, &zone, false)
	for i, ip := range intentsPointers {
		suite.Equal(intents[i], *ip)
	}

	suite.T().Logf("intents:\n%+v\n", intentsPointers)
}

func (suite *KeeperTestSuite) TestAggregateIntent() {
	tc := []struct {
		name     string
		intents  func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent
		balances func() map[string]int64
		expected func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents
	}{
		{
			name: "empty intents; returns equal weighting",
			intents: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				return out
			},
			balances: func() map[string]int64 { return map[string]int64{} },
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId) {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			name: "single intent; zero balance, returns default equal weighting",
			intents: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.OneDec()}}})
				return out
			},
			balances: func() map[string]int64 { return map[string]int64{} },
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId) {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			// name: "single intent; with balance, returns single weighting",
			name: "single intent; with balance, returns default equal weighting",
			intents: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.OneDec()}}})
				return out
			},
			balances: func() map[string]int64 {
				return map[string]int64{
					user1.String(): 1,
				}
			},
			// expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
			// 	out := icstypes.ValidatorIntents{}
			// 	out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})

			// 	return out.Sort()
			// },
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId) {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}

				return out.Sort()
			},
		},
		{
			// name: "two intents; with equal balances, same val, single weighting",
			name: "two intents; with equal balances, same val, returns default equal weighting",
			intents: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out,
					icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.OneDec()}}},
					icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.OneDec()}}},
				)
				return out
			},
			balances: func() map[string]int64 {
				return map[string]int64{
					user1.String(): 1,
					user2.String(): 1,
				}
			},
			// expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
			// 	out := icstypes.ValidatorIntents{}
			// 	out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})

			// 	return out.Sort()
			// },
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId) {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			// name: "two intents; with equal balances, diff val, equal weighting",
			name: "two intents; with equal balances, diff val, returns default equal weighting",
			intents: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out,
					icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.OneDec()}}},
					icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[1], Weight: sdk.OneDec()}}},
				)
				return out
			},
			balances: func() map[string]int64 {
				return map[string]int64{
					user1.String(): 1,
					user2.String(): 1,
				}
			},
			// expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
			// 	out := icstypes.ValidatorIntents{}
			// 	out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})
			// 	out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})

			// 	return out.Sort()
			// },
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId) {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			name: "two intents; with zer0 balances, diff val, returns default equal weights ",
			intents: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out,
					icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.ZeroDec()}}},
					icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[1], Weight: sdk.OneDec()}}},
				)
				return out
			},
			balances: func() map[string]int64 {
				return map[string]int64{
					user1.String(): 0,
					user2.String(): 0,
				}
			},
			// expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
			// 	out := icstypes.ValidatorIntents{}
			// 	out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})
			// 	out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})

			// 	return out.Sort()
			// },
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId) {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
	}

	for _, tt := range tc {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			// give each user some funds
			for addrString, balance := range tt.balances() {
				suite.giveFunds(ctx, zone.LocalDenom, balance, addrString)
			}

			for _, intent := range tt.intents(ctx, quicksilver, zone) {
				icsKeeper.SetDelegatorIntent(ctx, &zone, intent, false)
			}

			// If the supply is zero, mint some coins to avoid zero ordializedSum
			if quicksilver.BankKeeper.GetSupply(ctx, zone.LocalDenom).IsZero() {
				err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1000))))
				suite.NoError(err)
			}
			err := icsKeeper.AggregateDelegatorIntents(ctx, &zone)
			suite.NoError(err)

			// refresh zone to pull new aggregate
			zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			actual, err := icsKeeper.GetAggregateIntentOrDefault(ctx, &zone)
			suite.NoError(err)
			suite.Equal(tt.expected(ctx, quicksilver, zone), actual)
		})
	}
}

func (suite *KeeperTestSuite) TestAggregateIntentWithPRClaims() {
	tc := []struct {
		name           string
		intents        func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent
		claims         func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim
		unclaimedRatio sdkmath.Int
		expected       func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents
	}{
		{
			name: "single intent and claim, claims equals 100% supply, returns single weighting",
			intents: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.OneDec()}}})
				return out
			},
			claims: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.OneInt(), SourceChainId: "osmosis-1"},
				}
			},
			unclaimedRatio: sdkmath.ZeroInt(),
			expected: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.OneDec()})
				return out.Sort()
			},
		},
		{
			name: "one single-val intent, multiple claims, 100% supply, returns single weighting",
			intents: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.OneDec()}}})
				return out
			},
			claims: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(1000), SourceChainId: "osmosis-1"},
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(1000), SourceChainId: "osmosis-1"},
				}
			},
			unclaimedRatio: sdkmath.ZeroInt(),
			expected: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.OneDec()})
				return out.Sort()
			},
		},
		{
			name: "two single-val intents, and one claim, 100% supply, returns equal weighting",
			intents: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.OneDec()}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.OneDec()}}})
				return out
			},
			claims: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(1000), SourceChainId: "osmosis-1"},
				}
			},
			unclaimedRatio: sdkmath.ZeroInt(),
			expected: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.OneDec()})
				// only remote assets are considered, thus user2 balance is ignored...
				// out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.OneDec().Quo(sdk.NewDec(2))})
				return out.Sort()
			},
		},
		{
			name: "two single-val intents, one claim, 50% supply, returns equal weighting for unintended validators",
			intents: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(1)}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(1)}}})
				return out
			},
			claims: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(500), SourceChainId: "osmosis-1"},
				}
			},
			unclaimedRatio: sdkmath.NewInt(50),
			expected: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(5).Quo(sdk.NewDec(8))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(8))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[2].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(8))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[3].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(8))})
				return out.Sort()
			},
		},
		{
			name: "two single-val intents, two claim, 100% supply, returns equal weighting for unintended validators",
			intents: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(1)}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(1)}}})
				return out
			},
			claims: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(500), SourceChainId: "osmosis-1"},
					user2.String(): {UserAddress: user2.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(500), SourceChainId: "osmosis-1"},
				}
			},
			unclaimedRatio: sdkmath.NewInt(50),
			expected: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(3).Quo(sdk.NewDec(8))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(3).Quo(sdk.NewDec(8))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[2].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(8))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[3].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(8))})
				return out.Sort()
			},
		},
		{
			name: "two single-val intents, two claims, 50% supply, returns equal weighting for unintended validators",
			intents: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(1)}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(1)}}})
				return out
			},
			claims: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(500), SourceChainId: "osmosis-1"},
					user2.String(): {UserAddress: user2.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(1000), SourceChainId: "osmosis-1"},
				}
			},
			unclaimedRatio: sdkmath.NewInt(50),
			expected: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(7).Quo(sdk.NewDec(24))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(11).Quo(sdk.NewDec(24))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[2].ValoperAddress, Weight: sdk.NewDec(3).Quo(sdk.NewDec(24))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[3].ValoperAddress, Weight: sdk.NewDec(3).Quo(sdk.NewDec(24))})
				return out.Sort()
			},
		},
		{
			name: "two multi-val intents, two claims, 50% supply, returns equal weighting for unintended validators",
			intents: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				// user 1 split intent between two first validators
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(2))}, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(2))}}})
				// user 2 split intent between three first validators
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(3))}, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(3))}, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[2].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(3))}}})
				return out
			},
			claims: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					// user 1 claims 500, which is 1/6 of the total supply
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(500), SourceChainId: "osmosis-1"},
					// user 2 claims 1000, which is 1/3 of the total supply
					user2.String(): {UserAddress: user2.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: sdkmath.NewInt(1000), SourceChainId: "osmosis-1"},
				}
			},
			unclaimedRatio: sdkmath.NewInt(50),
			expected: func(ctx sdk.Context, app *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// aggregated_intent = sum(claim_weight * intent_weight_vector) + unclaimed_weight*default_weight_vector
				// claim_weight = claim_amount / total_supply
				// unclaimed_weight = unclaimed_supply / total_supply
				// 1/6*(1/2, 1/2,0,0) + 1/3*(1/3,1/3,1/3,0) + 1/2(1/4,1/4,1/4,1/4) = (23/72, 23/72, 17/72, 1/8)
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[0].ValoperAddress, Weight: sdk.NewDec(23).Quo(sdk.NewDec(72))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[1].ValoperAddress, Weight: sdk.NewDec(23).Quo(sdk.NewDec(72))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[2].ValoperAddress, Weight: sdk.NewDec(17).Quo(sdk.NewDec(72))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)[3].ValoperAddress, Weight: sdk.NewDec(1).Quo(sdk.NewDec(8))})
				return out.Sort()
			},
		},
	}

	for _, tt := range tc {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			// fund users based on claims
			for _, claim := range tt.claims(ctx, quicksilver, zone) {
				userAddr := claim.UserAddress
				amount := sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, claim.Amount))

				err := quicksilver.MintKeeper.MintCoins(ctx, amount)
				suite.NoError(err)
				addr, err := addressutils.AccAddressFromBech32(userAddr, zone.AccountPrefix)
				suite.NoError(err)
				err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amount)
				suite.NoError(err)
			}

			// add additional supply
			currSupply := quicksilver.BankKeeper.GetSupply(ctx, zone.LocalDenom)
			additionalAmount := currSupply.Amount.Mul(tt.unclaimedRatio).Quo(sdk.NewInt(100).Sub(tt.unclaimedRatio))
			err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, additionalAmount)))
			suite.NoError(err)

			for _, intent := range tt.intents(ctx, quicksilver, zone) {
				icsKeeper.SetDelegatorIntent(ctx, &zone, intent, false)
			}

			for _, claim := range tt.claims(ctx, quicksilver, zone) {
				claim := claim
				quicksilver.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &claim)
			}

			err = icsKeeper.AggregateDelegatorIntents(ctx, &zone)
			suite.NoError(err)

			// refresh zone to pull new aggregate
			zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			actual, err := icsKeeper.GetAggregateIntentOrDefault(ctx, &zone)
			suite.NoError(err)
			suite.Equal(tt.expected(ctx, quicksilver, zone), actual)
		})
	}
}

func (suite *KeeperTestSuite) TestDefaultIntent() {
	suite.SetupTest()
	suite.setupTestZones()
	t := suite.T()
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	zone := icstypes.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators,
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	// Need to set validators in the store to calculate the aggregate intent properly
	for _, val := range zone.Validators {
		err := icsKeeper.SetValidator(ctx, zone.ChainId, *val)
		require.NoError(t, err)
	}

	out := icsKeeper.DefaultAggregateIntents(ctx, zone.ChainId)
	require.Equal(t, len(out), 5)
	for _, v := range out {
		if !v.Weight.Equal(sdk.NewDecWithPrec(2, 1)) {
			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(2, 1), v.Weight)
		}
	}
}

func (suite *KeeperTestSuite) TestDefaultIntentWithJailed() {
	suite.SetupTest()
	suite.setupTestZones()
	t := suite.T()
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	zone := icstypes.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators,
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded, Jailed: true},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	// Need to set validators in the store to calculate the aggregate intent properly
	for _, val := range zone.Validators {
		err := icsKeeper.SetValidator(ctx, zone.ChainId, *val)
		require.NoError(t, err)
	}
	out := icsKeeper.DefaultAggregateIntents(ctx, zone.ChainId)
	require.Equal(t, len(out), 4)

	for _, v := range out {
		if !v.Weight.Equal(sdk.NewDecWithPrec(25, 2)) {
			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(25, 2), v.Weight)
		}
	}
}

func (suite *KeeperTestSuite) TestDefaultIntentWithTombstoned() {
	suite.SetupTest()
	suite.setupTestZones()
	t := suite.T()
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	zone := icstypes.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators,
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded, Tombstoned: true},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	// Need to set validators in the store to calculate the aggregate intent properly
	for _, val := range zone.Validators {
		err := icsKeeper.SetValidator(ctx, zone.ChainId, *val)
		suite.NoError(err)
	}
	out := icsKeeper.DefaultAggregateIntents(ctx, zone.ChainId)
	require.Equal(t, len(out), 4)

	for _, v := range out {
		if !v.Weight.Equal(sdk.NewDecWithPrec(25, 2)) {
			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(25, 2), v.Weight)
		}
	}
}

func (suite *KeeperTestSuite) TestDefaultIntentWithCommission100() {
	suite.SetupTest()
	suite.setupTestZones()
	t := suite.T()
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	zone := icstypes.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators,
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	// Need to set validators in the store to calculate the aggregate intent properly
	for _, val := range zone.Validators {
		err := icsKeeper.SetValidator(ctx, zone.ChainId, *val)
		require.NoError(t, err)
	}
	out := icsKeeper.DefaultAggregateIntents(ctx, zone.ChainId)
	require.Equal(t, len(out), 4)

	for _, v := range out {
		if !v.Weight.Equal(sdk.NewDecWithPrec(25, 2)) {
			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(25, 2), v.Weight)
		}
	}
}

func (suite *KeeperTestSuite) TestDefaultIntentWithOneUnbondedOneUnbonding() {
	suite.SetupTest()
	suite.setupTestZones()
	t := suite.T()
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	zone := icstypes.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators,
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusUnbonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusUnbonding},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
		&icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
	)

	// Need to set validators in the store to calculate the aggregate intent properly
	for _, val := range zone.Validators {
		err := icsKeeper.SetValidator(ctx, zone.ChainId, *val)
		require.NoError(t, err)
	}

	out := icsKeeper.DefaultAggregateIntents(ctx, zone.ChainId)
	require.Equal(t, len(out), 3)

	for _, v := range out {
		if !v.Weight.Equal(sdk.OneDec().Quo(sdk.NewDec(3))) {
			t.Errorf("Expected %v, got %v", sdk.OneDec().Quo(sdk.NewDec(3)), v.Weight)
		}
	}
}
