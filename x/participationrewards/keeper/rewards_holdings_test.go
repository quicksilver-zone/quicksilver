package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (s *KeeperTestSuite) TestCalcUserHoldingsAllocations() {
	user1 := utils.GenerateAccAddressForTest()
	user2 := utils.GenerateAccAddressForTest()

	tests := []struct {
		name      string
		malleate  func(ctx sdk.Context, appA *app.Quicksilver)
		want      []types.UserAllocation
		remainder math.Int
		wantErr   string
	}{
		{
			"zero claims; no allocation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				zone.HoldingsAllocation = 0
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			[]types.UserAllocation{},
			sdk.ZeroInt(),
			"",
		},
		{
			"zero relevant claims; 64k allocation, all returned",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				zone.HoldingsAllocation = 64000
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 500})
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 1000})
			},
			[]types.UserAllocation{},
			sdk.NewInt(64000),
			"",
		},
		{
			"valid claims - equal claims",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(5000)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 2500})
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 2500})
			},
			[]types.UserAllocation{
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
		{
			"valid claims - inequal claims, less than 100%, truncation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(2500)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 500})
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 1000})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewInt(1000), // 500 / 2500 (0.2) * 5000 = 1000
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewInt(2000), // 1000 / 2500 (0.4) * 5000 = 2000
				},
			},
			sdk.NewInt(2000),
			"",
		},
		{
			"valid claims - inequal claims, 100%, truncation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(1500)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 500})
				appA.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: s.chainA.ChainID, Amount: 1000})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewInt(1666), // 500/1500 (0.33333) * 5000 == 1666
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewInt(3333), // 1000/1500 (0.66666) * 5000 = 3333
				},
			},
			sdk.OneInt(),
			"",
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

			tt.malleate(s.chainA.GetContext(), appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			s.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			allocations, remainder := appA.ParticipationRewardsKeeper.CalcUserHoldingsAllocations(ctx, &zone)
			s.Require().ElementsMatch(tt.want, allocations)
			s.Require().True(tt.remainder.Equal(remainder))
		})
	}
}
