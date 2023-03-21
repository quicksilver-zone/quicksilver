package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) TestCalcUserHoldingsAllocations() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	user1 := utils.GenerateAccAddressForTest()
	user2 := utils.GenerateAccAddressForTest()

	tests := []struct {
		name      string
		malleate  func(ctx sdk.Context, appA *app.Quicksilver)
		want      []keeper.UserAllocation
		remainder math.Int
		wantErr   string
	}{
		{
			"zero claims; no allocation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
			},
			[]keeper.UserAllocation{},
			sdk.ZeroInt(),
			"",
		},
		{
			"zero relevant claims; no allocation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: 500})
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: 1000})
			},
			[]keeper.UserAllocation{},
			sdk.ZeroInt(),
			"",
		},
		{
			"valid claims",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: 500})
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: 1000})
			},
			[]keeper.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewInt(2500),
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewInt(2500),
				},
			},
			sdk.ZeroInt(),
			"",
		},
	}
	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			ctx := suite.chainA.GetContext()
			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(suite.chainA.GetContext(), appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			suite.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			suite.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			allocations, remainder := appA.ParticipationRewardsKeeper.CalcUserHoldingsAllocations(ctx, &zone)
			suite.Require().ElementsMatch(tt.want, allocations)
			suite.Require().ElementsMatch(tt.remainder, remainder)
		})
	}
}
