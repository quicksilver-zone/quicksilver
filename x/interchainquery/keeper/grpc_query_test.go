package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

func (suite *KeeperTestSuite) TestQueries() {
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

	icqsrvSrv := icqtypes.QuerySrvrServer(suite.GetSimApp(suite.chainA).InterchainQueryKeeper)

	res, err := icqsrvSrv.Queries(sdk.WrapSDKContext(suite.chainA.GetContext()), &icqtypes.QueryRequestsRequest{ChainId: suite.chainB.ChainID})
	suite.NoError(err)
	suite.Len(res.Queries, 1)
	suite.Equal(suite.path.EndpointB.ConnectionID, res.Queries[0].ConnectionId)
	suite.Equal(suite.chainB.ChainID, res.Queries[0].ChainId)
	suite.Equal("cosmos.staking.v1beta1.Query/Validators", res.Queries[0].QueryType)
	suite.Equal(sdk.NewInt(200), res.Queries[0].Period)
	suite.Equal("", res.Queries[0].CallbackId)
}
