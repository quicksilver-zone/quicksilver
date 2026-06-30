package upgrades

import (
	"fmt"
	"sort"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

var exploitCleanupAccounts = map[string]struct{}{
	"quick1sz5dgde3h3cs2u6zn6a7r5exkcenqhqxn2hhzx":                      {},
	"quick1qsr3pdhzt4tkqx5x502wf27dklfptxdhr3h0ff":                      {},
	"quick1g9q34rnewunqeu6zhdj6vj0lt9d5qh8zkg6pmg":                      {},
	"quick1q0km4k0yvzwwmyqzr2egyqrrc7kedncl3z6gjegt9l3pyqasj0xu2fmqw49": {},
	"quick1sqxhlr25v8e9lc4ejk3jres9c58uxzy9zlac0n":                      {},
	"quick1eud05ql9ftvqe80esls6mnwc534e8gm4n0h6c8":                      {},
	"quick1dt70tuyr3fn0pwdenerxvulv5sq0l35q0zmygf":                      {},
}

var exploitCleanupDenoms = []string{
	"uqatom",
	"uqosmo",
	"uqck",
}

var exploitCleanupDenomSet = map[string]struct{}{
	"uqatom": {},
	"uqosmo": {},
	"uqck":   {},
}

var exploitRemintDenomSet = map[string]math.Int{
	"uqatom2": math.NewInt(13642764531),
	"uqosmo2": math.NewInt(793710966168),
}

// V011000UpgradeHandler performs the Quicksilver sunset exploit cleanup. It
// cancels exploit-related redemptions by deleting their withdrawal records and
// burning the escrowed qAssets, then burns liquid qAssets still held by the known
// exploit-related Quicksilver accounts.
//
// For non-exploit withdrawal records, it cancels records that have not reached
// WithdrawStatusSend and refunds the original BurnAmount to the delegator. Status
// SEND records are left in place for manual verification that the remote-chain
// recipient did not already receive the funds.
//
// It then mints qAtom2 and qOsmo2 to the exploitRemintDenomSet amounts, and sends
// them to an account for distribution.
func V011000UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting v1.11.0 upgrade (sunset exploit cleanup)")

		if isMainnet(ctx) || isTest(ctx) {
			markAllZonesOffboarding(ctx, appKeepers)
			if err := cancelExploitRedemptionsAndBurnEscrow(ctx, appKeepers); err != nil {
				return nil, fmt.Errorf("cancel exploit redemptions: %w", err)
			}
			if err := burnExploitAccountBalances(ctx, appKeepers); err != nil {
				return nil, fmt.Errorf("burn exploit account balances: %w", err)
			}
			if err := cancelOutstandingRedemptionsAndRefundEscrow(ctx, appKeepers); err != nil {
				return nil, fmt.Errorf("cancel outstanding redemptions: %w", err)
			}
			if err := remintExploitDenoms(ctx, appKeepers); err != nil {
				return nil, fmt.Errorf("remint exploit denoms: %w", err)
			}
		}

		ctx.Logger().Info("Upgrade v1.11.0 complete")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func markAllZonesOffboarding(ctx sdk.Context, appKeepers *keepers.AppKeepers) {
	k := appKeepers.InterchainstakingKeeper
	count := 0

	k.IterateZones(ctx, func(_ int64, zone *icstypes.Zone) bool {
		zone.IsOffboarding = true
		zone.DepositsEnabled = false
		zone.UnbondingEnabled = false
		k.SetZone(ctx, zone)
		count++
		return false
	})

	ctx.Logger().Info("marked all zones as offboarding", "count", count)
}

func cancelExploitRedemptionsAndBurnEscrow(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	k := appKeepers.InterchainstakingKeeper
	statuses := []int32{
		icstypes.WithdrawStatusTokenize,
		icstypes.WithdrawStatusQueued,
		icstypes.WithdrawStatusUnbond,
		icstypes.WithdrawStatusSend,
	}

	var records []icstypes.WithdrawalRecord
	for _, status := range statuses {
		k.IterateWithdrawalRecords(ctx, func(_ int64, record icstypes.WithdrawalRecord) bool {
			if record.Status != status || !isExploitAccount(record.Delegator) || !isExploitDenom(record.BurnAmount.Denom) {
				return false
			}
			records = append(records, record)
			return false
		})
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].ChainId != records[j].ChainId {
			return records[i].ChainId < records[j].ChainId
		}
		return records[i].Txhash < records[j].Txhash
	})

	burned := sdk.NewCoins()
	for _, record := range records {
		if !record.BurnAmount.IsValid() || record.BurnAmount.IsZero() {
			ctx.Logger().Error("deleting malformed exploit withdrawal record without burn",
				"chain_id", record.ChainId,
				"txhash", record.Txhash,
				"delegator", record.Delegator,
				"status", record.Status,
				"burn_amount", record.BurnAmount.String(),
			)
			k.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
			continue
		}

		coins := sdk.NewCoins(record.BurnAmount)
		if err := appKeepers.BankKeeper.BurnCoins(ctx, icstypes.EscrowModuleAccount, coins); err != nil {
			return fmt.Errorf("burn escrowed %s for %s/%s: %w", coins.String(), record.ChainId, record.Txhash, err)
		}
		k.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
		burned = burned.Add(record.BurnAmount)
	}

	ctx.Logger().Info("cancelled exploit redemptions and burned escrow",
		"count", len(records),
		"burned", burned.String(),
	)
	return nil
}

func cancelOutstandingRedemptionsAndRefundEscrow(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	k := appKeepers.InterchainstakingKeeper

	var records []icstypes.WithdrawalRecord
	obligations := sdk.NewCoins()

	k.IterateWithdrawalRecords(ctx, func(_ int64, record icstypes.WithdrawalRecord) bool {
		if isExploitAccount(record.Delegator) {
			return false
		}

		switch record.Status {
		case icstypes.WithdrawStatusTokenize, icstypes.WithdrawStatusQueued, icstypes.WithdrawStatusUnbond:
			// These records have not reached SEND, so the remote-chain recipient
			// should not have received funds from the withdrawal flow.
		default:
			return false
		}

		if !record.BurnAmount.IsValid() || record.BurnAmount.IsZero() {
			ctx.Logger().Error("malformed withdrawal record burn amount; leaving for manual follow-up",
				"chain_id", record.ChainId,
				"txhash", record.Txhash,
				"delegator", record.Delegator,
				"status", record.Status,
				"burn_amount", record.BurnAmount.String(),
			)
			return false
		}

		records = append(records, record)
		obligations = obligations.Add(record.BurnAmount)
		return false
	})

	sort.Slice(records, func(i, j int) bool {
		if records[i].ChainId != records[j].ChainId {
			return records[i].ChainId < records[j].ChainId
		}
		if records[i].Status != records[j].Status {
			return records[i].Status < records[j].Status
		}
		return records[i].Txhash < records[j].Txhash
	})

	escrowAddr := appKeepers.AccountKeeper.GetModuleAddress(icstypes.EscrowModuleAccount)
	for _, obligation := range obligations {
		escrowBalance := appKeepers.BankKeeper.GetBalance(ctx, escrowAddr, obligation.Denom)
		if escrowBalance.Amount.GTE(obligation.Amount) {
			continue
		}

		shortfall := sdk.NewCoin(obligation.Denom, obligation.Amount.Sub(escrowBalance.Amount))
		shortfallCoins := sdk.NewCoins(shortfall)
		if err := appKeepers.BankKeeper.MintCoins(ctx, icstypes.ModuleName, shortfallCoins); err != nil {
			return fmt.Errorf("mint refund shortfall %s: %w", shortfallCoins.String(), err)
		}
		if err := appKeepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, shortfallCoins); err != nil {
			return fmt.Errorf("transfer refund shortfall %s to escrow: %w", shortfallCoins.String(), err)
		}

		ctx.Logger().Info("minted refund shortfall into escrow",
			"denom", obligation.Denom,
			"escrow_before", escrowBalance.Amount.String(),
			"obligation", obligation.Amount.String(),
			"minted", shortfall.Amount.String(),
		)
	}

	refunded := sdk.NewCoins()
	for _, record := range records {
		addr, err := sdk.AccAddressFromBech32(record.Delegator)
		if err != nil {
			return fmt.Errorf("parse delegator %s: %w", record.Delegator, err)
		}

		coins := sdk.NewCoins(record.BurnAmount)
		if err := appKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.EscrowModuleAccount, addr, coins); err != nil {
			return fmt.Errorf("refund %s to %s for %s/%s: %w", coins.String(), record.Delegator, record.ChainId, record.Txhash, err)
		}

		k.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
		refunded = refunded.Add(record.BurnAmount)
	}

	ctx.Logger().Info("cancelled outstanding non-exploit redemptions and refunded escrow",
		"count", len(records),
		"refunded", refunded.String(),
		"skipped_status", icstypes.WithdrawStatusSend,
	)
	return nil
}

func burnExploitAccountBalances(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	burned := sdk.NewCoins()
	accounts := make([]string, 0, len(exploitCleanupAccounts))
	for account := range exploitCleanupAccounts {
		accounts = append(accounts, account)
	}
	sort.Strings(accounts)

	for _, account := range accounts {
		addr, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return fmt.Errorf("parse exploit account %s: %w", account, err)
		}

		for _, denom := range exploitCleanupDenoms {
			balance := appKeepers.BankKeeper.GetBalance(ctx, addr, denom)
			if balance.IsZero() {
				continue
			}

			coins := sdk.NewCoins(balance)
			if err := appKeepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, icstypes.EscrowModuleAccount, coins); err != nil {
				return fmt.Errorf("move %s from %s to escrow: %w", coins.String(), account, err)
			}
			if err := appKeepers.BankKeeper.BurnCoins(ctx, icstypes.EscrowModuleAccount, coins); err != nil {
				return fmt.Errorf("burn %s from %s: %w", coins.String(), account, err)
			}
			burned = burned.Add(balance)
		}
	}

	ctx.Logger().Info("burned exploit account balances", "burned", burned.String())
	return nil
}

func remintExploitDenoms(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	for denom, amount := range exploitRemintDenomSet {
		if err := appKeepers.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin(denom, amount))); err != nil {
			return fmt.Errorf("mint %s: %w", denom, err)
		}
		if err := appKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, sdk.MustAccAddressFromBech32("quick12zqrwrgp4ntqswk8rn9dcv2ykwqa3ml2gak7yr"), sdk.NewCoins(sdk.NewCoin(denom, amount))); err != nil {
			return fmt.Errorf("send %s to %s: %w", denom, "quick12zqrwrgp4ntqswk8rn9dcv2ykwqa3ml2gak7yr", err)
		}
		ctx.Logger().Info("reminted exploit denoms and sent to distribution account", "denom", denom, "amount", amount.String())
	}

	return nil
}

func isExploitAccount(address string) bool {
	_, ok := exploitCleanupAccounts[address]
	return ok
}

func isExploitDenom(denom string) bool {
	_, ok := exploitCleanupDenomSet[denom]
	return ok
}
