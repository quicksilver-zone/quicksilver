package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
)

func (*Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) error {
	return nil
}

// AfterEpochEnd is called after any registered epoch ends.
// calls:
//
//	k.AggregateDelegatorIntents
//	k.HandleQueuedUnbondings
//	k.Rebalance
//
// and re-queries icq for new zone info.
func (k *Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for interchainstaking keeper.
type Hooks struct {
	k *Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
