package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// GetRemoteAddressMap retrieves a remote address using a local address.
func (k *Keeper) GetRemoteAddressMap(ctx sdk.Context, localAddress []byte, chainID string) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRemoteAddressKey(localAddress, chainID)
	value := store.Get(key)

	return value, value != nil
}

// IterateUserMappedAccounts iterates over all the user mapped accounts.
func (k Keeper) IterateUserMappedAccounts(ctx sdk.Context, localAddress []byte, fn func(index int64, chainID string, remoteAddressBytes []byte) (stop bool)) {
	// noop
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetRemoteAddressPrefix(localAddress))
	defer iterator.Close()

	i := int64(0)
	remoteAddrPrefixLen := len(types.GetRemoteAddressPrefix(localAddress))
	for ; iterator.Valid(); iterator.Next() {
		value := iterator.Value()
		key := iterator.Key()
		chainIDBytes := key[remoteAddrPrefixLen:]
		stop := fn(i, string(chainIDBytes), value)
		if stop {
			break
		}
		i++
	}
}

// SetRemoteAddressMap sets a remote address using a local address as a map.
func (k *Keeper) SetRemoteAddressMap(ctx sdk.Context, localAddress, remoteAddress []byte, chainID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRemoteAddressKey(localAddress, chainID)
	store.Set(key, remoteAddress)
}

// GetLocalAddressMap retrieves a local address using a remote address.
func (k *Keeper) GetLocalAddressMap(ctx sdk.Context, remoteAddress []byte, chainID string) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLocalAddressKey(remoteAddress, chainID)
	value := store.Get(key)

	return value, value != nil
}

// SetLocalAddressMap sets a local address using a remote address as a map.
func (k *Keeper) SetLocalAddressMap(ctx sdk.Context, localAddress, remoteAddress []byte, chainID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLocalAddressKey(remoteAddress, chainID)
	store.Set(key, localAddress)
}

// SetAddressMapPair sets forward and reverse maps for localAddress => remoteAddress and remoteAddress => localAddress.
func (k *Keeper) SetAddressMapPair(ctx sdk.Context, localAddress, remoteAddress []byte, chainID string) {
	k.SetLocalAddressMap(ctx, localAddress, remoteAddress, chainID)
	k.SetRemoteAddressMap(ctx, localAddress, remoteAddress, chainID)
}
