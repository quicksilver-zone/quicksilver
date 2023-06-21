package keeper_test

import (
	"github.com/ingenuity-build/quicksilver/utils/randomutils"
)

const (
	testChainID = "test-1"
)

func (suite *KeeperTestSuite) TestKeeper_RemoteAddressStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	localAddress := randomutils.GenerateRandomBytes(28)
	remoteAddress := randomutils.GenerateRandomBytes(40)

	suite.Run("not found", func() {
		_, found := icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		suite.Require().False(found)
	})

	suite.Run("set", func() {
		icsKeeper.SetRemoteAddressMap(ctx, localAddress, remoteAddress, testChainID)
	})

	suite.Run("found", func() {
		addr, found := icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		suite.Require().True(found)
		suite.Require().Equal(remoteAddress, addr)
	})
}

func (suite *KeeperTestSuite) TestKeeper_LocalAddressStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	localAddress := randomutils.GenerateRandomBytes(28)
	remoteAddress := randomutils.GenerateRandomBytes(40)

	suite.Run("not found", func() {
		_, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.Require().False(found)
	})

	suite.Run("set", func() {
		icsKeeper.SetLocalAddressMap(ctx, localAddress, remoteAddress, testChainID)
	})

	suite.Run("found", func() {
		addr, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.Require().True(found)
		suite.Require().Equal(localAddress, addr)
	})
}

func (suite *KeeperTestSuite) TestKeeper_AddressMapPair() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	localAddress := randomutils.GenerateRandomBytes(28)
	remoteAddress := randomutils.GenerateRandomBytes(40)

	suite.Run("not found", func() {
		_, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.Require().False(found)
		_, found = icsKeeper.GetRemoteAddressMap(ctx, remoteAddress, testChainID)
		suite.Require().False(found)
	})

	suite.Run("set", func() {
		icsKeeper.SetAddressMapPair(ctx, localAddress, remoteAddress, testChainID)
	})

	suite.Run("found", func() {
		addr, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.Require().True(found)
		suite.Require().Equal(localAddress, addr)
		addr, found = icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		suite.Require().True(found)
		suite.Require().Equal(remoteAddress, addr)
	})
}
