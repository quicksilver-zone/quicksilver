package keeper_test

import (
	"encoding/json"
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
					Type:        "",
					Data:        nil,
					Key:         "",
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
					Type:        "testtype",
					Data:        []byte("{}"),
					Key:         "",
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
					Type:        types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
					Data:        []byte("{}"),
					Key:         "",
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
					Type:        types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
					Data:        []byte(connpdstr),
					Key:         "",
				}
			},
			true,
		},
		{
			"valid_prop",
			func() {
				connpdstr := fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d, \"prefix\": \"cosmos\"}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0)

				prop = types.AddProtocolDataProposal{
					Title:       "Add connection protocol for test chain B",
					Description: "A connection protocol for testing connection protocols",
					Type:        types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
					Data:        []byte(connpdstr),
					Key:         "",
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

func (suite *KeeperTestSuite) TestHandleRemoveProtocolDataProposal() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	pd := types.ConnectionProtocolData{
		ConnectionID: suite.path.EndpointB.ConnectionID,
		ChainID:      suite.chainB.ChainID,
		LastEpoch:    0,
		Prefix:       "cosmos",
	}

	pdString, err := json.Marshal(pd)
	suite.Require().NoError(err)

	ctx := suite.chainA.GetContext()

	prop := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
		Data: pdString,
	}

	k := appA.ParticipationRewardsKeeper

	k.SetProtocolData(ctx, pd.GenerateKey(), &prop)
	// set the protocol data

	_, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, string(pd.GenerateKey()))
	suite.Require().True(found)

	msgServer := keeper.NewMsgServerImpl(k)

	// submit proposal

	proposalMsg := types.MsgGovRemoveProtocolData{
		Title:       "remove chain B connection string",
		Description: "remove the protocol data",
		Key:         string(pd.GenerateKey()),
		Authority:   k.GetGovAuthority(ctx),
	}

	_, err = msgServer.GovRemoveProtocolData(ctx, &proposalMsg)

	suite.Require().NoError(err)

	_, found = k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, string(pd.GenerateKey()))
	suite.Require().True(found)
}
