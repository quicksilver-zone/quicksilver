package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const (
	ProductionChainID = "quicksilver-2"
	InnuendoChainID   = "innuendo-5"
	DevnetChainID     = "quicktest-1"

	v010204UpgradeName = "v1.2.4"
	v010207UpgradeName = "v1.2.7"
	v010208UpgradeName = "v1.2.8"
	v010300UpgradeName = "v1.3.0" // retained for testy
)

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v010300UpgradeName, noOpUpgradeHandler(app)) // retained for testy
	app.UpgradeKeeper.SetUpgradeHandler(v010204UpgradeName, v010204UpgradeHandler(app))
	app.UpgradeKeeper.SetUpgradeHandler(v010207UpgradeName, v010207UpgradeHandler(app))
	app.UpgradeKeeper.SetUpgradeHandler(v010208UpgradeName, v010208UpgradeHandler(app))

	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name { //nolint:gocritic
	// case v001000UpgradeName:

	// 	storeUpgrades = &storetypes.StoreUpgrades{
	// 		Added: []string{claimsmanagertypes.ModuleName},
	// 	}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func noOpUpgradeHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func v010204UpgradeHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// upgrade receipts
		time := ctx.BlockTime()
		for _, r := range app.InterchainstakingKeeper.AllReceipts(ctx) {
			r.FirstSeen = &time
			r.Completed = &time
			app.InterchainstakingKeeper.SetReceipt(ctx, r)
		}

		// remove failed redelegation records
		for _, r := range app.InterchainstakingKeeper.AllRedelegationRecords(ctx) {
			if r.CompletionTime.IsZero() {
				app.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, r.ChainId, r.Source, r.Destination, r.EpochNumber)
			}
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func v010207UpgradeHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// update minter epoch-provisions
		minter := app.MintKeeper.GetMinter(ctx)
		minter.EpochProvisions = sdk.NewDec(50_000_000_000_000).Quo(sdk.NewDec(365))
		app.MintKeeper.SetMinter(ctx, minter)

		// update params
		params := app.MintKeeper.GetParams(ctx)
		params.DistributionProportions = minttypes.DistributionProportions{
			Staking:              sdk.NewDecWithPrec(80, 2),
			PoolIncentives:       sdk.NewDecWithPrec(17, 2),
			ParticipationRewards: sdk.NewDec(0),
			CommunityPool:        sdk.NewDecWithPrec(3, 2),
		}
		app.MintKeeper.SetParams(ctx, params)
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func v010208UpgradeHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// set messages per tx a maximum of (max_gas per tx/block divided by 1m); be conservative for now and we can tweak later.

		// cosmos
		zone, found := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
		if !found {
			panic("unable to find zone cosmoshub-4")
		}
		// block gas is 40m
		zone.MessagesPerTx = 30
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		// stargaze
		zone, found = app.InterchainstakingKeeper.GetZone(ctx, "stargaze-1")
		if !found {
			panic("unable to find zone stargaze-1")
		}
		// block gas is 75m
		zone.MessagesPerTx = 50
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		// osmosis
		zone, found = app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
		if !found {
			panic("unable to find zone osmosis-1")
		}
		// block gas is 120m, but gas per tx is 25m
		zone.MessagesPerTx = 20
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		// regen
		zone, found = app.InterchainstakingKeeper.GetZone(ctx, "regen-1")
		if !found {
			panic("unable to find zone regen-1")
		}
		// block gas is 2m
		zone.MessagesPerTx = 2
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		if _, err := app.InterchainstakingKeeper.FlushOutstandingDelegations(ctx, &zone); err != nil {
			panic("unable to emit delegation messages")
		}
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
