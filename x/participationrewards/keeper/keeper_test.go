package keeper_test

import (
	"fmt"
	"testing"
	"time"

	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/app"
	icqkeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

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
	app, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return app
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
		ConnectionId:    suite.path.EndpointA.ConnectionID,
		ChainId:         suite.chainB.ChainID,
		AccountPrefix:   "bcosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       true,
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
		MultiSend:       true,
		LiquidityModule: true,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

	// osmosis zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       true,
		LiquidityModule: true,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

	// connection type
	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0),
		suite.chainB.ChainID,
	)
	// osmosis
	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisParams,
		fmt.Sprintf("{\"ChainID\": %q}", "osmosis-1"),
		"",
	)

	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", "connection-77002", "osmosis-1", 0),
		"osmosis-1",
	)

	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisPool,
		fmt.Sprintf(
			"{\"poolid\":%d,\"poolname\":%q,\"zones\":{%q:%q,%q:%q}}",
			1,
			"atom/osmo",
			"cosmoshub-4",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"osmosis-1",
			"ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159",
		),
		"1",
	)

	// ensure that epoch boundaries are reached...
	suite.coordinator.CommitNBlocks(suite.chainA, 1)
	suite.coordinator.CommitNBlocks(suite.chainB, 1)

	suite.executeOsmosisPoolUpdateCallback()
	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.AfterEpochEnd(suite.chainA.GetContext(), "epoch", 3)
}

func (suite *KeeperTestSuite) addProtocolData(Type types.ProtocolDataType, Data string, Key string) {
	pd := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(Type)],
		Data: []byte(Data),
	}

	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.SetProtocolData(suite.chainA.GetContext(), Key, &pd)
}

func (suite *KeeperTestSuite) executeOsmosisPoolUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper

	osm := &keeper.OsmosisModule{}
	qid := icqkeeper.GenerateQueryHash("connection-77002", "osmosis-1", "store/gamm/key", osm.GetKeyPrefixPools(1), types.ModuleName)

	query, found := prk.IcqKeeper.GetQuery(suite.chainA.GetContext(), qid)
	suite.Require().True(found, "qid: %s", qid)

	resp := []byte{10, 26, 47, 111, 115, 109, 111, 115, 105, 115, 46, 103, 97, 109, 109, 46, 118, 49, 98, 101, 116, 97, 49, 46, 80, 111, 111, 108, 18, 202, 2, 10, 63, 111, 115, 109, 111, 49, 109, 119, 48, 97, 99, 54, 114, 119, 108, 112, 53, 114, 56, 119, 97, 112, 119, 107, 51, 122, 115, 54, 103, 50, 57, 104, 56, 102, 99, 115, 99, 120, 113, 97, 107, 100, 122, 119, 57, 101, 109, 107, 110, 101, 54, 99, 56, 119, 106, 112, 57, 113, 48, 116, 51, 118, 56, 116, 16, 1, 26, 6, 10, 1, 48, 18, 1, 48, 34, 4, 49, 54, 56, 104, 42, 43, 10, 11, 103, 97, 109, 109, 47, 112, 111, 111, 108, 47, 49, 18, 28, 49, 48, 48, 48, 48, 48, 48, 50, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 48, 48, 50, 94, 10, 80, 10, 68, 105, 98, 99, 47, 49, 53, 69, 57, 67, 53, 67, 70, 53, 57, 54, 57, 48, 56, 48, 53, 51, 57, 68, 66, 51, 57, 53, 70, 65, 55, 68, 57, 67, 48, 56, 54, 56, 50, 54, 53, 50, 49, 55, 69, 70, 67, 53, 50, 56, 52, 51, 51, 54, 55, 49, 65, 65, 70, 57, 66, 49, 57, 49, 50, 68, 49, 53, 57, 18, 8, 49, 48, 48, 48, 48, 48, 48, 51, 18, 10, 49, 48, 55, 51, 55, 52, 49, 56, 50, 52, 50, 94, 10, 80, 10, 68, 105, 98, 99, 47, 51, 48, 50, 48, 57, 50, 50, 66, 55, 53, 55, 54, 70, 67, 55, 53, 66, 66, 69, 48, 53, 55, 65, 48, 50, 57, 48, 65, 57, 65, 69, 69, 70, 70, 52, 56, 57, 66, 66, 49, 49, 49, 51, 69, 54, 69, 51, 54, 53, 67, 69, 52, 55, 50, 68, 52, 66, 70, 66, 55, 70, 70, 65, 51, 18, 8, 49, 48, 48, 48, 48, 48, 48, 51, 18, 10, 49, 48, 55, 51, 55, 52, 49, 56, 50, 52, 58, 10, 50, 49, 52, 55, 52, 56, 51, 54, 52, 56}

	err := keeper.OsmosisPoolUpdateCallback(
		prk,
		suite.chainA.GetContext(),
		resp,
		query,
	)
	suite.Require().NoError(err)
}
