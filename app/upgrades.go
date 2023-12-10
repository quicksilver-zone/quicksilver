package app

import (
	"fmt"

	packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v5/router/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/upgrades"
)

func (app *Quicksilver) setUpgradeHandlers() {
	for _, upgrade := range upgrades.Upgrades() {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.mm,
				app.configurator,
				&app.AppKeepers,
			),
		)
	}
}

func (app *Quicksilver) setUpgradeStoreLoaders() {
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

	switch upgradeInfo.Name { // nolint:gocritic
	// case v001000UpgradeName:

	// 	storeUpgrades = &storetypes.StoreUpgrades{
	// 		Added: []string{claimsmanagertypes.ModuleName},
	// 	}

	case upgrades.V010402rc3UpgradeName:

		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{packetforwardtypes.ModuleName},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
