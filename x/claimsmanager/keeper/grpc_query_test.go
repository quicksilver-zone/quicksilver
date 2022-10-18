package keeper_test

import (
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

func (suite *KeeperTestSuite) TestKeeper_Params() {
	suite.Run("Params", func() {
		k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper
		want := types.QueryParamsResponse{
			Params: types.DefaultParams(),
		}
		got, err := k.Params(suite.chainA.GetContext(), &types.QueryParamsRequest{})
		suite.Require().NoError(err)
		suite.Require().NotNil(got)
		suite.Require().Equal(want, *got)
	})
}
