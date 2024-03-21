package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// =========== PRODUCTION UPGRADE HANDLER ===========

func V010600UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Update dust threshold configuration for all zones
		thresholds := map[string]int64{
			"osmosis-1":   1_000_000,
			"cosmoshub-4": 1_000_000,
			"stargaze-1":  5_000_000,
			"juno-1":      2_000_000,
			"sommelier-3": 5_000_000,
			"regen-1":     5_000_000,
			"dydx-mainnet-1":      500_000_000_000_000_000,
		}
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			threshold, ok := thresholds[zone.ChainId]
			// if threshold not exist => get default value
			if !ok {
				threshold = 1_000_000
			}
			zone.DustThreshold = threshold
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
