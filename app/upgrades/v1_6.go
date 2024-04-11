package upgrades

import (
	"cosmossdk.io/math"
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
		// no action yet.
		// migrate killer queen incentives
		migrations := map[string]string{
			"quick1qfyntnmlvznvrkk9xqppmcxqcluv7wd74nmyus": "quick1898x2jpjfelg4jvl4hqm9a9vugyctfdcl9t64x",
		}
		err := migrateVestingAccounts(ctx, appKeepers, migrations, migratePeriodicVestingAccount)
		if err != nil {
			panic(err)
		}

		// Update dust threshold configuration for all zones
		thresholds := map[string]math.Int{
			"osmosis-1":      math.NewInt(1_000_000),
			"cosmoshub-4":    math.NewInt(1_000_000),
			"stargaze-1":     math.NewInt(5_000_000),
			"juno-1":         math.NewInt(2_000_000),
			"sommelier-3":    math.NewInt(5_000_000),
			"regen-1":        math.NewInt(5_000_000),
			"dydx-mainnet-1": math.NewInt(500_000_000_000_000_000),
		}
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			threshold, ok := thresholds[zone.ChainId]
			// if threshold not exist => get default value
			if !ok {
				threshold = math.NewInt(1_000_000)
			}
			zone.DustThreshold = threshold
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010600rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// iterate over all withdrawal records with zero BurnAmount and delete them
		appKeepers.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(_ int64, record icstypes.WithdrawalRecord) (stop bool) {
			if record.BurnAmount.IsZero() {
				appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
			}
			return false
		})
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
