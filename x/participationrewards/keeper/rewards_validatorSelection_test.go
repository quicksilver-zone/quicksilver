package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"strings"
)

var (
	user1, _ = addressutils.AccAddressFromBech32("cosmos1kcgprgjxntc5w4ygfsgvjnnypeptf3vw6gyv0z77h27cx23vg5rsptlw4a", "")
	user2, _ = addressutils.AccAddressFromBech32("cosmos1u4ln57y5m2qyna7aq09u3r05waf74ad9rsk4hzr79acapar6lhqqumtd5d", "")
)

func (s *KeeperTestSuite) TestCalcUserValidatorSelectionAllocations() {

	tests := []struct {
		name            string
		malleate        func(sdk.Context, *app.Quicksilver)
		validatorScores func(sdk.Context, *app.Quicksilver, string) map[string]*types.Validator
		want            []types.UserAllocation
	}{
		{
			name: "no allocation",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 0
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(context sdk.Context, quicksilver *app.Quicksilver, s string) map[string]*types.Validator {
				return nil
			},
			want: []types.UserAllocation{},
		},
		{
			name: "zero weight intents, no user allocation",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 5000

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)

				validatorIntents := icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.NewDec(0)}}

				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: validatorIntents}, true)
				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: validatorIntents}, true)

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, chainId) {
					val.Score = sdk.NewDec(1)
					validatorScores[val.ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &val,
					}
				}
				return validatorScores
			},
			want: []types.UserAllocation{},
		},
		{
			name: "unit weight intents - default validator scores - same valopaddress",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 5000

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)

				validatorIntents := icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.NewDec(1)}}

				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: validatorIntents}, true)
				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: validatorIntents}, true)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, chainId) {
					val.Score = sdk.NewDec(1)
					validatorScores[val.ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &val,
					}
				}
				return validatorScores
			},
			want: []types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewInt(2500),
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewInt(2500),
				},
			},
		},
		{
			name: "unit weight intents - default validator scores - different validators",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 5000

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)

				validatorIntentsA := icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.NewDec(1)}}
				validatorIntentsB := icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[1], Weight: sdk.NewDec(1)}}

				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: validatorIntentsA}, true)
				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: validatorIntentsB}, true)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, chainId) {
					val.Score = sdk.NewDec(1)
					validatorScores[val.ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &val,
					}
				}
				return validatorScores
			},
			want: []types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewInt(2500),
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewInt(2500),
				},
			},
		},
		{
			name: "weighted intents - default validator scores - same validators",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 5000

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)

				validatorIntentsA := icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.NewDec(10)}}
				validatorIntentsB := icstypes.ValidatorIntents{&icstypes.ValidatorIntent{ValoperAddress: appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0], Weight: sdk.NewDec(100)}}

				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user1.String(), Intents: validatorIntentsA}, true)
				appA.InterchainstakingKeeper.SetDelegatorIntent(ctx, &zone, icstypes.DelegatorIntent{Delegator: user2.String(), Intents: validatorIntentsB}, true)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, chainId) {
					val.Score = sdk.NewDec(1)
					validatorScores[val.ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &val,
					}
				}
				return validatorScores
			},
			want: []types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewInt(454),
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewInt(4545),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()

			appA := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(ctx, appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			s.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			validatorScores := tt.validatorScores(ctx, appA, zone.ChainId)

			zs := types.ZoneScore{
				ZoneID:           zone.ChainId,
				TotalVotingPower: sdk.NewInt(0),
				ValidatorScores:  validatorScores,
			}

			userAllocations := appA.ParticipationRewardsKeeper.CalcUserValidatorSelectionAllocations(ctx, &zone, zs)
			s.Require().Equal(tt.want, userAllocations)
		})
	}
}
func (s *KeeperTestSuite) TestCalcDistributionScores() {

	tests := []struct {
		name            string
		malleate        func(sdk.Context, *app.Quicksilver)
		validatorScores func(sdk.Context, *app.Quicksilver, string) map[string]*types.Validator
		verify          func(zs types.ZoneScore)
		wantErr         bool
	}{
		{
			name: "zero validators",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 5000

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId) {
					address, _ := val.GetAddressBytes()
					appA.InterchainstakingKeeper.DeleteValidator(ctx, zone.ChainId, address)
				}

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				return nil
			},
			verify: func(zs types.ZoneScore) {
			},
			wantErr: true,
		},
		{
			name: "zero voting power",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 5000

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId) {
					val.VotingPower = sdk.NewInt(0)
					appA.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val)
				}

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, chainId) {
					val.Score = sdk.NewDec(1)
					validatorScores[val.ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &val,
					}
				}
				return validatorScores
			},
			verify: func(zs types.ZoneScore) {
			},
			wantErr: true,
		},
		{
			name: "valid zonescore, same power",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId) {
					val.VotingPower = sdk.NewInt(10)
					appA.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val)
				}
				zone.ValidatorSelectionAllocation = 5000

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				for _, val := range appA.InterchainstakingKeeper.GetValidators(ctx, chainId) {
					validatorScores[val.ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &val,
					}
				}
				return validatorScores
			},
			verify: func(zs types.ZoneScore) {
				s.Require().Equal(zs.TotalVotingPower, sdk.NewInt(40))

				for _, val := range zs.ValidatorScores {
					s.Require().Equal(strings.TrimRight(val.PowerPercentage.String(), "0"), "0.25")
					s.Require().Equal(val.DistributionScore, sdk.NewDec(1))
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()

			appA := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(ctx, appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			s.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			validatorScores := tt.validatorScores(ctx, appA, zone.ChainId)

			zs := types.ZoneScore{
				ZoneID:           zone.ChainId,
				TotalVotingPower: sdk.NewInt(0),
				ValidatorScores:  validatorScores,
			}

			err := appA.ParticipationRewardsKeeper.CalcDistributionScores(ctx, zone, &zs)
			s.Require().Equal(err != nil, tt.wantErr)
			tt.verify(zs)
		})
	}
}
