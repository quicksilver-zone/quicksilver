package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker of module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
}

// EndBlocker of module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	for _, zd := range k.UnconcludedAirdrops(ctx) {
		if err := k.EndZoneDrop(ctx, zd.ChainId); err != nil {
			// failure in EndBlocker should NOT panic
			k.Logger(ctx).Error(err.Error())
		}
	}
}
