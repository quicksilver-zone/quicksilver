package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func V010700UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isMainnet(ctx) || isTest(ctx) {

			hashes := []struct {
				Zone string
				Hash string
			}{
				{Zone: "cosmoshub-4", Hash: "0c8269f04109a55a152d3cdfd22937b4e5c2746111d579935eef4cd7ffa71f7f"},
			}
			for _, hashRecord := range hashes {
				// delete duplicate records.
				appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, hashRecord.Zone, hashRecord.Hash, icstypes.WithdrawStatusUnbond)
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("delete duplicate withdrawal record", "hash", hashRecord.Hash, "zone", hashRecord.Zone)
			}

			err := appKeepers.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(50699994))))
			if err != nil {
				panic(err)
			}
			err = appKeepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(50699994))))
			if err != nil {
				panic(err)
			}

			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
				appKeepers.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, zone)
				return false
			})

		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
