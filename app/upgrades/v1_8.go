package upgrades

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

func V010800UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting module migrations...")

		// Migrate Tendermint consensus parameters from x/params module to dedicated x/consensus module.
		baseAppLegacySS := appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &appKeepers.ConsensusKeeper)

		// explicitly update the IBC 02-client params, adding the localhost client type
		params := appKeepers.IBCKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		appKeepers.IBCKeeper.ClientKeeper.SetParams(ctx, params)

		ctx.Logger().Info("Upgrade v1.8.0 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
