package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

const TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	path   *ibctesting.Path
}

func (suite *KeeperTestSuite) GetSimApp(chain *ibctesting.TestChain) *app.Quicksilver {
	app, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return app
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = newSimAppPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)
}

func (suite *KeeperTestSuite) TestMakeRequest() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	suite.GetSimApp(suite.chainA).InterchainQueryKeeper.MakeRequest(
		suite.chainA.GetContext(),
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		"",
		0,
	)

	id := keeper.GenerateQueryHash(suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")
	query, found := suite.GetSimApp(suite.chainA).InterchainQueryKeeper.GetQuery(suite.chainA.GetContext(), id)
	suite.True(found)
	suite.Equal(suite.path.EndpointB.ConnectionID, query.ConnectionId)
	suite.Equal(suite.chainB.ChainID, query.ChainId)
	suite.Equal("cosmos.staking.v1beta1.Query/Validators", query.QueryType)
	suite.Equal(sdk.NewInt(200), query.Period)
	suite.Equal("", query.CallbackId)

	suite.GetSimApp(suite.chainA).InterchainQueryKeeper.MakeRequest(
		suite.chainA.GetContext(),
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		"",
		0,
	)
}

func (suite *KeeperTestSuite) TestSubmitQueryResponse() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: suite.GetSimApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext()),
	}

	tests := []struct {
		query       *icqtypes.Query
		setQuery    bool
		expectError error
	}{
		{
			suite.GetSimApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(200),
					"",
					0,
				),
			true,
			nil,
		},
		{
			suite.GetSimApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(200),
					"",
					10,
				),
			true,
			nil,
		},
		{
			suite.GetSimApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(-200),
					"",
					0,
				),
			true,
			nil,
		},
		{
			suite.GetSimApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(100),
					"",
					0,
				),
			false,
			nil,
		},
	}

	for _, tc := range tests {
		// set the query
		if tc.setQuery {
			suite.GetSimApp(suite.chainA).InterchainQueryKeeper.SetQuery(suite.chainA.GetContext(), *tc.query)
		}

		icqmsgSrv := keeper.NewMsgServerImpl(suite.GetSimApp(suite.chainA).InterchainQueryKeeper)

		qmsg := icqtypes.MsgSubmitQueryResponse{
			ChainId:     suite.chainB.ChainID,
			QueryId:     keeper.GenerateQueryHash(tc.query.ConnectionId, tc.query.ChainId, tc.query.QueryType, bz, ""),
			Result:      suite.GetSimApp(suite.chainB).AppCodec().MustMarshalJSON(&qvr),
			Height:      suite.chainB.CurrentHeader.Height,
			FromAddress: TestOwnerAddress,
		}

		_, err = icqmsgSrv.SubmitQueryResponse(sdk.WrapSDKContext(suite.chainA.GetContext()), &qmsg)
		suite.Equal(tc.expectError, err)
	}
}

func (suite *KeeperTestSuite) TestDataPoints() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: suite.GetSimApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext()),
	}

	id := keeper.GenerateQueryHash(suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")

	err = suite.GetSimApp(suite.chainA).InterchainQueryKeeper.SetDatapointForID(
		suite.chainA.GetContext(),
		id,
		suite.GetSimApp(suite.chainB).AppCodec().MustMarshalJSON(&qvr),
		sdk.NewInt(suite.chainB.CurrentHeader.Height),
	)
	suite.NoError(err)

	dataPoint, err := suite.GetSimApp(suite.chainA).InterchainQueryKeeper.GetDatapointForID(suite.chainA.GetContext(), id)
	suite.NoError(err)
	suite.NotNil(dataPoint)

	suite.GetSimApp(suite.chainA).InterchainQueryKeeper.DeleteDatapoint(suite.chainA.GetContext(), id)
}

func newSimAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}
