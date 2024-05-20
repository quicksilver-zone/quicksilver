package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/utils/randomutils"
)

const (
	testChainID = "test-1"
)

var (
	localAddress  = sdk.AccAddress(randomutils.GenerateRandomBytes(20))
	remoteAddress = sdk.AccAddress(randomutils.GenerateRandomBytes(32))
)

func (suite *KeeperTestSuite) TestKeeper_RemoteAddressStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	suite.Run("not found", func() {
		_, found := icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		suite.False(found)
	})

	suite.Run("set", func() {
		icsKeeper.SetRemoteAddressMap(ctx, localAddress, remoteAddress, testChainID)
	})

	suite.Run("found", func() {
		addr, found := icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		suite.True(found)
		suite.Equal(remoteAddress, addr)
	})
}

func (suite *KeeperTestSuite) TestKeeper_LocalAddressStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	suite.Run("not found", func() {
		_, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.False(found)
	})

	suite.Run("set", func() {
		icsKeeper.SetLocalAddressMap(ctx, localAddress, remoteAddress, testChainID)
	})

	suite.Run("found", func() {
		addr, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.True(found)
		suite.Equal(localAddress, addr)
	})
}

func (suite *KeeperTestSuite) TestKeeper_AddressMapPair() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	suite.Run("not found", func() {
		_, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.False(found)
		_, found = icsKeeper.GetRemoteAddressMap(ctx, remoteAddress, testChainID)
		suite.False(found)
	})

	suite.Run("set", func() {
		icsKeeper.SetAddressMapPair(ctx, localAddress, remoteAddress, testChainID)
	})

	suite.Run("found", func() {
		addr, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		suite.True(found)
		suite.Equal(localAddress, addr)
		addr, found = icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		suite.True(found)
		suite.Equal(remoteAddress, addr)
	})
}
