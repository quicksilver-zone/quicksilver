package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

func (k Keeper) NewClaim(ctx sdk.Context, address string, chainID string, module types.ClaimType, srcChainID string, amount uint64) types.Claim {
	return types.Claim{UserAddress: address, ChainId: chainID, Module: module, SourceChainId: srcChainID, Amount: amount}
}

// GetClaim returns claim
func (k Keeper) GetClaim(ctx sdk.Context, chainID string, address string, module types.ClaimType, srcChainID string) (types.Claim, bool) {
	data := types.Claim{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	key := types.GetKeyClaim(chainID, address, module, srcChainID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return data, false
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

// GetLastEpochClaim returns claim from last epoch
func (k Keeper) GetLastEpochClaim(ctx sdk.Context, chainID string, address string, module types.ClaimType, srcChainID string) (types.Claim, bool) {
	data := types.Claim{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	key := types.GetKeyLastEpochClaim(chainID, address, module, srcChainID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return data, false
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

// SetClaim sets claim
func (k Keeper) SetClaim(ctx sdk.Context, claim *types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := k.cdc.MustMarshal(claim)
	store.Set(types.GetKeyClaim(claim.ChainId, claim.UserAddress, claim.Module, claim.SourceChainId), bz)
}

// SetLastEpochClaim sets claim for last epoch
func (k Keeper) SetLastEpochClaim(ctx sdk.Context, claim *types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := k.cdc.MustMarshal(claim)
	store.Set(types.GetKeyLastEpochClaim(claim.ChainId, claim.UserAddress, claim.Module, claim.SourceChainId), bz)
}

// DeleteClaim deletes claim
func (k Keeper) DeleteClaim(ctx sdk.Context, claim *types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	store.Delete(types.GetKeyClaim(claim.ChainId, claim.UserAddress, claim.Module, claim.SourceChainId))
}

// DeleteLastEpochClaim deletes claim for last epoch
func (k Keeper) DeleteLastEpochClaim(ctx sdk.Context, claim *types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	store.Delete(types.GetKeyLastEpochClaim(claim.ChainId, claim.UserAddress, claim.Module, claim.SourceChainId))
}

// IterateClaims iterates through zone claims.
func (k Keeper) IterateClaims(ctx sdk.Context, chainID string, fn func(index int64, data types.Claim) (stop bool)) {
	// noop
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixClaim(chainID))
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.Claim{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, data)
		if stop {
			break
		}
		i++
	}
}

// IterateUserClaims iterates through zone claims for a given address.
func (k Keeper) IterateUserClaims(ctx sdk.Context, chainID string, address string, fn func(index int64, data types.Claim) (stop bool)) {
	// noop
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixUserClaim(chainID, address))
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.Claim{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, data)
		if stop {
			break
		}
		i++
	}
}

// IterateLastEpochClaims iterates through zone claims from last epoch.
func (k Keeper) IterateLastEpochClaims(ctx sdk.Context, chainID string, fn func(index int64, data types.Claim) (stop bool)) {
	// noop
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixLastEpochClaim(chainID))
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.Claim{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, data)
		if stop {
			break
		}
		i++
	}
}

// IterateLastEpochUserClaims iterates through zone claims from last epoch for a given user.
func (k Keeper) IterateLastEpochUserClaims(ctx sdk.Context, chainID string, address string, fn func(index int64, data types.Claim) (stop bool)) {
	// noop
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixLastEpochUserClaim(chainID, address))
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.Claim{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, data)
		if stop {
			break
		}
		i++
	}
}

// IterateLastEpochUserClaims iterates through zone claims from last epoch for a given user.
func (k Keeper) IterateAllLastEpochClaims(ctx sdk.Context, fn func(index int64, key []byte, data types.Claim) (stop bool)) {
	// noop
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixLastEpochClaim)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.Claim{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, iterator.Key(), data)
		if stop {
			break
		}
		i++
	}
}

// IterateAllClaims iterates through all claims.
func (k Keeper) IterateAllClaims(ctx sdk.Context, fn func(index int64, key []byte, data types.Claim) (stop bool)) {
	// noop
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixClaim)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.Claim{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, iterator.Key(), data)
		if stop {
			break
		}
		i++
	}
}

// AllClaims returns a slice containing all claims from the store.
func (k Keeper) AllClaims(ctx sdk.Context) []*types.Claim {
	claims := []*types.Claim{}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixClaim)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		claim := types.Claim{}
		k.cdc.MustUnmarshal(iterator.Value(), &claim)

		claims = append(claims, &claim)
	}

	return claims
}

func (k Keeper) AllZoneClaims(ctx sdk.Context, chainID string) []*types.Claim {
	claims := []*types.Claim{}
	k.IterateClaims(ctx, chainID, func(_ int64, claim types.Claim) (stop bool) {
		claims = append(claims, &claim)
		return false
	})
	return claims
}

func (k Keeper) AllZoneUserClaims(ctx sdk.Context, chainID string, address string) []*types.Claim {
	claims := []*types.Claim{}
	k.IterateUserClaims(ctx, chainID, address, func(_ int64, claim types.Claim) (stop bool) {
		claims = append(claims, &claim)
		return false
	})
	return claims
}

func (k Keeper) AllZoneLastEpochClaims(ctx sdk.Context, chainID string) []*types.Claim {
	claims := []*types.Claim{}
	k.IterateLastEpochClaims(ctx, chainID, func(_ int64, claim types.Claim) (stop bool) {
		claims = append(claims, &claim)
		return false
	})
	return claims
}

func (k Keeper) AllZoneLastEpochUserClaims(ctx sdk.Context, chainID string, address string) []*types.Claim {
	claims := []*types.Claim{}
	k.IterateLastEpochUserClaims(ctx, chainID, address, func(_ int64, claim types.Claim) (stop bool) {
		claims = append(claims, &claim)
		return false
	})
	return claims
}

// ClearClaims deletes all the current epoch claims of the given zone.
func (k Keeper) ClearClaims(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixClaim(chainID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

// ClearLastEpochClaims deletes all the last epoch claims of the given zone.
func (k Keeper) ClearLastEpochClaims(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixLastEpochClaim(chainID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

// ArchiveAndGarbageCollectClaims deletes all the last epoch claims and moves the current epoch claims to the last epoch store.
func (k Keeper) ArchiveAndGarbageCollectClaims(ctx sdk.Context, chainID string) {
	k.ClearLastEpochClaims(ctx, chainID)

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixClaim(chainID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
		newKey := types.KeyPrefixLastEpochClaim
		newKey = append(newKey, key[1:]...) // update prefix from KeyPrefixClaim to KeyPrefixLastEpochClaim
		store.Set(newKey, iterator.Value())
	}
}
