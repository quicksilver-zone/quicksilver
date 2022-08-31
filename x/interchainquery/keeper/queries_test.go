package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
)

func (suite *KeeperTestSuite) TestQuery() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	query := suite.GetSimApp(suite.chainA).InterchainQueryKeeper.NewQuery(
		suite.chainA.GetContext(),
		"",
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
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
	suite.Equal(sdk.NewInt(200), getQuery.Period)
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
