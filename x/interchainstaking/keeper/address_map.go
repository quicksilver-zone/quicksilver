package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetRemoteAddress retrieves a remote address using a local address.
func (k *Keeper) GetRemoteAddress(ctx sdk.Context, localAddress []byte, chainID string) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRemoteAddressKey(localAddress, chainID)
	value := store.Get(key)

	return value, value != nil
}

// SetRemoteAddress sets a remote address using a local address as a map.
func (k *Keeper) SetRemoteAddress(ctx sdk.Context, localAddress, remoteAddress []byte, chainID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRemoteAddressKey(localAddress, chainID)
	store.Set(key, remoteAddress)
}

// GetLocalAddress retrieves a local address using a remote address.
func (k *Keeper) GetLocalAddress(ctx sdk.Context, remoteAddress []byte, chainID string) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLocalAddressKey(remoteAddress, chainID)
	value := store.Get(key)

	return value, value != nil
}

// SetLocalAddress sets a local address using a remote address as a map.
func (k *Keeper) SetLocalAddress(ctx sdk.Context, localAddress, remoteAddress []byte, chainID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLocalAddressKey(remoteAddress, chainID)
	store.Set(key, localAddress)
}

// SetAddressMapPair sets forward and reverse maps for localAddress => remoteAddress and remoteAddress => localAddress.
func (k *Keeper) SetAddressMapPair(ctx sdk.Context, localAddress, remoteAddress []byte, chainID string) {
	k.SetLocalAddress(ctx, localAddress, remoteAddress, chainID)
	k.SetRemoteAddress(ctx, localAddress, remoteAddress, chainID)
}
