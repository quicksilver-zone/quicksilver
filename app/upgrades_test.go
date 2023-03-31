package app

import (
<<<<<<< HEAD
	"testing"
	"time"
=======
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
>>>>>>> origin/main

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/app/upgrades"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

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

<<<<<<< HEAD
=======
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

>>>>>>> origin/main
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

<<<<<<< HEAD
	receipt := icstypes.Receipt{
=======
	reciept := icstypes.Receipt{
>>>>>>> origin/main
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

<<<<<<< HEAD
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetReceipt(suite.chainA.GetContext(), receipt)
=======
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetReceipt(suite.chainA.GetContext(), reciept)
>>>>>>> origin/main

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

<<<<<<< HEAD
func (s *AppTestSuite) TestV010400UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)
	handler := upgrades.V010400UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	ctx := s.chainA.GetContext()
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	osmosis, found := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
	s.Require().True(found)
	s.Require().Equal(int64(6), osmosis.Decimals)
	s.Require().Equal("osmo", osmosis.AccountPrefix)
	s.Require().Equal("connection-77002", osmosis.ConnectionId)
	s.Require().False(osmosis.UnbondingEnabled)
	s.Require().False(osmosis.ReturnToSender)
	s.Require().True(osmosis.LiquidityModule)

	cosmos, found := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.Require().True(found)
	s.Require().Equal(int64(6), cosmos.Decimals)
	s.Require().Equal("uatom", cosmos.BaseDenom)
	s.Require().Equal("uqatom", cosmos.LocalDenom)
	s.Require().False(cosmos.UnbondingEnabled)
	s.Require().False(cosmos.ReturnToSender)
	s.Require().False(cosmos.LiquidityModule)

	chainb, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(int64(6), chainb.Decimals)
	s.Require().False(chainb.UnbondingEnabled)
	s.Require().False(chainb.ReturnToSender)
	s.Require().True(chainb.LiquidityModule)

	juno, found := app.InterchainstakingKeeper.GetZone(ctx, "uni-5")
	s.Require().False(found)

	reciepts := app.InterchainstakingKeeper.AllReceipts(ctx)
	s.Require().Equal(0, len(reciepts))

	unbondings := app.InterchainstakingKeeper.AllZoneUnbondingRecords(ctx, "uni-5")
	s.Require().Equal(0, len(unbondings))

	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "uni-5")
	s.Require().Equal(0, len(redelegations))

	delegations := app.InterchainstakingKeeper.GetAllDelegations(ctx, &juno)
	s.Require().Equal(0, len(delegations))

	perfDelegations := app.InterchainstakingKeeper.GetAllPerformanceDelegations(ctx, &juno)
	s.Require().Equal(0, len(perfDelegations))

	_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "uni-5", "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D", icskeeper.WithdrawStatusQueued)
	s.Require().False(found)
}

func (s *AppTestSuite) TestV010400rc6UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)

	handler := upgrades.V010400rc6UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	ctx := s.chainA.GetContext()

	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
	s.Require().Equal(1, len(redelegations))

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	redelegations = app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
	s.Require().Equal(0, len(redelegations))
}

func (s *AppTestSuite) TestV010400rc8UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)

	handler := upgrades.V010400rc8UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	ctx := s.chainA.GetContext()

	zone, _ := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
	osmodels := []icstypes.Delegation{
		{
			Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000)),
			DelegationAddress: "osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44",
			Height:            10,
			ValidatorAddress:  "osmovaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep88n0y4",
			RedelegationEnd:   -62135596800,
		}, {
			Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17005)),
			DelegationAddress: "osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44",
			Height:            11,
			ValidatorAddress:  "osmovaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpf6t4agt",
			RedelegationEnd:   0,
		},
	}

	for _, dels := range osmodels {
		app.InterchainstakingKeeper.SetDelegation(ctx, &zone, dels)
	}

	zone, _ = app.InterchainstakingKeeper.GetZone(ctx, "uni-5")

	var negRedelEndsBefore []icstypes.Delegation
	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		app.InterchainstakingKeeper.IterateAllDelegations(ctx, zone, func(delegation icstypes.Delegation) (stop bool) {
			if delegation.RedelegationEnd < 0 {
				negRedelEndsBefore = append(negRedelEndsBefore, delegation)
			}
			return false
		})
		return false
	})

	s.Require().Equal(2, len(negRedelEndsBefore))

	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
	s.Require().Equal(1, len(redelegations))

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	var negRedelEndsAfter []icstypes.Delegation

	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		app.InterchainstakingKeeper.IterateAllDelegations(ctx, zone, func(delegation icstypes.Delegation) (stop bool) {
			if delegation.RedelegationEnd < 0 {
				negRedelEndsAfter = append(negRedelEndsAfter, delegation)
			}
			return false
		})
		return false
	})

	s.Require().Equal(0, len(negRedelEndsAfter))
	redelegations = app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
	s.Require().Equal(0, len(redelegations))
=======
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

>>>>>>> origin/main
}
