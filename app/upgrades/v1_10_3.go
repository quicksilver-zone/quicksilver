package upgrades

import (
	"fmt"
	"sort"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// V0101003UpgradeHandler re-attempts the v1.10.2 sunset that halted the chain.
// Strategy: clamp-ASC — records are refunded in ascending BurnAmount order until
// the next record would overspend the escrow balance. Remaining records are
// left in place (still in their Queued/Unbond status) for a follow-up
// governance-driven reconciliation. The handler never panics on insufficient
// funds and never creates new qAsset supply; it only spends what exists.
func V0101003UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting v1.10.3 upgrade (clamp-ASC sunset)")

		if isMainnet(ctx) || isTest(ctx) {
			for _, chainID := range []string{"stargaze-1", "omniflixhub-1"} {
				if err := sunsetZoneClampASC(ctx, appKeepers, chainID); err != nil {
					return nil, fmt.Errorf("failed to sunset zone %s: %w", chainID, err)
				}
			}
		}

		ctx.Logger().Info("Upgrade v1.10.3 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// sunsetZoneClampASC locks the zone and refunds pending redemptions in
// ascending BurnAmount order, stopping cleanly when the escrow balance cannot
// cover the next record. Skipped records are logged and left intact for manual
// or governance-driven handling.
func sunsetZoneClampASC(ctx sdk.Context, appKeepers *keepers.AppKeepers, chainID string) error {
	k := appKeepers.InterchainstakingKeeper

	zone, found := k.GetZone(ctx, chainID)
	if !found {
		ctx.Logger().Info("zone not found; skipping sunset", "chain_id", chainID)
		return nil
	}

	zone.IsOffboarding = true
	zone.DepositsEnabled = false
	zone.UnbondingEnabled = false
	k.SetZone(ctx, &zone)

	var records []icstypes.WithdrawalRecord
	collect := func(_ int64, r icstypes.WithdrawalRecord) bool {
		// Defensive: legacy WDRs pre-v1.10.1 guards could have malformed BurnAmount
		// (empty denom, zero, or negative). sdk.NewCoins panics on those, so drop
		// them here. They remain in state under their original status for manual
		// follow-up.
		if !r.BurnAmount.IsValid() || r.BurnAmount.IsZero() {
			ctx.Logger().Error("malformed BurnAmount on WDR; skipping refund",
				"chain_id", chainID,
				"txhash", r.Txhash,
				"delegator", r.Delegator,
				"burn_amount", r.BurnAmount.String(),
			)
			return false
		}
		records = append(records, r)
		return false
	}
	k.IterateZoneStatusWithdrawalRecords(ctx, chainID, icstypes.WithdrawStatusQueued, collect)
	k.IterateZoneStatusWithdrawalRecords(ctx, chainID, icstypes.WithdrawStatusUnbond, collect)

	if len(records) == 0 {
		ctx.Logger().Info("no pending redemptions to cancel", "chain_id", chainID)
		return nil
	}

	// Ascending by BurnAmount so smaller refunds succeed first. Tie-break on
	// Txhash to keep ordering deterministic across nodes.
	sort.SliceStable(records, func(i, j int) bool {
		if c := records[i].BurnAmount.Amount.BigInt().Cmp(records[j].BurnAmount.Amount.BigInt()); c != 0 {
			return c < 0
		}
		return records[i].Txhash < records[j].Txhash
	})

	escrowAddr := appKeepers.AccountKeeper.GetModuleAddress(icstypes.EscrowModuleAccount)
	denom := zone.LocalDenom
	remaining := appKeepers.BankKeeper.GetBalance(ctx, escrowAddr, denom).Amount

	refunded := sdk.NewCoins()
	skipped := 0
	skippedAmount := math.ZeroInt()

	for _, record := range records {
		if record.BurnAmount.Amount.GT(remaining) {
			skipped++
			skippedAmount = skippedAmount.Add(record.BurnAmount.Amount)
			ctx.Logger().Error("insufficient escrow to refund record; leaving in place",
				"chain_id", chainID,
				"txhash", record.Txhash,
				"delegator", record.Delegator,
				"burn_amount", record.BurnAmount.String(),
				"escrow_remaining", remaining.String(),
			)
			continue
		}

		userAccAddress, err := addressutils.AddressFromBech32(record.Delegator, "")
		if err != nil {
			return fmt.Errorf("parse delegator %s: %w", record.Delegator, err)
		}
		if err := appKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.EscrowModuleAccount, userAccAddress, sdk.NewCoins(record.BurnAmount)); err != nil {
			return fmt.Errorf("refund %s to %s: %w", record.BurnAmount.String(), record.Delegator, err)
		}
		k.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
		remaining = remaining.Sub(record.BurnAmount.Amount)
		refunded = refunded.Add(record.BurnAmount)
	}

	ctx.Logger().Info("sunset refund summary",
		"chain_id", chainID,
		"refunded_count", len(records)-skipped,
		"refunded", refunded.String(),
		"skipped_count", skipped,
		"skipped_amount", skippedAmount.String(),
		"escrow_remaining", remaining.String(),
	)
	return nil
}
