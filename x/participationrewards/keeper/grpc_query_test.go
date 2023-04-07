package keeper_test

import (
	encoding_json "encoding/json"
	"fmt"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (s *KeeperTestSuite) TestKeeper_Params() {
	s.Run("Params", func() {
		k := s.GetQuicksilverApp(s.chainA).ParticipationRewardsKeeper
		want := types.QueryParamsResponse{
			Params: types.DefaultParams(),
		}
		got, err := k.Params(s.chainA.GetContext(), &types.QueryParamsRequest{})
		s.Require().NoError(err)
		s.Require().NotNil(got)
		s.Require().Equal(want, *got)
	})
}

func (s *KeeperTestSuite) TestKeeper_ProtocolData() {
	connpdstr := fmt.Sprintf("{\"ConnectionID\":%q,\"ChainID\":%q,\"LastEpoch\":%d,\"Prefix\":\"\"}", s.path.EndpointB.ConnectionID, s.chainB.ChainID, 90767)
	s.Run("ProtocolData", func() {
		k := s.GetQuicksilverApp(s.chainA).ParticipationRewardsKeeper
		want := types.QueryProtocolDataResponse{
			Data: []encoding_json.RawMessage{
				[]byte(connpdstr),
			},
		}
		got, err := k.ProtocolData(
			s.chainA.GetContext(),
			&types.QueryProtocolDataRequest{
				Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
				Key:  s.chainB.ChainID,
			},
		)
		s.Require().NoError(err)
		s.Require().NotNil(got)
		s.Require().Equal(want, *got)
	})
}
