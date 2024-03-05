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
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			appKeepers.InterchainstakingKeeper.SetLocalDenomZoneMapping(ctx, zone)
			return false
		})
		migrations := map[string]string{
			"quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l": "quick1h0sqndv2y4xty6uk0sv4vckgyc5aa7n5at7fll",
			"quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f": "quick1n4g6037cjm0e0v2nvwj2ngau7pk758wtwk6lwq",
		}

		if err := migrateVestingAccounts(ctx, appKeepers, migrations, migrateVestingAccountWithActions); err != nil {
			panic(err)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
