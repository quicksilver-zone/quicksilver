package keeper_test

import (
	"testing"
	"time"

	"github.com/quicksilver-zone/quicksilver/app"
	cmdcfg "github.com/quicksilver-zone/quicksilver/cmd/config"
	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/keeper"
	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

var (
	SecondaryDenom  = "ura"
	SecondaryAmount = sdk.NewInt(100000000)
)

type KeeperTestSuite struct {
	suite.Suite

	App *app.Quicksilver
	Ctx sdk.Context

	queryClient types.QueryClient
	msgServer   types.MsgServer
	QueryHelper *baseapp.QueryServiceTestHelper
	// defaultDenom is on the suite, as it depends on the creator test address.
	defaultDenom string
	TestAccs     []sdk.AccAddress
}

// Setup sets up basic environment for suite (App, Ctx, and test accounts).
func (s *KeeperTestSuite) Setup() {
	cmdcfg.SetBech32Prefixes(sdk.GetConfig())
	s.App = app.Setup(s.T(), false)
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "quick-1", Time: time.Now().UTC()})
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = CreateRandomAccounts(3)
}

// CreateRandomAccounts is a function return a list of randomly generated AccAddresses.
func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// FundAcc funds target address with specified amount.
func (s *KeeperTestSuite) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, amounts)
	s.NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, acc, amounts)
	s.NoError(err)
}

func (s *KeeperTestSuite) SetupTestForInitGenesis() {
	// Setting to True, leads to init genesis not running
	s.App = app.Setup(s.T(), true)
	s.Ctx = s.App.BaseApp.NewContext(true, tmproto.Header{})
}

// AssertEventEmitted asserts that ctx's event manager has emitted the given number of events
// of the given type.
func (s *KeeperTestSuite) AssertEventEmitted(ctx sdk.Context, eventTypeExpected string, numEventsExpected int) {
	allEvents := ctx.EventManager().Events()
	// filter out other events
	actualEvents := make([]sdk.Event, 0)
	for _, event := range allEvents {
		if event.Type == eventTypeExpected {
			actualEvents = append(actualEvents, event)
		}
	}
	s.Equal(numEventsExpected, len(actualEvents))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.Setup()
	// Fund every TestAcc with two denoms, one of which is the denom creation fee
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)), sdk.NewCoin(SecondaryDenom, SecondaryAmount))
	for _, acc := range s.TestAccs {
		s.FundAcc(acc, fundAccsAmount)
	}

	s.queryClient = types.NewQueryClient(s.QueryHelper)
	s.msgServer = keeper.NewMsgServerImpl(s.App.TokenFactoryKeeper)
}

func (s *KeeperTestSuite) CreateDefaultDenom() {
	s.T().Helper()

	res, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
	s.Require().NoError(err)
	s.defaultDenom = res.GetNewTokenDenom()
}

func (s *KeeperTestSuite) TestCreateModuleAccount() {
	quicksilver := s.App

	// remove module account
	tokenfactoryModuleAccount := quicksilver.AccountKeeper.GetAccount(s.Ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
	quicksilver.AccountKeeper.RemoveAccount(s.Ctx, tokenfactoryModuleAccount)

	// ensure module account was removed
	s.Ctx = quicksilver.BaseApp.NewContext(false, tmproto.Header{})
	tokenfactoryModuleAccount = quicksilver.AccountKeeper.GetAccount(s.Ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	// create module account
	quicksilver.TokenFactoryKeeper.CreateModuleAccount(s.Ctx)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = quicksilver.AccountKeeper.GetAccount(s.Ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)
}
