package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
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
	suite.Require().True(found)
	zoneValidatorAddresses := icsKeeper.GetValidators(ctx, zone.ChainID())

	// check that there are no intents
	intents := icsKeeper.AllDelegatorIntents(ctx, &zone, false)
	suite.Require().Len(intents, 0)

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
	suite.Require().Len(intents, 3)

	// delete intent for testAddress
	icsKeeper.DeleteDelegatorIntent(ctx, &zone, testAddress, false)

	// check intents
	intents = icsKeeper.AllDelegatorIntents(ctx, &zone, false)
	suite.Require().Len(intents, 2)

	suite.T().Logf("intents:\n%+v\n", intents)

	// update intent for user1
	err := icsKeeper.UpdateDelegatorIntent(
		ctx,
		user2,
		&zone,
		sdk.NewCoins(
			sdk.NewCoin(
				zone.BaseDenom,
				math.NewInt(7313913),
			),
		),
		nil,
	)
	suite.Require().NoError(err)

	// load and match pointers
	intentsPointers := icsKeeper.AllDelegatorIntentsAsPointer(ctx, &zone, false)
	for i, ip := range intentsPointers {
		suite.Require().Equal(intents[i], *ip)
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
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainID()) {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			name: "single intent; zero balance, returns default equal weighting",
			intents: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainID())[0], Weight: sdk.OneDec()}}})
				return out
			},
			balances: func() map[string]int64 { return map[string]int64{} },
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainID()) {
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
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainID())[0], Weight: sdk.OneDec()}}})
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
			suite.Require().True(found)

			// give each user some funds
			for addrString, balance := range tt.balances() {
				suite.giveFunds(ctx, zone.LocalDenom, balance, addrString)
			}

			for _, intent := range tt.intents(ctx, quicksilver, zone) {
				icsKeeper.SetDelegatorIntent(ctx, &zone, intent, false)
			}

			_ = icsKeeper.AggregateDelegatorIntents(ctx, &zone)

			// refresh zone to pull new aggregate
			zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			actual, err := icsKeeper.GetAggregateIntentOrDefault(ctx, &zone)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.expected(ctx, quicksilver, zone), actual)
		})
	}
}

// func (suite *KeeperTestSuite) TestAggregateIntentWithPRClaims() {
// 	tc := []struct {
// 		name     string
// 		intents  func(zone icstypes.Zone) []icstypes.DelegatorIntent
// 		balances func(denom string) map[string]sdk.Coins
// 		claims   func(zone icstypes.Zone) map[string]cmtypes.Claim
// 		expected func(zone icstypes.Zone) icstypes.ValidatorIntents
// 	}{
// 		{
// 			name: "single intent; zero balance, but claim, returns single weighting",
// 			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
// 				out := make([]icstypes.DelegatorIntent, 0)
// 				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
// 				return out
// 			},
// 			claims: func(zone icstypes.Zone) map[string]cmtypes.Claim {
// 				return map[string]cmtypes.Claim{
// 					user1.String(): {UserAddress: user1.String(), ChainID: zone.ChainID, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
// 				}
// 			},
// 			balances: func(denom string) map[string]sdk.Coins { return map[string]sdk.Coins{} },
// 			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
// 				// four delegators each at 25%
// 				out := icstypes.ValidatorIntents{}
// 				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})
// 				return out.Sort()
// 			},
// 		},
// 		{
// 			name: "single intent; with balance and claim, returns single weighting",
// 			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
// 				out := make([]icstypes.DelegatorIntent, 0)
// 				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
// 				return out
// 			},
// 			claims: func(zone icstypes.Zone) map[string]cmtypes.Claim {
// 				return map[string]cmtypes.Claim{
// 					user1.String(): {UserAddress: user1.String(), ChainID: zone.ChainID, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
// 				}
// 			},
// 			balances: func(denom string) map[string]sdk.Coins {
// 				return map[string]sdk.Coins{user1.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.OneInt()))}
// 			},
// 			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
// 				out := icstypes.ValidatorIntents{}
// 				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})
// 				return out.Sort()
// 			},
// 		},
// 		{
// 			name: "two intents; one balance and one claim, returns equal weighting",
// 			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
// 				out := make([]icstypes.DelegatorIntent, 0)
// 				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
// 				// next intent cannot be set, as local asset balance does not qualify
// 				// out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec()}}})
// 				return out
// 			},
// 			claims: func(zone icstypes.Zone) map[string]cmtypes.Claim {
// 				return map[string]cmtypes.Claim{
// 					user1.String(): {UserAddress: user1.String(), ChainID: zone.ChainID, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
// 				}
// 			},
// 			balances: func(denom string) map[string]sdk.Coins {
// 				return map[string]sdk.Coins{user2.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000)))}
// 			},
// 			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
// 				out := icstypes.ValidatorIntents{}
// 				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})
// 				// only remote assets are considered, thus user2 balance is ignored...
// 				// out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})
// 				return out.Sort()
// 			},
// 		},
// 	}

// 	for _, tt := range tc {
// 		suite.Run(tt.name, func() {
// 			suite.SetupTest()
// 			suite.setupTestZones()

// 			quicksilver := suite.GetQuicksilverApp(suite.chainA)
// 			ctx := suite.chainA.GetContext()
// 			icsKeeper := quicksilver.InterchainstakingKeeper
// 			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
// 			suite.Require().True(found)

// 			// give each user some funds
// 			for addrString, balance := range tt.balances(zone.LocalDenom) {
// 				quicksilver.MintKeeper.MintCoins(ctx, balance)
// 				addr, err := utils.AccAddressFromBech32(addrString, "")
// 				suite.Require().NoError(err)
// 				quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, balance)
// 			}

// 			for _, intent := range tt.intents(zone) {
// 				icsKeeper.SetDelegatorIntent(ctx, zone, intent, false)
// 			}

// 			for _, claim := range tt.claims(zone) {
// 				quicksilver.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &claim)
// 			}

// 			icsKeeper.AggregateDelegatorIntents(ctx, zone)

// 			// refresh zone to pull new aggregate
// 			zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
// 			suite.Require().True(found)

// 			actual := zone.GetAggregateIntentOrDefault()
// 			suite.Require().Equal(tt.expected(zone), actual)
// 		})
// 	}
// }

// TODO: convert to keeper tests

// func TestDefaultIntent(t *testing.T) {
// 	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
// 	zone.Validators = append(zone.Validators,
// 		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 	)

// 	out := zone.DefaultAggregateIntents()
// 	require.Equal(t, len(out), 5)
// 	for _, v := range out {
// 		if !v.Weight.Equal(sdk.NewDecWithPrec(2, 1)) {
// 			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(2, 1), v.Weight)
// 		}
// 	}
// }

// func TestDefaultIntentWithJailed(t *testing.T) {
// 	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
// 	zone.Validators = append(zone.Validators,
// 		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded, Jailed: true},
// 		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 	)

// 	out := zone.DefaultAggregateIntents()
// 	require.Equal(t, len(out), 4)

// 	for _, v := range out {
// 		if !v.Weight.Equal(sdk.NewDecWithPrec(25, 2)) {
// 			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(25, 2), v.Weight)
// 		}
// 	}
// }

// func TestDefaultIntentWithTombstoned(t *testing.T) {
// 	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
// 	zone.Validators = append(zone.Validators,
// 		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded, Tombstoned: true},
// 		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 	)

// 	out := zone.DefaultAggregateIntents()
// 	require.Equal(t, len(out), 4)

// 	for _, v := range out {
// 		if !v.Weight.Equal(sdk.NewDecWithPrec(25, 2)) {
// 			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(25, 2), v.Weight)
// 		}
// 	}
// }

// func TestDefaultIntentWithCommission100(t *testing.T) {
// 	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
// 	zone.Validators = append(zone.Validators,
// 		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 	)

// 	out := zone.DefaultAggregateIntents()
// 	require.Equal(t, len(out), 4)

// 	for _, v := range out {
// 		if !v.Weight.Equal(sdk.NewDecWithPrec(25, 2)) {
// 			t.Errorf("Expected %v, got %v", sdk.NewDecWithPrec(25, 2), v.Weight)
// 		}
// 	}
// }

// func TestDefaultIntentWithOneUnbondedOneUnbonding(t *testing.T) {
// 	zone := types.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
// 	zone.Validators = append(zone.Validators,
// 		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusUnbonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusUnbonding},
// 		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded},
// 	)

// 	out := zone.DefaultAggregateIntents()
// 	require.Equal(t, len(out), 3)

// 	for _, v := range out {
// 		if !v.Weight.Equal(sdk.OneDec().Quo(sdk.NewDec(3))) {
// 			t.Errorf("Expected %v, got %v", sdk.OneDec().Quo(sdk.NewDec(3)), v.Weight)
// 		}
// 	}
// }
