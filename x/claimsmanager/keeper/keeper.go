package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcKeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	cdc       codec.BinaryCodec
	storeKey  storetypes.StoreKey
	IBCKeeper ibcKeeper.Keeper
}

// NewKeeper returns a new instance of participationrewards Keeper.
// This function will panic on failure.
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	ibcKeeper ibcKeeper.Keeper,
) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		IBCKeeper: ibcKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
