package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibcchanneltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	"github.com/ingenuity-build/quicksilver/utils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func closeChannel(ctx sdk.Context, k *Keeper, connectionID string, portID string) {
	channelID, found := k.ICAControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		panic("unable to fetch channelID for closing")
	}
	channel, found := k.IBCKeeper.ChannelKeeper.GetChannel(ctx, portID, channelID)
	if !found {
		panic("unable to fetch channel for closing")
	}
	channel.State = ibcchanneltypes.CLOSED
	k.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, channelID, channel)
}

// this is not a conventional upgrade handler, as it needs to be run by the begin blocker of the ICS module immediately on restart.
// it needs to be exported/public so can be called from the appropriate begin blocker.
func V010100UpgradeHandler(ctx sdk.Context, k *Keeper) {

	// do NOT remove receipts, as reregistering the chain with the same connection will re-open the same deposit account, and existing txs would be re-processed.
	// no delegation records, withdrawal, unbonding or rebalancing records exist, because the delegations never happened.

	// and delete the zone.
	k.DeleteZone(ctx, "cosmoshub-4")

	// close channels 3,4,5 and 6 on connection-4.
	closeChannel(ctx, k, "connection-4", "icacontroller-cosmoshub-4.deposit")
	closeChannel(ctx, k, "connection-4", "icacontroller-cosmoshub-4.withdrawal")
	closeChannel(ctx, k, "connection-4", "icacontroller-cosmoshub-4.performance")
	closeChannel(ctx, k, "connection-4", "icacontroller-cosmoshub-4.delegate")

	// move all qAssets into interchainstakingtypes.EscrowModuleAccount, and burn. Original funds are returned to foundation address by the following
	// Gaia migration: https://github.com/cosmos/gaia/pull/1976
	balances := map[string]int64{
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
	}

	for addressStr, balance := range balances {
		addressBytes, err := utils.AccAddressFromBech32(addressStr, "quick")
		if err != nil {
			panic(err)
		}
		k.Logger(ctx).Info(fmt.Sprintf("Moving %d uqatom from %s to module account", balance, addressStr))
		if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, addressBytes, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewInt64Coin("uqatom", balance))); err != nil {
			panic(err)
		}
	}

	k.Logger(ctx).Info(fmt.Sprintf("Burning %d uqatom in module account", 3460036988))
	if err := k.BankKeeper.BurnCoins(ctx, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewInt64Coin("uqatom", 3460036988))); err != nil {
		panic(err)
	}

	k.Logger(ctx).Info(fmt.Sprintf("%d uqatom burned and removed from supply; current supply (should be zero!): %v", 3460036988, k.BankKeeper.GetSupply(ctx, "uqatom")))
}
