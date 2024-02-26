package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibctesting "github.com/cosmos/ibc-go/v5/testing"

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

func (s *AppTestSuite) InitV146TestZones() {
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

	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1e22za5qrqqp488h5p7vw2pfx8v0y4u444ufeuw", 10, 864000, 5000000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1qlckz3nplj3sf323n4ma7n75fmv60lpclq5ccc", 20, 864000, 200000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1edavtxhdfs8luyvedgkjcxjc9dtvks3ve7etku", 5, 86400, 50000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1pajjuywnj6w3y6pclp4tj55a7ngz9tp2z4pgep", 3, 31536000, 100000000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1vhd4n5u8rsmsdgs4h7zsn4h4klsej6n8spvsl3", 10, 864000, 5000000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1rufya429ss9nlhdram0xkcu0jejsz5atap0xan", 10, 864000, 5000000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1f8jp5tr86gn5yvwecr7a4a9zypqf2mg85p96rw", 10, 864000, 5000000000)
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

func (s *AppTestSuite) TestV010406UpgradeHandler() {
	s.InitV146TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	handler := upgrades.V010406UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)
	ctx := s.chainA.GetContext()

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	osmoZone, ok := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
	s.True(ok)
	s.True(osmoZone.Is_118)
	s.True(osmoZone.UnbondingEnabled)
	s.False(osmoZone.SupportLsm())

	cosmosZone, ok := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.True(ok)
	s.True(cosmosZone.Is_118)
	s.True(cosmosZone.UnbondingEnabled)
	s.True(cosmosZone.SupportLsm())

	caps, ok := app.InterchainstakingKeeper.GetLsmCaps(ctx, "cosmoshub-4")
	s.True(ok)
	s.Equal(sdk.NewDecWithPrec(25, 2), caps.GlobalCap)
	s.Equal(sdk.NewDecWithPrec(100, 2), caps.ValidatorCap)
	s.Equal(sdk.NewDecWithPrec(250, 0), caps.ValidatorBondCap)

	_, ok = app.InterchainstakingKeeper.GetLsmCaps(ctx, "juno-1")
	s.False(ok)

	// original account ought to no longer exist.
	account := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1e22za5qrqqp488h5p7vw2pfx8v0y4u444ufeuw", ""))
	s.Nil(account)

	// replacement account is PVA and all fields are good.
	account = app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1gxrks2rcj9gthzfgrkjk5lnk0g00cg0cpyntlm", ""))
	pva, ok := account.(*vestingtypes.PeriodicVestingAccount)
	s.True(ok)
	s.Equal(int64(1712880000), pva.EndTime)
	s.Equal(10, len(pva.VestingPeriods))

	ctestZone, found := app.InterchainstakingKeeper.GetZoneForAccount(ctx, cosmosZone.DepositAddress.Address)
	s.True(found)
	s.Equal(ctestZone.ChainId, cosmosZone.ChainId)

	otestZone, found := app.InterchainstakingKeeper.GetZoneForAccount(ctx, osmoZone.PerformanceAddress.Address)
	s.True(found)
	s.Equal(otestZone.ChainId, osmoZone.ChainId)

	noTestZone, found := app.InterchainstakingKeeper.GetZoneForAccount(ctx, addressutils.GenerateAddressForTestWithPrefix("cosmos"))
	s.False(found)
	s.Nil(noTestZone)
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

	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l", 10, 864000, 5000000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f", 10, 864000, 5000000000)
}

func (s *AppTestSuite) TestV010600UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()
	validators := app.StakingKeeper.GetAllValidators(ctx)
	// Setting up
	accountA := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l", ""))
	err := app.BankKeeper.SendCoins(ctx, s.chainA.SenderAccount.GetAddress(), accountA.GetAddress(), sdk.Coins{
		sdk.NewInt64Coin("stake", 100),
	})
	s.Require().NoError(err)

	accountB := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f", ""))
	err = app.BankKeeper.SendCoins(ctx, s.chainA.SenderAccount.GetAddress(), accountB.GetAddress(), sdk.Coins{
		sdk.NewInt64Coin("stake", 300),
	})
	s.Require().NoError(err)

	// Stake old account
	amountA, _ := sdk.NewIntFromString("100")
	_, err = app.StakingKeeper.Delegate(ctx, accountA.GetAddress(), amountA, stakingtypes.Unbonded, validators[0], true)
	s.Require().NoError(err)

	amountB, _ := sdk.NewIntFromString("100")
	_, err = app.StakingKeeper.Delegate(ctx, accountB.GetAddress(), amountB, stakingtypes.Unbonded, validators[0], true)
	s.Require().NoError(err)

	amountB, _ = sdk.NewIntFromString("200")
	_, err = app.StakingKeeper.Delegate(ctx, accountB.GetAddress(), amountB, stakingtypes.Unbonded, validators[1], true)
	s.Require().NoError(err)

	handler := upgrades.V010600UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)

	_, err = handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.Require().NoError(err)

	_, found := app.StakingKeeper.GetDelegation(ctx, accountA.GetAddress(), validators[0].GetOperator())
	s.Require().False(found)

	migratedA := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1h0sqndv2y4xty6uk0sv4vckgyc5aa7n5at7fll", ""))
	stakeBalanceA := app.BankKeeper.GetBalance(ctx, migratedA.GetAddress(), "stake")
	s.Require().Equal(int64(100), stakeBalanceA.Amount.Int64())

	migratedB := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1n4g6037cjm0e0v2nvwj2ngau7pk758wtwk6lwq", ""))
	stakeBalanceB := app.BankKeeper.GetBalance(ctx, migratedB.GetAddress(), "stake")
	s.Require().Equal(int64(300), stakeBalanceB.Amount.Int64())

	// Check the vest period of new account
	vestMigratedA, ok := migratedA.(*vestingtypes.PeriodicVestingAccount)
	s.Require().True(ok)
	s.Require().Equal(int64(5000000000), vestMigratedA.OriginalVesting.AmountOf("uqck").Int64())
	s.Require().Equal(float64(864000), vestMigratedA.VestingPeriods[0].Duration().Seconds())

	vestMigratedB, ok := migratedB.(*vestingtypes.PeriodicVestingAccount)
	s.Require().True(ok)
	s.Require().Equal(int64(5000000000), vestMigratedB.OriginalVesting.AmountOf("uqck").Int64())
	s.Require().Equal(float64(864000), vestMigratedB.VestingPeriods[0].Duration().Seconds())
}
