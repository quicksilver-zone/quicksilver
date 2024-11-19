package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

func Upgrades() []Upgrade {
	return []Upgrade{
		// testnet upgrades
		{UpgradeName: V010600beta0UpgradeName, CreateUpgradeHandler: V010600beta0UpgradeHandler},
		{UpgradeName: V010600beta1UpgradeName, CreateUpgradeHandler: V010600beta1UpgradeHandler},
		{UpgradeName: V010600rc0UpgradeName, CreateUpgradeHandler: V010600rc0UpgradeHandler},
		{UpgradeName: V010601rc0UpgradeName, CreateUpgradeHandler: V010601rc0UpgradeHandler},
		{UpgradeName: V010601rc2UpgradeName, CreateUpgradeHandler: V010601rc0UpgradeHandler}, // this name mismatch is intentional, as we want to rerun the upgrade after resolving some issues.
		{UpgradeName: V010601rc3UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010601rc4UpgradeName, CreateUpgradeHandler: NoOpHandler},

		// v1.5: this needs to be present to support upgrade on mainnet
		{UpgradeName: V010601UpgradeName, CreateUpgradeHandler: V010601UpgradeHandler},
		{UpgradeName: V010602UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010603UpgradeName, CreateUpgradeHandler: V010603UpgradeHandler},
		{UpgradeName: V010604UpgradeName, CreateUpgradeHandler: V010604UpgradeHandler},
		{UpgradeName: V010700UpgradeName, CreateUpgradeHandler: V010700UpgradeHandler},
		{UpgradeName: V010702UpgradeName, CreateUpgradeHandler: V010702UpgradeHandler},
	}
}

// NoOpHandler no-op handler for upgrades with no state manipulation.
func NoOpHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
