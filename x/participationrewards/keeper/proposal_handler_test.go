package keeper_test

import (
	"fmt"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) TestHandleAddProtocolDataProposal() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	prop := types.AddProtocolDataProposal{}
	tests := []struct {
		name     string
		malleate func()
		wantErr  bool
	}{
		{
			"blank",
			func() {},
			true,
		},
		{
			"invalid_prop",
			func() {
				prop = types.AddProtocolDataProposal{
					Title:       "Add connection protocol for test chain B",
					Description: "A connection protocol for testing connection protocols",
					Protocol:    "",
					Type:        "",
					Key:         "",
					Data:        nil,
				}
			},
			true,
		},
		{
			"invalid_prop_data_type",
			func() {
				prop = types.AddProtocolDataProposal{
					Title:       "Add connection protocol for test chain B",
					Description: "A connection protocol for testing connection protocols",
					Protocol:    fmt.Sprintf("connection/%s", suite.chainB.ChainID),
					Type:        "testtype",
					Key:         "testkey",
					Data:        []byte("{}"),
				}
			},
			true,
		},
		{
			"invalid_prop_data_empty",
			func() {
				prop = types.AddProtocolDataProposal{
					Title:       "Add connection protocol for test chain B",
					Description: "A connection protocol for testing connection protocols",
					Protocol:    fmt.Sprintf("connection/%s", suite.chainB.ChainID),
					Type:        types.ProtocolDataType_name[types.ProtocolDataConnection],
					Key:         "connection",
					Data:        []byte("{}"),
				}
			},
			true,
		},
		{
			"invalid_prop_data",
			func() {
				connpdstr := fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", "", "", 100)

				prop = types.AddProtocolDataProposal{
					Title:       "Add connection protocol for test chain B",
					Description: "A connection protocol for testing connection protocols",
					Protocol:    fmt.Sprintf("connection/%s", suite.chainB.ChainID),
					Type:        types.ProtocolDataType_name[types.ProtocolDataConnection],
					Key:         "connection",
					Data:        []byte(connpdstr),
				}
			},
			true,
		},
		{
			"valid_prop",
			func() {
				connpdstr := fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0)

				prop = types.AddProtocolDataProposal{
					Title:       "Add connection protocol for test chain B",
					Description: "A connection protocol for testing connection protocols",
					Protocol:    fmt.Sprintf("connection/%s", suite.chainB.ChainID),
					Type:        types.ProtocolDataType_name[types.ProtocolDataConnection],
					Key:         "connection",
					Data:        []byte(connpdstr),
				}
			},
			false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			k := appA.ParticipationRewardsKeeper
			err := keeper.HandleAddProtocolDataProposal(suite.chainA.GetContext(), k, &prop)
			if tt.wantErr {
				suite.Require().Error(err)
				suite.T().Logf("Error: %v", err)
				return
			}

			suite.Require().NoError(err)
		})
	}
}
