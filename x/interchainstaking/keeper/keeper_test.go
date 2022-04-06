package keeper_test

import (
	"context"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	qapp "github.com/ingenuity-build/quicksilver/app"
	icqkeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/suite"
)

var (
	TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"
)

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
}

func (s *KeeperTestSuite) SetupRegisteredZones() {
	path := NewQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(path)

	zonemsg := icstypes.MsgRegisterZone{
		Identifier:   "cosmos",
		ConnectionId: path.EndpointA.ConnectionID,
		LocalDenom:   "uqatom",
		BaseDenom:    "uatom",
		FromAddress:  TestOwnerAddress,
	}

	msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
	ctx := s.chainA.GetContext()
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), "TEST", "TEST"))
	_, err := msgSrv.RegisterZone(sdktypes.WrapSDKContext(ctx), &zonemsg)
	s.Require().NoError(err)

	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext()),
	}

	icqmsgSrv := icqkeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainQueryKeeper)
	qmsg := icqtypes.MsgSubmitQueryResponse{
		// target or source chain_id?
		ChainId: s.chainB.ChainID,
		QueryId: icqkeeper.GenerateQueryHash(
			path.EndpointA.ConnectionID,
			s.chainB.ChainID,
			"cosmos.staking.v1beta1.Query/Validators",
			map[string]string{"status": stakingtypes.BondStatusBonded},
		),
		Result:      s.GetQuicksilverApp(s.chainB).AppCodec().MustMarshalJSON(&qvr),
		Height:      s.chainB.CurrentHeader.Height,
		FromAddress: TestOwnerAddress,
	}
	_, err = icqmsgSrv.SubmitQueryResponse(sdktypes.WrapSDKContext(ctx), &qmsg)
	s.Require().NoError(err)

	s.coordinator.CommitNBlocks(s.chainA, 25)
	s.coordinator.CommitNBlocks(s.chainB, 25)
}

func NewQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}
