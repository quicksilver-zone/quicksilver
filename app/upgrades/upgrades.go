package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

func Upgrades() []Upgrade {
	return []Upgrade{
		{UpgradeName: V010700UpgradeName, CreateUpgradeHandler: V010700UpgradeHandler},
		{UpgradeName: V010702UpgradeName, CreateUpgradeHandler: V010702UpgradeHandler},
		{UpgradeName: V010704UpgradeName, CreateUpgradeHandler: V010704UpgradeHandler},
		{UpgradeName: V010705UpgradeName, CreateUpgradeHandler: V010705UpgradeHandler},
		{UpgradeName: V010706UpgradeName, CreateUpgradeHandler: V010706UpgradeHandler},
		{UpgradeName: V010707UpgradeName, CreateUpgradeHandler: NoOpHandler},

		{UpgradeName: V010800r1UpgradeName, CreateUpgradeHandler: NoOpHandler},

		{UpgradeName: V010800UpgradeName, CreateUpgradeHandler: V010800UpgradeHandler},
		{UpgradeName: V010801UpgradeName, CreateUpgradeHandler: NoOpHandler},

		{UpgradeName: V010900UpgradeName, CreateUpgradeHandler: NoOpHandler},

		// v1.10.0 - Zone Offboarding
		{UpgradeName: V0101000rc0UpgradeName, CreateUpgradeHandler: V0101000UpgradeHandler},
		{UpgradeName: V0101000UpgradeName, CreateUpgradeHandler: V0101000UpgradeHandler},
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
