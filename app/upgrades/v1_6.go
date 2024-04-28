package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v6migration "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/migrations/v6"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// ============ TESTNET UPGRADE HANDLERS ============

func V010600beta0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) {
			appKeepers.UpgradeKeeper.Logger(ctx).Info("migrating capabilities")
			err := v6migration.MigrateICS27ChannelCapability(
				ctx,
				appKeepers.IBCKeeper.Codec(),
				appKeepers.GetKey(capabilitytypes.StoreKey),
				appKeepers.CapabilityKeeper,
				icstypes.ModuleName,
			)
			if err != nil {
				panic(err)
			}

			appKeepers.UpgradeKeeper.Logger(ctx).Info("removing defunct zones")
			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "agoric-3")
			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "archway-1")
			appKeepers.InterchainQueryKeeper.SetLatestHeight(ctx, "provider", 6209948)
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010600rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// no action yet.
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// =========== PRODUCTION UPGRADE HANDLER ===========

func V010600UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// no action yet.

		// TODO: remove incorrect ProtocolDataLiquidTokens
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
