package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// GetProtocolData returns protocol data.
func (k *Keeper) GetProtocolData(ctx sdk.Context, pdType types.ProtocolDataType, key string) (types.ProtocolData, bool) {
	data := types.ProtocolData{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	bz := store.Get(types.GetProtocolDataKey(pdType, []byte(key)))
	if len(bz) == 0 {
		return data, false
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

// SetProtocolData set protocol data info.
func (k Keeper) SetProtocolData(ctx sdk.Context, key []byte, data *types.ProtocolData) {
	if data == nil {
		k.Logger(ctx).Error("protocol data not set; value is nil")
		return
	}

	pdType, exists := types.ProtocolDataType_value[data.Type]
	if !exists {
		k.Logger(ctx).Error("protocol data not set; type does not exist", "type", data.Type)
		return
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	bz := k.cdc.MustMarshal(data)
	store.Set(types.GetProtocolDataKey(types.ProtocolDataType(pdType), key), bz)
}

// DeleteProtocolData deletes protocol data info.
func (k *Keeper) DeleteProtocolData(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	store.Delete(key)
}

// IteratePrefixedProtocolDatas iterate through protocol data with the given prefix and perform the provided function.
func (k *Keeper) IteratePrefixedProtocolDatas(ctx sdk.Context, key []byte, fn func(index int64, key []byte, data types.ProtocolData) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtocolData)
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.ProtocolData{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, iterator.Key(), data)
		if stop {
			break
		}
		i++
	}
}

// IterateAllProtocolDatas iterates through protocol data and perform the provided function.
func (k *Keeper) IterateAllProtocolDatas(ctx sdk.Context, fn func(index int64, key string, data types.ProtocolData) (stop bool)) {
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
func (k *Keeper) AllKeyedProtocolDatas(ctx sdk.Context) []*types.KeyedProtocolData {
	out := make([]*types.KeyedProtocolData, 0)
	k.IterateAllProtocolDatas(ctx, func(_ int64, key string, data types.ProtocolData) (stop bool) {
		out = append(out, &types.KeyedProtocolData{Key: key, ProtocolData: &data})
		return false
	})
	return out
}

/*func GetProtocolDataKey(protocol string, key string) string {
	return fmt.Sprintf("%s/%s", protocol, key)
}*/
