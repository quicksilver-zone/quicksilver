package keeper_test

import "github.com/ingenuity-build/quicksilver/utils"

const (
	testChainID = "test-1"
)

func (s *KeeperTestSuite) TestKeeper_RemoteAddressStore() {
	s.SetupTest()
	s.setupTestZones()

	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	localAddress, err := utils.GenerateRandomBytes(28)
	s.Require().NoError(err)
	remoteAddress, err := utils.GenerateRandomBytes(40)
	s.Require().NoError(err)

	s.Run("not found", func() {
		_, found := icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		s.Require().False(found)
	})

	s.Run("set", func() {
		icsKeeper.SetRemoteAddressMap(ctx, localAddress, remoteAddress, testChainID)
	})

	s.Run("found", func() {
		addr, found := icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		s.Require().True(found)
		s.Require().Equal(remoteAddress, addr)
	})
}

func (s *KeeperTestSuite) TestKeeper_LocalAddressStore() {
	s.SetupTest()
	s.setupTestZones()

	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	localAddress, err := utils.GenerateRandomBytes(28)
	s.Require().NoError(err)
	remoteAddress, err := utils.GenerateRandomBytes(40)
	s.Require().NoError(err)

	s.Run("not found", func() {
		_, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		s.Require().False(found)
	})

	s.Run("set", func() {
		icsKeeper.SetLocalAddressMap(ctx, localAddress, remoteAddress, testChainID)
	})

	s.Run("found", func() {
		addr, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		s.Require().True(found)
		s.Require().Equal(localAddress, addr)
	})
}

func (s *KeeperTestSuite) TestKeeper_AddressMapPair() {
	s.SetupTest()
	s.setupTestZones()

	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	localAddress, err := utils.GenerateRandomBytes(28)
	s.Require().NoError(err)
	remoteAddress, err := utils.GenerateRandomBytes(40)
	s.Require().NoError(err)

	s.Run("not found", func() {
		_, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		s.Require().False(found)
		_, found = icsKeeper.GetRemoteAddressMap(ctx, remoteAddress, testChainID)
		s.Require().False(found)
	})

	s.Run("set", func() {
		icsKeeper.SetAddressMapPair(ctx, localAddress, remoteAddress, testChainID)
	})

	s.Run("found", func() {
		addr, found := icsKeeper.GetLocalAddressMap(ctx, remoteAddress, testChainID)
		s.Require().True(found)
		s.Require().Equal(localAddress, addr)
		addr, found = icsKeeper.GetRemoteAddressMap(ctx, localAddress, testChainID)
		s.Require().True(found)
		s.Require().Equal(remoteAddress, addr)
	})
}
