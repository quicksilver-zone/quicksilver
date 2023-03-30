package keeper_test

import (
	"time"

	"github.com/ingenuity-build/quicksilver/x/epochs/types"
)

func (s *KeeperTestSuite) TestEpochLifeCycle() {
	s.SetupTest()

	epochInfo := types.EpochInfo{
		Identifier:            "monthly",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24 * 30,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}
	s.app.EpochsKeeper.SetEpochInfo(s.ctx, epochInfo)
	epochInfoSaved := s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
	s.Require().Equal(epochInfo, epochInfoSaved)

	allEpochs := s.app.EpochsKeeper.AllEpochInfos(s.ctx)
	s.Require().Len(allEpochs, 4)
	s.Require().Equal(allEpochs[0].Identifier, "day") // alphabetical order
	s.Require().Equal(allEpochs[1].Identifier, "epoch")
	s.Require().Equal(allEpochs[2].Identifier, "monthly")
	s.Require().Equal(allEpochs[3].Identifier, "week")
}
