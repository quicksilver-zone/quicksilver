package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) allocateLockupRewards(ctx sdk.Context, allocation sdk.Coins) error {
	k.Logger(ctx).Info("allocateLockupRewards", "allocation", allocation)

	// allocate staking incentives into fee collector account to be moved to on next begin blocker by staking module
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, allocation); err != nil {
		return err
	}

	return nil
}
