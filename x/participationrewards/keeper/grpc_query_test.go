package keeper_test

import (
	encoding_json "encoding/json"
	"fmt"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) TestKeeper_Params() {
	suite.Run("Params", func() {
		k := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
		want := types.QueryParamsResponse{
			Params: types.DefaultParams(),
		}
		got, err := k.Params(suite.chainA.GetContext(), &types.QueryParamsRequest{})
		suite.Require().NoError(err)
		suite.Require().NotNil(got)
		suite.Require().Equal(want, *got)
	})
}

func (suite *KeeperTestSuite) TestKeeper_ProtocolData() {
	connstr := fmt.Sprintf("connection/%s", suite.chainB.ChainID)
	connpdstr := fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0)
	suite.Run("ProtocolData", func() {
		k := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
		want := types.QueryProtocolDataResponse{
			Data: []encoding_json.RawMessage{
				[]byte(connpdstr),
			},
		}
		got, err := k.ProtocolData(suite.chainA.GetContext(), &types.QueryProtocolDataRequest{Protocol: connstr})
		suite.Require().NoError(err)
		suite.Require().NotNil(got)
		suite.Require().Equal(want, *got)
	})
}
