package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/supply/keeper"
	"github.com/quicksilver-zone/quicksilver/x/supply/types"
)

func (suite *KeeperTestSuite) TestKeeper_Supply_Disabled() {
	suite.Run("Params", func() {
		ctx := suite.chainA.GetContext()
		k := suite.GetQuicksilverApp(suite.chainA).SupplyKeeper
		querier := keeper.NewQuerier(k)
		_, err := querier.Supply(ctx, &types.QuerySupplyRequest{})
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestKeeper_Supply() {
	suite.Run("Params", func() {
		ctx := suite.chainA.GetContext()

		k := suite.GetQuicksilverApp(suite.chainA).SupplyKeeper
		k.Enable(ctx, true)
		s, ok := math.NewIntFromString("100000000000058098352")
		suite.True(ok)
		cs, ok := math.NewIntFromString("100000000000004000000")
		suite.True(ok)
		want := types.QuerySupplyResponse{
			Supply:            s,
			CirculatingSupply: cs,
		}
		querier := keeper.NewQuerier(k)
		got, err := querier.Supply(ctx, &types.QuerySupplyRequest{})
		suite.NoError(err)
		suite.NotNil(got)
		fmt.Println("got", got)
		fmt.Println("want", want)
		suite.Equal(want, *got)
	})
}

func (suite *KeeperTestSuite) TestKeeper_Supply_Excluded_Account() {
	suite.Run("Params", func() {
		ctx := suite.chainA.GetContext()

		k := suite.GetQuicksilverApp(suite.chainA).SupplyKeeper
		bk := suite.GetQuicksilverApp(suite.chainA).BankKeeper
		stk := suite.GetQuicksilverApp(suite.chainA).StakingKeeper
		mk := suite.GetQuicksilverApp(suite.chainA).MintKeeper
		k.Enable(ctx, true)
		s, ok := math.NewIntFromString("100000000000063098352") // this included the 5stake minted below.
		suite.True(ok)
		cs, ok := math.NewIntFromString("100000000000004000000") // this does not include the new stake, as it is excluded.
		suite.True(ok)
		addr := addressutils.MustAccAddressFromBech32("quick1yxe3vmd2ypjf0fs4cejnmv2559tqq5x5cc5nyh", "")
		amount, ok := math.NewIntFromString("5000000")
		suite.True(ok)
		// add coins to an excluded account.
		err := mk.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(stk.BondDenom(ctx), amount)))
		suite.NoError(err)
		err = bk.SendCoinsFromModuleToAccount(ctx, "mint", addr, sdk.NewCoins(sdk.NewCoin(stk.BondDenom(ctx), amount)))
		suite.NoError(err)
		want := types.QuerySupplyResponse{
			Supply:            s,
			CirculatingSupply: cs,
		}
		querier := keeper.NewQuerier(k)
		got, err := querier.Supply(ctx, &types.QuerySupplyRequest{})
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(want, *got)
	})
}

func (suite *KeeperTestSuite) TestKeeper_TopN_Disabled() {
	suite.Run("Params", func() {
		ctx := suite.chainA.GetContext()
		k := suite.GetQuicksilverApp(suite.chainA).SupplyKeeper
		querier := keeper.NewQuerier(k)
		_, err := querier.TopN(ctx, &types.QueryTopNRequest{N: 5})
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestKeeper_TopN() {
	suite.Run("Params", func() {
		ctx := suite.chainA.GetContext()

		k := suite.GetQuicksilverApp(suite.chainA).SupplyKeeper
		k.Enable(ctx, true)
		bk := suite.GetQuicksilverApp(suite.chainA).BankKeeper
		stk := suite.GetQuicksilverApp(suite.chainA).StakingKeeper
		mk := suite.GetQuicksilverApp(suite.chainA).MintKeeper
		querier := keeper.NewQuerier(k)
		// all accounts (random addrs) have 10000000000000000000 +/- 4 stake. Lets create an account that'll sit at the top of the topN.
		addr := addressutils.GenerateAccAddressForTest()
		amount, ok := math.NewIntFromString("20000000000000000004")
		suite.True(ok)
		err := mk.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(stk.BondDenom(ctx), amount)))
		suite.NoError(err)
		err = bk.SendCoinsFromModuleToAccount(ctx, "mint", addr, sdk.NewCoins(sdk.NewCoin(stk.BondDenom(ctx), amount)))
		suite.NoError(err)
		topN, err := querier.TopN(ctx, &types.QueryTopNRequest{N: 5})
		suite.NoError(err)
		suite.Equal(5, len(topN.Accounts))
		suite.Equal(addr.String(), topN.Accounts[0].Address)
		suite.Equal(amount, topN.Accounts[0].Balance)
	})
}
