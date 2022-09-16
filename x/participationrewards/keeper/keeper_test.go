package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/ibctesting"

	"github.com/stretchr/testify/suite"

	qapp "github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
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
	return chain.App
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
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("TEST"), "TEST"))

	err := icskeeper.HandleRegisterZoneProposal(ctx, s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper, proposal)
	s.Require().NoError(err)

	chainBVals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext())

	for _, val := range chainBVals {
		qvr := stakingtypes.QueryValidatorResponse{
			Validator: val,
		}

		addr, err := utils.ValAddressFromBech32(val.OperatorAddress, "")
		s.Require().NoError(err)

		data := stakingtypes.GetValidatorKey(addr)

		query := s.GetQuicksilverApp(s.chainA).InterchainQueryKeeper.NewQuery(
			ctx,
			icstypes.ModuleName,
			s.path.EndpointA.ConnectionID,
			s.chainB.ChainID,
			"store/staking/key",
			data,
			sdk.ZeroInt(),
			"validator",
			0,
		)
		err = icskeeper.ValidatorCallback(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper, ctx, s.GetQuicksilverApp(s.chainB).AppCodec().MustMarshal(&qvr), *query)
		s.Require().NoError(err)
	}

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
