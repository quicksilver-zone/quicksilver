package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	claimsmanagertypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const v001000UpgradeName = "v0.10.0"

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v001000UpgradeName, getv001000Upgrade(app))

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

	switch upgradeInfo.Name {
	case v001000UpgradeName:

		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{claimsmanagertypes.ModuleName},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func getv001000Upgrade(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.0; no state transitions to apply.")
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
