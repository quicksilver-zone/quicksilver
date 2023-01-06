package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

// GetSelfConsensusState returns consensus state stored every epoch
func (k Keeper) GetSelfConsensusState(ctx sdk.Context) (ibctmtypes.ConsensusState, bool) {
	store := ctx.KVStore(k.storeKey)

	var selfConsensusState ibctmtypes.ConsensusState
	k.cdc.MustUnmarshal(store.Get(types.KeySelfConsensusState), &selfConsensusState)

	return selfConsensusState, true
}

// SetSelfConsensusState sets the self consensus state
func (k Keeper) SetSelfConsensusState(ctx sdk.Context, consState ibctmtypes.ConsensusState) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeySelfConsensusState, k.cdc.MustMarshal(&consState))
}
