package keeper_test

import (
	"fmt"
	"testing"

	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	qapp "github.com/ingenuity-build/quicksilver/app"
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

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

func (s *KeeperTestSuite) TestEpochEnd() {
	fmt.Println("TestEpochEnd >>>")
	fmt.Println("<<< TestEpochEnd")
}
