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
	// cosmos zone
	zone := icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: false,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	// osmosis zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       false,
		LiquidityModule: false,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	// uni-5 zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77003",
		ChainId:         "juno-1",
		AccountPrefix:   "juno",
		LocalDenom:      "uqjuno",
		BaseDenom:       "ujuno",
		MultiSend:       false,
		LiquidityModule: false,
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
}
