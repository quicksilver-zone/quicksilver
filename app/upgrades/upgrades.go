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
		{UpgradeName: V010500rc0UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010500rc1UpgradeName, CreateUpgradeHandler: V010500rc1UpgradeHandler},
		{UpgradeName: V010503rc0UpgradeName, CreateUpgradeHandler: V010503rc0UpgradeHandler},
		{UpgradeName: V010600beta0UpgradeName, CreateUpgradeHandler: V010600beta0UpgradeHandler},
		{UpgradeName: V010600beta1UpgradeName, CreateUpgradeHandler: V010600beta1UpgradeHandler},
		{UpgradeName: V010600rc0UpgradeName, CreateUpgradeHandler: V010600rc0UpgradeHandler},
		{UpgradeName: V010601rc0UpgradeName, CreateUpgradeHandler: V010601rc0UpgradeHandler},
		{UpgradeName: V010601rc2UpgradeName, CreateUpgradeHandler: V010601rc0UpgradeHandler}, // this name mismatch is intentional, as we want to rerun the upgrade after resolving some issues.
		{UpgradeName: V010601rc3UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010601rc4UpgradeName, CreateUpgradeHandler: NoOpHandler},

		// v1.5: this needs to be present to support upgrade on mainnet
		{UpgradeName: V010500UpgradeName, CreateUpgradeHandler: V010500UpgradeHandler},
		{UpgradeName: V010501UpgradeName, CreateUpgradeHandler: V010501UpgradeHandler},
		{UpgradeName: V010503UpgradeName, CreateUpgradeHandler: V010503UpgradeHandler},
		{UpgradeName: V010504UpgradeName, CreateUpgradeHandler: V010504UpgradeHandler},
		{UpgradeName: V010505UpgradeName, CreateUpgradeHandler: V010505UpgradeHandler},
		{UpgradeName: V010601UpgradeName, CreateUpgradeHandler: V010601UpgradeHandler},
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
