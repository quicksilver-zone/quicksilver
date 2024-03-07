package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

// =========== PRODUCTION UPGRADE HANDLER ===========

func V010600UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Update dust threshold configuration for all zones
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			zone.DustThreshold = 1_000_000
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
