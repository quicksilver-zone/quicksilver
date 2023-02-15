package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

// GetSelfConsensusState returns consensus state stored every epoch
func (k Keeper) GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool) {
	store := ctx.KVStore(k.storeKey)
	var selfConsensusState ibctmtypes.ConsensusState

	bz := store.Get(append(types.KeySelfConsensusState, []byte(key)...))
	if bz == nil {
		return selfConsensusState, false
	}
	k.cdc.MustUnmarshal(bz, &selfConsensusState)
	return selfConsensusState, true
}

// SetSelfConsensusState sets the self consensus state
func (k Keeper) SetSelfConsensusState(ctx sdk.Context, key string, consState *ibctmtypes.ConsensusState) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.KeySelfConsensusState, []byte(key)...), k.cdc.MustMarshal(consState))
}

// DeleteSelfConsensusState deletes the self consensus state
func (k Keeper) DeleteSelfConsensusState(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.KeySelfConsensusState, []byte(key)...))
}
