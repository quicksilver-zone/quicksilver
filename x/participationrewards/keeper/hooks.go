package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	k.Logger(ctx).Info("Distribute participation rewards...")

	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalances := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
	k.Logger(ctx).Info("module account", "address", moduleAddress, "balances", moduleBalances)
	// DEVTEST:
	fmt.Printf("AfterEpochEnd >>>\n\tAddress = %v\n\tBalance = %v\n", moduleAddress, moduleBalances)

	if moduleBalances.Empty() {
		k.Logger(ctx).Info("nothing to distribute...")
		// DEVTEST:
		fmt.Println("nothing to distribute...")

		return
	}

	// get distribution proportions (params)
	params := k.GetParams(ctx)
	k.Logger(ctx).Info("module parameters", "params", params)
	// DEVTEST:
	fmt.Printf("\nParams:\n%v\n", params)

	// split participation rewards allocations
	validatorSelectionAllocation := sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalances[0],
			params.DistributionProportions.ValidatorSelectionAllocation,
		),
	)
	holdingsAllocation := sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalances[0],
			params.DistributionProportions.HoldingsAllocation,
		),
	)
	lockupAllocation := sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalances[0],
			params.DistributionProportions.LockupAllocation,
		),
	)

	// use sum to check total distribution to collect and allocate dust
	total := moduleBalances[0]
	sum := lockupAllocation.Add(validatorSelectionAllocation...).Add(holdingsAllocation...)
	dust := total.Sub(sum[0])
	k.Logger(ctx).Info(
		"rewards distribution",
		"total", total,
		"validatorSelectionAllocation", validatorSelectionAllocation,
		"holdingsAllocation", holdingsAllocation,
		"lockupAllocation", lockupAllocation,
		"sum", sum,
		"dust", dust,
	)
	// DEVTEST:
	fmt.Printf("\tTotal:\t\t\t\t%v\n", total)
	fmt.Printf("\tValidator Selection Allocation:\t%v\n", validatorSelectionAllocation)
	fmt.Printf("\tHoldings Allocation:\t\t%v\n", holdingsAllocation)
	fmt.Printf("\tLockup Allocation:\t\t%v\n", lockupAllocation)
	fmt.Printf("\tSum:\t\t\t\t%v\n", sum)
	fmt.Printf("\tDust:\t\t\t\t%v\n", dust)

	// Add dust to validator choice allocation (favors decentralization)
	validatorSelectionAllocation = validatorSelectionAllocation.Add(dust)
	k.Logger(ctx).Info("add dust to validatorSelectionAllocation...")
	// DEVTEST:
	fmt.Printf("\n\tAdd dust to validatorSelectionAllocation...\n\n")

	if err := k.allocateValidatorChoiceRewards(ctx, validatorSelectionAllocation); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	if err := k.allocateHoldingsRewards(ctx, holdingsAllocation); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	if err := k.allocateLockupRewards(ctx, lockupAllocation); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	// DEVTEST:
	fmt.Printf("\n<<<\n")
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
