package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func NewProtocolData(datatype string, protocol string, data json.RawMessage) *types.ProtocolData {
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
	store.Set([]byte(key), bz)
}

// DeleteProtocolData delete protocol data info
func (k Keeper) DeleteProtocolData(ctx sdk.Context, key string, protocol string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	store.Delete([]byte(key))
}

// IteratePrefixedProtocolDatas iterate through protocol datas with the given prefix and perform the provided function
func (k Keeper) IteratePrefixedProtocolDatas(ctx sdk.Context, keyPrefix string, fn func(index int64, data types.ProtocolData) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	iterator := sdk.KVStorePrefixIterator(store, []byte(keyPrefix))
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

// IterateAllProtocolDatas iterate through protocol data and perform the provided function
func (k Keeper) IterateAllProtocolDatas(ctx sdk.Context, fn func(index int64, key string, data types.ProtocolData) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.ProtocolData{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, string(iterator.Key()), data)
		if stop {
			break
		}
		i++
	}
}

// AllKeyedProtocolDatas returns a slice containing all protocol datas and their keys from the store.
func (k Keeper) AllKeyedProtocolDatas(ctx sdk.Context) []*types.KeyedProtocolData {
	out := make([]*types.KeyedProtocolData, 0)
	k.IterateAllProtocolDatas(ctx, func(_ int64, key string, data types.ProtocolData) (stop bool) {
		out = append(out, &types.KeyedProtocolData{Key: key, ProtocolData: &data})
		return false
	})
	return out
}

func GetProtocolDataKey(protocol string, key string) string {
	return fmt.Sprintf("%s/%s", protocol, key)
}
