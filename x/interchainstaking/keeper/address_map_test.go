package keeper_test

import (
	"github.com/ingenuity-build/quicksilver/utils/randomutils"
)

const (
	testChainID = "test-1"
)

func (s *KeeperTestSuite) TestKeeper_RemoteAddressStore() {
	s.SetupTest()
	s.setupTestZones()

	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	localAddress := randomutils.GenerateRandomBytes(28)
	remoteAddress := randomutils.GenerateRandomBytes(40)

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

	localAddress := randomutils.GenerateRandomBytes(28)
	remoteAddress := randomutils.GenerateRandomBytes(40)

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

	localAddress := randomutils.GenerateRandomBytes(28)
	remoteAddress := randomutils.GenerateRandomBytes(40)

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
