package upgrades

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// V0101002UpgradeHandler sunsets stargaze-1 and omniflixhub-1 via the
// mint-to-cover strategy: if the escrow balance is short of the total qAsset
// refund obligation for a zone, the handler mints the exact shortfall into the
// interchainstaking module account and transfers it into the escrow module
// account before refunding users. This guarantees every WDR is refunded in
// full, avoiding the escrow-underflow panic that halted the original v1.10.2
// attempt on mainnet.
//
// The newly-minted qAssets are economically backed by the base-denom balances
// held across the zone's ICA accounts on the host chain (delegated, unbonding,
// rewards, and liquid). Those host-side tokens will be reclaimed via a
// post-sunset governance action to cover the minted supply.
func V0101002UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting v1.10.2 upgrade (mint-to-cover sunset)")

		if isMainnet(ctx) || isTest(ctx) {
			for _, chainID := range []string{"stargaze-1", "omniflixhub-1"} {
				if err := sunsetZoneMintToCover(ctx, appKeepers, chainID); err != nil {
					return nil, fmt.Errorf("failed to sunset zone %s: %w", chainID, err)
				}
			}
		}

		ctx.Logger().Info("Upgrade v1.10.2 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// sunsetZoneMintToCover locks the zone and refunds all pending redemptions,
// minting into escrow if the recorded obligations exceed the escrow balance.
func sunsetZoneMintToCover(ctx sdk.Context, appKeepers *keepers.AppKeepers, chainID string) error {
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
	totalObligation := math.ZeroInt()

	collect := func(_ int64, record icstypes.WithdrawalRecord) bool {
		// Defensive: legacy WDRs pre-v1.10.1 guards could have malformed BurnAmount
		// (empty denom, zero, or negative). sdk.NewCoins panics on those, so drop
		// them from the refund set. Excluding them here also prevents the mint
		// branch from over-minting against obligations we can't actually honor.
		// Malformed records remain in state under their original status for
		// manual follow-up.
		if !record.BurnAmount.IsValid() || record.BurnAmount.IsZero() {
			ctx.Logger().Error("malformed BurnAmount on WDR; skipping refund",
				"chain_id", chainID,
				"txhash", record.Txhash,
				"delegator", record.Delegator,
				"burn_amount", record.BurnAmount.String(),
			)
			return false
		}
		records = append(records, record)
		totalObligation = totalObligation.Add(record.BurnAmount.Amount)
		return false
	}
	k.IterateZoneStatusWithdrawalRecords(ctx, chainID, icstypes.WithdrawStatusQueued, collect)
	k.IterateZoneStatusWithdrawalRecords(ctx, chainID, icstypes.WithdrawStatusUnbond, collect)

	if len(records) == 0 {
		ctx.Logger().Info("no pending redemptions to cancel", "chain_id", chainID)
		return nil
	}

	localDenom := zone.LocalDenom
	escrowAddr := appKeepers.AccountKeeper.GetModuleAddress(icstypes.EscrowModuleAccount)
	escrowBalance := appKeepers.BankKeeper.GetBalance(ctx, escrowAddr, localDenom).Amount

	if escrowBalance.LT(totalObligation) {
		shortfall := totalObligation.Sub(escrowBalance)
		shortfallCoins := sdk.NewCoins(sdk.NewCoin(localDenom, shortfall))

		if err := appKeepers.BankKeeper.MintCoins(ctx, icstypes.ModuleName, shortfallCoins); err != nil {
			return fmt.Errorf("mint shortfall %s: %w", shortfallCoins.String(), err)
		}
		if err := appKeepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, shortfallCoins); err != nil {
			return fmt.Errorf("transfer shortfall to escrow %s: %w", shortfallCoins.String(), err)
		}
		ctx.Logger().Info("minted shortfall into escrow",
			"chain_id", chainID,
			"denom", localDenom,
			"escrow_before", escrowBalance.String(),
			"obligation", totalObligation.String(),
			"minted", shortfall.String(),
		)
	}

	refunded := sdk.NewCoins()
	for _, record := range records {
		userAccAddress, err := addressutils.AddressFromBech32(record.Delegator, "")
		if err != nil {
			return fmt.Errorf("parse delegator %s: %w", record.Delegator, err)
		}
		if err := appKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.EscrowModuleAccount, userAccAddress, sdk.NewCoins(record.BurnAmount)); err != nil {
			return fmt.Errorf("refund %s to %s: %w", record.BurnAmount.String(), record.Delegator, err)
		}
		k.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
		refunded = refunded.Add(record.BurnAmount)
	}

	ctx.Logger().Info("cancelled pending redemptions",
		"chain_id", chainID, "count", len(records), "refunded", refunded.String())
	return nil
}
