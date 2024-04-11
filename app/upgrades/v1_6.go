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
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			zone.DustThreshold = 1_000_000
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
