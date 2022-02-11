package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

// GetSingleQuery returns query
func (k Keeper) GetSingleQuery(ctx sdk.Context, id string) (types.SingleQuery, bool) {
	query := types.SingleQuery{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixData)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return query, false
	}

	k.cdc.MustUnmarshal(bz, &query)
	return query, true
}

func (k Keeper) NewPeriodicQuery(ctx sdk.Context, connection_id string, chain_id string, query_type string, query_params map[string]string, period sdk.Int) *types.PeriodicQuery {
	generated_id := fmt.Sprintf("%x", crypto.Sha256([]byte("p"+connection_id+chain_id+query_type+string(ctx.BlockHeight()))))
	return &types.PeriodicQuery{Id: generated_id, ConnectionId: connection_id, ChainId: chain_id, QueryType: query_type, QueryParameters: query_params, Period: period, LastHeight: sdk.NewInt(ctx.BlockHeight())}
}

func (k Keeper) NewSingleQuery(ctx sdk.Context, connection_id string, chain_id string, query_type string, query_params map[string]string) *types.SingleQuery {
	generated_id := fmt.Sprintf("%x", crypto.Sha256([]byte("s"+connection_id+chain_id+query_type+string(ctx.BlockHeight()))))
	return &types.SingleQuery{Id: generated_id, ConnectionId: connection_id, ChainId: chain_id, QueryType: query_type, QueryParameters: query_params}
}

// SetSingleQuery set query info
func (k Keeper) SetSingleQuery(ctx sdk.Context, query types.SingleQuery) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSingleQuery)
	bz := k.cdc.MustMarshal(&query)
	ctx.Logger().Error(fmt.Sprintf("%v", query))
	store.Set([]byte(query.ChainId), bz)
}

// DeleteSingleQuery delete query info
func (k Keeper) DeleteSingleQuery(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSingleQuery)
	ctx.Logger().Error(fmt.Sprintf("Removing query: %s", id))
	store.Delete([]byte(id))
}

// IterateSingleQueries iterate through querys
func (k Keeper) IterateSingleQueries(ctx sdk.Context, fn func(index int64, queryInfo types.SingleQuery) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSingleQuery)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		query := types.SingleQuery{}
		k.cdc.MustUnmarshal(iterator.Value(), &query)

		stop := fn(i, query)

		if stop {
			break
		}
		i++
	}
}

// AllSingleQueries returns every queryInfo in the store
func (k Keeper) AllSingleQueries(ctx sdk.Context) []types.SingleQuery {
	querys := []types.SingleQuery{}
	k.IterateSingleQueries(ctx, func(_ int64, queryInfo types.SingleQuery) (stop bool) {
		querys = append(querys, queryInfo)
		return false
	})
	return querys
}
