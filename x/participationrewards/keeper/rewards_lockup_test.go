package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/app"
)

func (suite *KeeperTestSuite) TestAllocateLockupRewards() {
	tests := []struct {
		name          string
		getAllocation func(types.Context, *app.Quicksilver) int64
		wantErr       bool
	}{
		{
			"valid unit allocation",
			func(ctx types.Context, appA *app.Quicksilver) int64 {
				return 1
			},
			false,
		},
		{
			"invalid -ve allocation",
			func(ctx types.Context, appA *app.Quicksilver) int64 {
				return -1
			},
			true,
		},
		{
			"redundant allocation",
			func(ctx types.Context, appA *app.Quicksilver) int64 {
				return 0
			},
			false,
		},
		{
			"valid total module balance allocation",
			func(ctx types.Context, appA *app.Quicksilver) int64 {
				return appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64()
			},
			false,
		},
		{
			"invalid greater than module balance allocation",
			func(ctx types.Context, appA *app.Quicksilver) int64 {
				return appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64() + 1
			},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()

			appA := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			initialBalance := appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64()
			allocation := tt.getAllocation(ctx, appA)

			if err := appA.ParticipationRewardsKeeper.AllocateLockupRewards(ctx, math.NewInt(allocation)); err != nil {
				suite.Require().True((err != nil) == tt.wantErr)
				return
			}
			finalBalance := appA.ParticipationRewardsKeeper.GetModuleBalance(ctx).Int64()
			suite.Require().Equal(initialBalance-finalBalance, allocation)
		})
	}
}
