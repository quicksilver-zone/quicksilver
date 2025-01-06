package keeper

import (
	"encoding/json"
	"fmt"

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

func GetAndUnmarshalProtocolData[T any](ctx sdk.Context, k *Keeper, key string, pdType types.ProtocolDataType) (dt types.ProtocolData, tt T, err error) {
	data, ok := k.GetProtocolData(ctx, pdType, key)
	if !ok {
		return dt, tt, fmt.Errorf("unable to find protocol data for %q", key)
	}
	pd, err := types.UnmarshalProtocolData(pdType, data.Data)
	if err != nil {
		return dt, tt, err
	}
	asType, ok := pd.(T)
	if !ok {
		return dt, tt, fmt.Errorf("could not retrieve type of %T, actual type: %T", (*T)(nil), pd)
	}
	return data, asType, nil
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

// AllKeyedProtocolDatas returns a slice containing all protocol data and their keys from the store.
func (k *Keeper) AllKeyedProtocolDatas(ctx sdk.Context) []*types.KeyedProtocolData {
	out := make([]*types.KeyedProtocolData, 0)
	k.IterateAllProtocolDatas(ctx, func(_ int64, key string, data types.ProtocolData) (stop bool) {
		out = append(out, &types.KeyedProtocolData{Key: key, ProtocolData: &data})
		return false
	})
	return out
}

// MarshalAndSetProtocolData marshals and sets protocol data given a protocol data type and protocol data.
// It returns an error if the protocol data cannot be marshalled, and panic if can not set the protocol data.
func MarshalAndSetProtocolData(ctx sdk.Context, k *Keeper, datatype types.ProtocolDataType, pd types.ProtocolDataI) error {
	pdString, err := json.Marshal(pd)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling protocol data", "error", err)
		return err
	}
	storedProtocolData := types.NewProtocolData(datatype.String(), pdString)
	k.SetProtocolData(ctx, pd.GenerateKey(), storedProtocolData)
	return nil
}
