package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/quicksilver-zone/quicksilver/app"
	cmdcfg "github.com/quicksilver-zone/quicksilver/cmd/config"
	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/keeper"
	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/types"
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
func (suite *KeeperTestSuite) Setup() {
	cmdcfg.SetBech32Prefixes(sdk.GetConfig())
	suite.App = app.Setup(suite.T(), false)
	suite.Ctx = suite.App.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "quick-1", Time: time.Now().UTC()})
	suite.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: suite.App.GRPCQueryRouter(),
		Ctx:             suite.Ctx,
	}

	suite.TestAccs = CreateRandomAccounts(3)
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
func (suite *KeeperTestSuite) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := suite.App.BankKeeper.MintCoins(suite.Ctx, minttypes.ModuleName, amounts)
	suite.NoError(err)
	err = suite.App.BankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, minttypes.ModuleName, acc, amounts)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) SetupTestForInitGenesis() {
	// Setting to True, leads to init genesis not running
	suite.App = app.Setup(suite.T(), true)
	suite.Ctx = suite.App.BaseApp.NewContext(true, tmproto.Header{})
}

// AssertEventEmitted asserts that ctx's event manager has emitted the given number of events
// of the given type.
func (suite *KeeperTestSuite) AssertEventEmitted(ctx sdk.Context, eventTypeExpected string, numEventsExpected int) {
	allEvents := ctx.EventManager().Events()
	// filter out other events
	actualEvents := make([]sdk.Event, 0)
	for _, event := range allEvents {
		if event.Type == eventTypeExpected {
			actualEvents = append(actualEvents, event)
		}
	}
	suite.Equal(numEventsExpected, len(actualEvents))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()
	// Fund every TestAcc with two denoms, one of which is the denom creation fee
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)), sdk.NewCoin(SecondaryDenom, SecondaryAmount))
	for _, acc := range suite.TestAccs {
		suite.FundAcc(acc, fundAccsAmount)
	}

	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.App.TokenFactoryKeeper)
}

func (suite *KeeperTestSuite) CreateDefaultDenom() {
	suite.T().Helper()

	res, err := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.Ctx), types.NewMsgCreateDenom(suite.TestAccs[0].String(), "bitcoin"))
	suite.Require().NoError(err)
	suite.defaultDenom = res.GetNewTokenDenom()
}

func (suite *KeeperTestSuite) TestCreateModuleAccount() {
	quicksilver := suite.App

	// remove module account
	tokenfactoryModuleAccount := quicksilver.AccountKeeper.GetAccount(suite.Ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
	quicksilver.AccountKeeper.RemoveAccount(suite.Ctx, tokenfactoryModuleAccount)

	// ensure module account was removed
	suite.Ctx = quicksilver.BaseApp.NewContext(false, tmproto.Header{})
	tokenfactoryModuleAccount = quicksilver.AccountKeeper.GetAccount(suite.Ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().Nil(tokenfactoryModuleAccount)

	// create module account
	quicksilver.TokenFactoryKeeper.CreateModuleAccount(suite.Ctx)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = quicksilver.AccountKeeper.GetAccount(suite.Ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().NotNil(tokenfactoryModuleAccount)
}
