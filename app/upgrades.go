package app

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

func GetInnuendo1Upgrade(app *Quicksilver) types.UpgradeHandler {
	return func(ctx sdk.Context, _ types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ReplaceZoneDropChain(ctx, app, "osmotestnet-4", "osmo-test-4", ctx.BlockHeader().Time)

		// update unbonding time to 48h for innuendo-1 testnet to avoid ibc client expiry.
		// this only applies to innuendo-1; production value will be 21 days.
		stakeParams := app.StakingKeeper.GetParams(ctx)
		stakeParams.UnbondingTime = 48 * time.Hour
		app.StakingKeeper.SetParams(ctx, stakeParams)

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

// replaces zonedrop and claimrecords for a given chain, with another chain and update start time.
// this function will panic if zonedrop for the given chainId is not found, or claim records fail to be set or deleted as expected.
func ReplaceZoneDropChain(ctx sdk.Context, app *Quicksilver, chainIDFrom string, chainIDTo string, start time.Time) {
	ad, found := app.AirdropKeeper.GetZoneDrop(ctx, chainIDFrom)
	if !found {
		panic(chainIDFrom + " zonedrop not found")
	}
	// update chainid for chainIdFrom airdrop and reset start time.
	ad.ChainId = chainIDTo
	ad.StartTime = start

	app.AirdropKeeper.SetZoneDrop(ctx, ad)
	app.AirdropKeeper.IterateClaimRecords(ctx, chainIDFrom, func(index int64, cr airdroptypes.ClaimRecord) (stop bool) {
		ctx.Logger().Info("migrating claimdrop record", "address", cr.Address)
		cr.ChainId = chainIDTo
		err := app.AirdropKeeper.SetClaimRecord(ctx, cr)
		if err != nil {
			panic(err)
		}
		err = app.AirdropKeeper.DeleteClaimRecord(ctx, chainIDFrom, cr.Address)
		if err != nil {
			panic(err)
		}
		return false
	})

	zonedropOldAddress := app.AirdropKeeper.GetZoneDropAccountAddress(chainIDFrom)
	zonedropNewAddress := app.AirdropKeeper.GetZoneDropAccountAddress(chainIDTo)

	coinsToMove := sdk.NewCoins(
		sdk.NewCoin(
			app.AirdropKeeper.BondDenom(ctx),
			sdk.NewIntFromUint64(ad.Allocation),
		),
	)

	ctx.Logger().Info("migrating zonedrop bounty", "from", zonedropOldAddress, "to", zonedropNewAddress, "coins", coinsToMove)

	// migrate coins from old chain account to the new one - via the airdrop module.
	if err := app.BankKeeper.SendCoinsFromAccountToModule(
		ctx, zonedropOldAddress, airdroptypes.ModuleName, coinsToMove,
	); err != nil {
		panic(err)
	}

	if err := app.AirdropKeeper.SendCoinsFromModuleToAccount(
		ctx,
		airdroptypes.ModuleName, zonedropNewAddress, coinsToMove,
	); err != nil {
		panic(err)
	}

	app.AirdropKeeper.DeleteZoneDrop(ctx, chainIDFrom)
}
