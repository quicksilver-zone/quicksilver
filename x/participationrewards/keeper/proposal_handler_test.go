package keeper_test

import (
	"encoding/base64"
	"fmt"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/keeper"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
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
				connpdstr := fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d,}", "", "", 100)

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
				connpdstr := fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d, \"prefix\": \"cosmos\", \"transferchannel\": %q}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0, "channel-1")

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
				suite.Error(err)
				// suite.T().Logf("Error: %v", err)
				return
			}

			suite.NoError(err)
		})
	}
}

func (suite *KeeperTestSuite) TestHandleRemoveProtocolDataProposal() {
	appA := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	k := appA.ParticipationRewardsKeeper

	// set the protocol data
	pd := types.ConnectionProtocolData{
		ConnectionID: suite.path.EndpointB.ConnectionID,
		ChainID:      suite.chainB.ChainID,
		LastEpoch:    0,
		Prefix:       "cosmos",
	}
	err := keeper.MarshalAndSetProtocolData(ctx, k, types.ProtocolDataTypeConnection, &pd)
	suite.NoError(err)

	// check if the protocol data is set
	_, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, string(pd.GenerateKey()))
	suite.True(found)

	msgServer := keeper.NewMsgServerImpl(k)

	// submit proposal to remove the protocol data
	proposalMsg := types.MsgGovRemoveProtocolData{
		Title:       "remove chain B connection string",
		Description: "remove the protocol data",
		Key:         base64.StdEncoding.EncodeToString(types.GetProtocolDataKey(types.ProtocolDataTypeConnection, pd.GenerateKey())),
		Authority:   k.GetGovAuthority(ctx),
	}

	_, err = msgServer.GovRemoveProtocolData(ctx, &proposalMsg)
	suite.NoError(err)

	_, found = k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, string(pd.GenerateKey()))
	suite.False(found)
}
