package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

func Upgrades() []Upgrade {
	return []Upgrade{
		// v1.10.0 - Zone Offboarding
		{UpgradeName: V0101000rc0UpgradeName, CreateUpgradeHandler: V0101000UpgradeHandler},
		{UpgradeName: V0101000UpgradeName, CreateUpgradeHandler: V0101000UpgradeHandler},

		// v1.10.1 - HandleFailedUndelegate BurnAmount zero guard hotfix
		{UpgradeName: V0101001UpgradeName, CreateUpgradeHandler: NoOpHandler},

		// v1.10.2 - Zone sunset (stargaze-1, omniflixhub-1) with mint-to-cover refund
		{UpgradeName: V0101002UpgradeName, CreateUpgradeHandler: V0101002UpgradeHandler},

		// v1.11.0 - Quicksilver sunset exploit cleanup
		{UpgradeName: V011000UpgradeName, CreateUpgradeHandler: V011000UpgradeHandler},
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
