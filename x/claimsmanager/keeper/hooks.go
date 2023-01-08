package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	if epochIdentifier == "epoch" && epochNumber > 1 {
		if err := k.StoreSelfConsensusState(ctx, "epoch"); err != nil {
			k.Logger(ctx).Error("unable to store consensus state", "error", err)
		}
	}
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {}

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
