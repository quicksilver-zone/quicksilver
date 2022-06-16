package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) allocateHoldingsRewards(ctx sdk.Context) error {
	k.Logger(ctx).Info("allocateHoldingsRewards")

	return fmt.Errorf("allocateHoldingsRewards not implemented")
}
