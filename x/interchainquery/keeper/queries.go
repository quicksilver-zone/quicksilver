package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

func GenerateQueryHash(connection_id string, chain_id string, query_type string, query_params map[string]string) string {
	param_bytes, _ := json.Marshal(query_params)
	return fmt.Sprintf("%x", crypto.Sha256(append([]byte(connection_id+chain_id+query_type), param_bytes...)))
}

// ----------------------------------------------------------------
func (k Keeper) NewSingleQuery(ctx sdk.Context, connection_id string, chain_id string, query_type string, query_params map[string]string) *types.SingleQuery {
	return &types.SingleQuery{Id: GenerateQueryHash(connection_id, chain_id, query_type, query_params), ConnectionId: connection_id, ChainId: chain_id, QueryType: query_type, QueryParameters: query_params, EmitHeight: sdk.ZeroInt()}
}

// GetSingleQuery returns query
func (k Keeper) GetSingleQuery(ctx sdk.Context, id string) (types.SingleQuery, bool) {
	query := types.SingleQuery{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSingleQuery)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return query, false
	}

	k.cdc.MustUnmarshal(bz, &query)
	return query, true
}

// SetSingleQuery set query info
func (k Keeper) SetSingleQuery(ctx sdk.Context, query types.SingleQuery) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSingleQuery)
	bz := k.cdc.MustMarshal(&query)
	store.Set([]byte(query.Id), bz)
}

// DeleteSingleQuery delete query info
func (k Keeper) DeleteSingleQuery(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSingleQuery)
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

// ----------------------------------------------------------------

func (k Keeper) NewPeriodicQuery(ctx sdk.Context, connection_id string, chain_id string, query_type string, query_params map[string]string, period sdk.Int) *types.PeriodicQuery {
	return &types.PeriodicQuery{Id: GenerateQueryHash(connection_id, chain_id, query_type, query_params), ConnectionId: connection_id, ChainId: chain_id, QueryType: query_type, QueryParameters: query_params, Period: period, LastHeight: sdk.NewInt(ctx.BlockHeight())}
}

// GetPeriodicQuery returns query
func (k Keeper) GetPeriodicQuery(ctx sdk.Context, id string) (types.PeriodicQuery, bool) {
	query := types.PeriodicQuery{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPeriodicQuery)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return query, false
	}
	k.cdc.MustUnmarshal(bz, &query)
	return query, true
}

// SetPeriodicQuery set query info
func (k Keeper) SetPeriodicQuery(ctx sdk.Context, query types.PeriodicQuery) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPeriodicQuery)
	bz := k.cdc.MustMarshal(&query)
	k.Logger(ctx).Info("Created/updated query", "ID", query.Id)
	store.Set([]byte(query.Id), bz)
}

// DeletePeriodicQuery delete query info
func (k Keeper) DeletePeriodicQuery(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPeriodicQuery)
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
