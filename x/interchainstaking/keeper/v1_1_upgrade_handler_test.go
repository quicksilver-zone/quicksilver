package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v5/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	tmclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) Test_v11UpgradeHandler() {
	qApp := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	qApp.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", &tmclienttypes.ClientState{ChainId: "cosmoshub-4", TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
	qApp.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: ctx.BlockTime()})
	qApp.IBCKeeper.ConnectionKeeper.SetConnection(ctx, "connection-4", connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0", Versions: []*connectiontypes.Version{{Identifier: "1", Features: []string{"ORDER_ORDERED", "ORDER_UNORDERED"}}}})

	proposal := &icstypes.RegisterZoneProposal{
		Title:           "register zone A",
		Description:     "register zone A",
		ConnectionId:    "connection-4",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		AccountPrefix:   "cosmos",
		MultiSend:       true,
		LiquidityModule: true,
	}

	err := icskeeper.HandleRegisterZoneProposal(ctx, qApp.InterchainstakingKeeper, proposal)
	suite.Require().NoError(err)

	zone, found := qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), "cosmoshub-4")
	suite.Require().True(found)

	suite.Require().NoError(suite.setupChannelForICA(ctx, "cosmoshub-4", "connection-4", "deposit", zone.AccountPrefix))
	suite.Require().NoError(suite.setupChannelForICA(ctx, "cosmoshub-4", "connection-4", "withdrawal", zone.AccountPrefix))
	suite.Require().NoError(suite.setupChannelForICA(ctx, "cosmoshub-4", "connection-4", "performance", zone.AccountPrefix))
	suite.Require().NoError(suite.setupChannelForICA(ctx, "cosmoshub-4", "connection-4", "delegate", zone.AccountPrefix))

	for _, val := range suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext()) {
		// refetch the zone for each validator, else we end up with an empty valset each time!
		zone, found := qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), "cosmoshub-4")
		suite.Require().True(found)
		suite.Require().NoError(icskeeper.SetValidatorForZone(&qApp.InterchainstakingKeeper, suite.chainA.GetContext(), zone, app.DefaultConfig().Codec.MustMarshal(&val)))
	}

	qApp.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewInt64Coin("uqatom", 3460036988)))
	for addr, amount := range map[string]int64{
		// depositor accounts; already refunded.
		"quick1dpr6x7ggn32u2fgewtv08yvdrlh76mqwz4vva3": 4000000,
		"quick1t9t4ch0gvt2wtyejnzd03gtdawlhhut56ymqzh": 3000000,
		"quick133df0rz42d49jl3hzh6gtv6gf42x8vqas37vn3": 20000,
		"quick1afj8d8e6tfdnj4hwczjugl5hu0frfamxsjh030": 750000,
		"quick1kx3xxcpfaegq8nw0e46la339lgfmqahcph4nql": 2540000,
		"quick1cu89smpf8cmlctxz5y5eud6qhpsp308ckgs0fu": 250000000,
		"quick1nnhlu7j6r4e2efuaydyvdmqs5plp963es0ukex": 1000000,
		"quick1rdhss0e720yjqu28xxen30t7u55selqfzwdnf5": 17100000,
		"quick14hew5e5ua5pzhr6swnr0t2md6up7qmgpy8fe06": 10000000,
		"quick1k00n9wvapdkwcmucct8f5wwnrmytqqrwkjkrmk": 10000000,
		"quick1qgme8vlq4ly8tcye6xdgxnz4khzq9es03n2yzy": 10000000,
		"quick1n4c56vddeqg67ktukprkteqdmpph2ck0xtm8sw": 2000000,
		"quick1k8g0vlfmctyqtwahrxhudksz7rgrm6nsye3kpp": 66988,
		"quick1gjgjsh74w5trmfaaauq4qt2mwvh8p7gsp4rc5v": 760000,
		"quick1eslm0n3ypkd6f67n6ymf08m2t6kt9cypsf2w7d": 20000000,
		// quicksilver accounts; to refund after gaia returns funds.
		"quick1780znw95jjcdk4wtac44t5cjtrmxqfe9q49pej": 295000000,
		"quick1954q9apawr6kg8ez4ukx8jyuaxakz7ye4jvyk4": 800000,
		"quick16x03wcp37kx5e8ehckjxvwcgk9j0cqnhcccnty": 2833000000,
	} {

		addrBytes, _ := utils.AccAddressFromBech32(addr, "quick")
		if err := qApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", addrBytes, sdk.NewCoins(sdk.NewInt64Coin("uqatom", amount))); err != nil {
			panic(err)
		}
	}
	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)

	keeper.V010100UpgradeHandler(ctx, &qApp.InterchainstakingKeeper)

	for _, ica := range qApp.ICAControllerKeeper.GetAllInterchainAccounts(ctx) {
		channelId, found := qApp.ICAControllerKeeper.GetActiveChannelID(ctx, ica.ConnectionId, ica.PortId)
		suite.Require().True(found)
		channel, found := qApp.IBCKeeper.ChannelKeeper.GetChannel(ctx, ica.PortId, channelId)
		suite.Require().True(found)
		suite.Require().Equal(ibcchanneltypes.CLOSED, channel.State)
	}

	suite.Require().Equal(sdk.Coin{Denom: "uqatom", Amount: sdk.ZeroInt()}, qApp.BankKeeper.GetSupply(ctx, "uqatom"))
	_, found = qApp.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	suite.Require().False(found)
}
