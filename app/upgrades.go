package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const (
	ProductionChainID = "quicksilver-2"
	InnuendoChainID   = "innuendo-5"
	DevnetChainID     = "quicktest-1"

	v010300UpgradeName = "v1.3.0"
)

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v010300UpgradeName, v010300UpgradeHandler(app))

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

func v010300UpgradeHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
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
