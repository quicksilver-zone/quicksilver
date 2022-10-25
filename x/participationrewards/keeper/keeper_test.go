package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
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

	suite.coreTest()
}

func (suite *KeeperTestSuite) coreTest() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:    suite.path.EndpointA.ConnectionID,
		ChainId:         suite.chainB.ChainID,
		AccountPrefix:   "bcosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       true,
		LiquidityModule: true,
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           utils.GenerateAccAddressForTestWithPrefix("bcosmos"),
			PortName:          fmt.Sprintf("%s.performance", suite.chainB.ChainID),
			WithdrawalAddress: utils.GenerateAccAddressForTestWithPrefix("bcosmos"),
		},
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

	// cosmos zone
	zoneCosmos := icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       true,
		LiquidityModule: true,
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           utils.GenerateAccAddressForTestWithPrefix("cosmos"),
			PortName:          "cosmoshub-4.performance",
			WithdrawalAddress: utils.GenerateAccAddressForTestWithPrefix("cosmos"),
		},
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneCosmos)

	// osmosis zone
	zoneOsmosis := icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       true,
		LiquidityModule: true,
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           utils.GenerateAccAddressForTestWithPrefix("osmo"),
			PortName:          "cosmoshub-4.performance",
			WithdrawalAddress: utils.GenerateAccAddressForTestWithPrefix("osmo"),
		},
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneOsmosis)

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
			"{\"poolid\":%d,\"poolname\":%q,\"pooltype\":\"balancer\",\"zones\":{%q:%q,%q:%q}}",
			1,
			"atom/osmo",
			"cosmoshub-4",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"osmosis-1",
			"ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159",
		),
		"1",
	)

	// advance the chains
	suite.coordinator.CommitNBlocks(suite.chainA, 1)
	suite.coordinator.CommitNBlocks(suite.chainB, 1)

	// callback test
	suite.executeOsmosisPoolUpdateCallback()

	// add some deposits
	suite.addReceipt(
		&zoneCosmos,
		testAddress,
		"testTxHash",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(120000000))),
	)

	// add some claims
	suite.addClaim(
		testAddress,
		"cosmoshub-4",
		cmtypes.ClaimTypeLiquidToken,
		"osmosis-1",
		40000000,
	)

	// Epoch boundary
	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.AfterEpochEnd(suite.chainA.GetContext(), "epoch", 3)
}

func (suite *KeeperTestSuite) addProtocolData(Type types.ProtocolDataType, Data string, Key string) {
	pd := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(Type)],
		Data: []byte(Data),
	}

	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.SetProtocolData(suite.chainA.GetContext(), Key, &pd)
}

func (suite *KeeperTestSuite) addReceipt(zone *icstypes.Zone, sender string, hash string, coins sdk.Coins) {
	receipt := icstypes.Receipt{
		ChainId: zone.ChainId,
		Sender:  sender,
		Txhash:  hash,
		Amount:  coins,
	}

	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetReceipt(suite.chainA.GetContext(), receipt)

	delegationAddress := utils.GenerateAccAddressForTestWithPrefix("cosmos")
	validatorAddress := utils.GenerateValAddressForTestWithPrefix("cosmos")
	delegation := icstypes.Delegation{
		DelegationAddress: delegationAddress,
		ValidatorAddress:  validatorAddress,
		Amount:            coins[0],
		Height:            1,
		RedelegationEnd:   101,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetDelegation(suite.chainA.GetContext(), zone, delegation)
}

func (suite *KeeperTestSuite) addClaim(address string, chainID string, claimType cmtypes.ClaimType, sourceChainID string, amount uint64) {
	claim := cmtypes.Claim{
		UserAddress:   address,
		ChainId:       chainID,
		Module:        claimType,
		SourceChainId: sourceChainID,
		Amount:        amount,
	}
	suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper.SetClaim(suite.chainA.GetContext(), &claim)
}
