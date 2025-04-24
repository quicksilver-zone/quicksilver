package app

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	"github.com/quicksilver-zone/quicksilver/app/upgrades"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
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

func (*AppTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *Quicksilver {
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
}

func addVestingAccount(ctx sdk.Context, ak *authkeeper.AccountKeeper, address string, numPeriods int64, periodLength int64, total int64) {
	start := int64(1704240000)
	duration := numPeriods * periodLength
	perPeriod := total / numPeriods
	dust := total - (perPeriod * numPeriods)

	periods := make(vestingtypes.Periods, 0, numPeriods)
	for i := numPeriods; i > 0; i-- {
		periods = append(periods, vestingtypes.Period{Length: periodLength, Amount: sdk.NewCoins(sdk.NewCoin("uqck", math.NewInt(perPeriod)))})
	}
	periods[0].Amount.Add(sdk.NewCoin("uqck", math.NewInt(dust)))
	vest := vestingtypes.NewPeriodicVestingAccountRaw(
		vestingtypes.NewBaseVestingAccount(
			authtypes.NewBaseAccountWithAddress(addressutils.MustAccAddressFromBech32(address, "")),
			sdk.NewCoins(sdk.NewCoin("uqck", math.NewInt(total))),
			start+duration,
		),
		start,
		periods,
	)
	ak.SetAccount(ctx, vest)
}

func (s *AppTestSuite) InitV150TestZones() {
	// zone to match prod
	zone := icstypes.Zone{ConnectionId: "connection-1", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-0", ChainId: "stargaze-1", AccountPrefix: "stars", LocalDenom: "uqstars", BaseDenom: "ustars"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-50", ChainId: "juno-1", AccountPrefix: "juno", LocalDenom: "uqjuno", BaseDenom: "ujuno"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-2", ChainId: "osmosis-1", AccountPrefix: "osmo", LocalDenom: "uqosmo", BaseDenom: "uosmo"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-9", ChainId: "regen-1", AccountPrefix: "regen", LocalDenom: "uqregen", BaseDenom: "uregen"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-54", ChainId: "sommelier-3", AccountPrefix: "somm", LocalDenom: "uqsomm", BaseDenom: "usomm"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l", 10, 864000, 5000000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f", 10, 864000, 5000000000)

	// set counterparty channels to match prod so we can assert denoms
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-0", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-124")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-1", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-467")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-2", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-522")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-52", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-65")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-49", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-53")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-101", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-59")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-86", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-272")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-17", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-62")})
}

// Init a zone with some zero burnAmount withdrawal records

func (s *AppTestSuite) InitV160rc0TestZone() {
	cosmosWithdrawal := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosPerformance := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDeposit := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDelegate := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	// cosmos zone
	zone := icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           cosmosWithdrawal,
			PortName:          "icacontroller-cosmoshub-4.withdrawal",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           cosmosDelegate,
			PortName:          "icacontroller-cosmoshub-4.delegate",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           cosmosDeposit,
			PortName:          "icacontroller-cosmoshub-4.deposit",
			WithdrawalAddress: cosmosWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           cosmosPerformance,
			PortName:          "icacontroller-cosmoshub-4.performance",
			WithdrawalAddress: cosmosWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
}

// UncheckedSetWithdrawalRecord store the withdrawal record without checking the burnAmount.
// WARNING: This function is intended for testing purposes only and should not be used in production code.
func (s *AppTestSuite) UncheckedSetWithdrawalRecord(ctx sdk.Context, app *Quicksilver, record icstypes.WithdrawalRecord) {
	key, err := hex.DecodeString(record.Txhash)
	if err != nil {
		panic(err)
	}

	store := prefix.NewStore(ctx.KVStore(app.GetKey(icstypes.StoreKey)), icstypes.GetWithdrawalKey(record.ChainId, record.Status))
	bz := app.InterchainstakingKeeper.GetCodec().MustMarshal(&record)
	store.Set(key, bz)
}

func (s *AppTestSuite) InitV160TestZones() {
	cosmosWithdrawal := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosPerformance := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDeposit := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDelegate := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	// cosmos zone
	zone := icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           cosmosWithdrawal,
			PortName:          "icacontroller-cosmoshub-4.withdrawal",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           cosmosDelegate,
			PortName:          "icacontroller-cosmoshub-4.delegate",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           cosmosDeposit,
			PortName:          "icacontroller-cosmoshub-4.deposit",
			WithdrawalAddress: cosmosWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           cosmosPerformance,
			PortName:          "icacontroller-cosmoshub-4.performance",
			WithdrawalAddress: cosmosWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	osmoWithdrawal := addressutils.GenerateAddressForTestWithPrefix("osmo")
	osmoPerformance := addressutils.GenerateAddressForTestWithPrefix("osmo")
	osmoDeposit := addressutils.GenerateAddressForTestWithPrefix("osmo")
	osmoDelegate := addressutils.GenerateAddressForTestWithPrefix("osmo")
	// osmosis zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           osmoWithdrawal,
			PortName:          "icacontroller-osmosis-1.withdrawal",
			WithdrawalAddress: osmoWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           osmoDelegate,
			PortName:          "icacontroller-osmosis-1.delegate",
			WithdrawalAddress: osmoWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           osmoDeposit,
			PortName:          "icacontroller-osmosis-1.deposit",
			WithdrawalAddress: osmoWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           osmoPerformance,
			PortName:          "icacontroller-osmosis-1.performance",
			WithdrawalAddress: osmoWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	// uni-5 zone

	junoWithdrawal := addressutils.GenerateAddressForTestWithPrefix("juno")
	junoPerformance := addressutils.GenerateAddressForTestWithPrefix("juno")
	junoDeposit := addressutils.GenerateAddressForTestWithPrefix("juno")
	junoDelegate := addressutils.GenerateAddressForTestWithPrefix("juno")

	zone = icstypes.Zone{
		ConnectionId:    "connection-77003",
		ChainId:         "juno-1",
		AccountPrefix:   "juno",
		LocalDenom:      "uqjuno",
		BaseDenom:       "ujuno",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           junoWithdrawal,
			PortName:          "icacontroller-juno-1.withdrawal",
			WithdrawalAddress: junoWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           junoDelegate,
			PortName:          "icacontroller-juno-1.delegate",
			WithdrawalAddress: junoWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           junoDeposit,
			PortName:          "icacontroller-juno-1.deposit",
			WithdrawalAddress: junoWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           junoPerformance,
			PortName:          "icacontroller-juno-1.performance",
			WithdrawalAddress: junoWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1qfyntnmlvznvrkk9xqppmcxqcluv7wd74nmyus", 10, 864000, 5000000000)

	// set withdrawal records
	invalidWithdrawal := icstypes.WithdrawalRecord{
		ChainId:        zone.ChainId,
		Delegator:      cosmosDelegate,
		Recipient:      cosmosWithdrawal,
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(0)),
		Requeued:       true,
		Txhash:         fmt.Sprintf("%064d", 1),
		Acknowledged:   false,
		Status:         icstypes.WithdrawStatusQueued,
		EpochNumber:    1,
		CompletionTime: time.Time{},
	}
	s.UncheckedSetWithdrawalRecord(s.chainA.GetContext(), s.GetQuicksilverApp(s.chainA), invalidWithdrawal)
	validWithdrawal := icstypes.WithdrawalRecord{
		ChainId:        zone.ChainId,
		Delegator:      cosmosDelegate,
		Recipient:      cosmosWithdrawal,
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(1000)),
		Requeued:       true,
		Txhash:         fmt.Sprintf("%064d", 2),
		Acknowledged:   false,
		Status:         icstypes.WithdrawStatusQueued,
		EpochNumber:    1,
		CompletionTime: time.Time{},
	}
	s.UncheckedSetWithdrawalRecord(s.chainA.GetContext(), s.GetQuicksilverApp(s.chainA), validWithdrawal)

	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-2", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-522")})

	// agoric-3 zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-12312",
		ChainId:         "agoric-3",
		AccountPrefix:   "agoric",
		LocalDenom:      "uqbld",
		BaseDenom:       "ubld",
		MultiSend:       false,
		LiquidityModule: false,
		Is_118:          true,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
}

func (s *AppTestSuite) TestV010705UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	ctx := s.chainA.GetContext()
	completion, err := time.Parse("2006-01-02T15:04:05Z", "2024-12-20T17:00:47Z")
	s.NoError(err)

	record := icstypes.WithdrawalRecord{
		ChainId:        "cosmoshub-4",
		Delegator:      "quick1efpktthkfsuzqhsnqdyyjxv5fl9eemchrlmhyd",
		Recipient:      "cosmos1efpktthkfsuzqhsnqdyyjxv5fl9eemchgmt9al",
		Txhash:         "02c2d4bcb869b9ddf26540c2854c2ca09d70492a3831170da293f4101fda32b3",
		Status:         icstypes.WithdrawStatusUnbond,
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(3534090000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(4841699172))),
		Acknowledged:   true,
		CompletionTime: completion,
	}

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.True(found)

	delegation := icstypes.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  "cosmosvaloper1efpktthkfsuzqhsnqdyyjxv5fl9eemchd0ls3v",
		Amount:            sdk.NewCoin("uatom", math.NewInt(37848188596)),
	}

	app.InterchainstakingKeeper.SetDelegation(ctx, "cosmoshub-4", delegation)

	s.NoError(app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(30811987786)))))
	s.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(3534090000)))))
	s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record))

	handler := upgrades.V010705UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)

	_, err = handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	zone, found = app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.True(found)
	s.Equal(sdk.NewDecWithPrec(1387503864591246254, 18), zone.RedemptionRate)
	s.Equal(sdk.NewDecWithPrec(138, 2), zone.LastRedemptionRate)
}

func (s *AppTestSuite) TestV010706UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	ctx := s.chainA.GetContext()
	completion, err := time.Parse("2006-01-02T15:04:05Z", "2025-01-30T18:17:24Z")
	s.NoError(err)

	record := icstypes.WithdrawalRecord{
		ChainId:        "regen-1",
		Delegator:      "quick1nvyqrve35fpzgrewnaxmyxsqtaq4dxvydf5gpj",
		Recipient:      "regen1nvyqrve35fpzgrewnaxmyxsqtaq4dxvye00xwy",
		Txhash:         "ee0b5f5c423508c8dd6a501168a77a0b72d5a8aaf1702a64804e522334ff272b",
		Status:         icstypes.WithdrawStatusUnbond,
		BurnAmount:     sdk.NewCoin("uqregen", math.NewInt(31569450000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("uregen", math.NewInt(43107830342))),
		Acknowledged:   true,
		CompletionTime: completion,
		EpochNumber:    230,
		SendErrors:     187,
	}

	s.NoError(app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqregen", math.NewInt(31569450000)))))
	s.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqregen", math.NewInt(31569450000)))))
	s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record))

	completion, err = time.Parse("2006-01-02T15:04:05Z", "2025-01-30T16:45:43Z")
	s.NoError(err)

	record = icstypes.WithdrawalRecord{
		ChainId:        "sommelier-3",
		Delegator:      "quick1t5zgnfz0jrvflywjmgs95rey3un57n42jr7qd4",
		Recipient:      "somm1t5zgnfz0jrvflywjmgs95rey3un57n424mp79d",
		Txhash:         "a55f1f4deaa501ff5671ef96fbbb5b60e225d4b8db4825ae3706893bb94e052c",
		Status:         icstypes.WithdrawStatusUnbond,
		BurnAmount:     sdk.NewCoin("uqsomm", math.NewInt(1231350000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("usomm", math.NewInt(1320763919))),
		Acknowledged:   true,
		CompletionTime: completion,
		EpochNumber:    230,
		SendErrors:     150,
	}

	s.NoError(app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqsomm", math.NewInt(1231350000)))))
	s.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqsomm", math.NewInt(1231350000)))))
	s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record))

	// setup unbondign records
	// no ct
	ubr1 := icstypes.UnbondingRecord{
		ChainId:        "sommelier-3",
		Validator:      "sommvaloper10m5g48u53vss7xqmqw6089d02ua0d85d7s8gu0",
		EpochNumber:    133,
		CompletionTime: time.Time{},
		RelatedTxhash: []string{
			"0fc9c66af331cbb3b015f97a2257220e16c2d58f75f5cdcbce3e77cb90570834",
			"1b3c51013bc5753b728083703cd00dd236f14759dc1e13f20899e4e8863baa66",
			"204c56e09f176373d73a63be060f560eb4d2bda83ffedd838de814e6175c2bc4",
			"24fc99033051c0504b85e3da033a6aa006ac76c373ef60791f8046c7713814e4",
		},
		Amount: sdk.NewCoin("usomm", math.NewInt(2759679739)),
	}

	completion, err = time.Parse("2006-01-02T15:04:05Z", "2025-01-20T16:45:43Z")
	s.NoError(err)

	ubr2 := icstypes.UnbondingRecord{
		ChainId:        "sommelier-3",
		Validator:      "sommvaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d",
		EpochNumber:    238,
		CompletionTime: completion,
		RelatedTxhash: []string{
			"184365ad8ea78478cf49e7fb28c829f64aa9dac026905bef239f8fb16cc3dcb6",
		},
		Amount: sdk.NewCoin("usomm", math.NewInt(19803)),
	}

	ubr3 := icstypes.UnbondingRecord{
		ChainId:        "sommelier-3",
		Validator:      "sommvaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d",
		EpochNumber:    243,
		CompletionTime: ctx.BlockTime().Add(-time.Hour * 12),
		RelatedTxhash: []string{
			"bb09ae31a44f83a292c4ddf16870b033daccc1ea9fb1620ccf36c1e232af11c3",
		},
		Amount: sdk.NewCoin("usomm", math.NewInt(95016053)),
	}

	app.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr1)
	app.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr2)
	app.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr3)

	handler := upgrades.V010706UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)

	_, err = handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	_, found := app.InterchainstakingKeeper.GetUnbondingRecord(ctx, "sommelier-3", "sommvaloper10m5g48u53vss7xqmqw6089d02ua0d85d7s8gu0", 133)
	s.False(found)

	_, found = app.InterchainstakingKeeper.GetUnbondingRecord(ctx, "sommelier-3", "sommvaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d", 238)
	s.False(found)

	_, found = app.InterchainstakingKeeper.GetUnbondingRecord(ctx, "sommelier-3", "sommvaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d", 243)
	s.True(found)
}

func (s *AppTestSuite) TestV010800UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	ctx := s.chainA.GetContext()
	completion, err := time.Parse("2006-01-02T15:04:05Z", "2025-04-11T11:17:29Z")
	s.NoError(err)

	record := icstypes.WithdrawalRecord{
		ChainId:        "juno-1",
		Delegator:      "quick1npxuk4c30xm3q27lyx24vvrhqzm089wdxmhe93",
		Recipient:      "juno1npxuk4c30xm3q27lyx24vvrhqzm089wdmdysml",
		Txhash:         "564e8a6263763644bbe32e4bd0bf9f99619aaf68b938216fff2acef2dfb8aec6",
		Status:         icstypes.WithdrawStatusSend,
		BurnAmount:     sdk.NewCoin("uqjuno", math.NewInt(59560000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("ujuno", math.NewInt(82121494))),
		Requeued:       false,
		Acknowledged:   true,
		CompletionTime: completion,
		EpochNumber:    266,
		SendErrors:     0,
	}

	s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record))

	completion, err = time.Parse("2006-01-02T15:04:05Z", "2025-04-11T11:20:26Z")
	s.NoError(err)

	record = icstypes.WithdrawalRecord{
		ChainId:        "juno-1",
		Delegator:      "quick194dawsp29zcp4s9r6hdppdak0cy5kf3xqumt9x",
		Recipient:      "juno194dawsp29zcp4s9r6hdppdak0cy5kf3xa2gzmg",
		Txhash:         "c746ceba8da060f25a81f2e0cc6ed53fecd69dbbd89ff7c9aa8b5d0464302f84",
		Status:         icstypes.WithdrawStatusUnbond,
		BurnAmount:     sdk.NewCoin("uqjuno", math.NewInt(773810000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("ujuno", math.NewInt(1066930053))),
		Requeued:       false,
		Acknowledged:   true,
		CompletionTime: completion,
		EpochNumber:    266,
		SendErrors:     0,
	}

	s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record))

	// setup unbondign records
	// no ct
	ubr1 := icstypes.UnbondingRecord{
		ChainId:        "juno-1",
		Validator:      "junovaloper10m5g48u53vss7xqmqw6089d02ua0d85d7s8gu0",
		EpochNumber:    277,
		CompletionTime: time.Time{},
		RelatedTxhash: []string{
			"0fc9c66af331cbb3b015f97a2257220e16c2d58f75f5cdcbce3e77cb90570834",
			"1b3c51013bc5753b728083703cd00dd236f14759dc1e13f20899e4e8863baa66",
			"204c56e09f176373d73a63be060f560eb4d2bda83ffedd838de814e6175c2bc4",
			"24fc99033051c0504b85e3da033a6aa006ac76c373ef60791f8046c7713814e4",
		},
		Amount: sdk.NewCoin("usomm", math.NewInt(2759679739)),
	}

	completion, err = time.Parse("2006-01-02T15:04:05Z", "2025-01-20T16:45:43Z")
	s.NoError(err)

	ubr2 := icstypes.UnbondingRecord{
		ChainId:        "juno-1",
		Validator:      "junovaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d",
		EpochNumber:    277,
		CompletionTime: completion,
		RelatedTxhash: []string{
			"184365ad8ea78478cf49e7fb28c829f64aa9dac026905bef239f8fb16cc3dcb6",
		},
		Amount: sdk.NewCoin("usomm", math.NewInt(19803)),
	}

	ubr3 := icstypes.UnbondingRecord{
		ChainId:        "sommelier-3",
		Validator:      "sommvaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d",
		EpochNumber:    277,
		CompletionTime: ctx.BlockTime().Add(time.Hour * 36),
		RelatedTxhash: []string{
			"bb09ae31a44f83a292c4ddf16870b033daccc1ea9fb1620ccf36c1e232af11c3",
		},
		Amount: sdk.NewCoin("usomm", math.NewInt(95016053)),
	}

	app.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr1)
	app.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr2)
	app.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr3)

	handler := upgrades.V010800UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)

	_, err = handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	_, found := app.InterchainstakingKeeper.GetUnbondingRecord(ctx, "juno-1", "junovaloper10m5g48u53vss7xqmqw6089d02ua0d85d7s8gu0", 277)
	s.False(found)

	_, found = app.InterchainstakingKeeper.GetUnbondingRecord(ctx, "juno-1", "junovaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d", 277)
	s.False(found)

	_, found = app.InterchainstakingKeeper.GetUnbondingRecord(ctx, "sommelier-3", "sommvaloper1y0few0kgxa7vtq8nskjsdwdtyqglj3k5pv2c4d", 277)
	s.True(found)

	_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "juno-1", "quick1npxuk4c30xm3q27lyx24vvrhqzm089wdxmhe93", 266)
	s.False(found)

	_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "juno-1", "quick194dawsp29zcp4s9r6hdppdak0cy5kf3xqumt9x", 266)
	s.False(found)
}
