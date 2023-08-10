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
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(want, *got)
	})
}

func (suite *KeeperTestSuite) TestKeeper_ProtocolData() {
	connpdstr := fmt.Sprintf("{\"ConnectionID\":%q,\"ChainID\":%q,\"LastEpoch\":%d,\"Prefix\":\"\"}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 90767)
	suite.Run("ProtocolData", func() {
		k := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
		want := types.QueryProtocolDataResponse{
			Data: []encoding_json.RawMessage{
				[]byte(connpdstr),
			},
		}
		got, err := k.ProtocolData(
			suite.chainA.GetContext(),
			&types.QueryProtocolDataRequest{
				Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
				Key:  suite.chainB.ChainID,
			},
		)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(want, *got)
	})
}
