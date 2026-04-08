package upgrades

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// V0101002UpgradeHandler handles the v1.10.2 upgrade.
// - Sunset stargaze-1 and omniflixhub-1: lock zones and cancel all pending redemptions
func V0101002UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting v1.10.2 upgrade...")

		if isMainnet(ctx) || isTest(ctx) {
			zonesToSunset := []string{"stargaze-1", "omniflixhub-1"}

			for _, chainID := range zonesToSunset {
				if err := sunsetZone(ctx, appKeepers, chainID); err != nil {
					return nil, fmt.Errorf("failed to sunset zone %s: %w", chainID, err)
				}
			}
		}

		ctx.Logger().Info("Upgrade v1.10.2 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// sunsetZone sets a zone to offboarding mode and cancels all pending redemptions,
// refunding qAssets from escrow back to users.
func sunsetZone(ctx sdk.Context, appKeepers *keepers.AppKeepers, chainID string) error {
	k := appKeepers.InterchainstakingKeeper

	zone, found := k.GetZone(ctx, chainID)
	if !found {
		ctx.Logger().Info("zone not found; skipping sunset", "chain_id", chainID)
		return nil
	}

	// Step 1: Lock the zone — set offboarding mode, disable deposits and unbonding
	zone.IsOffboarding = true
	zone.DepositsEnabled = false
	zone.UnbondingEnabled = false
	k.SetZone(ctx, &zone)
	ctx.Logger().Info("zone set to offboarding mode", "chain_id", chainID)

	// Step 2: Cancel all pending redemptions — refund qAssets from escrow to users
	var recordsToDelete []icstypes.WithdrawalRecord
	amountsToRefund := sdk.NewCoins()

	k.IterateZoneStatusWithdrawalRecords(ctx, chainID, icstypes.WithdrawStatusQueued, func(_ int64, record icstypes.WithdrawalRecord) bool {
		recordsToDelete = append(recordsToDelete, record)
		amountsToRefund = amountsToRefund.Add(record.BurnAmount)
		return false
	})

	k.IterateZoneStatusWithdrawalRecords(ctx, chainID, icstypes.WithdrawStatusUnbond, func(_ int64, record icstypes.WithdrawalRecord) bool {
		recordsToDelete = append(recordsToDelete, record)
		amountsToRefund = amountsToRefund.Add(record.BurnAmount)
		return false
	})

	if len(recordsToDelete) == 0 {
		ctx.Logger().Info("no pending redemptions to cancel", "chain_id", chainID)
		return nil
	}

	for _, record := range recordsToDelete {
		userAccAddress, err := addressutils.AddressFromBech32(record.Delegator, "")
		if err != nil {
			return fmt.Errorf("failed to parse delegator address %s: %w", record.Delegator, err)
		}

		if err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.EscrowModuleAccount, userAccAddress, sdk.NewCoins(record.BurnAmount)); err != nil {
			return fmt.Errorf("failed to refund qAssets to %s: %w", record.Delegator, err)
		}

		k.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
	}

	ctx.Logger().Info("cancelled pending redemptions", "chain_id", chainID, "count", len(recordsToDelete), "refunded", amountsToRefund.String())
	return nil
}
