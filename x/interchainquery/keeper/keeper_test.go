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

func (s *KeeperTestSuite) GetSimApp(chain *ibctesting.TestChain) *app.Quicksilver {
	quicksilver, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return quicksilver
}

func (s *KeeperTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1))
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2))

	s.path = newSimAppPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)
}

func (s *KeeperTestSuite) TestMakeRequest() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	s.NoError(err)

	s.GetSimApp(s.chainA).InterchainQueryKeeper.MakeRequest(
		s.chainA.GetContext(),
		s.path.EndpointB.ConnectionID,
		s.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		"",
		0,
	)

	id := keeper.GenerateQueryHash(s.path.EndpointB.ConnectionID, s.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")
	query, found := s.GetSimApp(s.chainA).InterchainQueryKeeper.GetQuery(s.chainA.GetContext(), id)
	s.True(found)
	s.Equal(s.path.EndpointB.ConnectionID, query.ConnectionId)
	s.Equal(s.chainB.ChainID, query.ChainId)
	s.Equal("cosmos.staking.v1beta1.Query/Validators", query.QueryType)
	s.Equal(sdk.NewInt(200), query.Period)
	s.Equal("", query.CallbackId)

	s.GetSimApp(s.chainA).InterchainQueryKeeper.MakeRequest(
		s.chainA.GetContext(),
		s.path.EndpointB.ConnectionID,
		s.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		"",
		0,
	)
}

func (s *KeeperTestSuite) TestSubmitQueryResponse() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	s.NoError(err)

	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: s.GetSimApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext()),
	}

	tests := []struct {
		query       *icqtypes.Query
		setQuery    bool
		expectError error
	}{
		{
			s.GetSimApp(s.chainA).InterchainQueryKeeper.
				NewQuery(
					"",
					s.path.EndpointB.ConnectionID,
					s.chainB.ChainID,
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
			s.GetSimApp(s.chainA).InterchainQueryKeeper.
				NewQuery(
					"",
					s.path.EndpointB.ConnectionID,
					s.chainB.ChainID,
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
			s.GetSimApp(s.chainA).InterchainQueryKeeper.
				NewQuery(
					"",
					s.path.EndpointB.ConnectionID,
					s.chainB.ChainID,
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
			s.GetSimApp(s.chainA).InterchainQueryKeeper.
				NewQuery(
					"",
					s.path.EndpointB.ConnectionID,
					s.chainB.ChainID,
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
			s.GetSimApp(s.chainA).InterchainQueryKeeper.SetQuery(s.chainA.GetContext(), *tc.query)
		}

		icqmsgSrv := keeper.NewMsgServerImpl(s.GetSimApp(s.chainA).InterchainQueryKeeper)

		qmsg := icqtypes.MsgSubmitQueryResponse{
			ChainId:     s.chainB.ChainID,
			QueryId:     keeper.GenerateQueryHash(tc.query.ConnectionId, tc.query.ChainId, tc.query.QueryType, bz, ""),
			Result:      s.GetSimApp(s.chainB).AppCodec().MustMarshalJSON(&qvr),
			Height:      s.chainB.CurrentHeader.Height,
			FromAddress: TestOwnerAddress,
		}

		_, err = icqmsgSrv.SubmitQueryResponse(sdk.WrapSDKContext(s.chainA.GetContext()), &qmsg)
		s.Equal(tc.expectError, err)
	}
}

func (s *KeeperTestSuite) TestDataPoints() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	s.NoError(err)

	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: s.GetSimApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext()),
	}

	id := keeper.GenerateQueryHash(s.path.EndpointB.ConnectionID, s.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")

	err = s.GetSimApp(s.chainA).InterchainQueryKeeper.SetDatapointForID(
		s.chainA.GetContext(),
		id,
		s.GetSimApp(s.chainB).AppCodec().MustMarshalJSON(&qvr),
		sdk.NewInt(s.chainB.CurrentHeader.Height),
	)
	s.NoError(err)

	dataPoint, err := s.GetSimApp(s.chainA).InterchainQueryKeeper.GetDatapointForID(s.chainA.GetContext(), id)
	s.NoError(err)
	s.NotNil(dataPoint)

	s.GetSimApp(s.chainA).InterchainQueryKeeper.DeleteDatapoint(s.chainA.GetContext(), id)
}

func newSimAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}
