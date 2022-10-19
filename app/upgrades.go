package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const v001000UpgradeName = "v0.10.0"

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v001000UpgradeName, getv001000Upgrade(app))
}

func getv001000Upgrade(app *Quicksilver) types.UpgradeHandler {
	return func(ctx sdk.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		app.UpgradeKeeper.Logger(ctx).Info("UPGRADE: to v0.10.0; upgrading store.")
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(plan.Height, &storetypes.StoreUpgrades{
			Added: []string{cmtypes.ModuleName},
		}))
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
