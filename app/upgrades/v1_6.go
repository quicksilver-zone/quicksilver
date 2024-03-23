package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
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

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
