package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker of participationrewards module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// TODO: implement
	// - Calculate validator choice scores and allocations; for each zone independently;
	// - Calculate qAsset holdings scores and allocations;
	// - QCK staking allocations via x/distribution (using ModuleAccount and auth.FeeCollector);
}
