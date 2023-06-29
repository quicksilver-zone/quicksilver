package app

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/suite"

	"github.com/ingenuity-build/quicksilver/app/upgrades"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	prtypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
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
func (s *AppTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)         // initializes 2 test chains
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	s.path = newQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)

	s.coordinator.CurrentTime = time.Now().UTC()
	s.coordinator.UpdateTime()

	s.initTestZone()
}

func (s *AppTestSuite) initTestZone() {
	// test zone
	zone := icstypes.Zone{
		ConnectionId:    s.path.EndpointA.ConnectionID,
		ChainId:         s.chainB.ChainID,
		AccountPrefix:   "bcosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: true,
		Is_118:          true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	// cosmos zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: false,
		Is_118:          true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	// osmosis zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         upgrades.OsmosisTestnetChainID,
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       false,
		LiquidityModule: true,
		Is_118:          true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	// uni-5 zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77003",
		ChainId:         "uni-5",
		AccountPrefix:   "juno",
		LocalDenom:      "uqjunox",
		BaseDenom:       "ujunox",
		MultiSend:       false,
		LiquidityModule: true,
		Is_118:          true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	receipt := icstypes.Receipt{
		ChainId: "uni-5",
		Sender:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
		Txhash:  "TestDeposit01",
		Amount: sdk.NewCoins(
			sdk.NewCoin(
				"ujunox",
				sdk.NewIntFromUint64(2000000), // 20% deposit
			),
		),
	}

	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetReceipt(s.chainA.GetContext(), receipt)

	ubRecord := icstypes.UnbondingRecord{
		ChainId:       "uni-5",
		EpochNumber:   1,
		Validator:     "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn",
		RelatedTxhash: []string{"ABC012"},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetUnbondingRecord(s.chainA.GetContext(), ubRecord)

	rdRecord := icstypes.RedelegationRecord{
		ChainId:        "uni-5",
		EpochNumber:    1,
		Source:         "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn",
		Destination:    "junovaloper1z89utvygweg5l56fsk8ak7t6hh88fd0aa9ywed",
		Amount:         3000000,
		CompletionTime: s.chainA.GetContext().BlockTime().Add(time.Hour),
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetRedelegationRecord(s.chainA.GetContext(), rdRecord)

	rdRecord = icstypes.RedelegationRecord{
		ChainId:        upgrades.OsmosisTestnetChainID,
		EpochNumber:    1,
		Source:         "osmovaloper1zxavllftfx3a3y5ldfyze7jnu5uyuktsfx2jcc",
		Destination:    "osmovaloper13eq5c99ym05jn02e78l8cac2fagzgdhh4294zk",
		Amount:         3000000,
		CompletionTime: time.Time{},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetRedelegationRecord(s.chainA.GetContext(), rdRecord)

	delRecord := icstypes.Delegation{
		Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000)),
		DelegationAddress: "juno1z89utvygweg5l56fsk8ak7t6hh88fd0azcjpz5",
		Height:            10,
		ValidatorAddress:  "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn",
		RedelegationEnd:   -62135596800,
	}

	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetDelegation(s.chainA.GetContext(), &zone, delRecord)

	wRecord := icstypes.WithdrawalRecord{
		ChainId:   "uni-5",
		Delegator: addressutils.GenerateAccAddressForTest().String(),
		Distribution: []*icstypes.Distribution{
			{Valoper: "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn", Amount: 1000000},
			{Valoper: "junovaloper1z89utvygweg5l56fsk8ak7t6hh88fd0aa9ywed", Amount: 1000000},
		},
		Recipient:  "juno1z89utvygweg5l56fsk8ak7t6hh88fd0azcjpz5",
		Amount:     sdk.NewCoins(sdk.NewCoin("ujunox", sdk.NewInt(4000000))),
		BurnAmount: sdk.NewCoin("ujunox", sdk.NewInt(4000000)),
		Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
		Status:     icstypes.WithdrawStatusQueued,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), wRecord)

	err := s.GetQuicksilverApp(s.chainA).BankKeeper.MintCoins(s.chainA.GetContext(), tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(202000000))))
	if err != nil {
		return
	}
	addr1 := addressutils.GenerateAccAddressForTest()
	addr2 := addressutils.GenerateAccAddressForTest()

	err = s.GetQuicksilverApp(s.chainA).BankKeeper.SendCoinsFromModuleToAccount(s.chainA.GetContext(), tokenfactorytypes.ModuleName, addr1, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(1600000))))
	if err != nil {
		return
	}
	err = s.GetQuicksilverApp(s.chainA).BankKeeper.SendCoinsFromModuleToAccount(s.chainA.GetContext(), tokenfactorytypes.ModuleName, addr2, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(200000000))))
	if err != nil {
		return
	}
	err = s.GetQuicksilverApp(s.chainA).BankKeeper.SendCoinsFromModuleToModule(s.chainA.GetContext(), tokenfactorytypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(400000))))
	if err != nil {
		return
	}
}

func (s *AppTestSuite) TestV010402rc1UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)

	handler := upgrades.V010402rc1UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	ctx := s.chainA.GetContext()
	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		if zone.ChainId == "uni-5" {

			zone.Validators = []*icstypes.Validator{
				{
					ValoperAddress:  "junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn",
					CommissionRate:  sdk.MustNewDecFromStr("0.2"),
					DelegatorShares: sdk.NewDec(2000),
					VotingPower:     sdk.NewInt(2000),
					Score:           sdk.NewDec(0),
					Status:          stakingtypes.BondStatusBonded,
				},

				{
					ValoperAddress:  "junovaloper1z89utvygweg5l56fsk8ak7t6hh88fd0aa9ywed",
					CommissionRate:  sdk.MustNewDecFromStr("0.2"),
					DelegatorShares: sdk.NewDec(2000),
					VotingPower:     sdk.NewInt(3000),
					Score:           sdk.NewDec(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			}
			app.InterchainstakingKeeper.SetZone(ctx, zone)
		}
		if zone.ChainId == "osmosis-1" {
			zone.Validators = []*icstypes.Validator{
				{
					ValoperAddress:  "osmovaloper1zxavllftfx3a3y5ldfyze7jnu5uyuktsfx2jcc",
					CommissionRate:  sdk.MustNewDecFromStr("0.2"),
					DelegatorShares: sdk.NewDec(2000),
					VotingPower:     sdk.NewInt(2000),
					Score:           sdk.NewDec(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					ValoperAddress:  "osmovaloper13eq5c99ym05jn02e78l8cac2fagzgdhh4294zk",
					CommissionRate:  sdk.MustNewDecFromStr("0.2"),
					DelegatorShares: sdk.NewDec(2000),
					VotingPower:     sdk.NewInt(3000),
					Score:           sdk.NewDec(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			}

			app.InterchainstakingKeeper.SetZone(ctx, zone)
		}
		return false
	})

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		if zone.ChainId == "uni-5" {
			s.Require().Nil(zone.Validators)
			valAddrBytes, err := addressutils.ValAddressFromBech32("junovaloper185hgkqs8q8ysnc8cvkgd8j2knnq2m0ah6ae73gntv9ampgwpmrxqlfzywn", "junovaloper")
			s.Require().NoError(err)
			_, found := app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
			s.Require().True(found)
			valAddrBytes2, err := addressutils.ValAddressFromBech32("junovaloper1z89utvygweg5l56fsk8ak7t6hh88fd0aa9ywed", "junovaloper")
			s.Require().NoError(err)
			_, found = app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes2)
			s.Require().True(found)
		}
		if zone.ChainId == "osmosis-1" {
			s.Require().Nil(zone.Validators)
			valAddrBytes, err := addressutils.ValAddressFromBech32("osmovaloper1zxavllftfx3a3y5ldfyze7jnu5uyuktsfx2jcc", "osmovaloper")
			s.Require().NoError(err)
			_, found := app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
			s.Require().True(found)
			valAddrBytes2, err := addressutils.ValAddressFromBech32("osmovaloper13eq5c99ym05jn02e78l8cac2fagzgdhh4294zk", "osmovaloper")
			s.Require().NoError(err)
			_, found = app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes2)
			s.Require().True(found)
		}

		return false
	})
}

func (s *AppTestSuite) TestV010402rc3UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)

	handler := upgrades.V010402rc3UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	ctx := s.chainA.GetContext()
	pdType, exists := prtypes.ProtocolDataType_value["ProtocolDataTypeConnection"]
	s.Require().True(exists)

	prData := prtypes.ProtocolData{
		Type: "ProtocolDataTypeConnection",
		Data: []byte(`{"ConnectionID":"connection-2","ChainID":"regen-redwood-1","Prefix":"regen"}`),
	}

	app.ParticipationRewardsKeeper.SetProtocolData(ctx, prtypes.GetProtocolDataKey(prtypes.ProtocolDataType(pdType), []byte("rege-redwood-1")), &prData)
	val0 := icstypes.Validator{ValoperAddress: "osmovaloper1zxavllftfx3a3y5ldfyze7jnu5uyuktsfx2jcc", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	app.InterchainstakingKeeper.SetValidator(ctx, upgrades.OsmosisTestnetChainID, val0)
	val1 := icstypes.Validator{ValoperAddress: "osmovaloper13eq5c99ym05jn02e78l8cac2fagzgdhh4294zk", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	app.InterchainstakingKeeper.SetValidator(ctx, upgrades.OsmosisTestnetChainID, val1)
	vals := app.InterchainstakingKeeper.GetValidators(ctx, upgrades.OsmosisTestnetChainID)
	s.Require().Equal(2, len(vals))

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	_, found := app.InterchainstakingKeeper.GetZone(ctx, upgrades.OsmosisTestnetChainID)
	s.Require().False(found)
	_, found = app.ParticipationRewardsKeeper.GetProtocolData(ctx, prtypes.ProtocolDataType(pdType), "rege-redwood-1")
	s.Require().False(found)

	vals = app.InterchainstakingKeeper.GetValidators(ctx, upgrades.OsmosisTestnetChainID)
	s.Require().Equal(0, len(vals))
}

func (s *AppTestSuite) TestV010404beta0UpgradeHandler() {
	app := s.GetQuicksilverApp(s.chainA)
	// osmosis zone
	zone := icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         upgrades.OsmosisTestnetChainID,
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       false,
		LiquidityModule: true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	handler := upgrades.V010404beta0UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	ctx := s.chainA.GetContext()

	zone, _ = app.InterchainstakingKeeper.GetZone(ctx, upgrades.OsmosisTestnetChainID)
	s.Require().False(zone.Is_118)

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)
	zone, _ = app.InterchainstakingKeeper.GetZone(ctx, upgrades.OsmosisTestnetChainID)
	s.Require().True(zone.Is_118)
}

// func (s *AppTestSuite) TestV010400rc6UpgradeHandler() {
//	app := s.GetQuicksilverApp(s.chainA)
//
//	handler := upgrades.V010400rc6UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
//	ctx := s.chainA.GetContext()
//
//	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(1, len(redelegations))
//
//	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
//	s.Require().NoError(err)
//
//	redelegations = app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(0, len(redelegations))
// }
//
// func (s *AppTestSuite) TestV010400rc8UpgradeHandler() {
//	app := s.GetQuicksilverApp(s.chainA)
//
//	handler := upgrades.V010400rc8UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
//	ctx := s.chainA.GetContext()
//
//	zone, _ := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
//	osmodels := []icstypes.Delegation{
//		{
//			Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000)),
//			DelegationAddress: "osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44",
//			Height:            10,
//			ValidatorAddress:  "osmovaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep88n0y4",
//			RedelegationEnd:   -62135596800,
//		}, {
//			Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17005)),
//			DelegationAddress: "osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44",
//			Height:            11,
//			ValidatorAddress:  "osmovaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpf6t4agt",
//			RedelegationEnd:   0,
//		},
//	}
//
//	for _, dels := range osmodels {
//		app.InterchainstakingKeeper.SetDelegation(ctx, &zone, dels)
//	}
//
//	zone, _ = app.InterchainstakingKeeper.GetZone(ctx, "uni-5")
//
//	var negRedelEndsBefore []icstypes.Delegation
//	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
//		app.InterchainstakingKeeper.IterateAllDelegations(ctx, zone, func(delegation icstypes.Delegation) (stop bool) {
//			if delegation.RedelegationEnd < 0 {
//				negRedelEndsBefore = append(negRedelEndsBefore, delegation)
//			}
//			return false
//		})
//		return false
//	})
//
//	s.Require().Equal(2, len(negRedelEndsBefore))
//
//	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(1, len(redelegations))
//
//	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
//	s.Require().NoError(err)
//
//	var negRedelEndsAfter []icstypes.Delegation
//
//	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
//		app.InterchainstakingKeeper.IterateAllDelegations(ctx, zone, func(delegation icstypes.Delegation) (stop bool) {
//			if delegation.RedelegationEnd < 0 {
//				negRedelEndsAfter = append(negRedelEndsAfter, delegation)
//			}
//			return false
//		})
//		return false
//	})
//
//	s.Require().Equal(0, len(negRedelEndsAfter))
//	redelegations = app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(0, len(redelegations))
// }

// func (s *AppTestSuite) TestV010400UpgradeHandler() {
//	app := s.GetQuicksilverApp(s.chainA)
//	handler := upgrades.V010400UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
//	ctx := s.chainA.GetContext()
//	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
//	s.Require().NoError(err)
//
//	osmosis, found := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
//	s.Require().True(found)
//	s.Require().Equal(int64(6), osmosis.Decimals)
//	s.Require().Equal("osmo", osmosis.AccountPrefix)
//	s.Require().Equal("connection-77002", osmosis.ConnectionId)
//	s.Require().False(osmosis.UnbondingEnabled)
//	s.Require().False(osmosis.ReturnToSender)
//	s.Require().True(osmosis.LiquidityModule)
//
//	cosmos, found := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
//	s.Require().True(found)
//	s.Require().Equal(int64(6), cosmos.Decimals)
//	s.Require().Equal("uatom", cosmos.BaseDenom)
//	s.Require().Equal("uqatom", cosmos.LocalDenom)
//	s.Require().False(cosmos.UnbondingEnabled)
//	s.Require().False(cosmos.ReturnToSender)
//	s.Require().False(cosmos.LiquidityModule)
//
//	chainb, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
//	s.Require().True(found)
//	s.Require().Equal(int64(6), chainb.Decimals)
//	s.Require().False(chainb.UnbondingEnabled)
//	s.Require().False(chainb.ReturnToSender)
//	s.Require().True(chainb.LiquidityModule)
//
//	juno, found := app.InterchainstakingKeeper.GetZone(ctx, "uni-5")
//	s.Require().False(found)
//
//	reciepts := app.InterchainstakingKeeper.AllReceipts(ctx)
//	s.Require().Equal(0, len(reciepts))
//
//	unbondings := app.InterchainstakingKeeper.AllZoneUnbondingRecords(ctx, "uni-5")
//	s.Require().Equal(0, len(unbondings))
//
//	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "uni-5")
//	s.Require().Equal(0, len(redelegations))
//
//	delegations := app.InterchainstakingKeeper.GetAllDelegations(ctx, &juno)
//	s.Require().Equal(0, len(delegations))
//
//	perfDelegations := app.InterchainstakingKeeper.GetAllPerformanceDelegations(ctx, &juno)
//	s.Require().Equal(0, len(perfDelegations))
//
//	_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "uni-5", "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D", icskeeper.WithdrawStatusQueued)
//	s.Require().False(found)
// }
//
// func (s *AppTestSuite) TestV010400rc6UpgradeHandler() {
//	app := s.GetQuicksilverApp(s.chainA)
//
//	handler := upgrades.V010400rc6UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
//	ctx := s.chainA.GetContext()
//
//	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(1, len(redelegations))
//
//	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
//	s.Require().NoError(err)
//
//	redelegations = app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(0, len(redelegations))
// }
//
// func (s *AppTestSuite) TestV010400rc8UpgradeHandler() {
//	app := s.GetQuicksilverApp(s.chainA)
//
//	handler := upgrades.V010400rc8UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
//	ctx := s.chainA.GetContext()
//
//	zone, _ := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
//	osmodels := []icstypes.Delegation{
//		{
//			Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000)),
//			DelegationAddress: "osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44",
//			Height:            10,
//			ValidatorAddress:  "osmovaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep88n0y4",
//			RedelegationEnd:   -62135596800,
//		}, {
//			Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17005)),
//			DelegationAddress: "osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44",
//			Height:            11,
//			ValidatorAddress:  "osmovaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpf6t4agt",
//			RedelegationEnd:   0,
//		},
//	}
//
//	for _, dels := range osmodels {
//		app.InterchainstakingKeeper.SetDelegation(ctx, &zone, dels)
//	}
//
//	zone, _ = app.InterchainstakingKeeper.GetZone(ctx, "uni-5")
//
//	var negRedelEndsBefore []icstypes.Delegation
//	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
//		app.InterchainstakingKeeper.IterateAllDelegations(ctx, zone, func(delegation icstypes.Delegation) (stop bool) {
//			if delegation.RedelegationEnd < 0 {
//				negRedelEndsBefore = append(negRedelEndsBefore, delegation)
//			}
//			return false
//		})
//		return false
//	})
//
//	s.Require().Equal(2, len(negRedelEndsBefore))
//
//	redelegations := app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(1, len(redelegations))
//
//	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
//	s.Require().NoError(err)
//
//	var negRedelEndsAfter []icstypes.Delegation
//
//	app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
//		app.InterchainstakingKeeper.IterateAllDelegations(ctx, zone, func(delegation icstypes.Delegation) (stop bool) {
//			if delegation.RedelegationEnd < 0 {
//				negRedelEndsAfter = append(negRedelEndsAfter, delegation)
//			}
//			return false
//		})
//		return false
//	})
//
//	s.Require().Equal(0, len(negRedelEndsAfter))
//	redelegations = app.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "osmosis-1")
//	s.Require().Equal(0, len(redelegations))
// }
