package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker of airdrop module.
func (k *Keeper) EndBlocker(ctx sdk.Context) {
	for _, zd := range k.UnconcludedAirdrops(ctx) {
		if err := k.EndZoneDrop(ctx, zd.ChainId); err != nil {
			// failure in EndBlocker should NOT panic
			k.Logger(ctx).Error(err.Error())
		}
	}
}
