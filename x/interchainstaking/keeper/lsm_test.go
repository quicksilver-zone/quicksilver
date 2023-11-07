package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestLsmSetGetDelete() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	caps, found := icsKeeper.GetLsmCaps(ctx, suite.chainB.ChainID)
	suite.False(found)
	suite.Nil(caps)

	allCaps := icsKeeper.AllLsmCaps(ctx)
	suite.Equal(0, len(allCaps))

	icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID, types.LsmCaps{ValidatorCap: sdk.NewDecWithPrec(50, 2), GlobalCap: sdk.NewDecWithPrec(25, 2), ValidatorBondCap: sdk.NewDec(500)})

	caps, found = icsKeeper.GetLsmCaps(ctx, suite.chainB.ChainID)
	suite.True(found)
	suite.Equal(caps.ValidatorBondCap, sdk.NewDec(500))

	allCaps = icsKeeper.AllLsmCaps(ctx)
	suite.Equal(1, len(allCaps))

	icsKeeper.DeleteLsmCaps(ctx, suite.chainB.ChainID)

	caps, found = icsKeeper.GetLsmCaps(ctx, suite.chainB.ChainID)
	suite.False(found)
	suite.Nil(caps)

	allCaps = icsKeeper.AllLsmCaps(ctx)
	suite.Equal(0, len(allCaps))
}
