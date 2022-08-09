package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) NewClaim(ctx sdk.Context, address string, zone string, amount int64) *types.Claim {
	return &types.Claim{UserAddress: address, Zone: zone, HeldAmount: amount}
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
	store.Set(GetClaimKeyForClaim(claim), bz)
}

// DeleteClaim deletes claim
func (k Keeper) DeleteClaim(ctx sdk.Context, claim *types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixClaim)
	store.Delete(GetClaimKeyForClaim(claim))
}

// IterateQueries iterates through claims
func (k Keeper) IterateClaims(ctx sdk.Context, fn func(index int64, data types.Claim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixClaim)
	iterator := sdk.KVStorePrefixIterator(store, nil)
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
	out := make([]*types.Claim, 0)
	k.IterateClaims(ctx, func(_ int64, claim types.Claim) (stop bool) {
		out = append(out, &claim)
		return false
	})
	return out
}

// ClaimKey returns the key for storing a given claim.
func GetClaimKeyForClaim(claim *types.Claim) []byte {
	return GetClaimKey(claim.Zone, claim.UserAddress)
}

// ClaimKey returns the key for a given zone and user.
func GetClaimKey(zone string, userAddress string) []byte {
	return []byte(fmt.Sprintf("%s/%s", zone, userAddress))
}
