package keeper_test

import (
	"context"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	"github.com/stretchr/testify/suite"

<<<<<<< HEAD
	"github.com/ingenuity-build/quicksilver/app"
	qapp "github.com/ingenuity-build/quicksilver/app"
=======
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/ingenuity-build/quicksilver/app"
>>>>>>> 26bb442 (simplify quicksilver logic)
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var TestOwnerAddress = "quick17dtl0mjt3t77kpuhg2edqzjpszulwhgzhk4dtz"

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

func (s *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
	app, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
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

func (s *KeeperTestSuite) SetupZones() {
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

	// Simulate "cosmos.staking.v1beta1.Query/Validators" response

	qApp := s.GetQuicksilverApp(s.chainA)

	for _, val := range s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext()) {
		// refetch the zone for each validator, else we end up with an empty valset each time!
		zone, found := qApp.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
		s.Require().True(found)
		s.Require().NoError(icskeeper.SetValidatorForZone(&qApp.InterchainstakingKeeper, s.chainA.GetContext(), zone, app.DefaultConfig().Codec.MustMarshal(&val)))
	}

	// valsetInterval := uint64(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetParam(ctx, icstypes.KeyValidatorSetInterval))
	s.coordinator.CommitNBlocks(s.chainA, 2)
	s.coordinator.CommitNBlocks(s.chainB, 2)
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

func GetICSKeeper(t *testing.T) (*icskeeper.Keeper, sdk.Context) {
	app := app.Setup(t, false)
	keeper := app.InterchainstakingKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "mercury-1", Time: time.Now().UTC()})

	return &keeper, ctx
}
