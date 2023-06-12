package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/app"
)

func (s *KeeperTestSuite) TestAllocateLockupRewards() {

	tests := []struct {
		name          string
		getAllocation func(*app.Quicksilver, types.Context) int64
		wantErr       bool
		wantPanic     bool
	}{
		{
			"valid unit allocation",
			func(appA *app.Quicksilver, ctx types.Context) int64 {
				return 1
			},
			false,
			false,
		},
		{
			"panic case -ve allocation",
			func(appA *app.Quicksilver, ctx types.Context) int64 {
				return -1
			},
			false,
			true,
		},
		{
			"redundant allocation",
			func(appA *app.Quicksilver, ctx types.Context) int64 {
				return 0
			},
			false,
			false,
		},
		{
			"valid total module balance allocation",
			func(appA *app.Quicksilver, ctx types.Context) int64 {
				return appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64()
			},
			false,
			false,
		},
		{
			"invalid greater than module balance allocation",
			func(appA *app.Quicksilver, ctx types.Context) int64 {
				return appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64() + 1
			},
			true,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()

			appA := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			initialBalance := appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64()
			allocation := tt.getAllocation(appA, ctx)

			if tt.wantPanic {
				s.Require().True(s.Panics(func() {
					appA.ParticipationRewardsKeeper.AllocateLockupRewards(ctx, math.NewInt(allocation))
				}))
				return
			}
			if err := appA.ParticipationRewardsKeeper.AllocateLockupRewards(ctx, math.NewInt(allocation)); err != nil {
				s.Require().True((err != nil) == tt.wantErr)
				return
			}
			finalBalance := appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64()
			s.Require().Equal(initialBalance-finalBalance, allocation)
		})
	}
}
