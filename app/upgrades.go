package app

import (
	"fmt"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	v6 "github.com/cosmos/ibc-go/v6/testing/simapp/upgrades/v6"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v6/packetforward/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	v6 "github.com/cosmos/ibc-go/v6/testing/simapp/upgrades/v6" // nolint:revive

	"github.com/quicksilver-zone/quicksilver/app/upgrades"
	supplytypes "github.com/quicksilver-zone/quicksilver/x/supply/types"
)

const wasmModuleName = "wasm"

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

	kvStoreKeys := app.GetKVStoreKey()
	app.UpgradeKeeper.SetUpgradeHandler(
		upgrades.V010600rc1UpgradeName,
		v6.CreateUpgradeHandler(
			app.mm,
			app.configurator,
			app.appCodec,
			kvStoreKeys[capabilitytypes.ModuleName],
			app.CapabilityKeeper,
			icacontrollertypes.SubModuleName,
		),
	)
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

	case upgrades.V010405UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{packetforwardtypes.ModuleName, supplytypes.ModuleName},
		}
	case upgrades.V010600rc1UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Deleted: []string{wasmModuleName},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
