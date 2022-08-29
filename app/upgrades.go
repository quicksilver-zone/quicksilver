package app

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

func GetInnuendo1Upgrade(app *Quicksilver) types.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ReplaceZoneDropChain(ctx, app, "osmotestnet-4", "osmo-test-4", ctx.BlockHeader().Time)
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

// replaces zonedrop and claimrecords for a given chain, with another chain and update start time.
func ReplaceZoneDropChain(ctx sdk.Context, app *Quicksilver, chainIdFrom string, chainIdTo string, start time.Time) {
	ad, found := app.AirdropKeeper.GetZoneDrop(ctx, chainIdFrom)
	if !found {
		panic(chainIdFrom + " zonedrop not found")
	}
	// update chainid for chainIdFrom airdrop and reset start time.
	ad.ChainId = chainIdTo
	ad.StartTime = start

	app.AirdropKeeper.SetZoneDrop(ctx, ad)
	app.AirdropKeeper.IterateClaimRecords(ctx, chainIdFrom, func(index int64, cr airdroptypes.ClaimRecord) (stop bool) {
		ctx.Logger().Info("migrating claimdrop record", "address", cr.Address)
		cr.ChainId = chainIdTo
		app.AirdropKeeper.SetClaimRecord(ctx, cr)
		app.AirdropKeeper.DeleteClaimRecord(ctx, chainIdFrom, cr.Address)
		return false
	})

	zonedropOldAddress := app.AirdropKeeper.GetZoneDropAccountAddress(chainIdFrom)
	zonedropNewAddress := app.AirdropKeeper.GetZoneDropAccountAddress(chainIdTo)

	coinsToMove := sdk.NewCoins(
		sdk.NewCoin(
			app.AirdropKeeper.BondDenom(ctx),
			sdk.NewIntFromUint64(ad.Allocation),
		),
	)

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

	app.AirdropKeeper.DeleteZoneDrop(ctx, chainIdFrom)

	// update unbonding time to 48h.
	stakeParams := app.StakingKeeper.GetParams(ctx)
	stakeParams.UnbondingTime = 48 * time.Hour
	app.StakingKeeper.SetParams(ctx, stakeParams)
}
