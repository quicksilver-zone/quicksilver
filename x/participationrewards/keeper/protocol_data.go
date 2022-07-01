package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) NewProtocolData(ctx sdk.Context, datatype string, protocol string, data json.RawMessage) *types.ProtocolData {
	return &types.ProtocolData{Type: datatype, Protocol: protocol, Data: data}
}

// GetProtocolData returns data
func (k Keeper) GetProtocolData(ctx sdk.Context, key string) (types.ProtocolData, bool) {
	data := types.ProtocolData{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	bz := store.Get([]byte(key))
	if len(bz) == 0 {
		return data, false
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

// SetProtocolData set protocol data info
func (k Keeper) SetProtocolData(ctx sdk.Context, key string, data *types.ProtocolData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	bz := k.cdc.MustMarshal(data)
	store.Set([]byte(GetProtocolDataKey(data.Protocol, key)), bz)
}

// DeleteProtocolData delete protocol data info
func (k Keeper) DeleteProtocolData(ctx sdk.Context, key string, protocol string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	store.Delete([]byte(GetProtocolDataKey(protocol, key)))
}

// IterateQueries iterate through protocol datas
func (k Keeper) IterateProtocolDatas(ctx sdk.Context, protocol string, fn func(index int64, data types.ProtocolData) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	iterator := sdk.KVStorePrefixIterator(store, []byte(protocol))
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.ProtocolData{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, data)
		if stop {
			break
		}
		i++
	}
}

func GetProtocolDataKey(protocol string, key string) string {
	return fmt.Sprintf("%s/%s", protocol, key)
}
