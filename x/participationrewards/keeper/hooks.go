package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	k.Logger(ctx).Info("Distribute participation rewards...")

	allocation := k.getRewardsAllocations(ctx)

	// TODO: obtain zone tvl and calculate zone allocation proportions

	if err := k.allocateValidatorSelectionRewards(ctx, allocation.ValidatorSelection); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	if err := k.allocateHoldingsRewards(ctx, allocation.Holdings); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	if err := k.allocateLockupRewards(ctx, allocation.Lockup); err != nil {
		k.Logger(ctx).Error(err.Error())
	}
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
