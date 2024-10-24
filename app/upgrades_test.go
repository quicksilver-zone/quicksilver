package app

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/quicksilver-zone/quicksilver/app/upgrades"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/suite"
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

func (s *AppTestSuite) TestV010601UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	ctx := s.chainA.GetContext()

	handler := upgrades.V010601UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	junoZone, found := app.InterchainstakingKeeper.GetZone(ctx, "juno-1")
	s.True(found)
	s.Equal("juno-1", junoZone.ChainId)
	s.Equal("channel-86", junoZone.TransferChannel)

	atomZone, found := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.True(found)
	s.Equal("cosmoshub-4", atomZone.ChainId)
	s.Equal("channel-1", atomZone.TransferChannel)

	osmoZone, found := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
	s.True(found)
	s.Equal("osmosis-1", osmoZone.ChainId)
	s.Equal("channel-2", osmoZone.TransferChannel)

	agoricZone, found := app.InterchainstakingKeeper.GetZone(ctx, "agoric-3")
	s.True(found)
	s.Equal("agoric-3", agoricZone.ChainId)
	s.Equal("channel-125", agoricZone.TransferChannel)
	s.False(agoricZone.Is_118)

	// check block params
	consensusParams := app.GetConsensusParams(ctx)
	s.Equal(int64(2072576), consensusParams.Block.MaxBytes)
	s.Equal(int64(150000000), consensusParams.Block.MaxGas)
}

func (s *AppTestSuite) TestV010603UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	ctx := s.chainA.GetContext()

	record := icstypes.WithdrawalRecord{
		ChainId:    "cosmoshub-4",
		Delegator:  "quick1zyj57u72nwr23q2glz77jaana9kpvn8cxdp5gl",
		Recipient:  "cosmos1xnvuycukuex5eae336u7umrhfea9xndr0ksjlj",
		Txhash:     "ea0d86a3fb4b25fcb13a587e72542f99ebf8c7c3aa255a0922dfa7002a8ee861",
		Status:     icstypes.WithdrawStatusUnbond,
		BurnAmount: sdk.NewCoin("uqatom", math.NewInt(4754000000)),
	}

	s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record))

	s.NoError(app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(record.BurnAmount)))
	s.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount)))

	handler := upgrades.V010603UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)

	preBalance := app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(record.Delegator, ""), "uqatom")

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	// check this hash no longer exists.
	for _, status := range []int32{icstypes.WithdrawStatusQueued, icstypes.WithdrawStatusUnbond, icstypes.WithdrawStatusSend, icstypes.WithdrawStatusCompleted} {
		_, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, record.ChainId, record.Txhash, status)
		s.False(found)
	}

	// check balance of uatom increased by 4754000000
	postBalance := app.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(record.Delegator, ""), "uqatom")

	s.Equal(postBalance.Amount.Int64(), preBalance.Add(record.BurnAmount).Amount.Int64())
}

func (s *AppTestSuite) TestV010604UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	ctx := s.chainA.GetContext()

	hashes := []struct {
		Zone string
		Hash string
	}{
		{Zone: "cosmoshub-4", Hash: "6cc942b42150a43b45d56c39d05155206ffb40eb18268dbd0b3c1ce5248b2645"},
		{Zone: "stargaze-1", Hash: "10af0ee10a97f01467039a69cbfb8df05dc3111c975d955ca51adda201f36555"},
		{Zone: "juno-1", Hash: "627db4f106a8ef99053a0726f3f71d2f23bbfd4a9155b6d083ff7015bdfa44c0"},
		{Zone: "cosmoshub-4", Hash: "0c8269f04109a55a152d3cdfd22937b4e5c2746111d579935eef4cd7ffa71f7f"},
		{Zone: "cosmoshub-4", Hash: "677691e596338af42387cbafae9831c5e0fe4b7f31b683ad69d2cc3f17687bd8"},
		{Zone: "cosmoshub-4", Hash: "c8351fe7e6775b39b9f480182f9ea57c914ea566dd35912a4597f234b12405a6"},
		{Zone: "cosmoshub-4", Hash: "d750de16665edbfca2a889ccec7a16ce107987416a80304154453ff6e8e25c5d"},
		{Zone: "cosmoshub-4", Hash: "e5a44be995514d10cce7795a28d8a997e4eb95ba805d54cfaa9ce62e78a87a50"},
		{Zone: "cosmoshub-4", Hash: "fb73556a38faeffa4740923c585b609a002869d1a9006f660567166cd4f5a79b"},
		{Zone: "cosmoshub-4", Hash: "fd42b32563d8beecb64ae2aa47f9b38ddecd436ac4e8b84bf9d9c46f447439e6"},
	}

	for _, hashRecord := range hashes {
		record := icstypes.WithdrawalRecord{
			ChainId:        hashRecord.Zone,
			Delegator:      "quick1zyj57u72nwr23q2glz77jaana9kpvn8cxdp5gl",
			Recipient:      "cosmos1xnvuycukuex5eae336u7umrhfea9xndr0ksjlj",
			Txhash:         hashRecord.Hash,
			Status:         icstypes.WithdrawStatusSend,
			BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(4754000000)),
			Acknowledged:   true,
			CompletionTime: time.Now(),
		}

		s.NoError(app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record))
	}

	handler := upgrades.V010604UpgradeHandler(app.mm, app.configurator, &app.AppKeepers)

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	for _, hashRecord := range hashes {
		record, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, hashRecord.Zone, hashRecord.Hash, icstypes.WithdrawStatusQueued)
		s.True(found)
		s.False(record.Acknowledged)
		s.True(record.CompletionTime.IsZero())
	}
}
