package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestHandleUpdateZoneProposal() {
	testAccount, err := addressutils.AccAddressFromBech32(testAddress, "")
	suite.NoError(err)

	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, quicksilver *app.Quicksilver)
		proposals func(zone icstypes.Zone) []icstypes.UpdateZoneProposal
		expectErr string
		check     func()
	}{
		{
			name:      "valid - all changes except connection",
			expectErr: "",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "base_denom",
								Value: "uosmo",
							},
							{
								Key:   "local_denom",
								Value: "uqosmo",
							},
							{
								Key:   "liquidity_module",
								Value: "true",
							},
						},
					},
				}
			},
		},
		// {
		// 	name:      "valid - connection",
		// 	expectErr: "",
		// 	setup: func() {
		// 		suite.setupTestZones()
		// 	},
		// 	proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
		// 		return []icstypes.UpdateZoneProposal{
		// 			{
		// 				ChainId: zone.ChainId,
		// 				Changes: []*icstypes.UpdateZoneValue{
		// 					{
		// 						Key:   "connection_id",
		// 						Value: "connection-1",
		// 					},
		// 				},
		// 			},
		// 		}
		// 	},
		// },
		{
			name:      "valid - no changes",
			expectErr: "",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{},
					},
				}
			},
		},
		{
			name:      "invalid zone",
			expectErr: "unable to get registered zone for chain id",
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: "",
						Changes: []*icstypes.UpdateZoneValue{},
					},
				}
			},
		},
		{
			name:      "invalid change key",
			expectErr: "unexpected key",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "nothing",
								Value: "nothing",
							},
						},
					},
				}
			},
		},
		{
			name:      "invalid - base_denom not valid",
			expectErr: "invalid denom",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "base_denom",
								Value: "123456789",
							},
						},
					},
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "local_denom",
								Value: "!@#$",
							},
						},
					},
				}
			},
		},
		{
			name:      "invalid - zone has assets minted",
			expectErr: "zone has assets minted",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				proposal := &icstypes.RegisterZoneProposal{
					Title:            "register zone A",
					Description:      "register zone A",
					ConnectionId:     suite.path.EndpointA.ConnectionID,
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
				}
			
				err := quicksilver.InterchainstakingKeeper.HandleRegisterZoneProposal(ctx, proposal)
				suite.NoError(err)

				err = quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
				suite.NoError(err)
				err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
				suite.NoError(err)
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "base_denom",
								Value: "uosmo",
							},
						},
					},
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "local_denom",
								Value: "uqosmo",
							},
						},
					},
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "connection_id",
								Value: "connection-1",
							},
						},
					},
				}
			},
		},
		{
			name:      "invalid - parse bool",
			expectErr: "ParseBool",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "liquidity_module",
								Value: "no",
							},
						},
					},
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "return_to_sender",
								Value: "falSE",
							},
						},
					},
				}
			},
		},
		{
			name:      "invalid - atoi",
			expectErr: "parsing",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "messages_per_tx",
								Value: "one",
							},
						},
					},
				}
			},
		},
		{
			name:      "invalid - messages_per_tx",
			expectErr: "invalid value for messages_per_tx",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "messages_per_tx",
								Value: "0",
							},
						},
					},
				}
			},
		},
		{
			name:      "invalid - connection format",
			expectErr: "unexpected connection format",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				suite.setupTestZones()
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "connection_id",
								Value: "not a connection",
							},
						},
					},
				}
			},
		},
		{
			name: "invalid - zone intialised",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				proposal := &icstypes.RegisterZoneProposal{
					Title:            "register zone A",
					Description:      "register zone A",
					ConnectionId:     suite.path.EndpointA.ConnectionID,
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
				}
			
				err := quicksilver.InterchainstakingKeeper.HandleRegisterZoneProposal(ctx, proposal)
				suite.NoError(err)

				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)

				zone.DepositAddress = &icstypes.ICAAccount{}
				quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			expectErr: "zone already intialised, cannot update connection_id",
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "connection_id",
								Value: "connection-1",
							},
						},
					},
				}
			},
		},
		// {
		// 	name:      "invalid - unable to fetch",
		// 	expectErr: "unable to fetch",
		// 	proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
		// 		return []icstypes.UpdateZoneProposal{
		// 			{
		// 				ChainId: zone.ChainId,
		// 				Changes: []*icstypes.UpdateZoneValue{
		// 					{
		// 						Key:   "connection_id",
		// 						Value: "connection-1",
		// 					},
		// 				},
		// 			},
		// 		}
		// 	},
		// },
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx := suite.chainA.GetContext()
			quicksilver := suite.GetQuicksilverApp(suite.chainA)

			var zone icstypes.Zone
			var found bool
			if tc.setup != nil {
				tc.setup(ctx, quicksilver)
				zone, found = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)
			}

			proposals := tc.proposals(zone)
			for _, proposal := range proposals {
				err := quicksilver.InterchainstakingKeeper.HandleUpdateZoneProposal(ctx, &proposal)
				if tc.expectErr != "" {
					suite.ErrorContains(err, tc.expectErr)
					suite.T().Logf("Error: %v", err)
				} else {
					suite.NoError(err)
				}
			}
		})
	}
}
