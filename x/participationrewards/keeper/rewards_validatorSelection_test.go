package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

var (
	user1, _ = addressutils.AccAddressFromBech32("cosmos1kcgprgjxntc5w4ygfsgvjnnypeptf3vw6gyv0z77h27cx23vg5rsptlw4a", "")
	user2, _ = addressutils.AccAddressFromBech32("cosmos1u4ln57y5m2qyna7aq09u3r05waf74ad9rsk4hzr79acapar6lhqqumtd5d", "")
)

func (suite *KeeperTestSuite) TestCalcUserValidatorSelectionAllocations() {
	tests := []struct {
		name            string
		malleate        func(sdk.Context, *app.Quicksilver)
		validatorScores func(sdk.Context, *app.Quicksilver, string) map[string]*types.Validator
		want            func(denom string) []types.UserAllocation
	}{
		{
			name: "no allocation",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
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
			want: func(denom string) []types.UserAllocation { return []types.UserAllocation{} },
		},
		{
			name: "zero weight intents, no user allocation",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
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
				validators := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)
				for i := range validators {
					validators[i].Score = sdk.NewDec(1)
					validatorScores[validators[i].ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &validators[i],
					}
				}
				return validatorScores
			},
			want: func(denom string) []types.UserAllocation { return []types.UserAllocation{} },
		},
		{
			name: "unit weight intents - default validator scores - same valopaddress",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
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

				validators := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)
				for i := range validators {
					validators[i].Score = sdk.NewDec(1)
					validatorScores[validators[i].ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &validators[i],
					}
				}
				return validatorScores
			},
			want: func(denom string) []types.UserAllocation {
				return []types.UserAllocation{
					{
						Address: user1.String(),
						Amount:  sdk.NewCoin(denom, sdk.NewInt(2500)),
					},
					{
						Address: user2.String(),
						Amount:  sdk.NewCoin(denom, sdk.NewInt(2500)),
					},
				}
			},
		},
		{
			name: "unit weight intents - default validator scores - different validators",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
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

				validators := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)
				for i := range validators {
					validators[i].Score = sdk.NewDec(1)
					validatorScores[validators[i].ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &validators[i],
					}
				}
				return validatorScores
			},
			want: func(denom string) []types.UserAllocation {
				return []types.UserAllocation{
					{
						Address: user1.String(),
						Amount:  sdk.NewCoin(denom, sdk.NewInt(2500)),
					},
					{
						Address: user2.String(),
						Amount:  sdk.NewCoin(denom, sdk.NewInt(2500)),
					},
				}
			},
		},
		{
			name: "weighted intents - default validator scores - same validators",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
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

				validators := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)
				for i := range validators {
					validators[i].Score = sdk.NewDec(1)
					validatorScores[validators[i].ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &validators[i],
					}
				}
				return validatorScores
			},
			want: func(denom string) []types.UserAllocation {
				return []types.UserAllocation{
					{
						Address: user1.String(),
						Amount:  sdk.NewCoin(denom, sdk.NewInt(454)),
					},
					{
						Address: user2.String(),
						Amount:  sdk.NewCoin(denom, sdk.NewInt(4545)),
					},
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()

			appA := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(ctx, appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			suite.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			suite.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			validatorScores := tt.validatorScores(ctx, appA, zone.ChainId)

			zs := types.ZoneScore{
				ZoneID:           zone.ChainId,
				TotalVotingPower: sdk.NewInt(0),
				ValidatorScores:  validatorScores,
			}

			userAllocations := appA.ParticipationRewardsKeeper.CalcUserValidatorSelectionAllocations(ctx, &zone, zs)
			suite.Require().Equal(tt.want(appA.StakingKeeper.BondDenom(ctx)), userAllocations)
		})
	}
}

func (suite *KeeperTestSuite) TestCalcDistributionScores() {
	tests := []struct {
		name            string
		malleate        func(sdk.Context, *app.Quicksilver)
		validatorScores func(sdk.Context, *app.Quicksilver, string) map[string]*types.Validator
		verify          func(sdk.Context, *app.Quicksilver, types.ZoneScore)
		wantErr         bool
	}{
		{
			name: "zero validators",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
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
			verify: func(context sdk.Context, quicksilver *app.Quicksilver, score types.ZoneScore) {
			},
			wantErr: true,
		},
		{
			name: "zero voting power",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})
				zone.ValidatorSelectionAllocation = 5000

				validators := appA.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				for i := range validators {
					validators[i].VotingPower = sdk.NewInt(0)
					appA.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, validators[i])
				}

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				validators := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)
				for i := range validators {
					validators[i].Score = sdk.NewDec(1)
					validatorScores[validators[i].ValoperAddress] = &types.Validator{
						PowerPercentage:   sdk.NewDec(1),
						DistributionScore: sdk.NewDec(1),
						PerformanceScore:  sdk.NewDec(1),
						Validator:         &validators[i],
					}
				}
				return validatorScores
			},
			verify: func(context sdk.Context, quicksilver *app.Quicksilver, score types.ZoneScore) {
			},
			wantErr: true,
		},
		{
			name: "valid zonescore, different power",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})

				for i, val := range appA.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId) {
					val.VotingPower = sdk.NewInt(int64(10 + i*10))
					appA.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val)
				}
				zone.ValidatorSelectionAllocation = 5000

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				vals := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)
				(&vals[0]).VotingPower = sdk.NewInt(10)
				validatorScores[vals[0].ValoperAddress] = &types.Validator{
					PowerPercentage:   sdk.NewDec(1),
					DistributionScore: sdk.NewDec(1),
					PerformanceScore:  sdk.NewDec(1),
					Validator:         &vals[0],
				}
				(&vals[1]).VotingPower = sdk.NewInt(20)
				validatorScores[vals[1].ValoperAddress] = &types.Validator{
					PowerPercentage:   sdk.NewDec(1),
					DistributionScore: sdk.NewDec(1),
					PerformanceScore:  sdk.NewDec(1),
					Validator:         &vals[1],
				}
				(&vals[2]).VotingPower = sdk.NewInt(30)
				validatorScores[vals[2].ValoperAddress] = &types.Validator{
					PowerPercentage:   sdk.NewDec(1),
					DistributionScore: sdk.NewDec(1),
					PerformanceScore:  sdk.NewDec(1),
					Validator:         &vals[2],
				}
				return validatorScores
			},
			verify: func(ctx sdk.Context, appA *app.Quicksilver, zs types.ZoneScore) {
				suite.Require().Equal(zs.TotalVotingPower, sdk.NewInt(100))

				validators := appA.InterchainstakingKeeper.GetValidators(ctx, zs.ZoneID)

				suite.Require().Equal(strings.TrimRight(zs.ValidatorScores[validators[0].ValoperAddress].PowerPercentage.String(), "0"), "0.1")
				suite.Require().Equal(strings.TrimRight(zs.ValidatorScores[validators[1].ValoperAddress].PowerPercentage.String(), "0"), "0.2")
				suite.Require().Equal(strings.TrimRight(zs.ValidatorScores[validators[2].ValoperAddress].PowerPercentage.String(), "0"), "0.3")

				suite.Require().Equal(zs.ValidatorScores[validators[0].ValoperAddress].DistributionScore, sdk.NewDec(1))
				suite.Require().Equal(strings.TrimRight(zs.ValidatorScores[validators[1].ValoperAddress].DistributionScore.String(), "0"), "0.75")
				suite.Require().Equal(strings.TrimRight(zs.ValidatorScores[validators[2].ValoperAddress].DistributionScore.String(), "0"), "0.5")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()

			appA := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(ctx, appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			suite.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			suite.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			validatorScores := tt.validatorScores(ctx, appA, zone.ChainId)

			zs := types.ZoneScore{
				ZoneID:           zone.ChainId,
				TotalVotingPower: sdk.NewInt(0),
				ValidatorScores:  validatorScores,
			}

			err := appA.ParticipationRewardsKeeper.CalcDistributionScores(ctx, zone, &zs)
			suite.Require().Equal(err != nil, tt.wantErr)
			tt.verify(ctx, appA, zs)
		})
	}
}

func (suite *KeeperTestSuite) TestCalcOverallScores() {
	tests := []struct {
		name             string
		malleate         func(sdk.Context, *app.Quicksilver)
		validatorScores  func(sdk.Context, *app.Quicksilver, string) map[string]*types.Validator
		delegatorRewards func(sdk.Context, *app.Quicksilver, string) distributiontypes.QueryDelegationTotalRewardsResponse
		verify           func(types.ZoneScore, []distributiontypes.DelegationDelegatorReward, []icstypes.Validator)
		wantErr          bool
	}{
		{
			name: "nil delegation rewards",
			malleate: func(context sdk.Context, quicksilver *app.Quicksilver) {
			},
			validatorScores: func(context sdk.Context, quicksilver *app.Quicksilver, s string) map[string]*types.Validator {
				return nil
			},
			delegatorRewards: func(_ sdk.Context, _ *app.Quicksilver, _ string) distributiontypes.QueryDelegationTotalRewardsResponse {
				return distributiontypes.QueryDelegationTotalRewardsResponse{}
			},
			verify: func(_ types.ZoneScore, _ []distributiontypes.DelegationDelegatorReward, _ []icstypes.Validator) {
			},
			wantErr: true,
		},
		{
			name: "zero total rewards",
			malleate: func(context sdk.Context, quicksilver *app.Quicksilver) {
			},
			validatorScores: func(context sdk.Context, quicksilver *app.Quicksilver, s string) map[string]*types.Validator {
				return nil
			},
			delegatorRewards: func(ctx sdk.Context, appA *app.Quicksilver, chainID string) distributiontypes.QueryDelegationTotalRewardsResponse {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				valAddress := appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, chainID)[0]
				return distributiontypes.QueryDelegationTotalRewardsResponse{Rewards: []distributiontypes.DelegationDelegatorReward{{ValidatorAddress: valAddress, Reward: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(1)))}}, Total: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(0)))}
			},
			verify: func(_ types.ZoneScore, _ []distributiontypes.DelegationDelegatorReward, _ []icstypes.Validator) {
			},
			wantErr: true,
		},
		{
			name: "validator removed from active set - performance greater than limit",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)

				val := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)[1]
				validatorScores[val.ValoperAddress] = &types.Validator{
					PowerPercentage:   sdk.NewDec(1),
					DistributionScore: sdk.NewDec(3),
					PerformanceScore:  sdk.NewDec(619),
					Validator:         &val,
				}

				return validatorScores
			},
			delegatorRewards: func(ctx sdk.Context, appA *app.Quicksilver, chainID string) distributiontypes.QueryDelegationTotalRewardsResponse {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				inactiveAddress := appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, chainID)[0]
				activeAddress := appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, chainID)[1]
				return distributiontypes.QueryDelegationTotalRewardsResponse{Rewards: []distributiontypes.DelegationDelegatorReward{
					{ValidatorAddress: inactiveAddress, Reward: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(1)))},
					{ValidatorAddress: activeAddress, Reward: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(10)))},
				}, Total: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(11)))}
			},
			verify: func(zs types.ZoneScore, delegatorRewards []distributiontypes.DelegationDelegatorReward, validators []icstypes.Validator) {
				suite.Require().True(zs.ValidatorScores[delegatorRewards[0].ValidatorAddress] == nil)
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[1].ValidatorAddress].PerformanceScore, sdk.NewDec(1))
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[1].ValidatorAddress].Score, sdk.NewDec(3))
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[1].ValidatorAddress].Score, validators[1].Score)
			},
		},
		{
			name: "multiple validators rewarded",
			malleate: func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				appA.InterchainstakingKeeper.IterateDelegatorIntents(ctx, &zone, true, func(_ int64, di icstypes.DelegatorIntent) (stop bool) {
					appA.InterchainstakingKeeper.DeleteDelegatorIntent(ctx, &zone, di.Delegator, true)
					return false
				})

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			validatorScores: func(ctx sdk.Context, appA *app.Quicksilver, chainId string) map[string]*types.Validator {
				validatorScores := make(map[string]*types.Validator)
				vals := appA.InterchainstakingKeeper.GetValidators(ctx, chainId)

				validatorScores[vals[0].ValoperAddress] = &types.Validator{
					PowerPercentage:   sdk.NewDec(1),
					DistributionScore: sdk.NewDec(1),
					PerformanceScore:  sdk.NewDec(1),
					Validator:         &vals[0],
				}
				validatorScores[vals[1].ValoperAddress] = &types.Validator{
					PowerPercentage:   sdk.NewDec(1),
					DistributionScore: sdk.NewDec(5),
					PerformanceScore:  sdk.NewDec(1),
					Validator:         &vals[1],
				}
				validatorScores[vals[2].ValoperAddress] = &types.Validator{
					PowerPercentage:   sdk.NewDec(1),
					DistributionScore: sdk.NewDec(7),
					PerformanceScore:  sdk.NewDec(1),
					Validator:         &vals[2],
				}
				return validatorScores
			},
			delegatorRewards: func(ctx sdk.Context, appA *app.Quicksilver, chainID string) distributiontypes.QueryDelegationTotalRewardsResponse {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				validators := appA.InterchainstakingKeeper.GetValidatorAddresses(ctx, chainID)
				return distributiontypes.QueryDelegationTotalRewardsResponse{Rewards: []distributiontypes.DelegationDelegatorReward{
					{ValidatorAddress: validators[0], Reward: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(5)))},
					{ValidatorAddress: validators[1], Reward: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(10)))},
					{ValidatorAddress: validators[2], Reward: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(15)))},
				}, Total: sdk.NewDecCoins(sdk.NewDecCoin(zone.BaseDenom, sdk.NewInt(30)))}
			},
			verify: func(zs types.ZoneScore, delegatorRewards []distributiontypes.DelegationDelegatorReward, validators []icstypes.Validator) {
				suite.Require().Equal(strings.TrimRight(zs.ValidatorScores[delegatorRewards[0].ValidatorAddress].PerformanceScore.String(), "0"), "0.5")
				suite.Require().Equal(strings.TrimRight(zs.ValidatorScores[delegatorRewards[0].ValidatorAddress].Score.String(), "0"), "0.5")
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[0].ValidatorAddress].Score, validators[0].Score)

				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[1].ValidatorAddress].PerformanceScore, sdk.NewDec(1))
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[1].ValidatorAddress].Score, sdk.NewDec(5))
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[1].ValidatorAddress].Score, validators[1].Score)

				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[2].ValidatorAddress].PerformanceScore, sdk.NewDec(1))
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[2].ValidatorAddress].Score, sdk.NewDec(7))
				suite.Require().Equal(zs.ValidatorScores[delegatorRewards[2].ValidatorAddress].Score, validators[2].Score)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()

			appA := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(ctx, appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			suite.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			suite.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			validatorScores := tt.validatorScores(ctx, appA, zone.ChainId)

			zs := types.ZoneScore{
				ZoneID:           zone.ChainId,
				TotalVotingPower: sdk.NewInt(0),
				ValidatorScores:  validatorScores,
			}

			delegatorRewards := tt.delegatorRewards(ctx, appA, zone.ChainId)

			err := appA.ParticipationRewardsKeeper.CalcOverallScores(ctx, zone, delegatorRewards, &zs)

			suite.Require().Equal(err != nil, tt.wantErr)
			tt.verify(zs, delegatorRewards.Rewards, appA.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))
		})
	}
}
