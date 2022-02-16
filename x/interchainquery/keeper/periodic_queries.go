package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

// GetPeriodicQuery returns query
func (k Keeper) GetPeriodicQuery(ctx sdk.Context, id string) (types.PeriodicQuery, bool) {
	query := types.PeriodicQuery{}
	ctx.Logger().Error(fmt.Sprintf("Looking for query: %s", id))

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPeriodicQuery)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		fmt.Printf("BAILING")
		return query, false
	}
	k.cdc.MustUnmarshal(bz, &query)
	ctx.Logger().Error(fmt.Sprintf("Found query: %v", query))
	return query, true
}

// SetPeriodicQuery set query info
func (k Keeper) SetPeriodicQuery(ctx sdk.Context, query types.PeriodicQuery) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPeriodicQuery)
	bz := k.cdc.MustMarshal(&query)
	ctx.Logger().Error(fmt.Sprintf("Created/updated query: %v", query))
	store.Set([]byte(query.ChainId), bz)
}

// DeletePeriodicQuery delete query info
func (k Keeper) DeletePeriodicQuery(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPeriodicQuery)
	ctx.Logger().Error(fmt.Sprintf("Removing query: %s", id))
	store.Delete([]byte(id))
}

// IteratePeriodicQueries iterate through querys
func (k Keeper) IteratePeriodicQueries(ctx sdk.Context, fn func(index int64, queryInfo types.PeriodicQuery) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPeriodicQuery)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		query := types.PeriodicQuery{}
		k.cdc.MustUnmarshal(iterator.Value(), &query)
		ctx.Logger().Error(fmt.Sprintf("looking up %s", iterator.Value()))

		stop := fn(i, query)

		if stop {
			break
		}
		i++
	}
}

// AllPeriodicQueries returns every queryInfo in the store
func (k Keeper) AllPeriodicQueries(ctx sdk.Context) []types.PeriodicQuery {
	querys := []types.PeriodicQuery{}
	k.IteratePeriodicQueries(ctx, func(_ int64, queryInfo types.PeriodicQuery) (stop bool) {
		querys = append(querys, queryInfo)
		return false
	})
	return querys
}
