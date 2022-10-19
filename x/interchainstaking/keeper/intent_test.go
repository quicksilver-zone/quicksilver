package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var (
	user1 = utils.GenerateAccAddressForTest()
	user2 = utils.GenerateAccAddressForTest()
)

func (s *KeeperTestSuite) TestAggregateIntent() {
	tc := []struct {
		name     string
		intents  func(zone icstypes.Zone) []icstypes.DelegatorIntent
		balances func(denom string) map[string]sdk.Coins
		expected func(zone icstypes.Zone) icstypes.ValidatorIntents
	}{
		{
			name: "empty intents; returns equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				return out
			},
			balances: func(denom string) map[string]sdk.Coins { return map[string]sdk.Coins{} },
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
			balances: func(denom string) map[string]sdk.Coins { return map[string]sdk.Coins{} },
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
			name: "single intent; with balance, returns single weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				return out
			},
			balances: func(denom string) map[string]sdk.Coins {
				return map[string]sdk.Coins{user1.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.OneInt()))}
			},
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})

				return out.Sort()
			},
		},
		{
			name: "two intents; with equal balances, same val, single weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				return out
			},
			balances: func(denom string) map[string]sdk.Coins {
				return map[string]sdk.Coins{
					user1.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.OneInt())),
					user2.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.OneInt())),
				}
			},
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})

				return out.Sort()
			},
		},
		{
			name: "two intents; with equal balances, diff val, equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec()}}})
				return out
			},
			balances: func(denom string) map[string]sdk.Coins {
				return map[string]sdk.Coins{
					user1.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.OneInt())),
					user2.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.OneInt())),
				}
			},
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})

				return out.Sort()
			},
		},
	}

	for _, tt := range tc {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupZones()

			qapp := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			// give each user some funds
			for addrString, balance := range tt.balances(zone.LocalDenom) {
				qapp.MintKeeper.MintCoins(ctx, balance)
				addr, err := utils.AccAddressFromBech32(addrString, "")
				s.Require().NoError(err)
				qapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, balance)
			}

			for _, intent := range tt.intents(zone) {
				icsKeeper.SetIntent(ctx, zone, intent, false)
			}

			icsKeeper.AggregateIntents(ctx, zone)

			// refresh zone to pull new aggregate
			zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			actual := zone.GetAggregateIntentOrDefault()
			s.Require().Equal(tt.expected(zone), actual)
		})
	}
}

func (s *KeeperTestSuite) TestAggregateIntentWithPRClaims() {
	tc := []struct {
		name     string
		intents  func(zone icstypes.Zone) []icstypes.DelegatorIntent
		balances func(denom string) map[string]sdk.Coins
		claims   func(zone icstypes.Zone) map[string]cmtypes.Claim
		expected func(zone icstypes.Zone) icstypes.ValidatorIntents
	}{
		{
			name: "single intent; zero balance, but claim, returns single weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				return out
			},
			claims: func(zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
				}
			},
			balances: func(denom string) map[string]sdk.Coins { return map[string]sdk.Coins{} },
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				// four delegators each at 25%
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})
				return out.Sort()
			},
		},
		{
			name: "single intent; with balance and claim, returns single weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				return out
			},
			claims: func(zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
				}
			},
			balances: func(denom string) map[string]sdk.Coins {
				return map[string]sdk.Coins{user1.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.OneInt()))}
			},
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()})
				return out.Sort()
			},
		},
		{
			name: "two intents; one balance and one claim, returns equal weighting",
			intents: func(zone icstypes.Zone) []icstypes.DelegatorIntent {
				out := make([]icstypes.DelegatorIntent, 0)
				out = append(out, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec()}}})
				out = append(out, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec()}}})
				return out
			},
			claims: func(zone icstypes.Zone) map[string]cmtypes.Claim {
				return map[string]cmtypes.Claim{
					user1.String(): {UserAddress: user1.String(), ChainId: zone.ChainId, Module: cmtypes.ClaimTypeLiquidToken, Amount: 1000, SourceChainId: "osmosis-1"},
				}
			},
			balances: func(denom string) map[string]sdk.Coins {
				return map[string]sdk.Coins{user2.String(): sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000)))}
			},
			expected: func(zone icstypes.Zone) icstypes.ValidatorIntents {
				out := icstypes.ValidatorIntents{}
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})
				out = append(out, &icstypes.ValidatorIntent{ValoperAddress: zone.GetValidatorsAddressesAsSlice()[1], Weight: sdk.OneDec().Quo(sdk.NewDec(2))})
				return out.Sort()
			},
		},
	}

	for _, tt := range tc {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupZones()

			qapp := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			// give each user some funds
			for addrString, balance := range tt.balances(zone.LocalDenom) {
				qapp.MintKeeper.MintCoins(ctx, balance)
				addr, err := utils.AccAddressFromBech32(addrString, "")
				s.Require().NoError(err)
				qapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, balance)
			}

			for _, intent := range tt.intents(zone) {
				icsKeeper.SetIntent(ctx, zone, intent, false)
			}

			for _, claim := range tt.claims(zone) {
				qapp.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &claim)
			}

			icsKeeper.AggregateIntents(ctx, zone)

			// refresh zone to pull new aggregate
			zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			actual := zone.GetAggregateIntentOrDefault()
			s.Require().Equal(tt.expected(zone), actual)
		})
	}
}
