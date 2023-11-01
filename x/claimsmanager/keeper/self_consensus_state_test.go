package keeper_test

func (suite *KeeperTestSuite) TestGetSetDelete() {
	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper
	ctx := suite.chainA.GetContext()

	_, found := k.GetSelfConsensusState(ctx, "test")
	suite.Require().False(found)

	err := k.StoreSelfConsensusState(ctx, "test")
	suite.Require().NoError(err)

	_, found = k.GetSelfConsensusState(ctx, "test")
	suite.Require().True(found)

	k.DeleteSelfConsensusState(ctx, "test")

	_, found = k.GetSelfConsensusState(ctx, "test")
	suite.Require().False(found)
}
