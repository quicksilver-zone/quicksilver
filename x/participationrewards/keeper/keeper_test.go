package keeper_test

import (
	"testing"
	"time"

	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/app"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	path *ibctesting.Path
}

func (s *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
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

	suite.path = newQuicksilverPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)

	suite.coordinator.CurrentTime = time.Now().UTC()
	suite.coordinator.UpdateTime()

	suite.initTestZone()
}

func (suite *KeeperTestSuite) initTestZone() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:    suite.path.EndpointA.ConnectionID,
		ChainId:         suite.chainB.ChainID,
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       true,
		LiquidityModule: true,
	}

	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
}

/*func (s *KeeperTestSuite) SetupRegisteredZones() {
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
}*/

/*func (s *KeeperTestSuite) Test() {
	s.SetupRegisteredZones()
}*/
