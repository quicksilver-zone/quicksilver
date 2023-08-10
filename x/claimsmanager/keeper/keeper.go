package keeper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

type Keeper struct {
	cdc       codec.BinaryCodec
	storeKey  storetypes.StoreKey
	IBCKeeper *ibckeeper.Keeper
}

// NewKeeper returns a new instance of participationrewards Keeper.
// This function will panic on failure.
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	ibcKeeper *ibckeeper.Keeper,
) Keeper {
	if ibcKeeper == nil {
		panic("ibcKeeper is nil")
	}

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

func (k Keeper) StoreSelfConsensusState(ctx sdk.Context, key string) error {
	var height ibcclienttypes.Height
	if strings.Contains(ctx.ChainID(), "-") {
		revisionNum, err := strconv.ParseUint(strings.Split(ctx.ChainID(), "-")[1], 10, 64)
		if err != nil {
			k.Logger(ctx).Error("Error getting revision number for client ")
			return err
		}

		height = ibcclienttypes.Height{
			RevisionNumber: revisionNum,
			RevisionHeight: uint64(ctx.BlockHeight() - 1),
		}
	} else {
		// ONLY FOR TESTING - ibctesting module chains donot follow standard [chainname]-[num] structure
		height = ibcclienttypes.Height{
			RevisionNumber: 0, // revision number for testchain1-1 is 0 (because parseChainId splits on '-')
			RevisionHeight: uint64(ctx.BlockHeight() - 1),
		}
	}

	selfConsState, err := k.IBCKeeper.ClientKeeper.GetSelfConsensusState(ctx, height)
	if err != nil {
		k.Logger(ctx).Error("Error getting self consensus state of previous height")
		return err
	}

	state, _ := selfConsState.(*ibctmtypes.ConsensusState)
	k.SetSelfConsensusState(ctx, key, state)

	return nil
}
