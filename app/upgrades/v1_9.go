package upgrades

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	prtypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func V010900UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting module migrations...")

		if isMainnet(ctx) {
			// add MembraneParams protocol data
			membranePd := prtypes.MembraneProtocolData{
				ContractAddress: "osmo1gy5gpqqlth0jpm9ydxlmff6g5mpnfvrfxd3mfc8dhyt03waumtzqt8exxr",
			}

			blob, err := json.Marshal(membranePd)
			if err != nil {
				panic(fmt.Errorf("failed to marshal membrane protocol data: %w", err))
			}

			pd := prtypes.NewProtocolData(prtypes.ProtocolDataType_name[int32(prtypes.ProtocolDataTypeMembraneParams)], blob)
			appKeepers.ParticipationRewardsKeeper.SetProtocolData(ctx, membranePd.GenerateKey(), pd)
		}

		ctx.Logger().Info("Upgrade v1.9.0 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
