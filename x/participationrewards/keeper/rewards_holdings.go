package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) allocateHoldingsRewards(ctx sdk.Context, allocation sdk.Coins) error {
	k.Logger(ctx).Info("allocateHoldingsRewards", "allocation", allocation)
	// DEVTEST:
	fmt.Printf("\t\tAllocate Holdings Rewards:\t\t%v\n", allocation)

	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, allocation)
	if err != nil {
		return err
	}

	return fmt.Errorf("allocateHoldingsRewards not implemented")
}
