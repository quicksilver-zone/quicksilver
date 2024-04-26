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
		{UpgradeName: V010405rc6UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010405rc7UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010407rc0UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010407rc1UpgradeName, CreateUpgradeHandler: V010407rc1UpgradeHandler},
		{UpgradeName: V010407rc2UpgradeName, CreateUpgradeHandler: V010407rc2UpgradeHandler},
		{UpgradeName: V010500rc0UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010500rc1UpgradeName, CreateUpgradeHandler: V010500rc1UpgradeHandler},
		{UpgradeName: V010503rc0UpgradeName, CreateUpgradeHandler: V010503rc0UpgradeHandler},
		{UpgradeName: V010600beta0UpgradeName, CreateUpgradeHandler: V010600beta0UpgradeHandler},
		{UpgradeName: V010600rc0UpgradeName, CreateUpgradeHandler: V010600rc0UpgradeHandler},

		// v1.2: this needs to be present to support upgrade on mainnet
		{UpgradeName: V010217UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010405UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010406UpgradeName, CreateUpgradeHandler: V010406UpgradeHandler},
		{UpgradeName: V010407UpgradeName, CreateUpgradeHandler: V010407UpgradeHandler},
		{UpgradeName: V010500UpgradeName, CreateUpgradeHandler: V010500UpgradeHandler},
		{UpgradeName: V010501UpgradeName, CreateUpgradeHandler: V010501UpgradeHandler},
		{UpgradeName: V010503UpgradeName, CreateUpgradeHandler: V010503UpgradeHandler},
		{UpgradeName: V010504UpgradeName, CreateUpgradeHandler: V010504UpgradeHandler},
		{UpgradeName: V010505UpgradeName, CreateUpgradeHandler: V010505UpgradeHandler},
		{UpgradeName: V010600UpgradeName, CreateUpgradeHandler: V010600UpgradeHandler},
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
