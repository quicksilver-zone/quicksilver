package keeper_test

import (
	"testing"
	"time"

	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var testAddress = utils.GenerateAccAddressForTest().String()

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

// TestKeeperTestSuite runs all the tests within this package.
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
	quicksilver, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return quicksilver
}

// SetupTest creates a coordinator with 2 test chains.
func (s *KeeperTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)         // initializes 2 test chains
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	s.path = newQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)

	s.coordinator.CurrentTime = time.Now().UTC()
	s.coordinator.UpdateTime()

	s.initTestZone()
}

func (s *KeeperTestSuite) initTestZone() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:     s.path.EndpointA.ConnectionID,
		ChainId:          s.chainB.ChainID,
		AccountPrefix:    "bcosmos",
		LocalDenom:       "uqatom",
		BaseDenom:        "uatom",
		ReturnToSender:   false,
		UnbondingEnabled: false,
		LiquidityModule:  true,
		Decimals:         6,
		Is_118:           true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	// cosmos zone
	zone = icstypes.Zone{
		ConnectionId:     "connection-77001",
		ChainId:          "cosmoshub-4",
		AccountPrefix:    "cosmos",
		LocalDenom:       "uqatom",
		BaseDenom:        "uatom",
		ReturnToSender:   false,
		UnbondingEnabled: false,
		LiquidityModule:  true,
		Decimals:         6,
		Is_118:           true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	// osmosis zone
	zone = icstypes.Zone{
		ConnectionId:     "connection-77002",
		ChainId:          "osmosis-1",
		AccountPrefix:    "osmo",
		LocalDenom:       "uqosmo",
		BaseDenom:        "uosmo",
		ReturnToSender:   false,
		UnbondingEnabled: false,
		LiquidityModule:  true,
		Decimals:         6,
		Is_118:           true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
}
