package app

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func init() {
	ibctesting.DefaultTestingAppInit = SetupTestingApp
}

// TestKeeperTestSuite runs all the tests within this package.
func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

type AppTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	path *ibctesting.Path
}

func (s *AppTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *Quicksilver {
	app, ok := chain.App.(*Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return app
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *AppTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)         // initializes 2 test chains
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	suite.path = newQuicksilverPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)

	suite.coordinator.CurrentTime = time.Now().UTC()
	suite.coordinator.UpdateTime()

	suite.initTestZone()
}

func (suite *AppTestSuite) initTestZone() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:    suite.path.EndpointA.ConnectionID,
		ChainId:         suite.chainB.ChainID,
		AccountPrefix:   "bcosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: true,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

	// cosmos zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: false,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

	// osmosis zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       false,
		LiquidityModule: true,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
}

func (s *AppTestSuite) TestV010400UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)
	handler := v010400UpgradeHandler(app)
	ctx := s.chainA.GetContext()
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	osmosis, found := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
	s.Require().True(found)
	s.Require().Equal(int64(6), osmosis.Decimals)
	s.Require().Equal("osmo", osmosis.AccountPrefix)
	s.Require().Equal("connection-77002", osmosis.ConnectionId)
	s.Require().False(osmosis.UnbondingEnabled)
	s.Require().False(osmosis.ReturnToSender)
	s.Require().True(osmosis.LiquidityModule)

	cosmos, found := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.Require().True(found)
	s.Require().Equal(int64(6), cosmos.Decimals)
	s.Require().Equal("uatom", cosmos.BaseDenom)
	s.Require().Equal("uqatom", cosmos.LocalDenom)
	s.Require().False(cosmos.UnbondingEnabled)
	s.Require().False(cosmos.ReturnToSender)
	s.Require().False(cosmos.LiquidityModule)

	chainb, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(int64(6), chainb.Decimals)
	s.Require().False(chainb.UnbondingEnabled)
	s.Require().False(chainb.ReturnToSender)
	s.Require().True(chainb.LiquidityModule)
}
