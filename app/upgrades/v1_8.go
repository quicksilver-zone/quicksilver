package upgrades

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"

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

		_, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, appKeepers.AppCodec, appKeepers.IBCKeeper.ClientKeeper)
		if err != nil {
			panic(fmt.Errorf("failed to prune expired consensus states: %w", err))
		}

		if isMainnet(ctx) || isTest(ctx) {
			// juno-1 - ubrs where extra records were created.
			appKeepers.InterchainstakingKeeper.IteratePrefixedUnbondingRecords(ctx, []byte("juno-1"), func(index int64, record icstypes.UnbondingRecord) (stop bool) {
				if record.EpochNumber == 277 {
					appKeepers.InterchainstakingKeeper.DeleteUnbondingRecord(ctx, record.ChainId, record.Validator, record.EpochNumber)
				}
				return false
			})

			// juno-1 send that happened but were never acked, so were recreated.
			appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, "juno-1", "564e8a6263763644bbe32e4bd0bf9f99619aaf68b938216fff2acef2dfb8aec6", 4)
			appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, "juno-1", "c746ceba8da060f25a81f2e0cc6ed53fecd69dbbd89ff7c9aa8b5d0464302f84", 4)
		}

		ctx.Logger().Info("Upgrade v1.8.0 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
