package keeper_test

import (
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestHandleUpdateZoneProposal() {
	tests := []struct {
		name      string
		proposals func(zone icstypes.Zone) []icstypes.UpdateZoneProposal
		expectErr string
	}{
		{
			name:      "valid - all changes except connection",
			expectErr: "",
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
		{
			name:      "valid - connection",
			expectErr: "",
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
			name:      "valid - no changes",
			expectErr: "",
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
			name:      "invalid - zone intialised",
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
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			ctx := suite.chainA.GetContext()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

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
