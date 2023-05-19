package keeper_test

import (
	"crypto/rand"
)

const (
	testChainID = "test-1"
)

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *KeeperTestSuite) TestKeeper_RemoteAddressStore() {
	s.SetupTest()
	s.setupTestZones()

	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	localAddress, err := generateRandomBytes(28)
	s.Require().NoError(err)
	remoteAddress, err := generateRandomBytes(40)
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

	localAddress, err := generateRandomBytes(28)
	s.Require().NoError(err)
	remoteAddress, err := generateRandomBytes(40)
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

	localAddress, err := generateRandomBytes(28)
	s.Require().NoError(err)
	remoteAddress, err := generateRandomBytes(40)
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
