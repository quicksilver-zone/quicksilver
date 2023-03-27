package app

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

// TODO: this test runs in isolation, but fails as part of `make test`.
// In the `make test` context, MintCoins() seems to have no effect. Why is this?
// func TestReplaceZone(t *testing.T) {
// 	// set up zone drop record and claims.
// 	app := Setup(false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
// 	var err error
// 	denom := app.StakingKeeper.BondDenom(ctx)
// 	someCoins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000000)))
// 	// work around airdrop keeper can't mint :)

// 	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, someCoins)
// 	require.NoError(t, err)
// 	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, app.AirdropKeeper.GetZoneDropAccountAddress("osmotest-4"), someCoins)
// 	require.NoError(t, err)

// 	zd := airdroptypes.ZoneDrop{
// 		ChainId:    "osmotest-4",
// 		StartTime:  time.Now().AddDate(0, 0, -1),
// 		Duration:   time.Hour,
// 		Decay:      time.Hour,
// 		Allocation: someCoins.AmountOf(denom).Uint64(),
// 		Actions: []sdk.Dec{
// 			sdk.OneDec(),
// 		},
// 		IsConcluded: false,
// 	}

// 	app.AirdropKeeper.SetZoneDrop(ctx, zd)

// 	claim1 := airdroptypes.ClaimRecord{
// 		ChainId:          "osmotest-4",
// 		Address:          "quick1g035r8sl346ttxuj0555yxdwftr52t849t3q39",
// 		ActionsCompleted: make(map[int32]*airdroptypes.CompletedAction),
// 		MaxAllocation:    500000,
// 		BaseValue:        500000,
// 	}

// 	claim2 := airdroptypes.ClaimRecord{
// 		ChainId:          "osmotest-4",
// 		Address:          "quick1u53f8u6jjdpxquesk8tqxzv9hvqx7qyfzlkdrj",
// 		ActionsCompleted: make(map[int32]*airdroptypes.CompletedAction),
// 		MaxAllocation:    500000,
// 		BaseValue:        500000,
// 	}

// 	err = app.AirdropKeeper.SetClaimRecord(ctx, claim1)
// 	require.NoError(t, err)
// 	err = app.AirdropKeeper.SetClaimRecord(ctx, claim2)
// 	require.NoError(t, err)
// 	claims := app.AirdropKeeper.AllZoneClaimRecords(ctx, "osmotest-4")
// 	require.Equal(t, 2, len(claims))
// 	require.True(t, app.AirdropKeeper.GetZoneDropAccountBalance(ctx, "osmotest-4").Amount.Equal(sdk.NewInt(1000000)))
// 	require.NotPanics(t, func() { ReplaceZoneDropChain(ctx, app, "osmotest-4", "osmo-test-4", ctx.BlockHeader().Time) })
// 	claimsAfter := app.AirdropKeeper.AllZoneClaimRecords(ctx, "osmotest-4")
// 	require.Equal(t, 0, len(claimsAfter))
// 	claimsNew := app.AirdropKeeper.AllZoneClaimRecords(ctx, "osmo-test-4")
// 	require.Equal(t, 2, len(claimsNew))
// 	zoneDropsAfter := app.AirdropKeeper.AllZoneDrops(ctx)
// 	// check we don't suddenly have two airdrops.
// 	require.Equal(t, 1, len(zoneDropsAfter))
// 	// check the one aidrop we have has the expected values.
// 	require.Equal(t, zoneDropsAfter[0].ChainId, "osmo-test-4")
// 	require.Equal(t, zoneDropsAfter[0].StartTime, ctx.BlockHeader().Time)
// 	require.False(t, app.AirdropKeeper.GetZoneDropAccountBalance(ctx, "osmotest-4").Amount.Equal(sdk.NewInt(1000000)))
// 	require.True(t, app.AirdropKeeper.GetZoneDropAccountBalance(ctx, "osmo-test-4").Amount.Equal(sdk.NewInt(1000000)))
// }

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
	// uni-5 zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77003",
		ChainId:         "uni-5",
		AccountPrefix:   "juno",
		LocalDenom:      "uqjunox",
		BaseDenom:       "ujunox",
		MultiSend:       false,
		LiquidityModule: true,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

	receipt := icstypes.Receipt{
		ChainId: "uni-5",
		Sender:  utils.GenerateAccAddressForTest().String(),
		Txhash:  "TestDeposit01",
		Amount: sdk.NewCoins(
			sdk.NewCoin(
				"ujunox",
				sdk.NewIntFromUint64(2000000), // 20% deposit
			),
		),
	}

	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetReceipt(suite.chainA.GetContext(), receipt)

	ubRecord := icstypes.UnbondingRecord{
		ChainId:       "uni-5",
		EpochNumber:   1,
		Validator:     "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn",
		RelatedTxhash: []string{"ABC012"},
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetUnbondingRecord(suite.chainA.GetContext(), ubRecord)

	rdRecord := icstypes.RedelegationRecord{
		ChainId:        "uni-5",
		EpochNumber:    1,
		Source:         "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn",
		Destination:    "junovaloper1z89utvygweg5l56fsk8ak7t6hh88fd0aa9ywed",
		Amount:         3000000,
		CompletionTime: time.Time(suite.chainA.GetContext().BlockTime().Add(time.Hour)),
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetRedelegationRecord(suite.chainA.GetContext(), rdRecord)

	rdRecord = icstypes.RedelegationRecord{
		ChainId:        "osmosis-1",
		EpochNumber:    1,
		Source:         "osmovaloper1zxavllftfx3a3y5ldfyze7jnu5uyuktsfx2jcc",
		Destination:    "osmovaloper13eq5c99ym05jn02e78l8cac2fagzgdhh4294zk",
		Amount:         3000000,
		CompletionTime: time.Time{},
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetRedelegationRecord(suite.chainA.GetContext(), rdRecord)

	delRecord := icstypes.Delegation{
		Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000)),
		DelegationAddress: "juno1z89utvygweg5l56fsk8ak7t6hh88fd0azcjpz5",
		Height:            10,
		ValidatorAddress:  "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn",
		RedelegationEnd:   -62135596800,
	}

	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetDelegation(suite.chainA.GetContext(), &zone, delRecord)

	wRecord := icstypes.WithdrawalRecord{
		ChainId:   "uni-5",
		Delegator: utils.GenerateAccAddressForTest().String(),
		Distribution: []*icstypes.Distribution{
			{Valoper: "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn", Amount: 1000000},
			{Valoper: "junovaloper1z89utvygweg5l56fsk8ak7t6hh88fd0aa9ywed", Amount: 1000000},
		},
		Recipient:  "juno1z89utvygweg5l56fsk8ak7t6hh88fd0azcjpz5",
		Amount:     sdk.NewCoins(sdk.NewCoin("ujunox", sdk.NewInt(4000000))),
		BurnAmount: sdk.NewCoin("ujunox", sdk.NewInt(4000000)),
		Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
		Status:     icskeeper.WithdrawStatusQueued,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetWithdrawalRecord(suite.chainA.GetContext(), wRecord)

	err := suite.GetQuicksilverApp(suite.chainA).BankKeeper.MintCoins(suite.chainA.GetContext(), tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(202000000))))
	if err != nil {
		return
	}
	addr1, err := utils.AccAddressFromBech32("quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq", "quick")
	if err != nil {
		return
	}
	addr2, err := utils.AccAddressFromBech32("quick16x03wcp37kx5e8ehckjxvwcgk9j0cqnhcccnty", "quick")
	if err != nil {
		return
	}

	err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(suite.chainA.GetContext(), tokenfactorytypes.ModuleName, addr1, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(1600000))))
	if err != nil {
		return
	}
	err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(suite.chainA.GetContext(), tokenfactorytypes.ModuleName, addr2, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(200000000))))
	if err != nil {
		return
	}
	err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToModule(suite.chainA.GetContext(), tokenfactorytypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(400000))))
	if err != nil {
		return
	}
}

func (s *AppTestSuite) TestV010207UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)
	handler := v010207UpgradeHandler(app)
	ctx := s.chainA.GetContext()

	expectedVal := sdk.NewDec(50_000_000_000_000).Quo(sdk.NewDec(365))
	expectedProportions := minttypes.DistributionProportions{
		Staking:              sdk.NewDecWithPrec(80, 2),
		PoolIncentives:       sdk.NewDecWithPrec(17, 2),
		ParticipationRewards: sdk.NewDec(0),
		CommunityPool:        sdk.NewDecWithPrec(3, 2),
	}
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	// assert EpochProvisions
	minter := app.MintKeeper.GetMinter(ctx)
	s.Require().Equal(expectedVal, minter.EpochProvisions)

	// assert DistributionProportions
	params := app.MintKeeper.GetParams(ctx)
	s.Require().Equal(expectedProportions, params.DistributionProportions)
}
