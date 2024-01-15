package keeper_test

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctesting "github.com/cosmos/ibc-go/v8/testing"

	"github.com/quicksilver-zone/quicksilver/v7/app"
	"github.com/quicksilver-zone/quicksilver/v7/x/airdrop/types"
	icstypes "github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
	minttypes "github.com/quicksilver-zone/quicksilver/v7/x/mint/types"
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

func (*KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
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

	suite.coordinator.CommitNBlocks(suite.chainA, 10)
	suite.coordinator.CommitNBlocks(suite.chainB, 10)
}

func (suite *KeeperTestSuite) initTestZone() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:  suite.path.EndpointB.ConnectionID,
		ChainId:       suite.chainB.ChainID,
		AccountPrefix: "cosmos",
		LocalDenom:    "uqatom",
		BaseDenom:     "uatom",
		Is_118:        true,
	}

	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
}

func (suite *KeeperTestSuite) getZoneDrop() types.ZoneDrop {
	zd := types.ZoneDrop{
		ChainId:    suite.chainB.ChainID,
		StartTime:  time.Now().Add(-5 * time.Minute),
		Duration:   time.Hour,
		Decay:      30 * time.Minute,
		Allocation: 1000000000,
		Actions: []sdkmath.LegacyDec{
			0:  sdkmath.LegacyMustNewDecFromStr("0.15"), // 15%
			1:  sdkmath.LegacyMustNewDecFromStr("0.06"), // 21%
			2:  sdkmath.LegacyMustNewDecFromStr("0.07"), // 28%
			3:  sdkmath.LegacyMustNewDecFromStr("0.08"), // 36%
			4:  sdkmath.LegacyMustNewDecFromStr("0.09"), // 45%
			5:  sdkmath.LegacyMustNewDecFromStr("0.1"),  // 55%
			6:  sdkmath.LegacyMustNewDecFromStr("0.15"), // 70%
			7:  sdkmath.LegacyMustNewDecFromStr("0.05"), // 75%
			8:  sdkmath.LegacyMustNewDecFromStr("0.1"),  // 85%
			9:  sdkmath.LegacyMustNewDecFromStr("0.1"),  // 95%
			10: sdkmath.LegacyMustNewDecFromStr("0.05"), // 100%
		},
		IsConcluded: false,
	}

	return zd
}

func (suite *KeeperTestSuite) compressClaimRecords(crs []types.ClaimRecord) []byte {
	suite.T().Helper()

	bz, err := json.Marshal(&crs)
	suite.Require().NoError(err)

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	_, err = zw.Write(bz)
	suite.Require().NoError(err)

	err = zw.Close()
	suite.Require().NoError(err)

	return buf.Bytes()
}

func (suite *KeeperTestSuite) initTestZoneDrop() {
	zd := suite.getZoneDrop()
	suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.SetZoneDrop(suite.chainA.GetContext(), zd)
	suite.fundZoneDrop(zd.ChainId, zd.Allocation)
}

func (suite *KeeperTestSuite) fundZoneDrop(chainID string, amount uint64) {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	coins := sdk.NewCoins(
		sdk.NewCoin(
			quicksilver.StakingKeeper.BondDenom(ctx),
			sdkmath.NewIntFromUint64(amount),
		),
	)
	// fund zonedrop account
	zdacc := quicksilver.AirdropKeeper.GetZoneDropAccountAddress(chainID)

	err := quicksilver.MintKeeper.MintCoins(ctx, coins)
	suite.Require().NoError(err)

	err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, zdacc, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) setClaimRecord(cr types.ClaimRecord) {
	err := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.SetClaimRecord(suite.chainA.GetContext(), cr)
	if err != nil {
		suite.T().Logf("setClaimRecord error: %v", err)
	}
	suite.Require().NoError(err)
}
