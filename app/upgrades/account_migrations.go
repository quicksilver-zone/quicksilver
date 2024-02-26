package upgrades

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

type ProcessMigrateAccountStrategy func(ctx sdk.Context, appKeepers *keepers.AppKeepers, from sdk.AccAddress, to sdk.AccAddress) error

// Migrate the Ingenuity genesis allocation to Notional.
func migrateIngenuityMultisigToNotional(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	// migrate ingenuity multisig to notional multisig.
	migrations := map[string]string{
		"quick1e22za5qrqqp488h5p7vw2pfx8v0y4u444ufeuw": "quick1gxrks2rcj9gthzfgrkjk5lnk0g00cg0cpyntlm",
	}
	return migrateVestingAccounts(ctx, appKeepers, migrations, migratePeriodicVestingAccount)
}

// Migrate a map of address pairs and migrate from key -> value
func migrateVestingAccounts(ctx sdk.Context, appKeepers *keepers.AppKeepers, migrations map[string]string, strategy ProcessMigrateAccountStrategy) error {
	for fromBech32, toBech32 := range migrations {
		from, err := addressutils.AccAddressFromBech32(fromBech32, "quick")
		if err != nil {
			return err
		}
		to, err := addressutils.AccAddressFromBech32(toBech32, "quick")
		if err != nil {
			return err
		}
		err = strategy(ctx, appKeepers, from, to)
		if err != nil {
			return err
		}
	}
	return nil
}

// Migrate a PeriodicVestingAccount from address A to address B, maintaining periods, amounts and end date.
func migratePeriodicVestingAccount(ctx sdk.Context, appKeepers *keepers.AppKeepers, from sdk.AccAddress, to sdk.AccAddress) error {
	oldAccount := appKeepers.AccountKeeper.GetAccount(ctx, from)
	// if the new account already exists in the account keeper, we should fail.
	if newAccount := appKeepers.AccountKeeper.GetAccount(ctx, to); newAccount != nil {
		return fmt.Errorf("unable to migrate vesting account; destination is already an account")
	}

	oldPva, ok := oldAccount.(*vestingtypes.PeriodicVestingAccount)
	if !ok {
		return fmt.Errorf("from account is not a PeriodicVestingAccount")
	}

	// copy the existing PVA.
	newPva := *oldPva

	// create a new baseVesting account with the address provided.
	newBva := vestingtypes.NewBaseVestingAccount(authtypes.NewBaseAccountWithAddress(to), oldPva.OriginalVesting, oldPva.EndTime)
	// change vesting end time so we are able to negate the token lock.
	// if the endDate has passed, we circumvent the period checking logic.
	oldPva.BaseVestingAccount.EndTime = ctx.BlockTime().Unix() - 1
	newPva.BaseVestingAccount = newBva

	// set the old pva (with the altered date), so we can transfer assets.
	appKeepers.AccountKeeper.SetAccount(ctx, oldPva)
	// set the new pva with the correct period and end dates, and new address.
	appKeepers.AccountKeeper.SetAccount(ctx, &newPva)

	// send coins from old account to new.
	err := appKeepers.BankKeeper.SendCoins(ctx, from, to, appKeepers.BankKeeper.GetAllBalances(ctx, from))
	if err != nil {
		return err
	}

	// delete the old account from the account keeper.
	appKeepers.AccountKeeper.RemoveAccount(ctx, oldPva)
	return nil
}

// migrateVestingAccountWithActions migrate from A to B with actions before migration executed
func migrateVestingAccountWithActions(ctx sdk.Context, appKeepers *keepers.AppKeepers, from sdk.AccAddress, to sdk.AccAddress) error {
	// Complete all re-delegation before unbonding
	err := completeAllRedelegations(ctx, ctx.BlockTime(), appKeepers, from)
	if err != nil {
		fmt.Printf("processMigratePeriodicVestingAccount: complete all re-delegation for %s failed: %v", from.String(), err)
		return err
	}

	// Unbond all delegation of account
	err = unbondAllDelegation(ctx, ctx.BlockTime(), appKeepers, from)
	if err != nil {
		fmt.Printf("processMigratePeriodicVestingAccount: unbonded all delegation for %s failed: %v", from.String(), err)
		return err
	}

	return migratePeriodicVestingAccount(ctx, appKeepers, from, to)
}

func completeAllRedelegations(ctx sdk.Context, now time.Time,
	appKeepers *keepers.AppKeepers,
	accAddr sdk.AccAddress,
) error {

	for _, activeRedelegation := range appKeepers.StakingKeeper.GetRedelegations(ctx, accAddr, 100) {
		redelegationSrc, _ := sdk.ValAddressFromBech32(activeRedelegation.ValidatorSrcAddress)
		redelegationDst, _ := sdk.ValAddressFromBech32(activeRedelegation.ValidatorDstAddress)

		for i := range activeRedelegation.Entries {
			activeRedelegation.Entries[i].CompletionTime = now
		}

		appKeepers.StakingKeeper.SetRedelegation(ctx, activeRedelegation)
		_, err := appKeepers.StakingKeeper.CompleteRedelegation(ctx, accAddr, redelegationSrc, redelegationDst)
		if err != nil {
			return err
		}
	}

	return nil
}

func unbondAllDelegation(ctx sdk.Context, now time.Time, appKeepers *keepers.AppKeepers, accAddr sdk.AccAddress) error {
	// Undelegate all delegations from the account
	for _, delegation := range appKeepers.StakingKeeper.GetAllDelegatorDelegations(ctx, accAddr) {
		validatorValAddr := delegation.GetValidatorAddr()
		_, found := appKeepers.StakingKeeper.GetValidator(ctx, validatorValAddr)
		if !found {
			continue
		}

		_, err := appKeepers.StakingKeeper.Undelegate(ctx, accAddr, validatorValAddr, delegation.GetShares())
		if err != nil {
			return err
		}
	}

	// Complete unbonding of all account's delegations
	for _, unbondingDelegation := range appKeepers.StakingKeeper.GetAllUnbondingDelegations(ctx, accAddr) {
		validatorStringAddr := unbondingDelegation.ValidatorAddress
		validatorValAddr, _ := sdk.ValAddressFromBech32(validatorStringAddr)

		for i := range unbondingDelegation.Entries {
			unbondingDelegation.Entries[i].CompletionTime = now
		}

		appKeepers.StakingKeeper.SetUnbondingDelegation(ctx, unbondingDelegation)
		_, err := appKeepers.StakingKeeper.CompleteUnbonding(ctx, accAddr, validatorValAddr)
		if err != nil {
			return err
		}
	}

	return nil
}
