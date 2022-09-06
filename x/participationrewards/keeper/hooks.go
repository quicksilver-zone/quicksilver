package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
)

var epochsDeferred = int64(3)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	k.Logger(ctx).Info("distribute participation rewards...")

	allocation := k.getRewardsAllocations(ctx)

	k.Logger(ctx).Info("Triggering submodule hooks")
	for _, sub := range k.prSubmodules {
		sub.Hooks(ctx, k)
	}

	if epochNumber < epochsDeferred {
		k.Logger(ctx).Info("defer...", "epoch", epochNumber)

		// create snapshot of current intents for the next epoch boundary
		// requires intents to be set, no intents no snapshot...
		// further snapshots will be taken during
		// ValidatorSelectionRewardsCallback;
		for _, zone := range k.icsKeeper.AllZones(ctx) {
			zone := zone
			for _, di := range k.icsKeeper.AllIntents(ctx, zone, false) {
				k.icsKeeper.SetIntent(ctx, zone, di, true)
			}
		}

		return
	}

	tvs, err := k.calcTokenValues(ctx)
	if err != nil {
		k.Logger(ctx).Error("unable to calculate token values", "error", err.Error())
		return
	}

	// TODO: remove this when the above is implemented
	// >>>
	/*tvs := tokenValues{
		Tokens: map[string]tokenValue{
			"uatom": {
				Symbol:     "atom",
				Multiplier: 1000000,
				Value:      sdk.NewDec(10.0),
			},
			"uosmo": {
				Symbol:     "osmo",
				Multiplier: 1000000,
				Value:      sdk.NewDec(2.0),
			},
		},
	}*/
	// <<<

	if err := k.allocateZoneRewards(ctx, tvs, allocation); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	if !allocation.Lockup.IsZero() {
		// at genesis lockup will be disable, and enabled when ICS is used.
		if err := k.allocateLockupRewards(ctx, allocation.Lockup); err != nil {
			k.Logger(ctx).Error(err.Error())
		}
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
