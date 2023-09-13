package keeper_test

import (
	"github.com/quicksilver-zone/quicksilver/x/interchainquery/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *KeeperTestSuite) TestQuery() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	s.NoError(err)

	query := s.GetSimApp(s.chainA).InterchainQueryKeeper.NewQuery(
		"",
		s.path.EndpointB.ConnectionID,
		s.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		0,
	)

	// set the query
	s.GetSimApp(s.chainA).InterchainQueryKeeper.SetQuery(s.chainA.GetContext(), *query)

	// get the stored query
	id := keeper.GenerateQueryHash(query.ConnectionId, query.ChainId, query.QueryType, query.Request, "")
	getQuery, found := s.GetSimApp(s.chainA).InterchainQueryKeeper.GetQuery(s.chainA.GetContext(), id)
	s.True(found)
	s.Equal(s.path.EndpointB.ConnectionID, getQuery.ConnectionId)
	s.Equal(s.chainB.ChainID, getQuery.ChainId)
	s.Equal("cosmos.staking.v1beta1.Query/Validators", getQuery.QueryType)
	s.Equal(sdk.NewInt(200), getQuery.Period)
	s.Equal(uint64(0), getQuery.Ttl)
	s.Equal("", getQuery.CallbackId)

	// get all the queries
	queries := s.GetSimApp(s.chainA).InterchainQueryKeeper.AllQueries(s.chainA.GetContext())
	s.Len(queries, 1)

	// delete the query
	s.GetSimApp(s.chainA).InterchainQueryKeeper.DeleteQuery(s.chainA.GetContext(), id)

	// get query
	_, found = s.GetSimApp(s.chainA).InterchainQueryKeeper.GetQuery(s.chainA.GetContext(), id)
	s.False(found)

	queries = s.GetSimApp(s.chainA).InterchainQueryKeeper.AllQueries(s.chainA.GetContext())
	s.Len(queries, 0)
}
