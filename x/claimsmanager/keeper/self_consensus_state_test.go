package keeper_test

func (s *KeeperTestSuite) TestGetSetDelete() {
	k := s.GetQuicksilverApp(s.chainA).ClaimsManagerKeeper
	ctx := s.chainA.GetContext()

	_, found := k.GetSelfConsensusState(ctx, "test")
	s.Require().False(found)

	err := k.StoreSelfConsensusState(ctx, "test")
	s.Require().NoError(err)

	_, found = k.GetSelfConsensusState(ctx, "test")
	s.Require().True(found)

	k.DeleteSelfConsensusState(ctx, "test")

	_, found = k.GetSelfConsensusState(ctx, "test")
	s.Require().False(found)
}
