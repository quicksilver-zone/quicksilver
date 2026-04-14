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

		// v1.10.1 - HandleFailedUndelegate BurnAmount zero guard hotfix
		{UpgradeName: V0101001UpgradeName, CreateUpgradeHandler: NoOpHandler},

		// v1.10.2 - Sunset stargaze & omniflix + HandleFailedUndelegate fix
		{UpgradeName: V0101002UpgradeName, CreateUpgradeHandler: V0101002UpgradeHandler},

		// v1.10.3 - Zone sunset (re-attempt of failed v1.10.2) with mint-to-cover
		{UpgradeName: V0101003UpgradeName, CreateUpgradeHandler: V0101003UpgradeHandler},
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
