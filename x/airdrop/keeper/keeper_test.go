package keeper_test

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/stretchr/testify/suite"

	qapp "github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
)

func init() {
	ibctesting.DefaultTestingAppInit = qapp.SetupTestingApp
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

func (s *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *qapp.Quicksilver {
	app, ok := chain.App.(*qapp.Quicksilver)
	if !ok {
		panic("not Quicksilver app")
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

	suite.coordinator.CommitNBlocks(suite.chainA, 10)
	suite.coordinator.CommitNBlocks(suite.chainB, 10)
}

func (suite *KeeperTestSuite) initTestZone() {
	// osmosis zone
	zone := icstypes.Zone{
		ConnectionId:  suite.path.EndpointB.ConnectionID,
		ChainId:       suite.chainB.ChainID,
		AccountPrefix: "osmo",
		LocalDenom:    "uosmo",
		BaseDenom:     "stake",
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
		Actions: []sdk.Dec{
			0:  sdk.MustNewDecFromStr("0.15"), // 15%
			1:  sdk.MustNewDecFromStr("0.06"), // 21%
			2:  sdk.MustNewDecFromStr("0.07"), // 28%
			3:  sdk.MustNewDecFromStr("0.08"), // 36%
			4:  sdk.MustNewDecFromStr("0.09"), // 45%
			5:  sdk.MustNewDecFromStr("0.1"),  // 55%
			6:  sdk.MustNewDecFromStr("0.15"), // 70%
			7:  sdk.MustNewDecFromStr("0.05"), // 75%
			8:  sdk.MustNewDecFromStr("0.1"),  // 85%
			9:  sdk.MustNewDecFromStr("0.1"),  // 95%
			10: sdk.MustNewDecFromStr("0.05"), // 100%
		},
		IsConcluded: false,
	}

	return zd
}

func (suite *KeeperTestSuite) compressClaimRecords(crs []types.ClaimRecord) []byte {
	bz, err := json.Marshal(&crs)
	suite.Require().NoError(err)

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	zw.Write(bz)

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
	app := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	coins := sdk.NewCoins(
		sdk.NewCoin(
			app.StakingKeeper.BondDenom(ctx),
			sdk.NewIntFromUint64(amount),
		),
	)
	// fund zonedrop account
	zdacc := app.AirdropKeeper.GetZoneDropAccountAddress(chainID)

	err := app.MintKeeper.MintCoins(ctx, coins)
	suite.Require().NoError(err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, zdacc, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) setClaimRecord(cr types.ClaimRecord) {
	err := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.SetClaimRecord(suite.chainA.GetContext(), cr)
	if err != nil {
		suite.T().Logf("setClaimRecord error: %v", err)
	}
	suite.Require().NoError(err)
}
