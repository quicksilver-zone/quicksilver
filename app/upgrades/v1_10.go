package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

// V0101000UpgradeHandler handles the v1.10.0 upgrade.
// This upgrade introduces zone offboarding functionality:
// - MsgGovSetZoneOffboarding: Enable/disable offboarding mode for a zone
// - MsgGovCancelAllPendingRedemptions: Cancel all pending redemptions and refund qAssets
// - MsgGovForceUnbondAllDelegations: Force unbond all delegations via ICA
func V0101000UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting v1.10.0 upgrade...")
		ctx.Logger().Info("This upgrade introduces zone offboarding functionality")

		// No state migrations required for this upgrade.
		// The new governance messages are automatically available after the upgrade.

		ctx.Logger().Info("Upgrade v1.10.0 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
