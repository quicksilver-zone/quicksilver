package keeper_test

import (
	gocontext "context"
	"time"

	"github.com/quicksilver-zone/quicksilver/x/epochs/types"
)

func (s *KeeperTestSuite) TestQueryEpochInfos() {
	s.SetupTest()
	queryClient := s.queryClient

	chainStartTime := s.ctx.BlockTime()

	// Invalid param
	epochInfosResponse, err := queryClient.EpochInfos(gocontext.Background(), &types.QueryEpochsInfoRequest{})
	s.Require().NoError(err)
	s.Require().Len(epochInfosResponse.Epochs, 3)

	// check if EpochInfos are correct
	s.Require().Equal(epochInfosResponse.Epochs[0].Identifier, "day")
	s.Require().Equal(epochInfosResponse.Epochs[0].StartTime, chainStartTime)
	s.Require().Equal(epochInfosResponse.Epochs[0].Duration, time.Hour*24)
	s.Require().Equal(epochInfosResponse.Epochs[0].CurrentEpoch, int64(0))
	s.Require().Equal(epochInfosResponse.Epochs[0].CurrentEpochStartTime, chainStartTime)
	s.Require().Equal(epochInfosResponse.Epochs[0].EpochCountingStarted, false)
	s.Require().Equal(epochInfosResponse.Epochs[1].Identifier, "epoch")
	s.Require().Equal(epochInfosResponse.Epochs[1].StartTime, chainStartTime)
	s.Require().Equal(epochInfosResponse.Epochs[1].Duration, time.Second*240)
	s.Require().Equal(epochInfosResponse.Epochs[1].CurrentEpoch, int64(0))
	s.Require().Equal(epochInfosResponse.Epochs[1].CurrentEpochStartTime, chainStartTime)
	s.Require().Equal(epochInfosResponse.Epochs[1].EpochCountingStarted, false)
	s.Require().Equal(epochInfosResponse.Epochs[2].Identifier, "week")
	s.Require().Equal(epochInfosResponse.Epochs[2].StartTime, chainStartTime)
	s.Require().Equal(epochInfosResponse.Epochs[2].Duration, time.Hour*24*7)
	s.Require().Equal(epochInfosResponse.Epochs[2].CurrentEpoch, int64(0))
	s.Require().Equal(epochInfosResponse.Epochs[2].CurrentEpochStartTime, chainStartTime)
	s.Require().Equal(epochInfosResponse.Epochs[2].EpochCountingStarted, false)
}
