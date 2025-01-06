package keeper_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) TestCalcUserHoldingsAllocations() {
	user1 := addressutils.GenerateAccAddressForTest()
	user2 := addressutils.GenerateAccAddressForTest()
	appA := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	bondDenom := appA.StakingKeeper.BondDenom(ctx)
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
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
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
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 64000
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]types.UserAllocation{},
			sdk.NewInt(64000),
			"",
		},
		{
			"valid claims - equal claims",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(5000)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(2500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(2500)})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(2500)),
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(2500)),
				},
			},
			sdk.ZeroInt(),
			"",
		},
		{
			"valid claims - inequal claims, less than 100%, truncation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(2500)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(1000)), // 500 / 2500 (0.2) * 5000 = 1000
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(2000)), // 1000 / 2500 (0.4) * 5000 = 2000
				},
			},
			sdk.NewInt(2000),
			"",
		},
		{
			"valid claims - inequal claims, 100%, truncation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(1500)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(1666)), // 500/1500 (0.33333) * 5000 == 1666
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(3333)), // 1000/1500 (0.66666) * 5000 = 3333
				},
			},
			sdk.OneInt(),
			"",
		},
		{
			"valid claims - inequal claims, 100%, truncation, plus ics",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000
				icsAddress, _ := addressutils.AddressFromBech32(zone.WithdrawalAddress.Address, "")
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(1500)))))
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewIntFromUint64(900)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", icsAddress, sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewIntFromUint64(900)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(1666)), // 500/1500 (0.33333) * 5000 == 1666
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(3333)), // 1000/1500 (0.66666) * 5000 = 3333
				},
			},
			sdk.OneInt(),
			"",
		},
		{
			"valid claims - inequal claims, 100%, truncation, plus multiple ics + overflow",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000
				icsAddress, _ := addressutils.AddressFromBech32(zone.WithdrawalAddress.Address, "")
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(1500)))))
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewIntFromUint64(900)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", icsAddress, sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewIntFromUint64(900)))))
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("testcoin2", sdk.NewIntFromUint64(18002)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", icsAddress, sdk.NewCoins(sdk.NewCoin("testcoin2", sdk.NewIntFromUint64(18002)))))
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("testcoin3", sdk.NewIntFromUint64(150)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", icsAddress, sdk.NewCoins(sdk.NewCoin("testcoin3", sdk.NewIntFromUint64(150)))))

				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(1666)), // 500/1500 (0.33333) * 5000 == 1666
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(3333)), // 1000/1500 (0.66666) * 5000 = 3333
				},
			},
			sdk.OneInt(),
			"",
		},

		{
			"valid claims - inequal claims, less than 100%, truncation + ics + overflow",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				icsAddress, _ := addressutils.AddressFromBech32(zone.WithdrawalAddress.Address, "")
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(2500)))))
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewIntFromUint64(900)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", icsAddress, sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewIntFromUint64(900)))))
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("testcoin2", sdk.NewIntFromUint64(18002)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", icsAddress, sdk.NewCoins(sdk.NewCoin("testcoin2", sdk.NewIntFromUint64(18002)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]types.UserAllocation{
				{
					Address: user1.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(1000)), // 500 / 2500 (0.2) * 5000 = 1000
				},
				{
					Address: user2.String(),
					Amount:  sdk.NewCoin(bondDenom, sdk.NewInt(2000)), // 1000 / 2500 (0.4) * 5000 = 2000
				},
			},
			sdk.NewInt(2000),
			"",
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			appA := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(suite.chainA.GetContext(), appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			suite.NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			allocations, remainder := appA.ParticipationRewardsKeeper.CalcUserHoldingsAllocations(ctx, &zone)
			suite.ElementsMatch(tt.want, allocations)
			suite.True(tt.remainder.Equal(remainder))
		})
	}
}

func (suite *KeeperTestSuite) TestAllocateHoldingsRewards() {
	user1 := addressutils.GenerateAccAddressForTest()
	user2 := addressutils.GenerateAccAddressForTest()

	tests := []struct {
		name     string
		malleate func(ctx sdk.Context, appA *app.Quicksilver)
		balances []string
	}{
		{
			"zero claims; no allocation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 0
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
			},
			[]string{"0", "0"},
		},
		{
			"zero relevant claims; 64k allocation, all returned",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 64000
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: "otherchain-1", Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]string{"0", "0"},
		},
		{
			"valid claims - equal claims",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(5000)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(2500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(2500)})
			},
			[]string{"2500", "2500"},
		},
		{
			"valid claims - inequal claims, less than 100%, truncation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(2500)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]string{"1000", "2000"},
		},
		{
			"valid claims - inequal claims, 100%, truncation",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				zone, _ := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				zone.HoldingsAllocation = 5000

				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewIntFromUint64(1500)))))
				appA.InterchainstakingKeeper.SetZone(ctx, &zone)
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(500)})
				appA.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user2.String(), ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: suite.chainA.ChainID, Amount: math.NewInt(1000)})
			},
			[]string{"1666", "3333"},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			appA := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			tt.malleate(suite.chainA.GetContext(), appA)

			zone, found := appA.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
			suite.NoError(appA.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(appA.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))

			err := appA.ParticipationRewardsKeeper.AllocateHoldingsRewards(ctx)
			suite.True(err == nil)
			user1Balance := appA.BankKeeper.GetBalance(ctx, user1, appA.StakingKeeper.BondDenom(ctx)).Amount.String()
			user2Balance := appA.BankKeeper.GetBalance(ctx, user2, appA.StakingKeeper.BondDenom(ctx)).Amount.String()
			suite.Equal(tt.balances, []string{user1Balance, user2Balance})
		})
	}
}
