package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker of module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// this is hideous, but the voting period was set too low on innuendo-4 and needs to be fixed.
	// this is the only place we have access to the govKeeper.
	// it will run once at the selected blockheight.
	if ctx.ChainID() == "innuendo-4" && ctx.BlockHeight() == 334000 {
		k.govKeeper.Logger(ctx).Info("setting gov voting period to six hours")
		votingParams := k.govKeeper.GetVotingParams(ctx)
		sixHours := time.Hour * 6
		votingParams.VotingPeriod = &sixHours
		k.govKeeper.SetVotingParams(ctx, votingParams)
	}
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
