package keeper_test

import (
	"fmt"
	"time"

	"github.com/ingenuity-build/quicksilver/x/epochs"
	"github.com/ingenuity-build/quicksilver/x/epochs/types"
)

func (s *KeeperTestSuite) TestEpochInfoChangesBeginBlockerAndInitGenesis() {
	var epochInfo types.EpochInfo

	now := time.Now()

	testCases := []struct {
		expCurrentEpochStartTime   time.Time
		expCurrentEpochStartHeight int64
		expCurrentEpoch            int64
		expInitialEpochStartTime   time.Time
		fn                         func()
	}{
		{
			// Only advance 2 seconds, do not increment epoch
			expCurrentEpochStartHeight: 2,
			expCurrentEpochStartTime:   now,
			expCurrentEpoch:            1,
			expInitialEpochStartTime:   now,
			fn: func() {
				s.ctx = s.ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
			},
		},
		{
			expCurrentEpochStartHeight: 2,
			expCurrentEpochStartTime:   now,
			expCurrentEpoch:            1,
			expInitialEpochStartTime:   now,
			fn: func() {
				s.ctx = s.ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
			},
		},
		{
			expCurrentEpochStartHeight: 2,
			expCurrentEpochStartTime:   now,
			expCurrentEpoch:            1,
			expInitialEpochStartTime:   now,
			fn: func() {
				s.ctx = s.ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				s.ctx = s.ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 31))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
			},
		},
		// Test that incrementing _exactly_ 1 month increments the epoch count.
		{
			expCurrentEpochStartHeight: 3,
			expCurrentEpochStartTime:   now.Add(time.Hour * 24 * 31),
			expCurrentEpoch:            2,
			expInitialEpochStartTime:   now,
			fn: func() {
				s.ctx = s.ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				s.ctx = s.ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 32))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
			},
		},
		{
			expCurrentEpochStartHeight: 3,
			expCurrentEpochStartTime:   now.Add(time.Hour * 24 * 31),
			expCurrentEpoch:            2,
			expInitialEpochStartTime:   now,
			fn: func() {
				s.ctx = s.ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				s.ctx = s.ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 32))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				s.ctx.WithBlockHeight(4).WithBlockTime(now.Add(time.Hour * 24 * 33))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
			},
		},
		{
			expCurrentEpochStartHeight: 3,
			expCurrentEpochStartTime:   now.Add(time.Hour * 24 * 31),
			expCurrentEpoch:            2,
			expInitialEpochStartTime:   now,
			fn: func() {
				s.ctx = s.ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				s.ctx = s.ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 32))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				s.ctx.WithBlockHeight(4).WithBlockTime(now.Add(time.Hour * 24 * 33))
				s.app.EpochsKeeper.BeginBlocker(s.ctx)
				epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
			},
		},
	}

	for i, tc := range testCases {
		s.Run(fmt.Sprintf("Case %d", i), func() {
			s.SetupTest() // reset

			// On init genesis, default epochs information is set
			// To check init genesis again, should make it fresh status
			epochInfos := s.app.EpochsKeeper.AllEpochInfos(s.ctx)
			for _, epochInfo := range epochInfos {
				s.app.EpochsKeeper.DeleteEpochInfo(s.ctx, epochInfo.Identifier)
			}

			s.ctx = s.ctx.WithBlockHeight(1).WithBlockTime(now)

			// check init genesis
			epochs.InitGenesis(s.ctx, s.app.EpochsKeeper, types.GenesisState{
				Epochs: []types.EpochInfo{
					{
						Identifier:              "monthly",
						StartTime:               time.Time{},
						Duration:                time.Hour * 24 * 31,
						CurrentEpoch:            0,
						CurrentEpochStartHeight: s.ctx.BlockHeight(),
						CurrentEpochStartTime:   time.Time{},
						EpochCountingStarted:    false,
					},
				},
			})

			tc.fn()

			s.Require().Equal(epochInfo.Identifier, "monthly")
			s.Require().Equal(epochInfo.StartTime.UTC().String(), tc.expInitialEpochStartTime.UTC().String())
			s.Require().Equal(epochInfo.Duration, time.Hour*24*31)
			s.Require().Equal(epochInfo.CurrentEpoch, tc.expCurrentEpoch)
			s.Require().Equal(epochInfo.CurrentEpochStartHeight, tc.expCurrentEpochStartHeight)
			s.Require().Equal(epochInfo.CurrentEpochStartTime.UTC().String(), tc.expCurrentEpochStartTime.UTC().String())
			s.Require().Equal(epochInfo.EpochCountingStarted, true)
		})
	}
}

func (s *KeeperTestSuite) TestEpochStartingOneMonthAfterInitGenesis() {
	// On init genesis, default epochs information is set
	// To check init genesis again, should make it fresh status
	epochInfos := s.app.EpochsKeeper.AllEpochInfos(s.ctx)
	for _, epochInfo := range epochInfos {
		s.app.EpochsKeeper.DeleteEpochInfo(s.ctx, epochInfo.Identifier)
	}

	now := time.Now()
	week := time.Hour * 24 * 7
	month := time.Hour * 24 * 30
	initialBlockHeight := int64(1)
	s.ctx = s.ctx.WithBlockHeight(initialBlockHeight).WithBlockTime(now)

	epochs.InitGenesis(s.ctx, s.app.EpochsKeeper, types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               now.Add(month),
				Duration:                time.Hour * 24 * 30,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: s.ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    false,
			},
		},
	})

	// epoch not started yet
	epochInfo := s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
	s.Require().Equal(epochInfo.CurrentEpoch, int64(0))
	s.Require().Equal(epochInfo.CurrentEpochStartHeight, initialBlockHeight)
	s.Require().Equal(epochInfo.CurrentEpochStartTime, time.Time{})
	s.Require().Equal(epochInfo.EpochCountingStarted, false)

	// after 1 week
	s.ctx = s.ctx.WithBlockHeight(2).WithBlockTime(now.Add(week))
	s.app.EpochsKeeper.BeginBlocker(s.ctx)

	// epoch not started yet
	epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
	s.Require().Equal(epochInfo.CurrentEpoch, int64(0))
	s.Require().Equal(epochInfo.CurrentEpochStartHeight, initialBlockHeight)
	s.Require().Equal(epochInfo.CurrentEpochStartTime, time.Time{})
	s.Require().Equal(epochInfo.EpochCountingStarted, false)

	// after 1 month
	s.ctx = s.ctx.WithBlockHeight(3).WithBlockTime(now.Add(month))
	s.app.EpochsKeeper.BeginBlocker(s.ctx)

	// epoch started
	epochInfo = s.app.EpochsKeeper.GetEpochInfo(s.ctx, "monthly")
	s.Require().Equal(epochInfo.CurrentEpoch, int64(1))
	s.Require().Equal(epochInfo.CurrentEpochStartHeight, s.ctx.BlockHeight())
	s.Require().Equal(epochInfo.CurrentEpochStartTime.UTC().String(), now.Add(month).UTC().String())
	s.Require().Equal(epochInfo.EpochCountingStarted, true)
}
