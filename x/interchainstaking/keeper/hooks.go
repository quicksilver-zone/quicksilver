package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	// every epoch
	k.Logger(ctx).Info("Handling epoch end")
	if epochIdentifier == "epoch" {
		k.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
			for _, da := range zoneInfo.DelegationAddresses {
				k.WithdrawDelegationRewards(ctx, zoneInfo, da)
			}
			return false
		})
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
