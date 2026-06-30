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

// --- v1.10.2 sunset tests (mint-to-cover strategy) ---

// setupV0101002Zones writes stargaze-1 and omniflixhub-1 zones and returns three
// test user bech32 addresses (2 for stargaze, 1 for omniflix).
func (s *AppTestSuite) setupV0101002Zones(ctx sdk.Context, app *Quicksilver) (string, string, string) {
	app.InterchainstakingKeeper.SetZone(ctx, &icstypes.Zone{
		ConnectionId:     "connection-3",
		ChainId:          "stargaze-1",
		AccountPrefix:    "stars",
		LocalDenom:       "uqstars",
		BaseDenom:        "ustars",
		DepositsEnabled:  true,
		UnbondingEnabled: true,
	})
	app.InterchainstakingKeeper.SetZone(ctx, &icstypes.Zone{
		ConnectionId:     "connection-4",
		ChainId:          "omniflixhub-1",
		AccountPrefix:    "omniflix",
		LocalDenom:       "uqflix",
		BaseDenom:        "uflix",
		DepositsEnabled:  true,
		UnbondingEnabled: true,
	})
	return addressutils.GenerateAddressForTestWithPrefix("quick"),
		addressutils.GenerateAddressForTestWithPrefix("quick"),
		addressutils.GenerateAddressForTestWithPrefix("quick")
}

func (s *AppTestSuite) seedWDR(ctx sdk.Context, app *Quicksilver, chainID, delegator, denom, baseDenom string, burn int64, status int32, tx int) {
	s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
		ChainId:     chainID,
		Delegator:   delegator,
		Recipient:   addressutils.GenerateAddressForTestWithPrefix("stars"),
		BurnAmount:  sdk.NewCoin(denom, math.NewInt(burn)),
		Amount:      sdk.NewCoins(sdk.NewCoin(baseDenom, math.NewInt(burn))),
		Txhash:      fmt.Sprintf("%064d", tx),
		Status:      status,
		EpochNumber: 1,
	}))
}

func (s *AppTestSuite) fundEscrow(ctx sdk.Context, app *Quicksilver, coins sdk.Coins) {
	s.NoError(app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, coins))
	s.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, coins))
}

// Shortfall: escrow < sum(burn_amount). Mint-to-cover must mint the delta and
// refund every user in full.
func (s *AppTestSuite) TestV0101002UpgradeHandler_MintToCover_Shortfall() {
	s.SetupTest()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	u1, u2, u3 := s.setupV0101002Zones(ctx, app)

	// stargaze obligation: 5,000,000 + 3,000,000 = 8,000,000 uqstars
	s.seedWDR(ctx, app, "stargaze-1", u1, "uqstars", "ustars", 5_000_000, icstypes.WithdrawStatusQueued, 100)
	s.seedWDR(ctx, app, "stargaze-1", u2, "uqstars", "ustars", 3_000_000, icstypes.WithdrawStatusUnbond, 101)
	// omniflix obligation: 2,000,000 uqflix
	s.seedWDR(ctx, app, "omniflixhub-1", u3, "uqflix", "uflix", 2_000_000, icstypes.WithdrawStatusQueued, 102)

	// Escrow is short: stargaze has only 6,000,000 vs 8,000,000 needed (2M short);
	// omniflix has only 500,000 vs 2,000,000 needed (1.5M short).
	s.fundEscrow(ctx, app, sdk.NewCoins(
		sdk.NewCoin("uqstars", math.NewInt(6_000_000)),
		sdk.NewCoin("uqflix", math.NewInt(500_000)),
	))

	stargazeSupplyBefore := app.BankKeeper.GetSupply(ctx, "uqstars").Amount
	omniflixSupplyBefore := app.BankKeeper.GetSupply(ctx, "uqflix").Amount

	handler := upgrades.V0101002UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	// Zones offboarded.
	sz, ok := app.InterchainstakingKeeper.GetZone(ctx, "stargaze-1")
	s.True(ok)
	s.True(sz.IsOffboarding)
	s.False(sz.DepositsEnabled)
	s.False(sz.UnbondingEnabled)
	oz, ok := app.InterchainstakingKeeper.GetZone(ctx, "omniflixhub-1")
	s.True(ok)
	s.True(oz.IsOffboarding)

	// Records consumed.
	s.Equal(0, len(app.InterchainstakingKeeper.AllZoneWithdrawalRecords(ctx, "stargaze-1")))
	s.Equal(0, len(app.InterchainstakingKeeper.AllZoneWithdrawalRecords(ctx, "omniflixhub-1")))

	// Users refunded their full BurnAmount.
	s.Equal(math.NewInt(5_000_000), app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(u1, ""), "uqstars").Amount)
	s.Equal(math.NewInt(3_000_000), app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(u2, ""), "uqstars").Amount)
	s.Equal(math.NewInt(2_000_000), app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(u3, ""), "uqflix").Amount)

	// Escrow drained for both denoms.
	escrowAddr := app.AccountKeeper.GetModuleAddress(icstypes.EscrowModuleAccount)
	s.True(app.BankKeeper.GetBalance(ctx, escrowAddr, "uqstars").IsZero())
	s.True(app.BankKeeper.GetBalance(ctx, escrowAddr, "uqflix").IsZero())

	// Supply increased by exactly the shortfall (2M uqstars, 1.5M uqflix).
	s.Equal(stargazeSupplyBefore.Add(math.NewInt(2_000_000)), app.BankKeeper.GetSupply(ctx, "uqstars").Amount)
	s.Equal(omniflixSupplyBefore.Add(math.NewInt(1_500_000)), app.BankKeeper.GetSupply(ctx, "uqflix").Amount)
}

// Exact: escrow == sum(burn_amount). No mint, no leftover.
func (s *AppTestSuite) TestV0101002UpgradeHandler_MintToCover_Exact() {
	s.SetupTest()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	u1, _, _ := s.setupV0101002Zones(ctx, app)
	s.seedWDR(ctx, app, "stargaze-1", u1, "uqstars", "ustars", 7_000_000, icstypes.WithdrawStatusUnbond, 200)
	s.fundEscrow(ctx, app, sdk.NewCoins(sdk.NewCoin("uqstars", math.NewInt(7_000_000))))

	supplyBefore := app.BankKeeper.GetSupply(ctx, "uqstars").Amount

	handler := upgrades.V0101002UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	s.Equal(supplyBefore, app.BankKeeper.GetSupply(ctx, "uqstars").Amount, "no mint expected when escrow matches obligation")
	s.Equal(math.NewInt(7_000_000), app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(u1, ""), "uqstars").Amount)
	escrowAddr := app.AccountKeeper.GetModuleAddress(icstypes.EscrowModuleAccount)
	s.True(app.BankKeeper.GetBalance(ctx, escrowAddr, "uqstars").IsZero())
}

// Surplus: escrow > sum(burn_amount). No mint. Leftover remains in escrow
// (to be reclaimed by a follow-up governance action).
func (s *AppTestSuite) TestV0101002UpgradeHandler_MintToCover_Surplus() {
	s.SetupTest()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	u1, _, _ := s.setupV0101002Zones(ctx, app)
	s.seedWDR(ctx, app, "stargaze-1", u1, "uqstars", "ustars", 4_000_000, icstypes.WithdrawStatusQueued, 300)
	s.fundEscrow(ctx, app, sdk.NewCoins(sdk.NewCoin("uqstars", math.NewInt(10_000_000))))

	supplyBefore := app.BankKeeper.GetSupply(ctx, "uqstars").Amount

	handler := upgrades.V0101002UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	s.Equal(supplyBefore, app.BankKeeper.GetSupply(ctx, "uqstars").Amount, "no mint when escrow > obligation")
	s.Equal(math.NewInt(4_000_000), app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(u1, ""), "uqstars").Amount)
	escrowAddr := app.AccountKeeper.GetModuleAddress(icstypes.EscrowModuleAccount)
	s.Equal(math.NewInt(6_000_000), app.BankKeeper.GetBalance(ctx, escrowAddr, "uqstars").Amount, "surplus stays in escrow for later reclaim")
}

// No records: zone is still offboarded but nothing is refunded.
func (s *AppTestSuite) TestV0101002UpgradeHandler_NoRecords() {
	s.SetupTest()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	s.setupV0101002Zones(ctx, app)
	supplyBefore := app.BankKeeper.GetSupply(ctx, "uqstars").Amount

	handler := upgrades.V0101002UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	sz, ok := app.InterchainstakingKeeper.GetZone(ctx, "stargaze-1")
	s.True(ok)
	s.True(sz.IsOffboarding)
	s.Equal(supplyBefore, app.BankKeeper.GetSupply(ctx, "uqstars").Amount)
}

// Zone not found: handler silently skips (matches original v1.10.2 semantics).
func (s *AppTestSuite) TestV0101002UpgradeHandler_ZoneNotFound() {
	s.SetupTest()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	handler := upgrades.V0101002UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)
}

// Malformed BurnAmount: a legacy WDR with zero BurnAmount must be skipped
// (not panic, not counted toward the mint obligation) and must remain in
// state for manual follow-up.
func (s *AppTestSuite) TestV0101002UpgradeHandler_MintToCover_SkipsMalformedWDR() {
	s.SetupTest()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	u1, u2, _ := s.setupV0101002Zones(ctx, app)

	// Valid record.
	s.seedWDR(ctx, app, "stargaze-1", u1, "uqstars", "ustars", 2_000_000, icstypes.WithdrawStatusQueued, 600)

	// Bypass validation to inject a malformed (zero BurnAmount) record.
	s.UncheckedSetWithdrawalRecord(ctx, app, icstypes.WithdrawalRecord{
		ChainId:     "stargaze-1",
		Delegator:   u2,
		Recipient:   addressutils.GenerateAddressForTestWithPrefix("stars"),
		BurnAmount:  sdk.Coin{Denom: "uqstars", Amount: math.ZeroInt()},
		Amount:      sdk.NewCoins(sdk.NewCoin("ustars", math.NewInt(0))),
		Txhash:      fmt.Sprintf("%064d", 601),
		Status:      icstypes.WithdrawStatusUnbond,
		EpochNumber: 1,
	})

	// Escrow short by 500_000 relative to the VALID obligation (2_000_000).
	s.fundEscrow(ctx, app, sdk.NewCoins(sdk.NewCoin("uqstars", math.NewInt(1_500_000))))

	supplyBefore := app.BankKeeper.GetSupply(ctx, "uqstars").Amount

	handler := upgrades.V0101002UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)
	s.NotPanics(func() {
		_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
		s.NoError(err)
	})

	// Mint should target only the valid obligation (2M − 1.5M = 500K), NOT the
	// malformed record.
	s.Equal(supplyBefore.Add(math.NewInt(500_000)), app.BankKeeper.GetSupply(ctx, "uqstars").Amount,
		"mint amount must ignore malformed BurnAmount")

	// Valid user refunded in full.
	s.Equal(math.NewInt(2_000_000), app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(u1, ""), "uqstars").Amount)
	// Malformed record kept in state.
	remaining := app.InterchainstakingKeeper.AllZoneWithdrawalRecords(ctx, "stargaze-1")
	s.Equal(1, len(remaining))
	s.Equal(u2, remaining[0].Delegator)
	s.True(remaining[0].BurnAmount.Amount.IsZero())
}
