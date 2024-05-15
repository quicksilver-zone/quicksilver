package keeper_test

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	minttypes "github.com/quicksilver-zone/quicksilver/x/mint/types"
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

func (suite *KeeperTestSuite) TestDeleteClaimRecord() {
	suite.initTestZoneDrop()

	address := addressutils.GenerateAccAddressForTest().String()

	cr := types.ClaimRecord{
		ChainId:   suite.chainA.ChainID,
		Address:   address,
		BaseValue: 100,
	}
	suite.setClaimRecord(cr)

	// delete claim record
	err := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.DeleteClaimRecord(suite.chainA.GetContext(), cr.ChainId, cr.Address)
	suite.Require().NoError(err)

	// check if claim record is deleted
	_, err = suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.GetClaimRecord(suite.chainA.GetContext(), cr.ChainId, cr.Address)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestIterateClaimRecords() {
	suite.initTestZoneDrop()

	addresses := []string{
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
	}

	suite.setDefaultClaimRecords(addresses)

	// iterate claim records
	var count int
	suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.IterateClaimRecords(suite.chainA.GetContext(), suite.chainA.ChainID, func(_ int64, cr types.ClaimRecord) (stop bool) {
		count++
		return false
	})
	suite.Require().Equal(len(addresses), count)
}

func (suite *KeeperTestSuite) TestAllZoneClaimRecords() {
	suite.initTestZoneDrop()

	addresses := []string{
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
	}

	suite.setDefaultClaimRecords(addresses)

	// get all claim records
	allCRs := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.AllClaimRecords(suite.chainA.GetContext())
	suite.Require().Equal(len(addresses), len(allCRs))
}

func (suite *KeeperTestSuite) TestClearClaimRecords() {
	suite.initTestZoneDrop()

	addresses := []string{
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
	}

	suite.setDefaultClaimRecords(addresses)

	// clear all claim records
	suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ClearClaimRecords(suite.chainA.GetContext(), suite.chainA.ChainID)

	// get all claim records
	allCRs := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.AllClaimRecords(suite.chainA.GetContext())
	suite.Require().Equal(0, len(allCRs))
}

func (suite *KeeperTestSuite) TestZoneDrop() {
	suite.initTestZoneDrop()

	req := types.QueryZoneDropRequest{
		ChainId: suite.chainB.ChainID,
	}

	res, err := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ZoneDrop(suite.chainA.GetContext(), &req)
	suite.Require().NoError(err)

	zoneDrop := res.ZoneDrop
	suite.Require().Equal(suite.chainB.ChainID, zoneDrop.ChainId)
	suite.Require().EqualValues(1000000000, zoneDrop.Allocation)
}

func (suite *KeeperTestSuite) TestZoneDrops() {
	suite.initTestZoneDrop()

	req := types.QueryZoneDropsRequest{Status: types.StatusActive}

	res, err := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ZoneDrops(suite.chainA.GetContext(), &req)
	suite.Require().NoError(err)

	zoneDrops := res.ZoneDrops
	suite.Require().Len(zoneDrops, 1)
	suite.Require().Equal(suite.chainB.ChainID, zoneDrops[0].ChainId)
	suite.Require().EqualValues(1000000000, zoneDrops[0].Allocation)

	// empty status request
	req.Status = types.StatusFuture

	res, err = suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ZoneDrops(suite.chainA.GetContext(), &req)
	suite.Require().NoError(err)
	suite.Require().Len(res.ZoneDrops, 0)
}

func (suite *KeeperTestSuite) TestClaimRecord() {
	suite.initTestZoneDrop()

	address := addressutils.GenerateAccAddressForTest().String()

	cr := types.ClaimRecord{
		ChainId:   suite.chainA.ChainID,
		Address:   address,
		BaseValue: 100,
	}
	suite.setClaimRecord(cr)

	req := types.QueryClaimRecordRequest{
		ChainId: cr.ChainId,
		Address: cr.Address,
	}

	res, err := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ClaimRecord(suite.chainA.GetContext(), &req)
	suite.Require().NoError(err)

	claimRecord := res.ClaimRecord
	suite.Require().Equal(cr.ChainId, claimRecord.ChainId)
	suite.Require().Equal(cr.Address, claimRecord.Address)
	suite.Require().Equal(cr.BaseValue, claimRecord.BaseValue)

	// invalid address
	req.Address = "invalid"
	_, err = suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ClaimRecord(suite.chainA.GetContext(), &req)
	suite.Require().Error(err)

	// invalid chain id
	req.ChainId = "invalid"
	req.Address = cr.Address
	_, err = suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ClaimRecord(suite.chainA.GetContext(), &req)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestClaimRecords() {
	suite.initTestZoneDrop()

	addresses := []string{
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
		addressutils.GenerateAccAddressForTest().String(),
	}

	suite.setDefaultClaimRecords(addresses)

	req := types.QueryClaimRecordsRequest{
		ChainId: suite.chainA.ChainID,
	}

	res, err := suite.GetQuicksilverApp(suite.chainA).AirdropKeeper.ClaimRecords(suite.chainA.GetContext(), &req)
	suite.Require().NoError(err)

	claimRecords := res.ClaimRecords
	suite.Require().Len(claimRecords, len(addresses))
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

func (suite *KeeperTestSuite) setDefaultClaimRecords(addresses []string) {
	for _, address := range addresses {
		cr := types.ClaimRecord{
			ChainId:   suite.chainA.ChainID,
			Address:   address,
			BaseValue: 100,
		}
		suite.setClaimRecord(cr)
	}
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
			sdk.NewIntFromUint64(amount),
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
