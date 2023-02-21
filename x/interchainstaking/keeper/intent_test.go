package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var (
	user1 = utils.GenerateAccAddressForTest()
	user2 = utils.GenerateAccAddressForTest()
)

func (suite *KeeperTestSuite) TestKeeper_IntentStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	zoneValidatorAddresses := zone.GetValidatorsAddressesAsSlice()

	// check that there are no intents
	intents := icsKeeper.AllIntents(ctx, zone, false)
	suite.Require().Len(intents, 0)

	// set intents for testAddress
	icsKeeper.SetIntent(
		ctx,
		zone,
		icstypes.DelegatorIntent{
			Delegator: testAddress,
			Intents: icstypes.ValidatorIntents{
				{
					ValoperAddress: zoneValidatorAddresses[0],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[1],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[2],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[3],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
			},
		},
		false,
	)
	// set intents for user1
	icsKeeper.SetIntent(
		ctx,
		zone,
		icstypes.DelegatorIntent{
			Delegator: user1.String(),
			Intents: icstypes.ValidatorIntents{
				{
					ValoperAddress: zoneValidatorAddresses[0],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[1],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[2],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[3],
					Weight:         sdk.MustNewDecFromStr("0.25"),
				},
			},
		},
		false,
	)
	// set intents for user2
	icsKeeper.SetIntent(
		ctx,
		zone,
		icstypes.DelegatorIntent{
			Delegator: user2.String(),
			Intents: icstypes.ValidatorIntents{
				{
					ValoperAddress: zoneValidatorAddresses[0],
					Weight:         sdk.MustNewDecFromStr("0.5"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[1],
					Weight:         sdk.MustNewDecFromStr("0.3"),
				},
				{
					ValoperAddress: zoneValidatorAddresses[2],
					Weight:         sdk.MustNewDecFromStr("0.2"),
				},
			},
		},
		false,
	)

	// check for intents set above
	intents = icsKeeper.AllIntents(ctx, zone, false)
	suite.Require().Len(intents, 3)

	// delete intent for testAddress
	icsKeeper.DeleteIntent(ctx, zone, testAddress, false)

	// check intents
	intents = icsKeeper.AllIntents(ctx, zone, false)
	suite.Require().Len(intents, 2)

	suite.T().Logf("intents:\n%+v\n", intents)

	// update intent for user1
	err := icsKeeper.UpdateIntent(
		ctx,
		user2,
		zone,
		sdk.NewCoins(
			sdk.NewCoin(
				zone.BaseDenom,
				math.NewInt(7313913),
			),
		),
		"",
	)
	suite.Require().NoError(err)

	// load and match pointers
	intentsPointers := icsKeeper.AllIntentsAsPointer(ctx, zone, false)
	for i, ip := range intentsPointers {
		suite.Require().Equal(intents[i], *ip)
	}

	suite.T().Logf("intents:\n%+v\n", intentsPointers)
}

func (suite *KeeperTestSuite) TestAggregateIntent() {
	tc := []struct {
		name     string
		intents  func(zone icstypes.Zone) []icstypes.DelegatorIntent
		balances func() map[string]int64
		expected func(zone icstypes.Zone) icstypes.ValidatorIntents
	}{
		{
			name: "empty intents; returns equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				return out
			},
			balances: func() map[string]int64 { return map[string]int64{} },
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range zone.GetValidatorsAddressesAsSlice() {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			name: "single intent; zero balance, returns default equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				return out
			},
			balances: func() map[string]int64 { return map[string]int64{} },
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range zone.GetValidatorsAddressesAsSlice() {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			// name: "single intent; with balance, returns single weighting",
			name: "single intent; with balance, returns default equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
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
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range zone.GetValidatorsAddressesAsSlice() {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			// name: "two intents; with equal balances, same val, single weighting",
			name: "two intents; with equal balances, same val, returns default equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
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
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range zone.GetValidatorsAddressesAsSlice() {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			// name: "two intents; with equal balances, diff val, equal weighting",
			name: "two intents; with equal balances, diff val, returns default equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec()}}})
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
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range zone.GetValidatorsAddressesAsSlice() {
					out = append(out, &icstypes.ValidatorIntent{ValoperAddress: val, Weight: sdk.OneDec().Quo(sdk.NewDec(4))})
				}
				return out.Sort()
			},
		},
		{
			name: "two intents; with zer0 balances, diff val, returns default equal weights ",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.ZeroDec()}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec()}}})
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
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				for _, val := range zone.GetValidatorsAddressesAsSlice() {
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

			qapp := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			// give each user some funds
			for addrString, balance := range tt.balances() {
				suite.giveFunds(ctx, zone.LocalDenom, balance, addrString)
			}

			for _, intent := range tt.intents(zone) {
				icsKeeper.SetIntent(ctx, zone, intent, false)
			}

			icsKeeper.AggregateIntents(ctx, &zone)

			// refresh zone to pull new aggregate
			zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			actual := zone.GetAggregateIntentOrDefault()
			suite.Require().Equal(tt.expected(zone), actual)
		})
	}
}

// func (s *KeeperTestSuite) TestAggregateIntentWithPRClaims() {
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
// 					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
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
// 					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
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
// 					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
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
// 		s.Run(tt.name, func() {
// 			s.SetupTest()
// 			s.setupTestZones()

// 			qapp := s.GetQuicksilverApp(s.chainA)
// 			ctx := s.chainA.GetContext()
// 			icsKeeper := qapp.InterchainstakingKeeper
// 			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
// 			s.Require().True(found)

// 			// give each user some funds
// 			for addrString, balance := range tt.balances(zone.LocalDenom) {
// 				qapp.MintKeeper.MintCoins(ctx, balance)
// 				addr, err := utils.AccAddressFromBech32(addrString, "")
// 				s.Require().NoError(err)
// 				qapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, balance)
// 			}

// 			for _, intent := range tt.intents(zone) {
// 				icsKeeper.SetIntent(ctx, zone, intent, false)
// 			}

// 			for _, claim := range tt.claims(zone) {
// 				qapp.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &claim)
// 			}

// 			icsKeeper.AggregateIntents(ctx, zone)

// 			// refresh zone to pull new aggregate
// 			zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
// 			s.Require().True(found)

// 			actual := zone.GetAggregateIntentOrDefault()
// 			s.Require().Equal(tt.expected(zone), actual)
// 		})
// 	}
// }
