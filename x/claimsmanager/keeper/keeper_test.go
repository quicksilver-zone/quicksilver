package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

var testAddress = addressutils.GenerateAccAddressForTest().String()

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

func (suite *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
	quicksilver, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return quicksilver
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)         // initializes 2 test chains
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	suite.path = newQuicksilverPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)

	suite.coordinator.CurrentTime = time.Now().UTC()
	suite.coordinator.UpdateTime()

	suite.initTestZone()
}

func (suite *KeeperTestSuite) initTestZone() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:     suite.path.EndpointA.ConnectionID,
		ChainId:          suite.chainB.ChainID,
		AccountPrefix:    "bcosmos",
		LocalDenom:       "uqatom",
		BaseDenom:        "uatom",
		ReturnToSender:   false,
		UnbondingEnabled: false,
		LiquidityModule:  true,
		Decimals:         6,
		Is_118:           true,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

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
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

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
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
}
