package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) NewClaim(ctx sdk.Context, address string, chainID string, amount uint64) *types.Claim {
	return &types.Claim{UserAddress: address, ChainId: chainID, Amount: amount}
}

// GetClaim returns claim
func (k Keeper) GetClaim(ctx sdk.Context, key []byte) (types.Claim, bool) {
	data := types.Claim{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixClaim)
	bz := store.Get(key)
	if len(bz) == 0 {
		return data, false
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

// SetClaim sets claim
func (k Keeper) SetClaim(ctx sdk.Context, claim *types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixClaim)
	bz := k.cdc.MustMarshal(claim)
	store.Set(types.GetKeyClaim(claim.ChainId, claim.UserAddress), bz)
}

// DeleteClaim deletes claim
func (k Keeper) DeleteClaim(ctx sdk.Context, claim *types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixClaim)
	store.Delete(types.GetKeyClaim(claim.ChainId, claim.UserAddress))
}

// IterateClaims iterates through zone claims.
func (k Keeper) IterateClaims(ctx sdk.Context, chainID string, fn func(index int64, data types.Claim) (stop bool)) {
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

// ClearClaims deletes all the claims of the given zone.
func (k Keeper) ClearClaims(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixClaim(chainID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}
