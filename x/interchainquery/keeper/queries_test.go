package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/interchainquery/keeper"
)

func (suite *KeeperTestSuite) TestQuery() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	query := suite.GetSimApp(suite.chainA).InterchainQueryKeeper.NewQuery(
		"",
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdkmath.NewInt(200),
		"",
		0,
	)

	// set the query
	suite.GetSimApp(suite.chainA).InterchainQueryKeeper.SetQuery(suite.chainA.GetContext(), *query)

	// get the stored query
	id := keeper.GenerateQueryHash(query.ConnectionId, query.ChainId, query.QueryType, query.Request, "")
	getQuery, found := suite.GetSimApp(suite.chainA).InterchainQueryKeeper.GetQuery(suite.chainA.GetContext(), id)
	suite.True(found)
	suite.Equal(suite.path.EndpointB.ConnectionID, getQuery.ConnectionId)
	suite.Equal(suite.chainB.ChainID, getQuery.ChainId)
	suite.Equal("cosmos.staking.v1beta1.Query/Validators", getQuery.QueryType)
	suite.Equal(sdkmath.NewInt(200), getQuery.Period)
	suite.Equal(uint64(0), getQuery.Ttl)
	suite.Equal("", getQuery.CallbackId)

	// get all the queries
	queries := suite.GetSimApp(suite.chainA).InterchainQueryKeeper.AllQueries(suite.chainA.GetContext())
	suite.Len(queries, 1)

	// delete the query
	suite.GetSimApp(suite.chainA).InterchainQueryKeeper.DeleteQuery(suite.chainA.GetContext(), id)

	// get query
	_, found = suite.GetSimApp(suite.chainA).InterchainQueryKeeper.GetQuery(suite.chainA.GetContext(), id)
	suite.False(found)

	queries = suite.GetSimApp(suite.chainA).InterchainQueryKeeper.AllQueries(suite.chainA.GetContext())
	suite.Len(queries, 0)
}
