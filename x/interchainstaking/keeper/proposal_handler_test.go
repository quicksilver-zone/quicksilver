package keeper_test

import (
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	tmclienttypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	"github.com/quicksilver-zone/quicksilver/v7/app"
	"github.com/quicksilver-zone/quicksilver/v7/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestHandleUpdateZoneProposal() {
	testAccount, err := addressutils.AccAddressFromBech32(testAddress, "")
	suite.NoError(err)

	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, quicksilver *app.Quicksilver)
		proposals func(zone icstypes.Zone) []icstypes.UpdateZoneProposal
		expectErr string
		check     func(ctx sdk.Context, quicksilver *app.Quicksilver, prevZone icstypes.Zone)
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
							{
								Key:   "return_to_sender",
								Value: "F",
							},
							{
								Key:   "messages_per_tx",
								Value: "2",
							},
							{
								Key:   "account_prefix",
								Value: "osmo",
							},
						},
					},
				}
			},
			check: func(ctx sdk.Context, quicksilver *app.Quicksilver, prevZone icstypes.Zone) {
				newZone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)

				suite.Equal(newZone.BaseDenom, "uosmo")
				suite.Equal(newZone.LocalDenom, "uqosmo")
				suite.True(newZone.LiquidityModule)
				suite.False(newZone.ReturnToSender)
				suite.Equal(newZone.MessagesPerTx, int64(2))
				suite.Equal(newZone.AccountPrefix, "osmo")
			},
		},
		{
			name:      "valid - connection",
			expectErr: "",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				proposal := &icstypes.RegisterZoneProposal{
					Title:            "register zone A",
					Description:      "register zone A",
					ConnectionId:     suite.path.EndpointB.ConnectionID,
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
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "connection_id",
								Value: suite.path.EndpointA.ConnectionID,
							},
						},
					},
				}
			},
			check: func(ctx sdk.Context, quicksilver *app.Quicksilver, prevZone icstypes.Zone) {
				newZone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)

				suite.Equal(newZone.ConnectionId, suite.path.EndpointA.ConnectionID)
			},
		},
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
			check: func(ctx sdk.Context, quicksilver *app.Quicksilver, prevZone icstypes.Zone) {
				newZone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)

				suite.Equal(prevZone, newZone)
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

				quicksilver.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", &tmclienttypes.ClientState{ChainId: suite.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
				quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: ctx.BlockTime()})
				quicksilver.IBCKeeper.ConnectionKeeper.SetConnection(ctx, suite.path.EndpointA.ConnectionID, connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0"})
				suite.NoError(suite.setupChannelForICA(ctx, suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "deposit", zone.AccountPrefix))
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
		{
			name:      "invalid - unable to fetch",
			expectErr: "unable to fetch",
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
			},
			proposals: func(zone icstypes.Zone) []icstypes.UpdateZoneProposal {
				return []icstypes.UpdateZoneProposal{
					{
						ChainId: zone.ChainId,
						Changes: []*icstypes.UpdateZoneValue{
							{
								Key:   "connection_id",
								Value: "connection-10",
							},
						},
					},
				}
			},
		},
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
			for i := range proposals {
				err := quicksilver.InterchainstakingKeeper.HandleUpdateZoneProposal(ctx, &proposals[i])
				if tc.expectErr != "" {
					suite.ErrorContains(err, tc.expectErr)
					suite.T().Logf("Error: %v", err)
				} else {
					suite.NoError(err)
				}
			}

			if tc.expectErr == "" {
				tc.check(ctx, quicksilver, zone)
			}
		})
	}
}
