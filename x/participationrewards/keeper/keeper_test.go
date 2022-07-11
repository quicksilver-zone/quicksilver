package keeper_test

import (
	"context"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctesting "github.com/cosmos/ibc-go/v4/testing"
	qapp "github.com/ingenuity-build/quicksilver/app"
	icqkeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/suite"
)

var TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

func init() {
	ibctesting.DefaultTestingAppInit = qapp.SetupTestingApp
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

func (s *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *qapp.Quicksilver {
	app, ok := chain.App.(*qapp.Quicksilver)
	if !ok {
		panic("not Quicksilver app")
	}

	return app
}

func (s *KeeperTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1))
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2))

	s.path = newQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)
}

func (s *KeeperTestSuite) SetupRegisteredZones() {
	proposal := &icstypes.RegisterZoneProposal{
		Title:           "register zone A",
		Description:     "register zone A",
		ConnectionId:    s.path.EndpointA.ConnectionID,
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		AccountPrefix:   "cosmos",
		MultiSend:       true,
		LiquidityModule: true,
	}

	ctx := s.chainA.GetContext()

	// Set special testing context (e.g. for test / debug output)
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), "TEST", "TEST"))

	err := icskeeper.HandleRegisterZoneProposal(ctx, s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper, proposal)
	s.Require().NoError(err)

	// Simulate "cosmos.staking.v1beta1.Query/Validators" response
	// - this is not working anymore;
	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext()),
	}
	icqmsgSrv := icqkeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainQueryKeeper)

	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := s.GetQuicksilverApp(s.chainA).AppCodec().Marshal(&bondedQuery)
	s.Require().NoError(err)

	qmsg := icqtypes.MsgSubmitQueryResponse{
		ChainId: s.chainB.ChainID,
		QueryId: icqkeeper.GenerateQueryHash(
			s.path.EndpointA.ConnectionID,
			s.chainB.ChainID,
			"cosmos.staking.v1beta1.Query/Validators",
			bz,
			icstypes.ModuleName,
		),
		Result:      s.GetQuicksilverApp(s.chainB).AppCodec().MustMarshalJSON(&qvr),
		Height:      s.chainB.CurrentHeader.Height,
		FromAddress: TestOwnerAddress,
	}
	_, err = icqmsgSrv.SubmitQueryResponse(sdktypes.WrapSDKContext(ctx), &qmsg)
	s.Require().NoError(err)

	valsetInterval := uint64(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetParam(ctx, icstypes.KeyValidatorSetInterval))
	s.coordinator.CommitNBlocks(s.chainA, valsetInterval)
	s.coordinator.CommitNBlocks(s.chainB, valsetInterval)
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

func (s *KeeperTestSuite) Test() {
	s.SetupRegisteredZones()
}
