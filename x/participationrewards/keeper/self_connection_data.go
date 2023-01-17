package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) SetSelfProtocolData(ctx sdk.Context, data *types.ProtocolData) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeySelfProtocolData, k.cdc.MustMarshal(data))
}

func (k Keeper) GetSelfProtocolData(ctx sdk.Context) (types.ProtocolData, bool) {
	store := ctx.KVStore(k.storeKey)

	var selfProtocolData types.ProtocolData
	k.cdc.MustUnmarshal(store.Get(types.KeySelfProtocolData), &selfProtocolData)

	return selfProtocolData, true
}
