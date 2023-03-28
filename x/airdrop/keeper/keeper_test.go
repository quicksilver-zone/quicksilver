package keeper_test

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
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
func (s *KeeperTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)         // initializes 2 test chains
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	s.path = newQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)

	s.coordinator.CurrentTime = time.Now().UTC()
	s.coordinator.UpdateTime()

	s.initTestZone()

	s.coordinator.CommitNBlocks(s.chainA, 10)
	s.coordinator.CommitNBlocks(s.chainB, 10)
}

func (s *KeeperTestSuite) initTestZone() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:  s.path.EndpointB.ConnectionID,
		ChainId:       s.chainB.ChainID,
		AccountPrefix: "cosmos",
		LocalDenom:    "uqatom",
		BaseDenom:     "uatom",
	}

	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
}

func (s *KeeperTestSuite) getZoneDrop() types.ZoneDrop {
	zd := types.ZoneDrop{
		ChainId:    s.chainB.ChainID,
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

func (s *KeeperTestSuite) compressClaimRecords(crs []types.ClaimRecord) []byte {
	bz, err := json.Marshal(&crs)
	s.Require().NoError(err)

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	zw.Write(bz)

	err = zw.Close()
	s.Require().NoError(err)

	return buf.Bytes()
}

func (s *KeeperTestSuite) initTestZoneDrop() {
	zd := s.getZoneDrop()
	s.GetQuicksilverApp(s.chainA).AirdropKeeper.SetZoneDrop(s.chainA.GetContext(), zd)
	s.fundZoneDrop(zd.ChainId, zd.Allocation)
}

func (s *KeeperTestSuite) fundZoneDrop(chainID string, amount uint64) {
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()
	coins := sdk.NewCoins(
		sdk.NewCoin(
			app.StakingKeeper.BondDenom(ctx),
			sdk.NewIntFromUint64(amount),
		),
	)
	// fund zonedrop account
	zdacc := app.AirdropKeeper.GetZoneDropAccountAddress(chainID)

	err := app.MintKeeper.MintCoins(ctx, coins)
	s.Require().NoError(err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, zdacc, coins)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) setClaimRecord(cr types.ClaimRecord) {
	err := s.GetQuicksilverApp(s.chainA).AirdropKeeper.SetClaimRecord(s.chainA.GetContext(), cr)
	if err != nil {
		s.T().Logf("setClaimRecord error: %v", err)
	}
	s.Require().NoError(err)
}
